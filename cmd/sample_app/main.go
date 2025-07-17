package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	rand "math/rand/v2"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status" 
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
  	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
  	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"

	
)

// TODO: 妥当なサービス名
var serviceName = semconv.ServiceNameKey.String("observability-with-otel-sample")

func main() {
	// TODO : 妥当なコンテキストの扱いとシャットダウン処理
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// 環境変数を読み込む
	envVars, err := LoadEnvironmentVariables()
	if err != nil {
		log.Fatalf("failed to LoadEnvironmentVariables(): %v", err)
	}

	// Otel SDK のセットアップ
	otelShutdownFuncs, err := setupOtelSDK(ctx, envVars.OtelCollectorHost, envVars.UseOtelStdoutExporter)
	if err != nil {
		log.Fatalf("failed to setupOtelSDK(): %v", err)
	}

	// 処理中のテレメトリーデータのケア
	defer func() {
		err = errors.Join(err, otelShutdownFuncs(context.Background()))
		log.Fatalln(err)
	}()

	// インフラ層のインスタンスを生成
	infrastructures, err := createInfrastructuresWithGORMSpanner(
		envVars.GCPProjectID,
		envVars.SpannerInstanceID,
		envVars.SpannerDatabaseID,
	)
	if err != nil {
		log.Fatalf("failed to createInfrastructuresWithGORMSpanner(): %v", err)
	}

	// アプリケーション層のインスタンスを生成
	applications, err := createApplications(infrastructures)
	if err != nil {
		log.Fatalf("failed to createApplications(): %v", err)
	}

	// プレゼンテーション層のインスタンスを生成
	presentations, err := createPresentations(applications)
	if err != nil {
		log.Fatalf("failed to createPresentations(): %v", err)
	}

	// gRPCサーバーの起動
	if err := startGrpcServer(envVars.Port, presentations); err != nil {
		log.Fatalf("failed to startGrpcServer(): %v", err)
	}
}

// startGrpcServer gRPCサーバーの起動処理
func startGrpcServer(port int, presentations *presentations) error {
	// 1. 指定したプロトコル・ポートのListenerを作成
	// (net.Listenerが返される。Listenerとはポートに対して聞き耳を立てる人である)。
	listenProtocol := "tcp"
	portListener, err := net.Listen(listenProtocol, fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("プロトコル:%v, ポート:%v のListener作成に失敗しました: %w", listenProtocol, port, err)
	}

	// 2. gRPCサーバのインスタンスを生成
	// (grpc.Serverインスタンスのポインタが返ってくる)
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(randProcStatusInterceptor))

	// 3. gRPCサーバにサービスを登録
	presentations.registerProtocServices(grpcServer)

	// 4. gRPCサーバのServer Reflectionを有効にする
	// (「grpcurl」コマンドで、gRPCサーバに登録したサービスのRPCメソッドをシリアライズなしで実行可能になる)
	reflection.Register(grpcServer)

	// 5. gRPCサーバーを起動(指定したプロトコル・ポートのListenも開始)
	if err := grpcServer.Serve(portListener); err != nil {
		return fmt.Errorf("failed to grpcServer.Serve(): %w", err)
	}
	return nil
}

// Otel SDK 周りの初期化
// useStdExporter: true の場合テレメトリーデータを標準出力へ出力、 false の場合コレクターへ出力
func setupOtelSDK(ctx context.Context, collectorHost string, useStdExporter bool) (shutdown func(context.Context) error, err error) {
	var shutdownFuncs []func(context.Context) error
	shutdown = func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	// プロパゲーターのセットアップをコール
	prop := newPropagator()
	otel.SetTextMapPropagator(prop)

	// Otel コレクターへのコネクションのセットアップ
	otelConn, err := newOtelCollectorConn(collectorHost)
	if err != nil {
		handleErr(err)
		return
	}

	// リソースのセットアップ
	res, err := newResource(ctx)
	if err != nil {
		handleErr(err)
		return
	}

	// トレーサープロバイダのセットアップ
	shutdownTracerProvider, err := initTracerProvider(ctx, res, otelConn, useStdExporter)
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, shutdownTracerProvider)

	// ロガープロバイダのセットアップ
	shutDownloggerProvider, err := initLoggerProvider(ctx, res, otelConn, useStdExporter)
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, shutDownloggerProvider)
	
	// メータープロバイダのセットアップ
	shutdownMeterProvider, err := initMeterProvider(ctx, res, otelConn, useStdExporter)
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, shutdownMeterProvider)

	return
}

// プロパゲーターのセットアップ
func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

// Otel コレクターへのコネクションのセットアップ (トレーサープロバイダ、ロガープロバイダ、メータープロバイダ共通で使う)
func newOtelCollectorConn(collectorHost string) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(collectorHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to grpc.NewClient(): %w", err)
	}
	return conn, nil
}

// リソースのセットアップ
func newResource(ctx context.Context) (*resource.Resource, error) {
	res, err := resource.New(ctx,
		resource.WithAttributes(serviceName),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to resource.WithAttributes(): %w", err)
	}
	return res, nil
}

// トレーサープロバイダのセットアップ
func initTracerProvider(ctx context.Context, res *resource.Resource, conn *grpc.ClientConn, useStdoutExporter bool) (func(context.Context) error, error) {
	var traceExporter sdktrace.SpanExporter
	var err error
	if useStdoutExporter {
		traceExporter, err = stdouttrace.New(stdouttrace.WithPrettyPrint())
		if err != nil {
			return nil, fmt.Errorf("failed to stdouttrace.New(): %w", err)
		}
	} else {
		traceExporter, err = otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
		if err != nil {
			return nil, fmt.Errorf("failed to otlptracegrpc.New(): %w", err)
		}
	}

	bsp := sdktrace.NewBatchSpanProcessor(traceExporter, sdktrace.WithBatchTimeout(time.Minute)) // TODO: デフォルトの 5秒から 1分にしてみた
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)

	return tracerProvider.Shutdown, nil
}

// ロガープロバイダのセットアップ
func initLoggerProvider(ctx context.Context, res *resource.Resource, conn *grpc.ClientConn, useStdoutExporter bool) (func(context.Context) error, error) {
	var logExporter sdklog.Exporter
	var err error
	if useStdoutExporter {
		logExporter, err = stdoutlog.New()
		if err != nil {
			return nil, fmt.Errorf("failed to stdoutlog.New(): %w", err)
		}
	} else {
		logExporter, err = otlploggrpc.New(ctx, otlploggrpc.WithGRPCConn(conn))
		if err != nil {
			return nil, fmt.Errorf("failed to otlploggrpc.New(): %w", err)
		}
	}

	loggerProvider := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(sdklog.NewBatchProcessor(logExporter)),
		sdklog.WithResource(res),
	)
		
	global.SetLoggerProvider(loggerProvider)

	return loggerProvider.Shutdown, nil
}

// メータープロバイダのセットアップ
func initMeterProvider(ctx context.Context, res *resource.Resource, conn *grpc.ClientConn, useStdoutExporter bool) (func(context.Context) error, error) {
	var metricExporter sdkmetric.Exporter
	var err error
	if useStdoutExporter {
		metricExporter, err = stdoutmetric.New()
		if err != nil {
			return nil, fmt.Errorf("failed to stdoutmetric.New(): %w", err)
		}
	} else {
		metricExporter, err = otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(conn))
		if err != nil {
			return nil, fmt.Errorf("failed to otlpmetricgrpc.New(): %w", err)
		}
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter)),
		sdkmetric.WithResource(res),
	)
	otel.SetMeterProvider(meterProvider)

	return meterProvider.Shutdown, nil
}

// 処理ステータス
type procState int
const (
	// エラー
	procStateError procState = iota
	// 高レイテンシー
	procStateHighLatency
	// 正常
	procStateNormal
	// 未定義
	procStateNone
)

// 等確率で処理ステータス(エラー/高レイテンシー/正常)
func randProcStatus() procState {
	switch rand.IntN(3) {
	case 0:
		return procStateError
	case 1:
		return procStateHighLatency
	case 2:
		return procStateNormal
	default:
		return procStateNone
	}
}

// リクエストされたメソッド実行前に処理ステータス(エラー/高レイテンシー/正常)を確率で決定し、処理ステータスに応じて以降の処理を制御するインターセプター
// 1. スパンをこの関数単位で作成する。
// 2. 処理ステータスをランダムに決定する。
// 3. 処理ステータスに応じて以下の処理を実行する。
//   - 処理ステータスがエラー: ログレベル Error で処理ステータスがエラーだったことをログ出力、メトリクスとしてカウントして Internal エラーでレスポンスを返す
//   - 処理ステータスが高レイテンシー: 一定時間待った後、ログレベル Info で処理ステータスが高レイテンシーだったことをログ出力、メトリクスとしてカウントして処理を続行する
//   - 処理ステータスが正常: ログレベル Info で処理ステータスが正常だったことをログ出力、メトリクスとしてカウントして処理を続行する
func randProcStatusInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (res interface{}, err error) {
	scopeName := "randProcStatusInterceptor" // TODO 計装スコープ名のセマンティック規約意識
	// テスト的に各テレメトリーデータに共通で仕込むプロパティ
	commonAttrs := []attribute.KeyValue{
		attribute.String("attrA", "chocolate"),
		attribute.String("attrB", "raspberry"),
		attribute.String("attrC", "vanilla"),
	}

	logger := otelslog.NewLogger(scopeName)	
	tracer := otel.Tracer(scopeName)
	meter := otel.Meter(scopeName)
	
	ctx, span := tracer.Start(
		ctx,
		"randProcStatusInterceptor", // TODO スパン名のセマンティック規約意識
		trace.WithAttributes(commonAttrs...))
	defer span.End()

	errCount, err := meter.Int64Counter("procStateError", metric.WithDescription("procStateError count"))
	if err != nil {
		err = status.Error(codes.Internal, "failed to meter.Int64Counter") // TODO エラーハンドリング
		return
	}

	highLatencyCount, err := meter.Int64Counter("procStateHighLatency", metric.WithDescription("procStateHighLatency count"))
	if err != nil {
		err = status.Error(codes.Internal, "failed to meter.Int64Counter") // TODO エラーハンドリング
		return
	}

	normalCount, err := meter.Int64Counter("procStateNormal", metric.WithDescription("procStateNormal count"))
	if err != nil {
		err = status.Error(codes.Internal, "failed to meter.Int64Counter") // TODO エラーハンドリング
		return
	}

	switch randProcStatus() {
	case procStateError:
		logger.ErrorContext(ctx, "procState is error") // otel 的には log.SeverityError 扱いになる
		errCount.Add(ctx, 1, metric.WithAttributes(commonAttrs...))
		err = status.Error(codes.Internal, "error occurred randomly")
	case procStateHighLatency:
		<-time.After(2*time.Second)
		logger.InfoContext(ctx, "procState is highLatency. 2 second waited")
		highLatencyCount.Add(ctx, 1, metric.WithAttributes(commonAttrs...))
		res, err = handler(ctx, req)		
	default:
		logger.InfoContext(ctx, "procState is normal")
		normalCount.Add(ctx, 1, metric.WithAttributes(commonAttrs...))
		res, err = handler(ctx, req)
	}
	return
}