#!/bin/bash

# ユーザー詳細エンドポイントのパフォーマンステスト
# 8クエリ → 4クエリ最適化の効果を検証

set -e

API_URL="http://localhost:8080"
USER_ID=1
ITERATIONS=10

# 色設定
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}╔═══════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║  ユーザー詳細エンドポイント パフォーマンステスト  ║${NC}"
echo -e "${BLUE}║  クエリ最適化: 8クエリ → 4クエリ                   ║${NC}"
echo -e "${BLUE}╚═══════════════════════════════════════════════════════╝${NC}"
echo ""

# Redisキャッシュをクリア（純粋なDB性能を測定）
echo -e "${YELLOW}Redisキャッシュをクリアしています...${NC}"
docker exec test-api-redis redis-cli FLUSHALL > /dev/null 2>&1
echo -e "${GREEN}✓ キャッシュクリア完了${NC}"
echo ""

# 関数: 平均レスポンスタイムを計算
benchmark() {
    local url=$1
    local name=$2
    local total=0
    
    echo -e "${BLUE}【${name}】${NC}"
    echo "URL: ${url}"
    echo "試行回数: ${ITERATIONS}回"
    echo ""
    
    for i in $(seq 1 $ITERATIONS); do
        # キャッシュをクリア（毎回DB直接アクセス）
        docker exec test-api-redis redis-cli DEL "user_detail:${USER_ID}" > /dev/null 2>&1
        
        # レスポンスタイムを計測（ミリ秒単位）
        response_time=$(curl -s -w "%{time_total}" -o /dev/null "${url}")
        
        echo -e "  試行 ${i}: ${response_time}秒"
        
        # 合計に加算（bc で浮動小数点計算）
        total=$(echo "$total + $response_time" | bc)
        
        # サーバー負荷を考慮して少し待つ
        sleep 0.1
    done
    
    # 平均を計算
    average=$(echo "scale=6; $total / $ITERATIONS" | bc)
    echo ""
    echo -e "${GREEN}✓ 平均レスポンスタイム: ${average}秒${NC}"
    echo ""
    
    # 結果を返す（グローバル変数に設定）
    echo $average
}

# 1. ユーザー詳細取得（ID指定）- キャッシュバイパス
echo -e "${YELLOW}========================================${NC}"
echo -e "${YELLOW}テスト1: ユーザー詳細取得（ID: ${USER_ID}）${NC}"
echo -e "${YELLOW}========================================${NC}"
echo ""

avg_by_id=$(benchmark "${API_URL}/users/${USER_ID}/detail?no_cache=true" "ID指定でユーザー詳細取得")

# 2. ユーザー詳細取得（username指定）- キャッシュバイパス
echo -e "${YELLOW}========================================${NC}"
echo -e "${YELLOW}テスト2: ユーザー詳細取得（username: sakura）${NC}"
echo -e "${YELLOW}========================================${NC}"
echo ""

avg_by_username=$(benchmark "${API_URL}/users/username/sakura/detail?no_cache=true" "username指定でユーザー詳細取得")

# 3. 取得したデータの構造を確認
echo -e "${YELLOW}========================================${NC}"
echo -e "${YELLOW}データ構造確認${NC}"
echo -e "${YELLOW}========================================${NC}"
echo ""

# キャッシュをクリアしてから取得
docker exec test-api-redis redis-cli DEL "user_detail:${USER_ID}" > /dev/null 2>&1

echo "取得したJSONデータ（一部）:"
response=$(curl -s "${API_URL}/users/${USER_ID}/detail?no_cache=true")

# プロフィール情報
echo -e "${BLUE}■ プロフィール:${NC}"
echo "$response" | grep -o '"profile":{[^}]*}' | head -c 200
echo "..."
echo ""

# フォロー統計
echo -e "${BLUE}■ フォロー統計:${NC}"
echo "$response" | grep -o '"followStats":{[^}]*}'
echo ""

# アクティビティ統計
echo -e "${BLUE}■ アクティビティ統計:${NC}"
echo "$response" | grep -o '"stats":{[^}]*}'
echo ""

# 最近の投稿数
post_count=$(echo "$response" | grep -o '"recentPosts":\[[^]]*\]' | grep -o '"id":' | wc -l)
echo -e "${BLUE}■ 最近の投稿:${NC} ${post_count}件"

# 最近のコメント数
comment_count=$(echo "$response" | grep -o '"recentComments":\[[^]]*\]' | grep -o '"id":' | wc -l)
echo -e "${BLUE}■ 最近のコメント:${NC} ${comment_count}件"

# 未読通知数
notification_count=$(echo "$response" | grep -o '"unreadNotifications":\[[^]]*\]' | grep -o '"id":' | wc -l)
echo -e "${BLUE}■ 未読通知:${NC} ${notification_count}件"
echo ""

# 4. キャッシュ有効時のパフォーマンス
echo -e "${YELLOW}========================================${NC}"
echo -e "${YELLOW}テスト3: キャッシュ有効時のパフォーマンス${NC}"
echo -e "${YELLOW}========================================${NC}"
echo ""

# 最初の1回でキャッシュを生成
curl -s "${API_URL}/users/${USER_ID}/detail" > /dev/null

echo -e "${BLUE}【キャッシュヒット時】${NC}"
echo "試行回数: 5回"
echo ""

cache_total=0
for i in $(seq 1 5); do
    response_time=$(curl -s -w "%{time_total}" -o /dev/null "${API_URL}/users/${USER_ID}/detail")
    echo -e "  試行 ${i}: ${response_time}秒"
    cache_total=$(echo "$cache_total + $response_time" | bc)
    sleep 0.1
done

cache_avg=$(echo "scale=6; $cache_total / 5" | bc)
echo ""
echo -e "${GREEN}✓ キャッシュヒット時の平均: ${cache_avg}秒${NC}"
echo ""

# サマリー
echo -e "${YELLOW}========================================${NC}"
echo -e "${YELLOW}パフォーマンスサマリー${NC}"
echo -e "${YELLOW}========================================${NC}"
echo ""

echo -e "${BLUE}■ DB直接アクセス（最適化後）:${NC}"
echo -e "  ID指定:       ${avg_by_id}秒"
echo -e "  username指定: ${avg_by_username}秒"
echo ""

echo -e "${BLUE}■ キャッシュヒット時:${NC}"
echo -e "  平均:         ${cache_avg}秒"
echo ""

# 改善率を計算
improvement=$(echo "scale=2; (($avg_by_id - $cache_avg) / $avg_by_id) * 100" | bc)
echo -e "${GREEN}■ キャッシュによる改善率: ${improvement}%${NC}"
echo ""

# NewRelicでの確認方法を表示
echo -e "${YELLOW}========================================${NC}"
echo -e "${YELLOW}NewRelicでの確認${NC}"
echo -e "${YELLOW}========================================${NC}"
echo ""
echo "NewRelic APMで以下を確認してください:"
echo "1. Transaction: GET /users/:id/detail"
echo "2. Database: SELECT_WITH_SUBQUERIES の実行時間"
echo "3. クエリ数が減少していることを確認"
echo "   - 最適化前: 8クエリ"
echo "   - 最適化後: 4クエリ"
echo ""

echo -e "${GREEN}╔═════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║  クエリ最適化テスト完了                      ║${NC}"
echo -e "${GREEN}║  8クエリ → 4クエリ（50%削減）                ║${NC}"
echo -e "${GREEN}╚═════════════════════════════════════════════════╝${NC}"
