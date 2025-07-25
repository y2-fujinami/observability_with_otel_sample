//go:generate mockgen -source=create_sample_usecase_interface.go -destination=create_sample_usecase_mock.go -package=sample -mock_names=ICreateSampleUseCase=MockCreateSampleUseCase

package sample

import (
	"context"
	application "modern-dev-env-app-sample/internal/sample_app/application/request/sample"
	application2 "modern-dev-env-app-sample/internal/sample_app/application/response/sample"
)

// ICreateSampleUseCase ユースケース: サンプルデータを追加 のインターフェース
// プレゼンテーション層は、ユースケース層で定義したこのインターフェースを介して処理を依頼する
type ICreateSampleUseCase interface {
	Run(ctx context.Context, req *application.CreateSampleRequest) (*application2.CreateSampleResponse, error)
}
