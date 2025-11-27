# Spanner Schema

이 디렉토리는 Google Cloud Spanner 데이터베이스의 전체 스키마를 관리합니다.

## 📁 파일 구조

```
schema/
└── schema.sql          # 메인 스키마 파일 (모든 테이블 + 인덱스)
```

## 🎯 사용 방법

### 1. 데이터베이스 생성 (schema.sql 기반)

```bash
make createdb
```

이 명령어는 `schema.sql` 파일을 사용하여 데이터베이스를 생성합니다 (hammer 사용).

### 2. 스키마 변경 적용

```bash
# 스키마 파일 수정 후
make db-apply
```

### 3. 현재 DB와 스키마 파일 비교

```bash
make db-diff
```

### 4. 현재 DB 스키마 내보내기

```bash
make db-export
# → exported_schema.sql 생성
```

## 📝 schema.sql 구조

### 두 가지 구현 방식

**방식 1: FOREIGN KEY (현재 사용 중)**
- 일반적인 참조 관계
- 테이블이 독립적으로 분산
- CASCADE 미지원 (애플리케이션에서 처리)

**방식 2: INTERLEAVE (주석으로 제공)**
- 강한 부모-자식 관계
- CASCADE DELETE 자동 지원
- 성능 최적화 (같은 노드에 저장)
- PRIMARY KEY에 부모 키 포함 필수

### INTERLEAVE 방식으로 전환하려면?

1. `schema.sql` 파일 열기
2. [방식 1] 부분 주석 처리
3. [방식 2] 주석 해제
4. `make resetdb` 실행

## ⚡ Spanner 주요 특징

### 1. DEFAULT 값

```sql
-- ❌ 잘못된 방법
published BOOL NOT NULL DEFAULT false

-- ✅ 올바른 방법
published BOOL NOT NULL DEFAULT (false)  -- 괄호 필수!
```

### 2. INTERLEAVE (CASCADE DELETE)

```sql
CREATE TABLE posts (
  user_id STRING(36) NOT NULL,      -- 부모 키
  id STRING(36) NOT NULL,           -- 자식 키
  -- ...
) PRIMARY KEY (user_id, id),        -- 복합 키
  INTERLEAVE IN PARENT users ON DELETE CASCADE;  -- 자동 CASCADE!
```

### 3. 타임스탬프 자동 관리

```sql
created_at TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true)
```

애플리케이션에서:
```go
spanner.CommitTimestamp  // 현재 시각으로 자동 설정
```

## 🔄 워크플로우

### 개발 중 스키마 변경

1. **schema.sql 수정**
2. **변경사항 확인**: `make db-diff`
3. **적용**: `make db-apply`
4. **모델 재생성**: `make generate-models`

### 완전 리셋

```bash
make resetdb  # DB 삭제 후 schema.sql로 재생성
```

## 🛠️ 실전 활용 시나리오

### 시나리오 1: 안전한 배포 (db-diff 활용)

```bash
# 1. schema.sql 수정 (새 컬럼 추가)
vi schema/schema.sql

# 2. 적용 전 차이 확인 (중요!)
make db-diff
# → 의도하지 않은 변경사항 확인
# → 인덱스 삭제 같은 위험 작업 발견

# 3. 안전하면 적용
make db-apply
make generate-models
```

**왜 필요한가?**
- ✅ 실수로 테이블/인덱스 삭제 방지
- ✅ 배포 전 변경사항 검증
- ✅ 팀원과 리뷰 시 근거 자료

### 시나리오 2: 긴급 수정 후 동기화 (db-export 활용)

```bash
# 운영 중 긴급하게 컬럼 추가함
gcloud spanner databases ddl update test-db \
  --ddl='ALTER TABLE posts ADD COLUMN view_count INT64'

# 1. 현재 DB 상태 내보내기
make db-export  # → exported_schema.sql 생성

# 2. schema.sql과 비교
diff exported_schema.sql schema/schema.sql

# 3. schema.sql 업데이트 (동기화)
vi schema/schema.sql  # 필요한 부분만 반영

# 4. 동기화 확인
make db-diff  # 차이 없어야 함
```

### 시나리오 3: 환경별 스키마 검증

```bash
# 개발 환경
export SPANNER_DATABASE_ID=dev-db
make db-export
mv exported_schema.sql /tmp/dev_schema.sql

# 운영 환경
export SPANNER_DATABASE_ID=prod-db
make db-export
mv exported_schema.sql /tmp/prod_schema.sql

# 차이 확인
diff /tmp/dev_schema.sql /tmp/prod_schema.sql
# → 환경 간 불일치 발견!
```

### 시나리오 4: 버전별 스키마 백업

```bash
# v1.0.0 릴리스 전
make db-export
mv exported_schema.sql docs/schemas/schema_v1.0.0.sql
git add docs/schemas/schema_v1.0.0.sql
git commit -m "docs: Add schema snapshot for v1.0.0"
```

## 📚 참고

- [Spanner 공식 문서](https://cloud.google.com/spanner/docs)
- [Spanner DDL 문법](https://cloud.google.com/spanner/docs/data-definition-language)
- [INTERLEAVE 가이드](https://cloud.google.com/spanner/docs/schema-and-data-model#creating-interleaved-tables)

