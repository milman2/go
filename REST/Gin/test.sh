#!/bin/bash

# Gin REST API 테스트 스크립트

BASE_URL="http://localhost:8080"
API_URL="$BASE_URL/api/v1"

echo "🧪 Gin REST API 테스트 시작"
echo "======================================"

# Health Check
echo ""
echo "1️⃣ Health Check"
curl -s $BASE_URL/health | jq .

# 아이템 생성
echo ""
echo "2️⃣ 아이템 생성"
ITEM1=$(curl -s -X POST $API_URL/items \
  -H "Content-Type: application/json" \
  -d '{
    "name": "노트북",
    "description": "고성능 노트북",
    "price": 1500000
  }')
echo $ITEM1 | jq .

ITEM2=$(curl -s -X POST $API_URL/items \
  -H "Content-Type: application/json" \
  -d '{
    "name": "마우스",
    "description": "무선 마우스",
    "price": 30000
  }')
echo $ITEM2 | jq .

# 모든 아이템 조회
echo ""
echo "3️⃣ 모든 아이템 조회"
curl -s $API_URL/items | jq .

# 특정 아이템 조회
echo ""
echo "4️⃣ 특정 아이템 조회 (ID: 1)"
curl -s $API_URL/items/1 | jq .

# 아이템 수정
echo ""
echo "5️⃣ 아이템 수정 (ID: 1)"
curl -s -X PUT $API_URL/items/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "노트북 (수정)",
    "description": "최신 고성능 노트북",
    "price": 1800000
  }' | jq .

# 수정 후 조회
echo ""
echo "6️⃣ 수정된 아이템 확인"
curl -s $API_URL/items/1 | jq .

# 아이템 삭제
echo ""
echo "7️⃣ 아이템 삭제 (ID: 2)"
curl -s -X DELETE $API_URL/items/2 | jq .

# 삭제 후 목록 확인
echo ""
echo "8️⃣ 삭제 후 목록 확인"
curl -s $API_URL/items | jq .

# 존재하지 않는 아이템 조회 (404 테스트)
echo ""
echo "9️⃣ 존재하지 않는 아이템 조회 (404 테스트)"
curl -s $API_URL/items/999 | jq .

# 잘못된 데이터로 생성 시도 (유효성 검증 테스트)
echo ""
echo "🔟 잘못된 데이터로 생성 시도 (유효성 검증 테스트)"
curl -s -X POST $API_URL/items \
  -H "Content-Type: application/json" \
  -d '{
    "description": "이름 없음",
    "price": -1000
  }' | jq .

echo ""
echo "======================================"
echo "✅ 테스트 완료"

