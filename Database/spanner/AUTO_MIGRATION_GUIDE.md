# 🤖 자동 마이그레이션 생성 가이드

스키마 diff를 분석해서 마이그레이션 파일을 자동으로 생성하는 방법을 설명합니다.

## 🎯 목표

```
현재 DB 스키마 → schema.sql 비교 → 차이점 발견 → 마이그레이션 파일 생성
```

## 🔧 방법 1: 내장 스크립트 사용 (권장)

### 사용법

```bash
# 1. schema.sql 수정
vi schema/schema.sql
# (예: users 테이블에 age 컬럼 추가)

# 2. 마이그레이션 템플릿 자동 생성
make generate-migration

# 3. 프롬프트에 따라 파일 이름 입력
# 파일 이름 (예: add_age_column): add_age_column

# 4. 생성된 템플릿 편집
vi migrations/production/20251127_143000_add_age_column.sql
```

### 생성되는 내용

```sql
-- Migration: add_age_column
-- Generated: 2025-11-27
-- 
-- 주의: 이 파일은 자동 생성된 템플릿입니다.
-- 실제 적용 전에 반드시 검토하고 수정하세요!
--
-- 현재 DB → 목표 스키마 차이:
-- [hammer diff 결과 요약]

-- TODO: 아래 SQL을 실제 변경사항에 맞게 수정하세요

-- 예시: 컬럼 추가
-- ALTER TABLE users ADD COLUMN age INT64;
```

### 워크플로우

```bash
# 1. 개발: schema.sql 수정
vi schema/schema.sql
make resetdb  # 로컬 테스트

# 2. 운영용 마이그레이션 생성
make generate-migration
# → migrations/production/xxx.sql 생성

# 3. 마이그레이션 파일 검토 및 수정
vi migrations/production/xxx.sql

# 4. 로컬에서 테스트
make resetdb
gcloud spanner databases ddl update test-db \
  --instance=test-instance \
  --ddl-file=migrations/production/xxx.sql

# 5. 테스트 통과 후 커밋
git add migrations/production/xxx.sql schema/schema.sql
git commit -m "feat: add age column to users"

# 6. 운영 배포
gcloud spanner databases ddl update production-db \
  --instance=production-instance \
  --ddl-file=migrations/production/xxx.sql
```

## 🛠️ 방법 2: 외부 도구 사용

### Atlas (추천)

**장점:**
- 완전 자동 diff 생성
- 안전성 검증 포함
- 다양한 DB 지원

**설치:**
```bash
# macOS
brew install ariga/tap/atlas

# Linux
curl -sSf https://atlasgo.sh | sh
```

**사용법:**
```bash
# 1. 스키마 파일 준비
# old_schema.sql: 현재 DB 상태
make db-export > old_schema.sql

# 2. 새 스키마
# schema/schema.sql: 목표 상태

# 3. Diff 생성
atlas schema diff \
  --from "file://old_schema.sql" \
  --to "file://schema/schema.sql" \
  > migration_001.sql

# 4. 생성된 마이그레이션 검토
cat migration_001.sql
```

**한계:**
- Spanner 지원 제한적 (PostgreSQL 호환 모드만)
- 일부 Spanner 전용 기능 미지원

### Bytebase (GUI 선호 시)

**장점:**
- GUI로 diff 확인
- 팀 협업 기능
- 승인 워크플로우

**사용법:**
1. https://bytebase.com 접속
2. Spanner 연결 설정
3. Schema Editor에서 변경
4. "Generate Migration" 클릭

**한계:**
- Spanner 지원은 유료 플랜
- 외부 서비스 의존

### Prisma Migrate

**장점:**
- ORM과 통합
- 완전 자동화

**한계:**
- ❌ Spanner 미지원
- 다른 DB로 전환 시에만 유용

## 💡 실전 팁

### 1. 변경 사항별 파일 분리

```bash
# 나쁜 예: 한 파일에 모든 변경
20251127_big_migration.sql  # 10개 테이블 변경

# 좋은 예: 기능별 분리
20251127_001_add_age_column.sql
20251127_002_create_index.sql
20251127_003_add_posts_table.sql
```

**이유:**
- 롤백 단위가 명확
- 문제 발생 시 어느 부분인지 쉽게 파악

### 2. 안전한 변경 우선

```sql
-- ✅ 안전 (데이터 손실 없음)
ALTER TABLE users ADD COLUMN age INT64;
CREATE INDEX users_age_idx ON users(age);

-- ⚠️ 주의 (데이터 손실 가능)
ALTER TABLE users DROP COLUMN legacy_field;
DROP TABLE old_table;
```

### 3. 테스트 데이터로 검증

```bash
# 1. 샘플 데이터 삽입
make seed-data

# 2. 마이그레이션 적용
gcloud spanner databases ddl update test-db \
  --ddl-file=migrations/xxx.sql

# 3. 데이터 확인
make test-query

# 4. 문제 없으면 운영 적용
```

### 4. 롤백 계획 수립

```sql
-- migration_001_add_age.sql (UP)
ALTER TABLE users ADD COLUMN age INT64;

-- 별도로 롤백 SQL 준비
-- rollback_001_add_age.sql (DOWN)
ALTER TABLE users DROP COLUMN age;
```

## 🔍 Diff 분석 전략

### 작은 변경부터 시작

```bash
# 1단계: 컬럼 추가만
ALTER TABLE users ADD COLUMN age INT64;

# 테스트

# 2단계: 인덱스 추가
CREATE INDEX users_age_idx ON users(age);

# 테스트

# 3단계: 다음 변경...
```

### 의존성 순서 고려

```sql
-- 올바른 순서
1. CREATE TABLE parent ...
2. CREATE TABLE child ... (parent 참조)
3. CREATE INDEX ...

-- 잘못된 순서
1. CREATE TABLE child ... (parent 참조) ← 에러!
2. CREATE TABLE parent ...
```

## 📊 비교: 수동 vs 자동

| 측면 | 수동 작성 | 자동 생성 (스크립트) | 자동 생성 (Atlas) |
|------|-----------|---------------------|-------------------|
| **정확도** | 사람 실수 가능 | 템플릿 제공 | 완전 자동 |
| **속도** | 느림 | 빠름 | 가장 빠름 |
| **학습 곡선** | 낮음 | 중간 | 높음 |
| **커스터마이징** | 완전 제어 | 검토 필요 | 제한적 |
| **Spanner 지원** | ✅ | ✅ | ⚠️ 제한적 |

## 🎯 권장 전략

### 개발 환경
```bash
# 빠른 반복
make resetdb  # schema.sql 직접 사용
```

### 스테이징 환경
```bash
# 마이그레이션 테스트
make generate-migration
# → 마이그레이션 파일 생성 및 테스트
```

### 운영 환경
```bash
# 검증된 마이그레이션만 적용
gcloud spanner databases ddl update prod-db \
  --ddl-file=migrations/production/xxx.sql
```

## ⚠️ 주의사항

### 자동 생성의 한계

1. **복잡한 변경**
```sql
-- 자동 생성 어려움: 데이터 변환 로직
-- 1. 새 컬럼 추가
-- 2. 기존 데이터 변환
-- 3. 기존 컬럼 삭제
-- → 수동 작성 필요
```

2. **Spanner 특수 기능**
```sql
-- INTERLEAVE 같은 Spanner 전용 기능
-- → 도구가 이해 못할 수 있음
```

3. **성능 고려**
```sql
-- 대용량 테이블에 인덱스 추가
-- → 시간이 오래 걸릴 수 있음
-- → 점진적 적용 전략 필요
```

## 📚 참고

- **내장 스크립트**: `scripts/generate_migration.sh`
- **운영 마이그레이션**: `migrations/production/`
- **스키마 히스토리**: `schema/SCHEMA_HISTORY.md`

## 🎉 결론

**가장 실용적인 접근:**
1. 개발: `schema.sql` 직접 수정 + `make resetdb`
2. 운영 배포용: `make generate-migration`으로 템플릿 생성
3. 템플릿 검토 및 수동 조정
4. 철저한 테스트 후 운영 적용

**완전 자동화는 어렵지만, 템플릿 생성으로 90% 시간 절약!** 🚀

