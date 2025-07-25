package sample

import (
	"errors"
	"fmt"

	usecase "modern-dev-env-app-sample/internal/sample_app/application/usecase/sample"
	"modern-dev-env-app-sample/internal/sample_app/presentation/pb"

	"go.opentelemetry.io/contrib/bridges/otelslog"
)

// このパッケージ内でのみ使う logger 
var logger = otelslog.NewLogger("modern-dev-env-app-sample/internal/sample_app/presentation/sample")	

// SampleServiceServer protocで自動生成されたSampleServiceServerのインターフェースをみたす構造体
type SampleServiceServer struct {
	iListSamplesUseCase  usecase.IListSamplesUseCase
	iCreateSampleUseCase usecase.ICreateSampleUseCase
	iUpdateSampleUseCase usecase.IUpdateSampleUseCase
	iDeleteSampleUseCase usecase.IDeleteSampleUseCase
	pb.UnimplementedSampleServiceServer
}

// NewSampleServiceServer SampleServiceServerのコンストラクタ
func NewSampleServiceServer(
	iListSamplesUseCase usecase.IListSamplesUseCase,
	iCreateSampleUseCase usecase.ICreateSampleUseCase,
	iUpdateSampleUseCase usecase.IUpdateSampleUseCase,
	iDeleteSampleUseCase usecase.IDeleteSampleUseCase,
) (*SampleServiceServer, error) {
	sampleServiceServer := &SampleServiceServer{
		iListSamplesUseCase:  iListSamplesUseCase,
		iCreateSampleUseCase: iCreateSampleUseCase,
		iUpdateSampleUseCase: iUpdateSampleUseCase,
		iDeleteSampleUseCase: iDeleteSampleUseCase,
	}
	if err := sampleServiceServer.validate(); err != nil {
		return nil, fmt.Errorf("failed to validate() :%w", err)
	}
	return sampleServiceServer, nil
}

// validate SampleServiceServerのバリデーション
func (s *SampleServiceServer) validate() error {
	if s.iListSamplesUseCase == nil {
		return errors.New("iListSamplesUseCase is nil")
	}
	if s.iCreateSampleUseCase == nil {
		return errors.New("iCreateSampleUseCase is nil")
	}
	if s.iUpdateSampleUseCase == nil {
		return errors.New("iUpdateSampleUseCase is nil")
	}
	if s.iDeleteSampleUseCase == nil {
		return errors.New("iDeleteSampleUseCase is nil")
	}
	return nil
}
