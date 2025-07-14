package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.26.0"
)

// TODO
var serviceName = semconv.ServiceNameKey.String("observability-with-otel-sample")

func main() {
	// TODO 
	ctx := context.Background()

	// 環境変数を読み込む
	envVars, err := LoadEnvironmentVariables()
	if err != nil {
		log.Fatalf("failed to LoadEnvironmentVariables(): %v", err)
	}

	// プロパゲーターのセットアップをコール
	prop := newPropagator()
	otel.SetTextMapPropagator(prop)

	// Otel コレクターへのコネクションのセットアップ
	otelConn, err := newOtelCollectorConn(envVars.OtelCollectorHost)
	if err != nil {
		log.Fatalf("failed to newOtelCollectorConn(): %v", err)
	}

	// リソースのセットアップ
	res, err := newResource(ctx)
	if err != nil {
		log.Fatalf("failed to newResource(): %v", err) // TODO 終了しなくていいんだっけ？
	}

	// トレーサープロバイダのセットアップ
	shutdownTracerProvider, err := initTracerProvider(ctx, res, otelConn)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdownTracerProvider(ctx); err != nil {
			log.Fatalf("failed to shutdown TracerProvider: %s", err)
		}
	}()

	// ロガープロバイダのセットアップ
	shutDownloggerProvider, err := initLoggerProvider(ctx, res, otelConn)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutDownloggerProvider(ctx); err != nil {
			log.Fatalf("failed to shutdown LoggerProvider: %s", err)
		}
	}()

	// メータープロバイダのセットアップ
	shutdownMeterProvider, err := initMeterProvider(ctx, res, otelConn)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdownMeterProvider(ctx); err != nil {
			log.Fatalf("failed to shutdown MeterProvider: %s", err)
		}
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
	grpcServer := grpc.NewServer()

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
func initTracerProvider(ctx context.Context, res *resource.Resource, conn *grpc.ClientConn) (func(context.Context) error, error) {
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)

	return tracerProvider.Shutdown, nil
}

// ロガープロバイダのセットアップ
func initLoggerProvider(ctx context.Context, res *resource.Resource, conn *grpc.ClientConn) (func(context.Context) error, error) {
	logExporter, err := otlploggrpc.New(ctx, otlploggrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create  exporter: %w", err)
	}

	loggerProvider := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(sdklog.NewBatchProcessor(logExporter)),
		sdklog.WithResource(res),
	)
		
	global.SetLoggerProvider(loggerProvider)

	return loggerProvider.Shutdown, nil
}


// メータープロバイダのセットアップ
func initMeterProvider(ctx context.Context, res *resource.Resource, conn *grpc.ClientConn) (func(context.Context) error, error) {
	metricExporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics exporter: %w", err)
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter)),
		sdkmetric.WithResource(res),
	)
	otel.SetMeterProvider(meterProvider)

	return meterProvider.Shutdown, nil
}
