# Test API

Go言語で構築されたREST APIです。MySQL、Redis、NewRelicを統合し、クリーンアーキテクチャで実装されています。

## 機能

- ✅ クリーンアーキテクチャ (Domain, Usecase, Infrastructure, Interfaces)
- ✅ OpenAPI 3.0による API定義とコード自動生成
- ✅ Echo Webフレームワーク
- ✅ MySQLデータベース
- ✅ Redisキャッシング（データとリクエスト結果）
- ✅ NewRelic APM統合 (API, MySQL, Redis)
- ✅ Docker Compose による環境構築
- ✅ Makefileによる簡単操作

## プロジェクト構造

```
.
├── .github/
│   └── copilot-instructions.md   # GitHub Copilot用インストラクション
├── api/                          # アプリケーションコード
│   ├── cmd/
│   │   └── main.go              # エントリーポイント
│   ├── domain/                  # ドメイン層
│   │   └── user.go              # エンティティとリポジトリインターフェース
│   ├── usecase/                 # ユースケース層
│   │   └── user_usecase.go      # ビジネスロジック
│   ├── infrastructure/          # インフラ層
│   │   ├── persistence/mysql/   # MySQL実装
│   │   │   └── user_repository.go
│   │   └── cache/redis/         # Redisキャッシュ実装
│   │       └── cached_user_repository.go
│   ├── interfaces/              # インターフェース層
│   │   └── handler/
│   │       └── user_handler.go  # HTTPハンドラー
│   ├── gen/                     # OpenAPIから自動生成されるコード
│   ├── go.mod                   # Go依存関係
│   └── go.sum                   # Go依存関係ロックファイル
├── resources/                   # リソースファイル
│   ├── docker/
│   │   ├── Dockerfile           # アプリケーションDockerfile
│   │   └── docker-compose.yml   # Docker Compose設定
│   ├── openapi/
│   │   └── openapi.yaml         # OpenAPI定義
│   └── database/
│       └── schema.sql           # データベーススキーマ
└── Makefile                     # 操作用Makefile
```

## 必要要件

- Docker
- Docker Compose
- Make
- Go 1.25.3以上（ローカル開発の場合）

## クイックスタート

### 1. リポジトリのクローン

```bash
git clone <repository-url>
cd test-api
```

### 2. 初期セットアップと起動

```bash
make setup
```

このコマンドは以下を実行します：
- OpenAPIコードの生成
- Dockerイメージのビルド
- 全サービスの起動

### 3. 動作確認

```bash
# ヘルスチェック
curl http://localhost:8080/health

# ユーザー一覧取得（シンプルAPI）
curl http://localhost:8080/users

# ユーザー作成
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Test User","email":"test@example.com","age":25}'

# 特定ユーザーの取得
curl http://localhost:8080/users/1

# 投稿一覧取得（複雑なJOIN）
curl 'http://localhost:8080/posts?page=1&pageSize=10'

# 特定投稿の詳細取得（タグとコメント付き）
curl http://localhost:8080/posts/1

# カテゴリー別投稿取得
curl http://localhost:8080/posts/category/programming

# タグ別投稿取得
curl http://localhost:8080/posts/tag/go

# 注目投稿取得
curl 'http://localhost:8080/posts/featured?limit=5'
```

## Makefileコマンド

### 本番環境

```bash
make help         # 利用可能なコマンド一覧を表示
make build        # Dockerイメージをビルド
make up           # 全サービスを起動
make down         # 全サービスを停止
make restart      # 全サービスを再起動
make logs         # 全サービスのログを表示
make logs-api     # APIサービスのログを表示
make logs-mysql   # MySQLのログを表示
make logs-redis   # Redisのログを表示
make clean        # サービス停止とボリューム削除
make test         # テストを実行
make generate     # OpenAPIコードをローカルで生成
make shell-api    # APIコンテナのシェルを開く
make mysql-cli    # MySQL CLIを開く
make redis-cli    # Redis CLIを開く
make load-test    # ロードテスト実行（詳細出力）
make load-test-simple # シンプルロードテスト実行
```

### 開発環境（Hot Reload対応）

```bash
make dev-build    # 開発用Dockerイメージをビルド（reflex込み）
make dev-up       # 開発環境を起動（ホットリロード有効）
make dev-down     # 開発環境を停止
make dev-restart  # 開発環境を再起動
make dev-logs     # 開発環境のログを表示
```

開発環境では、`api/`ディレクトリのGoファイルを編集すると、reflexが自動的に変更を検知してアプリケーションを再起動します（約1-2秒）。

## ロードテスト

秒間1リクエストでAPIに継続的にアクセスするスクリプトを提供しています。

### シンプル版

```bash
make load-test-simple
```

ユーザーAPI（シンプルなクエリ）のみをテスト：

```
[1] GET /users
✓ Success
--------
[2] GET /users/1
✓ Success
--------
[3] POST /users (name=LoadTest_3)
✓ Success
```

### 複雑版（推奨）

```bash
make load-test-complex
```

ユーザーAPIと投稿API（複雑なJOINクエリ）の両方をテスト：

```
[1] GET /posts/tag/go (JOIN: posts+users+post_tags+tags)
✓ Success
--------
[2] GET /posts (JOIN: posts+users+profiles+categories)
✓ Success
--------
[3] GET /posts/category/programming (JOIN: posts+users+categories+tags)
✓ Success
```

**テスト配分:**
- 20% - 投稿一覧（4テーブルJOIN）
- 10% - 投稿詳細（タグ・コメント含む、6テーブル以上）
- 10% - カテゴリー別投稿
- 10% - タグ別投稿（多対多JOIN）
- 5% - 注目投稿
- 20% - ユーザー一覧
- 15% - ユーザー詳細
- 10% - ユーザー作成

### 詳細版（JSON出力）

```bash
make load-test
```

JSON レスポンスを含む詳細な出力を表示します（jq が必要）。

Ctrl+Cで停止できます。NewRelicダッシュボードでリアルタイムに各クエリのパフォーマンスメトリクスを確認できます。

## 環境変数

NewRelicを使用する場合は、環境変数を設定してください：

```bash
export NEW_RELIC_LICENSE_KEY="your-license-key-here"
make up
```

または、`.env`ファイルを作成：

```env
NEW_RELIC_LICENSE_KEY=your-license-key-here
```

## データベース構造

このプロジェクトは、実際のブログ/SNSアプリケーションを想定した複雑なデータベース構造を採用しています。

### テーブル構成

1. **users** - ユーザー基本情報（認証情報）
2. **user_profiles** - ユーザープロフィール詳細
3. **categories** - 投稿カテゴリー（階層構造対応）
4. **posts** - 投稿本体（フルテキスト検索対応）
5. **tags** - タグマスタ
6. **post_tags** - 投稿-タグ関連（多対多）
7. **comments** - コメント（ネスト対応）
8. **user_follows** - フォロー関係（多対多）
9. **likes** - いいね（ポリモーフィック）
10. **notifications** - 通知

### 最適化とインデックス戦略

- **外部キー制約**: 全てのリレーションにFOREIGN KEY制約
- **複合インデックス**: 頻繁にJOINされるカラムにINDEX
- **フルテキスト検索**: posts.title, posts.contentにFULLTEXT INDEX
- **パフォーマンスカウンタ**: view_count, like_count等を非正規化
- **ソフトデリート**: statusカラムで論理削除

### 複雑なJOINクエリの例

#### 投稿一覧取得（4テーブルJOIN + サブクエリ）
```sql
SELECT 
    p.*, 
    u.username, 
    up.display_name, up.avatar_url,
    c.name as category_name, c.slug as category_slug
FROM posts p
INNER JOIN users u ON p.user_id = u.id
LEFT JOIN user_profiles up ON u.id = up.user_id
LEFT JOIN categories c ON p.category_id = c.id
WHERE p.status = 'published'
ORDER BY p.published_at DESC;

-- さらに各投稿のタグとコメントを効率的に取得
-- (N+1問題を回避するため、IN句で一括取得)
```

#### カテゴリー別・タグ別の投稿取得
- カテゴリーで絞り込み: `INNER JOIN categories`
- タグで絞り込み: `INNER JOIN post_tags + INNER JOIN tags`
- 著者情報を結合: `INNER JOIN users + LEFT JOIN user_profiles`

### NewRelicでのクエリトレーシング

全てのMySQLクエリは`newrelic.DatastoreSegment`でトレースされ、以下の情報を記録：
- クエリ実行時間
- テーブル名（Collection）
- 操作タイプ（SELECT, INSERT, UPDATE, DELETE, SELECT_WITH_JOIN等）
- スロークエリの検出

## API仕様

OpenAPI仕様は `resources/openapi/openapi.yaml` に定義されています。

### ユーザーAPI（シンプル）

- `GET /health` - ヘルスチェック
- `GET /users` - ユーザー一覧取得
- `GET /users/{id}` - ユーザー詳細取得
- `POST /users` - ユーザー作成
- `PUT /users/{id}` - ユーザー更新
- `DELETE /users/{id}` - ユーザー削除

### 投稿API（複雑なJOIN）

- `GET /posts?page=1&pageSize=20` - 投稿一覧取得（ページネーション）
- `GET /posts/{id}` - 投稿詳細取得（タグ、コメント、著者情報付き）
- `GET /posts/slug/{slug}` - スラッグで投稿取得
- `GET /posts/category/{slug}` - カテゴリー別投稿取得
- `GET /posts/tag/{slug}` - タグ別投稿取得
- `GET /posts/featured?limit=10` - 注目投稿取得

## キャッシング戦略

Redisは以下のようにキャッシュを管理します：

1. **読み取り操作**:
   - まずRedisキャッシュを確認
   - キャッシュミスの場合、MySQLから取得してキャッシュに保存
   - TTL: 5分

2. **書き込み操作**:
   - MySQLへの書き込み後、関連キャッシュを無効化
   - 次回読み取り時に新しいデータがキャッシュされる

## NewRelic統合

NewRelicが有効な場合、以下が監視されます：

- HTTPリクエスト/レスポンス
- MySQLクエリ
- Redisコマンド
- エラーとスタックトレース
- パフォーマンスメトリクス

## 開発

### ホットリロード開発環境

**推奨**: 開発時は`dev-*`コマンドを使用してください。reflexによるホットリロードで効率的な開発が可能です。

```bash
# 開発環境のビルドと起動
make dev-build
make dev-up

# ログを監視（別ターミナル）
make dev-logs

# コードを編集すると自動的に再起動されます
# 例: api/interfaces/handler/user_handler.go を編集
# → 約1-2秒後にreflexが変更を検知して自動再起動
```

**ホットリロードの仕組み:**
- `api/`ディレクトリがDockerコンテナにマウント
- reflexが`.go`ファイルの変更を監視（1秒ごと）
- 変更検知時に自動的に`go run`で再実行
- OpenAPIコードも起動時に自動生成

### OpenAPIの変更

1. `resources/openapi/openapi.yaml`を編集
2. 開発環境の場合: コンテナを再起動（自動でコード生成）
   ```bash
   make dev-restart
   ```
3. 本番環境の場合: 手動でコード生成してビルド
   ```bash
   make generate
   make build
   ```

### ローカル開発（Dockerなし）

```bash
# OpenAPIコード生成
make generate

# 依存関係のインストール
go mod download

# アプリケーション実行（MySQL/Redisは別途起動済み）
go run api/cmd/main.go

# またはreflexを使用（ローカルでホットリロード）
reflex -r '\.go$' -s -- go run api/cmd/main.go
```

### 技術スタック

- **言語**: Go 1.25.3
- **Webフレームワーク**: Echo v4.12.0
- **データベース**: MySQL 8.0 (utf8mb4)
- **キャッシュ**: Redis 7
- **APM**: NewRelic v3
- **ホットリロード**: Reflex
- **コード生成**: oapi-codegen
- **コンテナ**: Docker & Docker Compose

## トラブルシューティング

### データベース接続エラー

```bash
# MySQLの状態確認
make logs-mysql

# MySQLが起動するまで待つ
docker-compose up -d mysql
sleep 10
make up
```

### Redisキャッシュのクリア

```bash
make redis-cli
> FLUSHALL
> exit
```

### 完全なリセット

```bash
make clean
make setup
```

## ライセンス

MIT
