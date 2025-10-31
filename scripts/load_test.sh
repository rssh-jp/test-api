#!/bin/bash

# Load test script for test-api
# Sends requests to the API at ~1 request per second

API_URL="http://localhost:8080"
COUNTER=1

echo "Starting load test against $API_URL"
echo "Press Ctrl+C to stop"
echo "----------------------------------------"

# Function to handle Ctrl+C
trap 'echo -e "\n\nStopped load test. Total requests: $COUNTER"; exit 0' INT

while true; do
    # Randomly choose an operation
    RAND=$((RANDOM % 100))
    
    if [ $RAND -lt 60 ]; then
        # 60% GET all users
        echo "[Request #$COUNTER] GET /users"
        curl -s -X GET "$API_URL/users" -H "Content-Type: application/json" | jq -c '.' || echo "Failed"
        
    elif [ $RAND -lt 80 ]; then
        # 20% GET single user (ID 1-3)
        USER_ID=$((RANDOM % 3 + 1))
        echo "[Request #$COUNTER] GET /users/$USER_ID"
        curl -s -X GET "$API_URL/users/$USER_ID" -H "Content-Type: application/json" | jq -c '.' || echo "Failed"
        
    elif [ $RAND -lt 95 ]; then
        # 15% POST new user
        RANDOM_NAME="User_$(date +%s)"
        RANDOM_EMAIL="user_$(date +%s)@example.com"
        RANDOM_AGE=$((RANDOM % 50 + 20))
        echo "[Request #$COUNTER] POST /users"
        curl -s -X POST "$API_URL/users" \
            -H "Content-Type: application/json" \
            -d "{\"name\":\"$RANDOM_NAME\",\"email\":\"$RANDOM_EMAIL\",\"age\":$RANDOM_AGE}" | jq -c '.' || echo "Failed"
    
    else
        # 5% PUT update user
        USER_ID=$((RANDOM % 3 + 1))
        RANDOM_NAME="Updated_User_$(date +%s)"
        RANDOM_EMAIL="updated_$(date +%s)@example.com"
        RANDOM_AGE=$((RANDOM % 50 + 20))
        echo "[Request #$COUNTER] PUT /users/$USER_ID"
        curl -s -X PUT "$API_URL/users/$USER_ID" \
            -H "Content-Type: application/json" \
            -d "{\"name\":\"$RANDOM_NAME\",\"email\":\"$RANDOM_EMAIL\",\"age\":$RANDOM_AGE}" | jq -c '.' || echo "Failed"
    fi
    
    echo "----------------------------------------"
    COUNTER=$((COUNTER + 1))
    sleep 1
done
