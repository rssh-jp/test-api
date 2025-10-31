#!/bin/bash

# Complex load test script for test-api
# Tests both simple user API and complex JOIN queries for posts

API_URL="http://localhost:8080"
COUNTER=1

echo "Starting comprehensive load test against $API_URL"
echo "Testing both User API and Post API (with complex JOINs)"
echo "Press Ctrl+C to stop"
echo "========================================"

# Function to handle Ctrl+C
trap 'echo -e "\n\nStopped. Total requests: $COUNTER"; exit 0' INT

while true; do
    RAND=$((RANDOM % 100))
    
    if [ $RAND -lt 20 ]; then
        # 20% GET /posts (complex JOIN)
        echo "[$COUNTER] GET /posts (JOIN: posts+users+profiles+categories)"
        curl -s 'http://localhost:8080/posts?page=1&pageSize=5' > /dev/null && echo "✓ Success" || echo "✗ Failed"
        
    elif [ $RAND -lt 30 ]; then
        # 10% GET /posts/:id (complex JOIN with tags and comments)
        POST_ID=$((RANDOM % 4 + 1))
        echo "[$COUNTER] GET /posts/$POST_ID (JOIN: posts+users+profiles+tags+comments)"
        curl -s "http://localhost:8080/posts/$POST_ID" > /dev/null && echo "✓ Success" || echo "✗ Failed"
        
    elif [ $RAND -lt 40 ]; then
        # 10% GET /posts/category/:slug
        CATEGORIES=("programming" "devops" "technology")
        CAT=${CATEGORIES[$((RANDOM % ${#CATEGORIES[@]}))]}
        echo "[$COUNTER] GET /posts/category/$CAT (JOIN: posts+users+categories+tags)"
        curl -s "http://localhost:8080/posts/category/$CAT" > /dev/null && echo "✓ Success" || echo "✗ Failed"
        
    elif [ $RAND -lt 50 ]; then
        # 10% GET /posts/tag/:slug
        TAGS=("go" "docker" "kubernetes" "mysql" "redis")
        TAG=${TAGS[$((RANDOM % ${#TAGS[@]}))]}
        echo "[$COUNTER] GET /posts/tag/$TAG (JOIN: posts+users+post_tags+tags)"
        curl -s "http://localhost:8080/posts/tag/$TAG" > /dev/null && echo "✓ Success" || echo "✗ Failed"
        
    elif [ $RAND -lt 55 ]; then
        # 5% GET /posts/featured
        echo "[$COUNTER] GET /posts/featured (JOIN: posts+users+profiles+categories)"
        curl -s 'http://localhost:8080/posts/featured?limit=3' > /dev/null && echo "✓ Success" || echo "✗ Failed"
        
    elif [ $RAND -lt 75 ]; then
        # 20% GET /users (simple query)
        echo "[$COUNTER] GET /users (Simple query)"
        curl -s "$API_URL/users" > /dev/null && echo "✓ Success" || echo "✗ Failed"
        
    elif [ $RAND -lt 90 ]; then
        # 15% GET /users/:id (simple query)
        USER_ID=$((RANDOM % 3 + 1))
        echo "[$COUNTER] GET /users/$USER_ID (Simple query)"
        curl -s "$API_URL/users/$USER_ID" > /dev/null && echo "✓ Success" || echo "✗ Failed"
        
    else
        # 10% POST /users (write operation)
        NAME="LoadTest_$COUNTER"
        EMAIL="test_${COUNTER}_$(date +%s)@example.com"
        AGE=$((RANDOM % 40 + 20))
        echo "[$COUNTER] POST /users (Write operation)"
        curl -s -X POST "$API_URL/users" \
            -H "Content-Type: application/json" \
            -d "{\"name\":\"$NAME\",\"email\":\"$EMAIL\",\"age\":$AGE}" > /dev/null && echo "✓ Success" || echo "✗ Failed"
    fi
    
    echo "--------"
    COUNTER=$((COUNTER + 1))
    sleep 1
done
