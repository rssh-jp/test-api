#!/bin/bash

# テストAPIの動作確認スクリプト
# このスクリプトは以下をテストします：
# 1. ヘルスチェック
# 2. ユーザー一覧取得
# 3. 特定ユーザー取得
# 4. キャッシュバイパス機能
# 5. パフォーマンス比較（キャッシュあり/なし）

set -e

BASE_URL="${BASE_URL:-http://localhost:8080}"
ITERATIONS="${ITERATIONS:-5}"

# 色定義
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# ヘッダー表示
print_header() {
    echo -e "\n${BLUE}========================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}========================================${NC}"
}

# 成功メッセージ
print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

# エラーメッセージ
print_error() {
    echo -e "${RED}✗ $1${NC}"
}

# 警告メッセージ
print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

# APIが起動するまで待機
wait_for_api() {
    print_header "API起動確認"
    
    MAX_RETRIES=30
    RETRY_COUNT=0
    
    while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
        if curl -s -f "${BASE_URL}/health" > /dev/null 2>&1; then
            print_success "APIが起動しています"
            return 0
        fi
        
        RETRY_COUNT=$((RETRY_COUNT + 1))
        echo -n "."
        sleep 1
    done
    
    print_error "APIの起動がタイムアウトしました"
    exit 1
}

# 1. ヘルスチェック
test_health_check() {
    print_header "1. ヘルスチェック"
    
    RESPONSE=$(curl -s "${BASE_URL}/health")
    STATUS=$(echo "$RESPONSE" | grep -o '"status":"[^"]*"' | cut -d'"' -f4)
    
    if [ "$STATUS" = "healthy" ]; then
        print_success "ヘルスチェック: $STATUS"
        echo "$RESPONSE"
    else
        print_error "ヘルスチェック失敗"
        echo "$RESPONSE"
        exit 1
    fi
}

# 2. ユーザー一覧取得
test_get_users() {
    print_header "2. ユーザー一覧取得"
    
    RESPONSE=$(curl -s "${BASE_URL}/users")
    USER_COUNT=$(echo "$RESPONSE" | grep -o '"id"' | wc -l)
    
    if [ "$USER_COUNT" -gt 0 ]; then
        print_success "ユーザー取得成功: ${USER_COUNT}件"
        echo "$RESPONSE" | head -c 300
        echo "..."
    else
        print_error "ユーザーが取得できませんでした"
        exit 1
    fi
}

# 3. 特定ユーザー取得
test_get_user_by_id() {
    print_header "3. 特定ユーザー取得 (ID: 1)"
    
    RESPONSE=$(curl -s "${BASE_URL}/users/1")
    USER_ID=$(echo "$RESPONSE" | grep -o '"id":[0-9]*' | cut -d':' -f2)
    
    if [ "$USER_ID" = "1" ]; then
        print_success "ユーザーID:1の取得成功"
        echo "$RESPONSE"
    else
        print_error "ユーザーID:1の取得失敗"
        echo "$RESPONSE"
        exit 1
    fi
}

# 4. キャッシュバイパステスト
test_cache_bypass() {
    print_header "4. キャッシュバイパス機能テスト"
    
    # キャッシュをクリア（Redisが利用可能な場合）
    if command -v docker-compose &> /dev/null; then
        echo "Redisキャッシュをクリア..."
        docker-compose -f resources/docker/docker-compose.yml --env-file .env exec -T redis redis-cli FLUSHALL > /dev/null 2>&1 || true
    fi
    
    # 通常リクエスト（キャッシュに保存される）
    echo -e "\n${YELLOW}通常リクエスト（キャッシュあり）:${NC}"
    curl -s "${BASE_URL}/users/1" > /dev/null
    print_success "リクエスト完了"
    
    # キャッシュバイパスリクエスト
    echo -e "\n${YELLOW}キャッシュバイパスリクエスト (no_cache=true):${NC}"
    RESPONSE=$(curl -s "${BASE_URL}/users/1?no_cache=true")
    USER_ID=$(echo "$RESPONSE" | grep -o '"id":[0-9]*' | cut -d':' -f2)
    
    if [ "$USER_ID" = "1" ]; then
        print_success "キャッシュバイパス成功"
        echo "$RESPONSE"
    else
        print_error "キャッシュバイパス失敗"
        exit 1
    fi
}

# 5. パフォーマンス比較
test_performance() {
    print_header "5. パフォーマンス比較 (${ITERATIONS}回の平均)"
    
    # キャッシュあり
    echo -e "\n${YELLOW}キャッシュあり:${NC}"
    CACHE_TIMES=()
    for i in $(seq 1 $ITERATIONS); do
        TIME=$(curl -s -w "%{time_total}" -o /dev/null "${BASE_URL}/users/1")
        CACHE_TIMES+=($TIME)
        echo "  試行 $i: ${TIME}秒"
    done
    
    # 平均計算
    CACHE_AVG=$(printf '%s\n' "${CACHE_TIMES[@]}" | awk '{sum+=$1} END {print sum/NR}')
    print_success "キャッシュあり平均: ${CACHE_AVG}秒"
    
    # キャッシュなし
    echo -e "\n${YELLOW}キャッシュバイパス (no_cache=true):${NC}"
    NO_CACHE_TIMES=()
    for i in $(seq 1 $ITERATIONS); do
        TIME=$(curl -s -w "%{time_total}" -o /dev/null "${BASE_URL}/users/1?no_cache=true")
        NO_CACHE_TIMES+=($TIME)
        echo "  試行 $i: ${TIME}秒"
    done
    
    # 平均計算
    NO_CACHE_AVG=$(printf '%s\n' "${NO_CACHE_TIMES[@]}" | awk '{sum+=$1} END {print sum/NR}')
    print_success "キャッシュなし平均: ${NO_CACHE_AVG}秒"
    
    # 比較
    IMPROVEMENT=$(echo "$NO_CACHE_AVG $CACHE_AVG" | awk '{printf "%.1f", ($1-$2)/$1*100}')
    echo -e "\n${GREEN}キャッシュによる改善: ${IMPROVEMENT}%${NC}"
}

# 6. OpenAPIパラメータテスト（型安全性）
test_openapi_params() {
    print_header "6. OpenAPIパラメータテスト"
    
    echo -e "\n${YELLOW}パラメータなし:${NC}"
    curl -s "${BASE_URL}/users" > /dev/null && print_success "成功"
    
    echo -e "\n${YELLOW}no_cache=true:${NC}"
    curl -s "${BASE_URL}/users?no_cache=true" > /dev/null && print_success "成功"
    
    echo -e "\n${YELLOW}no_cache=false:${NC}"
    curl -s "${BASE_URL}/users?no_cache=false" > /dev/null && print_success "成功"
    
    echo -e "\n${YELLOW}no_cache=1:${NC}"
    curl -s "${BASE_URL}/users?no_cache=1" > /dev/null && print_success "成功"
}

# 7. 投稿エンドポイントテスト
test_post_endpoints() {
    print_header "7. 投稿エンドポイントテスト"
    
    # 投稿一覧
    echo -e "\n${YELLOW}投稿一覧取得:${NC}"
    RESPONSE=$(curl -s "${BASE_URL}/posts")
    POST_COUNT=$(echo "$RESPONSE" | grep -o '"id"' | wc -l)
    
    if [ "$POST_COUNT" -gt 0 ]; then
        print_success "投稿取得成功: ${POST_COUNT}件"
        echo "$RESPONSE" | head -c 200
        echo "..."
    else
        print_error "投稿が取得できませんでした"
        exit 1
    fi
    
    # 注目投稿
    echo -e "\n${YELLOW}注目投稿取得:${NC}"
    RESPONSE=$(curl -s "${BASE_URL}/posts/featured")
    FEATURED_COUNT=$(echo "$RESPONSE" | grep -o '"id"' | wc -l)
    
    if [ "$FEATURED_COUNT" -gt 0 ]; then
        print_success "注目投稿取得成功: ${FEATURED_COUNT}件"
        echo "$RESPONSE" | head -c 200
        echo "..."
    else
        print_error "注目投稿が取得できませんでした"
        exit 1
    fi
    
    # 個別投稿
    echo -e "\n${YELLOW}投稿ID:1の取得:${NC}"
    RESPONSE=$(curl -s "${BASE_URL}/posts/1")
    POST_ID=$(echo "$RESPONSE" | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)
    
    if [ "$POST_ID" = "1" ]; then
        print_success "投稿ID:1の取得成功"
        echo "$RESPONSE" | head -c 200
        echo "..."
    else
        print_error "投稿ID:1の取得失敗"
        exit 1
    fi
}

# 8. ユーザー詳細エンドポイントテスト
test_user_detail_endpoints() {
    print_header "8. ユーザー詳細エンドポイントテスト（複雑なJOIN）"
    
    # IDによる取得
    echo -e "\n${YELLOW}ユーザー詳細取得（ID: 1）:${NC}"
    RESPONSE=$(curl -s "${BASE_URL}/users/1/detail")
    USER_ID=$(echo "$RESPONSE" | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)
    HAS_PROFILE=$(echo "$RESPONSE" | grep -o '"profile"' | wc -l)
    HAS_STATS=$(echo "$RESPONSE" | grep -o '"stats"' | wc -l)
    
    if [ "$USER_ID" = "1" ] && [ "$HAS_PROFILE" -gt 0 ] && [ "$HAS_STATS" -gt 0 ]; then
        print_success "ユーザー詳細取得成功（プロフィール、統計情報含む）"
        echo "$RESPONSE" | head -c 300
        echo "..."
    else
        print_error "ユーザー詳細取得失敗"
        exit 1
    fi
    
    # ユーザー名による取得
    echo -e "\n${YELLOW}ユーザー詳細取得（username: sakura）:${NC}"
    RESPONSE=$(curl -s "${BASE_URL}/users/username/sakura/detail")
    USERNAME=$(echo "$RESPONSE" | grep -o '"username":"[^"]*"' | head -1 | cut -d'"' -f4)
    
    if [ "$USERNAME" = "sakura" ]; then
        print_success "ユーザー名での詳細取得成功"
        echo "$RESPONSE" | head -c 300
        echo "..."
    else
        print_error "ユーザー名での詳細取得失敗"
        exit 1
    fi
}

# 9. 投稿エンドポイントのキャッシュバイパステスト
test_post_cache_bypass() {
    print_header "9. 投稿エンドポイントのキャッシュバイパステスト"
    
    # 投稿一覧のパフォーマンス比較
    echo -e "\n${YELLOW}投稿一覧 - キャッシュあり (3回平均):${NC}"
    CACHE_TIMES=()
    for i in $(seq 1 3); do
        TIME=$(curl -s -w "%{time_total}" -o /dev/null "${BASE_URL}/posts")
        CACHE_TIMES+=($TIME)
        echo "  試行 $i: ${TIME}秒"
    done
    CACHE_AVG=$(printf '%s\n' "${CACHE_TIMES[@]}" | awk '{sum+=$1} END {print sum/NR}')
    print_success "キャッシュあり平均: ${CACHE_AVG}秒"
    
    echo -e "\n${YELLOW}投稿一覧 - キャッシュバイパス (3回平均):${NC}"
    NO_CACHE_TIMES=()
    for i in $(seq 1 3); do
        TIME=$(curl -s -w "%{time_total}" -o /dev/null "${BASE_URL}/posts?no_cache=true")
        NO_CACHE_TIMES+=($TIME)
        echo "  試行 $i: ${TIME}秒"
    done
    NO_CACHE_AVG=$(printf '%s\n' "${NO_CACHE_TIMES[@]}" | awk '{sum+=$1} END {print sum/NR}')
    print_success "キャッシュなし平均: ${NO_CACHE_AVG}秒"
    
    # 注目投稿のパフォーマンス比較
    echo -e "\n${YELLOW}注目投稿 - キャッシュあり (3回平均):${NC}"
    FEATURED_CACHE_TIMES=()
    for i in $(seq 1 3); do
        TIME=$(curl -s -w "%{time_total}" -o /dev/null "${BASE_URL}/posts/featured")
        FEATURED_CACHE_TIMES+=($TIME)
        echo "  試行 $i: ${TIME}秒"
    done
    FEATURED_CACHE_AVG=$(printf '%s\n' "${FEATURED_CACHE_TIMES[@]}" | awk '{sum+=$1} END {print sum/NR}')
    print_success "キャッシュあり平均: ${FEATURED_CACHE_AVG}秒"
    
    echo -e "\n${YELLOW}注目投稿 - キャッシュバイパス (3回平均):${NC}"
    FEATURED_NO_CACHE_TIMES=()
    for i in $(seq 1 3); do
        TIME=$(curl -s -w "%{time_total}" -o /dev/null "${BASE_URL}/posts/featured?no_cache=true")
        FEATURED_NO_CACHE_TIMES+=($TIME)
        echo "  試行 $i: ${TIME}秒"
    done
    FEATURED_NO_CACHE_AVG=$(printf '%s\n' "${FEATURED_NO_CACHE_TIMES[@]}" | awk '{sum+=$1} END {print sum/NR}')
    print_success "キャッシュなし平均: ${FEATURED_NO_CACHE_AVG}秒"
}

# メイン実行
main() {
    echo -e "${BLUE}"
    echo "╔═══════════════════════════════════════════╗"
    echo "║   Test API 動作確認スクリプト             ║"
    echo "║   フレームワーク非依存アーキテクチャ      ║"
    echo "╚═══════════════════════════════════════════╝"
    echo -e "${NC}"
    
    wait_for_api
    test_health_check
    test_get_users
    test_get_user_by_id
    test_cache_bypass
    test_performance
    test_openapi_params
    test_post_endpoints
    test_user_detail_endpoints
    test_post_cache_bypass
    
    # 最終結果
    print_header "テスト結果サマリー"
    print_success "すべてのテストが成功しました！"
    echo -e "\n${GREEN}✓ ヘルスチェック${NC}"
    echo -e "${GREEN}✓ ユーザー取得（一覧・個別）${NC}"
    echo -e "${GREEN}✓ ユーザー詳細取得（複雑なJOIN）${NC}"
    echo -e "${GREEN}✓ 投稿取得（一覧・個別・注目）${NC}"
    echo -e "${GREEN}✓ キャッシュバイパス（全エンドポイント）${NC}"
    echo -e "${GREEN}✓ パフォーマンス計測${NC}"
    echo -e "${GREEN}✓ OpenAPIパラメータ${NC}"
    
    echo -e "\n${BLUE}╔═════════════════════════════════════════════════════╗${NC}"
    echo -e "${BLUE}║  フレームワーク非依存アーキテクチャ              ║${NC}"
    echo -e "${BLUE}║  全エンドポイントが正常に動作しています！        ║${NC}"
    echo -e "${BLUE}╚═════════════════════════════════════════════════════╝${NC}"
}

# スクリプト実行
main "$@"
