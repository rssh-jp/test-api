#!/bin/bash

# Simple load test script for test-api (READ ONLY)
# Tests all GET endpoints at ~1 request per second

API_URL="http://localhost:8080"
COUNTER=1

echo "Starting READ-ONLY load test against $API_URL"
echo "Testing all GET endpoints"
echo "Press Ctrl+C to stop"
echo "========================================"

# Function to handle Ctrl+C
trap 'echo -e "\n\nStopped. Total requests: $COUNTER"; exit 0' INT

while true; do
    RAND=$((RANDOM % 100))
    
    if [ $RAND -lt 20 ]; then
        # 20% GET all users
        echo "[$COUNTER] GET /users"
        curl -s -X GET "$API_URL/users" > /dev/null && echo "✓ Success" || echo "✗ Failed"
        
    elif [ $RAND -lt 35 ]; then
        # 15% GET single user by ID
        USER_ID=$((RANDOM % 4 + 1))
        echo "[$COUNTER] GET /users/$USER_ID"
        curl -s -X GET "$API_URL/users/$USER_ID" > /dev/null && echo "✓ Success" || echo "✗ Failed"
        
    elif [ $RAND -lt 45 ]; then
        # 10% GET user detail by ID
        USER_ID=$((RANDOM % 4 + 1))
        echo "[$COUNTER] GET /users/$USER_ID/detail"
        curl -s -X GET "$API_URL/users/$USER_ID/detail" > /dev/null && echo "✓ Success" || echo "✗ Failed"
        
    elif [ $RAND -lt 55 ]; then
        # 10% GET user detail by username
        USERNAMES=("sakura" "takeshi" "yuki" "haruto")
        USERNAME=${USERNAMES[$((RANDOM % 4))]}
        echo "[$COUNTER] GET /users/username/$USERNAME/detail"
        curl -s -X GET "$API_URL/users/username/$USERNAME/detail" > /dev/null && echo "✓ Success" || echo "✗ Failed"
        
    elif [ $RAND -lt 65 ]; then
        # 10% GET all posts
        echo "[$COUNTER] GET /posts"
        curl -s -X GET "$API_URL/posts" > /dev/null && echo "✓ Success" || echo "✗ Failed"
        
    elif [ $RAND -lt 73 ]; then
        # 8% GET single post by ID
        POST_ID=$((RANDOM % 5 + 1))
        echo "[$COUNTER] GET /posts/$POST_ID"
        curl -s -X GET "$API_URL/posts/$POST_ID" > /dev/null && echo "✓ Success" || echo "✗ Failed"
        
    elif [ $RAND -lt 80 ]; then
        # 7% GET post by slug
        SLUGS=("clean-architecture-go" "docker-kubernetes-best-practices" "mysql-performance-optimization" "redis-caching-strategies")
        SLUG=${SLUGS[$((RANDOM % 4))]}
        echo "[$COUNTER] GET /posts/slug/$SLUG"
        curl -s -X GET "$API_URL/posts/slug/$SLUG" > /dev/null && echo "✓ Success" || echo "✗ Failed"
        
    elif [ $RAND -lt 87 ]; then
        # 7% GET featured posts
        echo "[$COUNTER] GET /posts/featured"
        curl -s -X GET "$API_URL/posts/featured" > /dev/null && echo "✓ Success" || echo "✗ Failed"
        
    elif [ $RAND -lt 93 ]; then
        # 6% GET posts by category
        CATEGORIES=("programming" "devops" "technology")
        CATEGORY=${CATEGORIES[$((RANDOM % 3))]}
        echo "[$COUNTER] GET /posts/category/$CATEGORY"
        curl -s -X GET "$API_URL/posts/category/$CATEGORY" > /dev/null && echo "✓ Success" || echo "✗ Failed"
        
    else
        # 7% GET posts by tag
        TAGS=("go" "docker" "kubernetes" "mysql" "redis" "clean-architecture")
        TAG=${TAGS[$((RANDOM % 6))]}
        echo "[$COUNTER] GET /posts/tag/$TAG"
        curl -s -X GET "$API_URL/posts/tag/$TAG" > /dev/null && echo "✓ Success" || echo "✗ Failed"
    fi
    
    echo "--------"
    COUNTER=$((COUNTER + 1))
    sleep 1
done
