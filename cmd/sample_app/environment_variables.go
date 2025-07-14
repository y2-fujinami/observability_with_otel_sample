package main

import (
	"errors"
	"fmt"
	"slices"

	"github.com/kelseyhightower/envconfig"
)

// EnvironmentVariables 環境変数
type EnvironmentVariables struct {
	// Port リスンポート番号
	Port int `envconfig:"PORT" default:"8080"`
	// GCPProjectID GCPプロジェクトID
	GCPProjectID string `envconfig:"GCP_PROJECT_ID"`
	// SpannerInstanceID SpannerインスタンスID
	SpannerInstanceID string `envconfig:"SPANNER_INSTANCE_ID"`
	// SpannerDatabaseID SpannerデータベースID
	SpannerDatabaseID string `envconfig:"SPANNER_DATABASE_ID"`
	// OtelCollectorHost OpenTelemetry Collector のホスト
	OtelCollectorHost string `envconfig:"OTEL_COLLECTOR_HOST"`
	// Environment 環境(ローカル / 本番)
	Environment string `envconfig:"ENVIRONMENT"`
	// UseOtelStdExporter 標準出力へテレメトリーデータを出力するか
	UseOtelStdoutExporter bool `envconfig:"USE_OTEL_STDOUT_EXPORTER"`
}

// LoadEnvironmentVariables 環境変数を読み込む
func LoadEnvironmentVariables() (*EnvironmentVariables, error) {
	envVars := &EnvironmentVariables{}
	if err := envconfig.Process("", envVars); err != nil {
		return nil, fmt.Errorf("環境変数の読み込みに失敗しました: %w", err)
	}
	if err := envVars.validate(); err != nil {
		return nil, fmt.Errorf("failed to validate(): %w", err)
	}
	return envVars, nil
}

// validate 環境変数のバリデーション
func (e *EnvironmentVariables) validate() error {
	if e.Port == 0 {
		return errors.New("environment variable Port is 0")
	}
	if e.GCPProjectID == "" {
		return errors.New("environment variable GCPProjectID is empty")
	}
	if e.SpannerInstanceID == "" {
		return errors.New("environment variable SpannerInstanceID is empty")
	}
	if e.SpannerDatabaseID == "" {
		return errors.New("environment variable SpannerDatabaseID is empty")
	}
	if e.OtelCollectorHost == "" {
		return errors.New("environment variable OtelCollectorHost is empty")
	}
	collectEnvs := []string{"local", "prod"}
	if !slices.Contains(collectEnvs, e.Environment) {
		return errors.New("environment variable Environment is not collect")
	}

	return nil
}
