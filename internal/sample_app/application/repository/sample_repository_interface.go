//go:generate  mockgen -source=sample_repository_interface.go -destination=sample_repository_mock.go -package=repository -mock_names=ISampleRepository=MockSampleRepository

package repository

import (
	"context"
	usecase "modern-dev-env-app-sample/internal/sample_app/application/repository/transaction"
	entity "modern-dev-env-app-sample/internal/sample_app/domain/entity/sample"
	"modern-dev-env-app-sample/internal/sample_app/domain/value"
)

// ISampleRepository Sample集約リポジトリのインターフェース
type ISampleRepository interface {
	// Save 1件のSampleエンティティを保存
	Save(ctx context.Context, sample *entity.Sample, iTx usecase.ITransaction) error
	// FindByIDs 指定したID群でSampleエンティティ群を取得
	FindByIDs(ctx context.Context, ids value.SampleIDs, iTx usecase.ITransaction) ([]*entity.Sample, error)
	// FindAll 全てのSampleエンティティ群を取得
	FindAll(ctx context.Context, iTx usecase.ITransaction) ([]*entity.Sample, error)
	// Delete 1件のSampleエンティティを論理削除
	Delete(ctx context.Context, sample *entity.Sample, iTx usecase.ITransaction) error
}
