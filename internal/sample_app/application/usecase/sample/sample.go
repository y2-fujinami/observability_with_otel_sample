package sample

import "go.opentelemetry.io/contrib/bridges/otelslog"

var logger = otelslog.NewLogger("modern-dev-env-app-sample/internal/sample_app/application/usecase/sample")	