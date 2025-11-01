# GitHub Copilot Instructions for test-api Project

## プロジェクト概要
このプロジェクトは、Go言語で実装されたREST APIで、クリーンアーキテクチャパターンを採用しています。MySQL、Redis、NewRelicを統合し、OpenAPI 3.0で定義されたAPIを提供します。

### 技術スタック
- **言語**: Go 1.25rc1
- **フレームワーク**: Echo v4.12.0
- **データベース**: MySQL 8.0
- **キャッシュ**: Redis 7-alpine
- **監視**: NewRelic APM v3.40.1
- **開発環境**: Docker + reflex（ホットリロード）
- **アーキテクチャ**: Clean Architecture + Decorator Pattern（キャッシュ層）

## アーキテクチャ原則

### クリーンアーキテクチャの層構造
コードを生成・編集する際は、以下の層構造を厳密に守ってください：

1. **Domain層** (`api/domain/`)
   - エンティティとリポジトリインターフェースのみ
   - 外部依存を一切持たない
   - ビジネスルールの中核
   - 他の層への依存は禁止

2. **Usecase層** (`api/usecase/`)
   - ビジネスロジックを実装
   - Domainのリポジトリインターフェースに依存
   - Infrastructureの具体実装には依存しない
   - 入力バリデーションとオーケストレーション

3. **Infrastructure層** (`api/infrastructure/`)
   - 外部システムとの接続実装
   - `persistence/mysql/`: データベースアクセス
   - `cache/redis/`: キャッシュ実装
   - Domainのリポジトリインターフェースを実装

4. **Interfaces層** (`api/interfaces/`)
   - HTTPハンドラー
   - OpenAPI生成コードを使用
   - リクエスト/レスポンスの変換
   - Usecaseを呼び出す

## コーディング規約

### Go言語の基本ルール
- **エラーハンドリング**: 全てのエラーを適切に処理し、ログに記録
- **命名規則**: 
  - パッケージ名: 小文字、単一単語
  - 関数名: キャメルケース（公開: 大文字始まり、非公開: 小文字始まり）
  - 定数: 大文字始まりのキャメルケース
- **コメント**: 公開関数には必ずドキュメントコメントを記述
- **構造体タグ**: JSONタグは必ず `json:"fieldName"` 形式で記述

### データベース操作
- **MySQL**: 
  - プリペアドステートメントを使用してSQLインジェクションを防止
  - `sql.NullXXX`型を使ってNULL値を適切に処理
  - トランザクションが必要な場合は明示的に開始/コミット/ロールバック
  - **複雑なJOIN**: 正規化されたテーブルからデータを集約する場合、複数のクエリを組み合わせる
    - 例: ユーザー詳細取得（users + user_profiles + posts + comments + follows + notifications）
    - NewRelic DatastoreSegmentで`Operation: "SELECT_WITH_JOINS"`を使用してパフォーマンス監視

### キャッシュ戦略
- **Redis**:
  - 読み取り: キャッシュ優先、ミス時はDBから取得してキャッシュ
  - 書き込み: DB更新後、関連キャッシュを無効化（Cache-Aside パターン）
  - TTL: デフォルト5分、用途に応じて調整可能
  - キーの命名: `リソース名:ID` 形式（例: `user:123`）
  - **Decorator Pattern**: `cached_*_repository.go`でリポジトリをラップ
    - 元のリポジトリを内部に保持（`api/infrastructure/cache/redis/cached_user_repository.go`参照）
    - クリーンアーキテクチャに準拠（Domain層のインターフェースに依存）
  - **キャッシュバイパス**: クエリパラメータ`no_cache=true`でキャッシュをスキップしてDB直接アクセス
    - 例: `GET /users/1?no_cache=true`
    - デバッグや最新データ確認時に使用
    - Handler層で`selectUsecase()`メソッドを使って切り替え

### OpenAPI統合
- **定義ファイル**: `resources/openapi/openapi.yaml`
- **コード生成**: 
  ```bash
  oapi-codegen -package gen -generate types,server,spec openapi.yaml > api/gen/openapi.gen.go
  ```
- **型変換**: OpenAPI生成型（`openapi_types.Email`など）と内部型を適切に変換
- **ハンドラー実装**: `gen.ServerInterface`を実装

## 新機能追加時の手順

### 1. OpenAPI定義の更新
```yaml
# resources/openapi/openapi.yaml
paths:
  /new-endpoint:
    get:
      summary: 新しいエンドポイント
      operationId: newEndpoint
      responses:
        '200':
          description: Success
```

### 2. コード生成
```bash
make generate
```

### 3. Domain層の実装
```go
// api/domain/entity.go
type NewEntity struct {
    ID   int64
    Name string
}

type NewEntityRepository interface {
    FindByID(id int64) (*NewEntity, error)
}
```

### 4. Infrastructure層の実装
```go
// api/infrastructure/persistence/mysql/new_entity_repository.go
type newEntityRepository struct {
    db *sql.DB
}

func NewNewEntityRepository(db *sql.DB) domain.NewEntityRepository {
    return &newEntityRepository{db: db}
}
```

### 5. Usecase層の実装
```go
// api/usecase/new_entity_usecase.go
type NewEntityUsecase interface {
    GetByID(id int64) (*domain.NewEntity, error)
}
```

### 6. Handler層の実装
```go
// api/interfaces/handler/new_entity_handler.go
func (h *NewEntityHandler) NewEndpoint(ctx echo.Context) error {
    // OpenAPI生成型と内部型の変換
    // Usecaseの呼び出し
    // レスポンスの返却
}
```

## テスト

### テストファイルの配置
- ファイル名: `*_test.go`
- 各層ごとにテストを記述
- モックを使用して依存を分離

### テスト例
```go
func TestUsecase(t *testing.T) {
    // モックリポジトリの作成
    mockRepo := &mockRepository{}
    usecase := NewUsecase(mockRepo)
    
    // テストケース実行
    result, err := usecase.Method()
    
    // アサーション
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
}
```

## 開発環境

### ホットリロード
- **reflex**: ファイル変更を検知して自動再ビルド・再起動
- **設定ファイル**: `api/reflex.conf`
- **監視対象**: `.go`ファイルと`openapi.yaml`
- **除外**: `api/gen/`ディレクトリ（生成コードの無限ループ防止）
- **動作**: OpenAPI定義変更時も自動でコード生成＆再起動

### Dockerfile構成
- **ベースイメージ**: `golang:1.25rc1-alpine`
- **開発モード**: ソースコードはボリュームマウント（`docker-compose.yml`）
- **本番ビルド**: `make prod-build`で最適化されたイメージ作成

## Docker & Makefile

### よく使うコマンド
```bash
# 開発用（デフォルト・ホットリロード有効）
make build         # Dockerイメージのビルド
make up            # サービス起動（開発モード）
make down          # サービス停止
make restart       # 再起動
make logs-api      # APIログ表示
make generate      # OpenAPIコード生成
make test          # テスト実行
make clean         # 完全クリーンアップ

# 本番用
make prod-build    # 本番用イメージビルド
make prod-up       # 本番モードで起動
make prod-down     # 本番モード停止

# 負荷テスト
make load-test-simple           # 全GETエンドポイントの負荷テスト（キャッシュあり）
./scripts/simple_load_test.sh --no-cache  # キャッシュバイパスモードで負荷テスト
```

### 環境変数
- `DB_USER`, `DB_PASSWORD`, `DB_HOST`, `DB_PORT`, `DB_NAME`: MySQL接続情報
- `REDIS_HOST`, `REDIS_PORT`, `REDIS_PASSWORD`: Redis接続情報
- `NEW_RELIC_LICENSE_KEY`: NewRelicライセンスキー（オプション）
- `NEW_RELIC_APP_NAME`: NewRelicアプリケーション名

## セキュリティ

### 注意事項
- **SQLインジェクション**: 必ずプリペアドステートメントを使用
- **パスワード**: 平文保存禁止、bcryptでハッシュ化
- **API Key**: 環境変数で管理、コードに直接記述しない
- **ログ**: 機密情報（パスワード、トークン等）をログに出力しない

## パフォーマンス

### 最適化のポイント
- **DB接続プール**: `SetMaxOpenConns`, `SetMaxIdleConns`を適切に設定
- **インデックス**: 頻繁に検索されるカラムにインデックスを作成
- **N+1問題**: JOINや一括取得で回避
- **キャッシュ**: 頻繁にアクセスされるデータはRedisでキャッシュ

## NewRelic統合

### 必須パッケージ
```go
import (
    "github.com/newrelic/go-agent/v3/newrelic"
    "github.com/newrelic/go-agent/v3/integrations/nrecho-v4"  // Echo統合
    "github.com/newrelic/go-agent/v3/integrations/nrmysql"     // MySQL統合
    "github.com/newrelic/go-agent/v3/integrations/nrredis-v8"  // Redis統合
)
```

### モニタリング対象
- **HTTPリクエスト**: nrecho-v4ミドルウェアで自動トレース
- **MySQLクエリ**: DatastoreSegmentで明示的にトレース（必須）
- **Redis操作**: nrredis-v8フックで自動トレース
- **カスタムトランザクション**: 必要に応じて追加

### Handler層でのコンテキスト伝播（必須パターン）
```go
func (h *Handler) GetUser(c echo.Context) error {
    // NewRelicトランザクションを取得
    txn := newrelic.FromContext(c.Request().Context())
    
    // 新しいコンテキストを作成（Usecase/Repositoryに渡す）
    ctx := newrelic.NewContext(c.Request().Context(), txn)
    
    // Usecaseを呼び出す
    user, err := h.usecase.GetUserByID(ctx, id)
    // ...
}
```

### MySQL Repository層でのトレース（必須パターン）
```go
func (r *userRepository) FindByID(ctx context.Context, id int64) (*domain.User, error) {
    // DatastoreSegmentを作成（NewRelicでクエリをトレース）
    segment := newrelic.DatastoreSegment{
        Product:    newrelic.DatastoreMySQL,
        Collection: "users",              // テーブル名
        Operation:  "SELECT",             // 操作種別
        StartTime:  newrelic.FromContext(ctx).StartSegmentNow(),
    }
    defer segment.End()
    
    // SQLクエリ実行
    row := r.db.QueryRowContext(ctx, "SELECT * FROM users WHERE id = ?", id)
    // ...
}
```

### Redis統合の初期化
```go
// Redisクライアント作成時にフックを追加
redisClient := redis.NewClient(&redis.Options{
    Addr: redisAddr,
})
redisClient.AddHook(nrredis.NewHook(redisClient.Options()))
```

### トランザクション作成（カスタム処理用）
```go
txn := nrApp.StartTransaction("CustomOperation")
defer txn.End()

// 処理...

if err != nil {
    txn.NoticeError(err)
}
```

### NewRelic統合チェックリスト
新しいハンドラーやリポジトリを追加する際は、以下を確認：

✅ **Handler層**:
- `newrelic.FromContext()`でトランザクション取得
- `newrelic.NewContext()`でコンテキスト作成
- Usecase呼び出し時にコンテキストを渡す

✅ **Repository層（MySQL）**:
- `newrelic.DatastoreSegment`を作成
- `Product`, `Collection`, `Operation`を適切に設定
- `defer segment.End()`を忘れずに

✅ **Repository層（Redis）**:
- クライアント初期化時に`nrredis.NewHook()`を追加（1回のみ）
- 個別の操作にコードは不要（自動トレース）

## コードレビューのチェックポイント

1. ✅ クリーンアーキテクチャの層構造を守っているか
2. ✅ エラーハンドリングが適切か
3. ✅ SQLインジェクション対策ができているか
4. ✅ NULL値の処理が適切か
5. ✅ キャッシュの無効化が適切に実装されているか
6. ✅ テストが書かれているか
7. ✅ ログ出力に機密情報が含まれていないか
8. ✅ OpenAPI定義とコードが一致しているか
9. ✅ ドキュメントコメントが書かれているか
10. ✅ 適切なHTTPステータスコードを返しているか

## 禁止事項

❌ Domain層からInfrastructure層への直接依存
❌ ビジネスロジックをHandler層に記述
❌ 生のSQL文字列連結（SQLインジェクションリスク）
❌ panic()の使用（エラーは戻り値で返す）
❌ グローバル変数の多用
❌ 機密情報のハードコード
❌ 生成コード（`api/gen/`）の手動編集

## 推奨プラクティス

✅ インターフェースを使った疎結合設計
✅ 依存性注入（DI）パターン
✅ コンテキストの適切な伝播
✅ 適切な粒度でのトランザクション管理
✅ ログレベルの適切な使い分け（Debug, Info, Warn, Error）
✅ テスタビリティを考慮した設計
✅ RESTful API設計原則の遵守
✅ セマンティックバージョニング

## 実装済みエンドポイント

### ユーザー関連
- `GET /users` - 全ユーザー取得
- `GET /users/:id` - ユーザーID取得
- `GET /users/:id/detail` - ユーザー詳細（JOIN集約）
  - ユーザー基本情報 + プロフィール
  - フォロー統計（フォロワー数、フォロー中数）
  - アクティビティ統計（投稿数、コメント数、総いいね数、総閲覧数）
  - 最近の投稿（最新5件）
  - 最近のコメント（最新5件、投稿タイトル付き）
  - 未読通知（最新10件）
- `GET /users/username/:username/detail` - ユーザー名で詳細取得
- `POST /users` - ユーザー作成
- `PUT /users/:id` - ユーザー更新
- `DELETE /users/:id` - ユーザー削除

### 投稿関連
- `GET /posts` - 全投稿取得
- `GET /posts/:id` - 投稿ID取得
- `GET /posts/slug/:slug` - スラッグで取得
- `GET /posts/featured` - 注目投稿取得
- `GET /posts/category/:category` - カテゴリー別取得
- `GET /posts/tag/:tag` - タグ別取得
- `POST /posts` - 投稿作成
- `PUT /posts/:id` - 投稿更新
- `DELETE /posts/:id` - 投稿削除

### 実装パターン例

#### 複雑なJOIN集約（user_detail_repository.go）
```go
// 正規化されたテーブルから全関連データを集約
func (r *userDetailRepository) FindDetailByID(ctx context.Context, id int64) (*domain.UserDetail, error) {
    // NewRelicトレース設定
    segment := newrelic.DatastoreSegment{
        Product:    newrelic.DatastoreMySQL,
        Collection: "users",
        Operation:  "SELECT_WITH_JOINS",
        StartTime:  newrelic.FromContext(ctx).StartSegmentNow(),
    }
    defer segment.End()
    
    // 1. ユーザー基本情報 + プロフィール（JOIN）
    // 2. フォロワー数集計
    // 3. フォロー中数集計
    // 4. 投稿統計集計（件数、閲覧数、いいね数）
    // 5. コメント数集計
    // 6. 最近の投稿取得（最新5件）
    // 7. 最近のコメント取得（投稿タイトルとJOIN、最新5件）
    // 8. 未読通知取得（最新10件）
    
    // 各クエリを実行してUserDetail構造体に集約
    // ...
}
```

## 参考資料

- [Go公式ドキュメント](https://golang.org/doc/)
- [Echo Framework](https://echo.labstack.com/)
- [OpenAPI Specification](https://swagger.io/specification/)
- [クリーンアーキテクチャ](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Effective Go](https://golang.org/doc/effective_go)

---

このインストラクションに従ってコードを生成・編集してください。不明点がある場合は、既存のコードパターンを参考にしてください。
