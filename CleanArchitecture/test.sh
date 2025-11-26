#!/bin/bash

# Clean Architecture API 테스트 스크립트

BASE_URL="http://localhost:8080"
API_URL="$BASE_URL/api/v1"

echo "🧪 Clean Architecture API 테스트 시작"
echo "======================================"

# Health Check
echo ""
echo "1️⃣ Health Check"
curl -s $BASE_URL/health | jq .

# 사용자 생성
echo ""
echo "2️⃣ 사용자 생성"
USER1=$(curl -s -X POST $API_URL/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "alice@example.com",
    "name": "Alice"
  }')
echo $USER1 | jq .
USER1_ID=$(echo $USER1 | jq -r '.id')

USER2=$(curl -s -X POST $API_URL/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "bob@example.com",
    "name": "Bob"
  }')
echo $USER2 | jq .
USER2_ID=$(echo $USER2 | jq -r '.id')

USER3=$(curl -s -X POST $API_URL/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "charlie@example.com",
    "name": "Charlie"
  }')
echo $USER3 | jq .

# 모든 사용자 조회
echo ""
echo "3️⃣ 모든 사용자 조회"
curl -s $API_URL/users | jq .

# 특정 사용자 조회
echo ""
echo "4️⃣ 특정 사용자 조회 (ID: $USER1_ID)"
curl -s $API_URL/users/$USER1_ID | jq .

# 사용자 수정
echo ""
echo "5️⃣ 사용자 수정 (ID: $USER1_ID)"
curl -s -X PUT $API_URL/users/$USER1_ID \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Alice Updated"
  }' | jq .

# 수정 확인
echo ""
echo "6️⃣ 수정된 사용자 확인"
curl -s $API_URL/users/$USER1_ID | jq .

# 사용자 삭제
echo ""
echo "7️⃣ 사용자 삭제 (ID: $USER2_ID)"
curl -s -X DELETE $API_URL/users/$USER2_ID -w "\nHTTP Status: %{http_code}\n"

# 삭제 후 목록 확인
echo ""
echo "8️⃣ 삭제 후 목록 확인"
curl -s $API_URL/users | jq .

# 존재하지 않는 사용자 조회 (404 테스트)
echo ""
echo "9️⃣ 존재하지 않는 사용자 조회 (404 테스트)"
curl -s $API_URL/users/non-existent-id | jq .

# 중복 이메일로 생성 시도 (409 테스트)
echo ""
echo "🔟 중복 이메일로 생성 시도 (409 테스트)"
curl -s -X POST $API_URL/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "alice@example.com",
    "name": "Another Alice"
  }' | jq .

# 잘못된 데이터로 생성 시도 (400 테스트)
echo ""
echo "1️⃣1️⃣ 잘못된 데이터로 생성 시도 - 이메일 없음"
curl -s -X POST $API_URL/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "No Email"
  }' | jq .

echo ""
echo "1️⃣2️⃣ 잘못된 데이터로 생성 시도 - 이름 없음"
curl -s -X POST $API_URL/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com"
  }' | jq .

echo ""
echo "======================================"
echo "✅ 테스트 완료"
echo ""
echo "📚 Clean Architecture 레이어:"
echo "   - Domain (Entities): internal/domain/"
echo "   - Use Cases: internal/usecase/"
echo "   - Interface Adapters: internal/repository/, internal/delivery/http/"
echo "   - Frameworks & Drivers: cmd/api/"

