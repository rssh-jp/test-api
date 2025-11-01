#!/bin/bash

# Simple load test script for test-api (READ ONLY)
# Tests all GET endpoints at ~1 request per second
# Usage: ./simple_load_test.sh [--no-cache]

API_URL="http://localhost:8080"
COUNTER=1
NO_CACHE_PARAM=""

# --no-cacheオプションをチェック
if [ "$1" == "--no-cache" ]; then
    NO_CACHE_PARAM="?no_cache=true"
    echo "Starting READ-ONLY load test against $API_URL (NO CACHE MODE)"
else
    echo "Starting READ-ONLY load test against $API_URL (WITH CACHE)"
fi

echo "Testing all GET endpoints"
echo "Press Ctrl+C to stop"
echo "========================================"

# Function to handle Ctrl+C
trap 'echo -e "\n\nStopped. Total requests: $COUNTER"; exit 0' INT

while true; do
    RAND=$((RANDOM % 100))
    
    if [ $RAND -lt 20 ]; then
        # 20% GET all users
        echo "[$COUNTER] GET /users$NO_CACHE_PARAM"
        curl -s -X GET "$API_URL/users$NO_CACHE_PARAM" > /dev/null && echo "✓ Success" || echo "✗ Failed"
        
    elif [ $RAND -lt 35 ]; then
        # 15% GET single user by ID
        USER_ID=$((RANDOM % 4 + 1))
        echo "[$COUNTER] GET /users/$USER_ID$NO_CACHE_PARAM"
        curl -s -X GET "$API_URL/users/$USER_ID$NO_CACHE_PARAM" > /dev/null && echo "✓ Success" || echo "✗ Failed"
        
    elif [ $RAND -lt 45 ]; then
        # 10% GET user detail by ID
        USER_ID=$((RANDOM % 4 + 1))
        echo "[$COUNTER] GET /users/$USER_ID/detail$NO_CACHE_PARAM"
        curl -s -X GET "$API_URL/users/$USER_ID/detail$NO_CACHE_PARAM" > /dev/null && echo "✓ Success" || echo "✗ Failed"
        
    elif [ $RAND -lt 55 ]; then
        # 10% GET user detail by username
        USERNAMES=("sakura" "takeshi" "yuki" "haruto")
        USERNAME=${USERNAMES[$((RANDOM % 4))]}
        echo "[$COUNTER] GET /users/username/$USERNAME/detail$NO_CACHE_PARAM"
        curl -s -X GET "$API_URL/users/username/$USERNAME/detail$NO_CACHE_PARAM" > /dev/null && echo "✓ Success" || echo "✗ Failed"
        
    elif [ $RAND -lt 65 ]; then
        # 10% GET all posts
        echo "[$COUNTER] GET /posts$NO_CACHE_PARAM"
        curl -s -X GET "$API_URL/posts$NO_CACHE_PARAM" > /dev/null && echo "✓ Success" || echo "✗ Failed"
        
    elif [ $RAND -lt 73 ]; then
        # 8% GET single post by ID
        POST_ID=$((RANDOM % 5 + 1))
        echo "[$COUNTER] GET /posts/$POST_ID$NO_CACHE_PARAM"
        curl -s -X GET "$API_URL/posts/$POST_ID$NO_CACHE_PARAM" > /dev/null && echo "✓ Success" || echo "✗ Failed"
        
    elif [ $RAND -lt 80 ]; then
        # 7% GET post by slug
        SLUGS=("clean-architecture-go" "docker-kubernetes-best-practices" "mysql-performance-optimization" "redis-caching-strategies")
        SLUG=${SLUGS[$((RANDOM % 4))]}
        echo "[$COUNTER] GET /posts/slug/$SLUG$NO_CACHE_PARAM"
        curl -s -X GET "$API_URL/posts/slug/$SLUG$NO_CACHE_PARAM" > /dev/null && echo "✓ Success" || echo "✗ Failed"
        
    elif [ $RAND -lt 87 ]; then
        # 7% GET featured posts
        echo "[$COUNTER] GET /posts/featured$NO_CACHE_PARAM"
        curl -s -X GET "$API_URL/posts/featured$NO_CACHE_PARAM" > /dev/null && echo "✓ Success" || echo "✗ Failed"
        
    elif [ $RAND -lt 93 ]; then
        # 6% GET posts by category
        CATEGORIES=("programming" "devops" "technology")
        CATEGORY=${CATEGORIES[$((RANDOM % 3))]}
        echo "[$COUNTER] GET /posts/category/$CATEGORY$NO_CACHE_PARAM"
        curl -s -X GET "$API_URL/posts/category/$CATEGORY$NO_CACHE_PARAM" > /dev/null && echo "✓ Success" || echo "✗ Failed"
        
    else
        # 7% GET posts by tag
        TAGS=("go" "docker" "kubernetes" "mysql" "redis" "clean-architecture")
        TAG=${TAGS[$((RANDOM % 6))]}
        echo "[$COUNTER] GET /posts/tag/$TAG$NO_CACHE_PARAM"
        curl -s -X GET "$API_URL/posts/tag/$TAG$NO_CACHE_PARAM" > /dev/null && echo "✓ Success" || echo "✗ Failed"
    fi
    
    echo "--------"
    COUNTER=$((COUNTER + 1))
    sleep 1
done
