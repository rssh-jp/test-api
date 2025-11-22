.PHONY: help build up down restart logs clean test vulncheck generate

# デフォルトターゲット
help:
	@echo "利用可能なコマンド:"
	@echo "  make build      - Dockerイメージをビルド"
	@echo "  make up         - 全サービスを起動（フォアグラウンド）"
	@echo "  make up-d       - 全サービスを起動（バックグラウンド）"
	@echo "  make down       - 全サービスを停止"
	@echo "  make restart    - 全サービスを再起動"
	@echo "  make logs       - 全サービスのログを表示"
	@echo "  make logs-api   - APIサービスのログを表示"
	@echo "  make logs-mysql - MySQLサービスのログを表示"
	@echo "  make logs-redis - Redisサービスのログを表示"
	@echo "  make clean      - サービスを停止してボリュームを削除"
	@echo "  make prune      - 未使用のDockerリソースを全て削除（注意: 破壊的操作）"
	@echo "  make test       - Goテストを実行"
	@echo "  make test-api   - API統合テストを実行"
	@echo "  make test-api-perf - APIパフォーマンステストを実行（10回反復）"
	@echo "  make vulncheck  - Go脆弱性チェックを実行（govulncheck）"
	@echo "  make vulncheck-verbose - 詳細な脆弱性チェックを実行"
	@echo "  make generate   - OpenAPIコードを生成"
	@echo "  make shell-api  - APIコンテナのシェルを開く"
	@echo "  make mysql-cli  - MySQL CLIを開く"
	@echo "  make redis-cli  - Redis CLIを開く"
	@echo "  make swagger    - Swagger UIをブラウザで開く"
	@echo "  make setup      - 初期セットアップ（generate, build, up）"
	@echo "  make load-test  - 負荷テストを実行（詳細出力）"
	@echo "  make load-test-simple - シンプルな負荷テストを実行"
	@echo "  make load-test-complex - 複雑な負荷テストを実行（JOINクエリのテスト）"

# Dockerイメージをビルド
build:
	docker-compose -f resources/docker/docker-compose.yml --env-file .env build

# 全サービスを起動（フォアグラウンド）
up:
	docker-compose -f resources/docker/docker-compose.yml --env-file .env up

# 全サービスを起動（バックグラウンド）
up-d:
	docker-compose -f resources/docker/docker-compose.yml --env-file .env up -d

# 全サービスを停止
down:
	docker-compose -f resources/docker/docker-compose.yml --env-file .env down

# 全サービスを再起動
restart: down up

# 全サービスのログを表示
logs:
	docker-compose -f resources/docker/docker-compose.yml --env-file .env logs -f

# 特定サービスのログを表示
logs-api:
	docker-compose -f resources/docker/docker-compose.yml --env-file .env logs -f api

logs-mysql:
	docker-compose -f resources/docker/docker-compose.yml --env-file .env logs -f mysql

logs-redis:
	docker-compose -f resources/docker/docker-compose.yml --env-file .env logs -f redis

# サービスを停止してボリュームを削除
clean:
	docker-compose -f resources/docker/docker-compose.yml --env-file .env down -v
	rm -rf api/gen/

# 全ての未使用Dockerリソースを削除（コンテナ、イメージ、ボリューム、ネットワーク）
prune:
	@echo "警告: 未使用のDockerリソースを全て削除します！"
	@bash -c 'read -p "本当に実行しますか？ [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		docker system prune -a --volumes -f; \
		echo "Dockerリソースを削除しました！"; \
	else \
		echo "削除をキャンセルしました。"; \
	fi'

# Goテストを実行
test:
	cd api && go test -v ./...

# API統合テストを実行
test-api:
	@echo "API統合テストを実行中..."
	@bash scripts/test_api.sh

# APIパフォーマンステストを実行（10回反復）
test-api-perf:
	@echo "APIパフォーマンステストを実行中（10回反復）..."
	@ITERATIONS=10 bash scripts/test_api.sh

# Go脆弱性チェックを実行
vulncheck:
	@echo "脆弱性チェックを実行中..."
	@echo "注意: 終了コード3は脆弱性が見つかったが間接的な依存関係の可能性があります"
	@cd api && go run golang.org/x/vuln/cmd/govulncheck@latest ./... || true

# 詳細な脆弱性チェックを実行
vulncheck-verbose:
	@echo "詳細な脆弱性チェックを実行中..."
	@cd api && go run golang.org/x/vuln/cmd/govulncheck@latest -show verbose ./...

# OpenAPIコードをローカルで生成
generate:
	@echo "oapi-codegenをインストール中..."
	@cd api && go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
	@echo "OpenAPIコードを生成中..."
	@mkdir -p api/gen
	@cd api && oapi-codegen -package gen -generate types,server,spec ../resources/openapi/openapi.yaml > gen/openapi.gen.go
	@echo "OpenAPIコードの生成が完了しました！"

# APIコンテナのシェルを開く
shell-api:
	docker-compose -f resources/docker/docker-compose.yml --env-file .env exec api sh

# MySQL CLIを開く
mysql-cli:
	docker-compose -f resources/docker/docker-compose.yml --env-file .env exec mysql mysql -uroot -ppassword testdb

# Redis CLIを開く
redis-cli:
	docker-compose -f resources/docker/docker-compose.yml --env-file .env exec redis redis-cli

# Swagger UIをブラウザで開く
swagger:
	@echo "Swagger UIを開きます: http://localhost:8081/swagger"
	@command -v xdg-open > /dev/null && xdg-open http://localhost:8081/swagger || \
	command -v open > /dev/null && open http://localhost:8081/swagger || \
	echo "ブラウザで http://localhost:8081/swagger を開いてください"

# 初期セットアップ
setup: generate build up
	@echo "サービスが起動するまで待機中..."
	@sleep 10
	@echo "セットアップ完了！APIは http://localhost:8080 で動作しています"
	@echo "試してみましょう: curl http://localhost:8080/health"

# 負荷テストを実行（詳細なJSON出力）
load-test:
	@echo "詳細出力付きで負荷テストを開始します..."
	@echo "停止するには Ctrl+C を押してください"
	@./scripts/load_test.sh

# シンプルな負荷テストを実行（コンパクト出力）
load-test-simple:
	@echo "シンプルな負荷テストを開始します..."
	@echo "停止するには Ctrl+C を押してください"
	@./scripts/simple_load_test.sh

# 複雑な負荷テストを実行（JOINクエリをテスト）
load-test-complex:
	@echo "JOINクエリを含む複雑な負荷テストを開始します..."
	@echo "停止するには Ctrl+C を押してください"
	@./scripts/complex_load_test.sh
