package main

import (
	"fmt"

	"modern-dev-env-app-sample/internal/sample_app/presentation/pb"
	"modern-dev-env-app-sample/internal/sample_app/presentation/sample"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

)

// presentations 全プレゼンテーション層のインスタンスをまとめた構造体
type presentations struct {
	iHealthServer healthpb.HealthServer
	iSampleServiceServer pb.SampleServiceServer
}

// newPresentations コンストラクタ
func newPresentations(
	healthServer healthpb.HealthServer,
	sampleServiceServer pb.SampleServiceServer,
) *presentations {
	return &presentations{
		iHealthServer: healthServer,
		iSampleServiceServer: sampleServiceServer,
	}
}

// createPresentations 全プレゼンテーション層のインスタンスのファクトリ
func createPresentations(
	applications *applications,
) (*presentations, error) {
	healthServer := health.NewServer()
	healthServer.SetServingStatus("SampleService", healthpb.HealthCheckResponse_SERVING)
	
	sampleServiceServer, err := sample.NewSampleServiceServer(
		applications.iListSamplesUseCase,
		applications.iCreateSampleUseCase,
		applications.iUpdateSampleUseCase,
		applications.iDeleteSampleUseCase,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to NewSampleServiceServer(): %w", err)
	}
	return newPresentations(
		healthServer, 
		sampleServiceServer,
	), nil
}

// registerProtocServices protoc都合のRPCサービス構造体を登録
// Register<サービス名>Server(grpc.Serverのポインタ, <サービス名>Serverインターフェース)関数は、
// protocコマンド実行で生成された「<元になった.proroファイル名>_grpc.pb.go」内に自動で定義されている
// ここで登録されたRPCサービスについてのみ、gRPC通信が可能になる
func (p *presentations) registerProtocServices(grpcServer *grpc.Server) {
	healthpb.RegisterHealthServer(grpcServer, p.iHealthServer)
	pb.RegisterSampleServiceServer(grpcServer, p.iSampleServiceServer)
}
