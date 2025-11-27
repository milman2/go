#!/bin/bash
# 스키마 diff 기반 마이그레이션 파일 생성 도우미

set -e

# 색상 정의
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

# 환경변수
PROJECT_ID=${SPANNER_PROJECT_ID:-"test-project"}
INSTANCE_ID=${SPANNER_INSTANCE_ID:-"test-instance"}
DATABASE_ID=${SPANNER_DATABASE_ID:-"test-db"}

echo -e "${GREEN}🔍 스키마 차이 분석 중...${NC}"

# 1. 현재 DB 스키마 export
echo "1. 현재 DB 스키마 추출 중..."
CURRENT_SCHEMA=$(mktemp)
SPANNER_EMULATOR_HOST=localhost:9010 \
  bin/ext/hammer export \
  spanner://projects/${PROJECT_ID}/instances/${INSTANCE_ID}/databases/${DATABASE_ID} \
  > "$CURRENT_SCHEMA"

echo "   저장됨: $CURRENT_SCHEMA"

# 2. 목표 스키마 (schema.sql)
TARGET_SCHEMA="schema/schema.sql"

if [ ! -f "$TARGET_SCHEMA" ]; then
  echo -e "${RED}❌ schema/schema.sql 파일이 없습니다${NC}"
  exit 1
fi

# 3. Diff 생성
echo ""
echo "2. 차이점 분석 중..."
DIFF_OUTPUT=$(mktemp)

SPANNER_EMULATOR_HOST=localhost:9010 \
  bin/ext/hammer diff \
  "$CURRENT_SCHEMA" \
  "$TARGET_SCHEMA" \
  > "$DIFF_OUTPUT" 2>&1 || true

# 4. 결과 확인
if [ -s "$DIFF_OUTPUT" ]; then
  echo -e "${YELLOW}📋 발견된 차이점:${NC}"
  echo ""
  cat "$DIFF_OUTPUT"
  echo ""
  
  # 5. 마이그레이션 파일 생성 제안
  echo -e "${GREEN}💡 마이그레이션 파일을 생성하시겠습니까?${NC}"
  read -p "파일 이름 (예: add_age_column): " MIGRATION_NAME
  
  if [ -n "$MIGRATION_NAME" ]; then
    TIMESTAMP=$(date +%Y%m%d_%H%M%S)
    MIGRATION_DIR="migrations/production"
    mkdir -p "$MIGRATION_DIR"
    
    MIGRATION_FILE="${MIGRATION_DIR}/${TIMESTAMP}_${MIGRATION_NAME}.sql"
    
    # DDL 변경사항을 파일로 저장 (수동 편집 필요)
    cat > "$MIGRATION_FILE" << EOF
-- Migration: ${MIGRATION_NAME}
-- Generated: $(date)
-- 
-- 주의: 이 파일은 자동 생성된 템플릿입니다.
-- 실제 적용 전에 반드시 검토하고 수정하세요!
--
-- 현재 DB → 목표 스키마 차이:
-- $(cat "$DIFF_OUTPUT" | head -10)

-- TODO: 아래 SQL을 실제 변경사항에 맞게 수정하세요

-- 예시: 컬럼 추가
-- ALTER TABLE users ADD COLUMN age INT64;

-- 예시: 인덱스 추가
-- CREATE INDEX users_age_idx ON users(age);

-- 예시: 테이블 생성
-- CREATE TABLE new_table (
--   id STRING(36) NOT NULL,
--   ...
-- ) PRIMARY KEY (id);

EOF
    
    echo ""
    echo -e "${GREEN}✅ 마이그레이션 템플릿 생성됨:${NC}"
    echo "   $MIGRATION_FILE"
    echo ""
    echo -e "${YELLOW}📝 다음 단계:${NC}"
    echo "   1. $MIGRATION_FILE 파일을 열어서 실제 SQL 작성"
    echo "   2. 로컬에서 테스트: make resetdb && gcloud spanner databases ddl update ..."
    echo "   3. 테스트 통과 후 운영 적용"
    echo ""
    echo -e "${YELLOW}💡 Diff 상세 내용:${NC}"
    echo "   cat $DIFF_OUTPUT"
  fi
else
  echo -e "${GREEN}✅ 차이점 없음! 스키마가 동기화되어 있습니다.${NC}"
fi

# 임시 파일 정리
# rm -f "$CURRENT_SCHEMA" "$DIFF_OUTPUT"
echo ""
echo "임시 파일은 검토를 위해 보관됩니다:"
echo "  현재 스키마: $CURRENT_SCHEMA"
echo "  Diff 결과: $DIFF_OUTPUT"

