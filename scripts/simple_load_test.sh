#!/bin/bash

# Simple load test script for test-api
# Sends requests to the API at ~1 request per second

API_URL="http://localhost:8080"
COUNTER=1

echo "Starting simple load test against $API_URL"
echo "Press Ctrl+C to stop"
echo "========================================"

# Function to handle Ctrl+C
trap 'echo -e "\n\nStopped. Total requests: $COUNTER"; exit 0' INT

while true; do
    RAND=$((RANDOM % 100))
    
    if [ $RAND -lt 70 ]; then
        # 70% GET all users
        echo "[$COUNTER] GET /users"
        curl -s -X GET "$API_URL/users" > /dev/null && echo "✓ Success" || echo "✗ Failed"
        
    elif [ $RAND -lt 90 ]; then
        # 20% GET single user
        USER_ID=$((RANDOM % 3 + 1))
        echo "[$COUNTER] GET /users/$USER_ID"
        curl -s -X GET "$API_URL/users/$USER_ID" > /dev/null && echo "✓ Success" || echo "✗ Failed"
        
    else
        # 10% POST new user
        NAME="LoadTest_$COUNTER"
        EMAIL="test_${COUNTER}_$(date +%s)@example.com"
        AGE=$((RANDOM % 40 + 20))
        echo "[$COUNTER] POST /users (name=$NAME)"
        curl -s -X POST "$API_URL/users" \
            -H "Content-Type: application/json" \
            -d "{\"name\":\"$NAME\",\"email\":\"$EMAIL\",\"age\":$AGE}" > /dev/null && echo "✓ Success" || echo "✗ Failed"
    fi
    
    echo "--------"
    COUNTER=$((COUNTER + 1))
    sleep 1
done
