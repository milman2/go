#!/bin/bash

echo "🧪 Spanner Emulator 종합 테스트"
echo "=================================="

# 환경 변수
export SPANNER_EMULATOR_HOST=localhost:9010
export PROJECT=test-project
export INSTANCE=test-instance
export DATABASE=test-database

# 색상
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo ""
echo -e "${YELLOW}1️⃣ Docker 상태 확인${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
docker ps | grep spanner | grep -v cli

echo ""
echo -e "${YELLOW}2️⃣ HTTP 연결 테스트${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
if curl -s http://localhost:9020 > /dev/null 2>&1; then
  echo -e "${GREEN}✅ HTTP OK (port 9020)${NC}"
else
  echo -e "${RED}❌ HTTP 실패${NC}"
  exit 1
fi

echo ""
echo -e "${YELLOW}3️⃣ gcloud 설정 확인${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
# gcloud 설정
gcloud config set auth/disable_credentials true --quiet
gcloud config set project $PROJECT --quiet
gcloud config set api_endpoint_overrides/spanner http://localhost:9020/ --quiet
echo -e "${GREEN}✅ gcloud 설정 완료${NC}"

echo ""
echo -e "${YELLOW}4️⃣ Instance 목록 조회${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
gcloud spanner instances list 2>/dev/null || echo "Instance 없음 (정상 - 마이그레이션 전)"

echo ""
echo -e "${YELLOW}5️⃣ Database 목록 조회${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
gcloud spanner databases list --instance=$INSTANCE 2>/dev/null || echo "Database 없음 (정상 - 마이그레이션 전)"

# Instance/Database가 있는지 확인
echo ""
echo -e "${YELLOW}6️⃣ Instance/Database 존재 확인${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
DB_EXISTS=$(gcloud spanner databases list --instance=$INSTANCE 2>/dev/null | grep -c "$DATABASE" || echo "0")

if [ "$DB_EXISTS" -eq "0" ]; then
  echo -e "${YELLOW}⚠️ Database가 없습니다. 'make init'을 먼저 실행하세요.${NC}"
  echo ""
  echo "실행 명령어:"
  echo "  make init"
  exit 0
fi

echo -e "${GREEN}✅ Database '$DATABASE' 존재${NC}"

echo ""
echo -e "${YELLOW}7️⃣ 테이블 DDL 조회${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
gcloud spanner databases ddl describe $DATABASE --instance=$INSTANCE

echo ""
echo -e "${YELLOW}8️⃣ 데이터 카운트 조회${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
gcloud spanner databases execute-sql $DATABASE \
  --instance=$INSTANCE \
  --sql="SELECT COUNT(*) as user_count FROM users"

echo ""
echo -e "${YELLOW}9️⃣ Go 연결 테스트${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
if [ -f test_connection.go ]; then
  go run test_connection.go
else
  echo -e "${RED}❌ test_connection.go 파일이 없습니다${NC}"
fi

echo ""
echo -e "${YELLOW}🔟 테이블 정보 조회 (Go)${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
if [ -f test_tables.go ]; then
  go run test_tables.go
else
  echo -e "${RED}❌ test_tables.go 파일이 없습니다${NC}"
fi

echo ""
echo -e "${YELLOW}1️⃣1️⃣ CRUD 테스트 (Go)${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
if [ -f test_crud.go ]; then
  go run test_crud.go
else
  echo -e "${RED}❌ test_crud.go 파일이 없습니다${NC}"
fi

echo ""
echo "=================================="
echo -e "${GREEN}✅ 종합 테스트 완료!${NC}"
echo ""
echo "추가 테스트:"
echo "  make test               # API 서버 테스트"
echo "  make spanner-cli        # Spanner CLI 접속"
echo "  make show-schema        # 스키마 확인"

