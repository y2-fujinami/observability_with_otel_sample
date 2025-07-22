variable "default_project_id" {
  type        = string
  description = "各リソースのデフォルトのプロジェクトID"
  default     = "y2-fujinami-study"
}
variable "default_region" {
  type        = string
  description = "各リージョンリソースのデフォルトのリージョン"
  default     = "us-central1"
}
variable "default_zone" {
  type        = string
  description = "各ゾーンリソースのデフォルトのゾーン"
  default     = "us-central1-a"
}

variable "cloud_run_api" {
  type = object({
    min_instance_count              = number
    max_instance_count              = number
    cpu_idle                        = bool
    startup_cpu_boost               = bool
    app_port                        = number
    app_limit_cpu                   = string
    app_limit_memory                = string
    otel_collector_port             = number
    otel_collector_healthcheck_port = number
    otel_collector_limit_cpu        = string
    otel_collector_limit_memory     = string
  })
  description = "Cloud Run services (サービス名:API)の設定値"
  default = {
    # リビジョンインスタンス数最小値
    min_instance_count = 0
    # リビジョンインスタンス数最大値
    max_instance_count = 1
    # アプリケーションコンテナでリクエストがあるときだけCPU割り当てるか(true:リクエストベースの課金、false:インスタンスベースの課金)
    cpu_idle = false
    # CPUブーストするか(コールドスタート時のレイテンシを低減する)
    startup_cpu_boost = false
    # アプリケーションコンテナのポート番号(外部から内部への転送先)
    app_port = 8080
    # アプリケーションコンテナのCPU上限
    app_limit_cpu = "0.8"
    # アプリケーションコンテナのメモリ上限
    app_limit_memory = "512Mi"
    # Otel コレクターコンテナのポート
    otel_collector_port = 4317
    # Otel コレクターコンテナのヘルスチェック用ポート
    otel_collector_healthcheck_port = 13133
    # Otel コレクターコンテナのCPU上限
    otel_collector_limit_cpu = "0.2"
    # Otel コレクターコンテナのメモリ上限
    otel_collector_limit_memory = "512Mi"
  }
}

variable "spanner_instance_dev" {
  description = "Cloud Spannerのインスタンスの設定値"
  type = object({
    name         = string
    config       = string
    display_name = string
    num_nodes    = number
  })
  default = {
    # インスタンス名(ID)
    name = "y2-fujinami-study-dev"
    # インスタンス構成
    config = "regional-us-central1"
    # 表示名
    display_name = "y2-fujinami-study-dev"
    # ノード数
    num_nodes = 1
  }
}

variable "spanner_database_dev" {
  description = "Cloud Spannerのデータベースの設定値"
  type = object({
    instance = string
    name     = string
  })
  default = {
    # インスタンス名(ID)
    instance = "y2-fujinami-study-dev"
    # データベース名(ID)
    name = "y2-fujinami-study-dev-1"
  }
}
