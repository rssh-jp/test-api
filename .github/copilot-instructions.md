# GitHub Copilot Instructions for test-api Project

## プロジェクト概要
このプロジェクトは、Go言語で実装されたREST APIで、クリーンアーキテクチャパターンを採用しています。MySQL、Redis、NewRelicを統合し、OpenAPI 3.0で定義されたAPIを提供します。

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

### キャッシュ戦略
- **Redis**:
  - 読み取り: キャッシュ優先、ミス時はDBから取得してキャッシュ
  - 書き込み: DB更新後、関連キャッシュを無効化（Cache-Aside パターン）
  - TTL: デフォルト5分、用途に応じて調整可能
  - キーの命名: `リソース名:ID` 形式（例: `user:123`）

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

## Docker & Makefile

### よく使うコマンド
```bash
make build      # Dockerイメージのビルド
make up         # サービス起動
make down       # サービス停止
make restart    # 再起動
make logs-api   # APIログ表示
make generate   # OpenAPIコード生成
make test       # テスト実行
make clean      # 完全クリーンアップ
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

### モニタリング対象
- HTTPリクエスト/レスポンス（自動）
- データベースクエリ（要ラッパー実装）
- Redisコマンド（要インテグレーション）
- カスタムトランザクション（必要に応じて）

### トランザクション作成
```go
txn := nrApp.StartTransaction("CustomOperation")
defer txn.End()

// 処理...

if err != nil {
    txn.NoticeError(err)
}
```

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

## 参考資料

- [Go公式ドキュメント](https://golang.org/doc/)
- [Echo Framework](https://echo.labstack.com/)
- [OpenAPI Specification](https://swagger.io/specification/)
- [クリーンアーキテクチャ](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Effective Go](https://golang.org/doc/effective_go)

---

このインストラクションに従ってコードを生成・編集してください。不明点がある場合は、既存のコードパターンを参考にしてください。
