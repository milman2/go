#!/bin/bash

# Spanner + yo API 테스트 스크립트

BASE_URL="http://localhost:8080"
API_URL="$BASE_URL/api/v1"

echo "🧪 Google Cloud Spanner + yo API 테스트"
echo "==========================================="
echo "📦 Database: Cloud Spanner Emulator"
echo "🔨 Code Generator: yo (go.mercari.io/yo)"
echo "🔧 Migration: Hammer + Wrench"
echo "==========================================="

# Health Check
echo ""
echo "1️⃣ Health Check"
curl -s $BASE_URL/health | jq .

# 사용자 생성
echo ""
echo "2️⃣ 사용자 생성 (Spanner INSERT)"
USER1=$(curl -s -X POST $API_URL/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "alice@spanner.com",
    "name": "Alice (Spanner)"
  }')
echo $USER1 | jq .
USER1_ID=$(echo $USER1 | jq -r '.id')

USER2=$(curl -s -X POST $API_URL/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "bob@spanner.com",
    "name": "Bob (Spanner)"
  }')
echo $USER2 | jq .

# 모든 사용자 조회
echo ""
echo "3️⃣ 모든 사용자 조회 (Spanner SELECT)"
curl -s $API_URL/users | jq .

# 특정 사용자 조회
echo ""
echo "4️⃣ 특정 사용자 조회 (yo FindUserByID)"
curl -s $API_URL/users/$USER1_ID | jq .

# 사용자 삭제
echo ""
echo "5️⃣ 사용자 삭제 (Spanner DELETE)"
curl -s -X DELETE $API_URL/users/$USER1_ID -w "\nHTTP Status: %{http_code}\n"

# 삭제 후 목록
echo ""
echo "6️⃣ 삭제 후 목록 확인"
curl -s $API_URL/users | jq .

echo ""
echo "==========================================="
echo "✅ 테스트 완료"
echo ""
echo "📚 yo 생성 코드 위치: models/"
echo "📋 마이그레이션 파일: migrations/"
echo ""
echo "🔧 유용한 명령어:"
echo "  make show-schema    # 스키마 확인"
echo "  make spanner-cli    # Spanner CLI 접속"
echo "  make generate-yo    # 코드 재생성"

