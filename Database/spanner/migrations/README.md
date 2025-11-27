# migrations 디렉토리

## 용도

이 디렉토리는 **DML(Data Manipulation Language) 마이그레이션**을 위한 것입니다.

> ⚠️ **중요**: DDL(테이블, 인덱스 등)은 `schema/schema.sql`을 사용합니다!

## DML 마이그레이션이란?

데이터를 조작하는 SQL 문:
- `INSERT`: 데이터 삽입
- `UPDATE`: 데이터 수정
- `DELETE`: 데이터 삭제

## 사용 예시

### 1. 샘플 데이터

```bash
# migrations/dml/001_sample_users.sql
INSERT INTO users (id, email, name, created_at, updated_at)
VALUES 
  ('user-1', 'alice@example.com', 'Alice', CURRENT_TIMESTAMP(), CURRENT_TIMESTAMP()),
  ('user-2', 'bob@example.com', 'Bob', CURRENT_TIMESTAMP(), CURRENT_TIMESTAMP());
```

### 2. 마스터 데이터

```bash
# migrations/dml/002_categories.sql
INSERT INTO categories (id, name, display_order)
VALUES
  ('cat-1', 'Technology', 1),
  ('cat-2', 'Business', 2),
  ('cat-3', 'Entertainment', 3);
```

### 3. 설정 데이터

```bash
# migrations/dml/003_config.sql
INSERT INTO system_config (key, value)
VALUES
  ('maintenance_mode', 'false'),
  ('max_upload_size', '10485760');
```

## 실행 방법

```bash
# 모든 DML 마이그레이션 실행
make migrate-dml
```

이 명령어는 `migrations/dml/` 디렉토리의 모든 `.sql` 파일을 `wrench apply --dml`로 실행합니다.

## 디렉토리 구조

```
migrations/
├── README.md           # 이 파일
└── dml/                # DML SQL 파일들
    ├── 001_sample_users.sql
    ├── 002_categories.sql
    └── 003_config.sql
```

## DDL vs DML

| 특징 | DDL | DML |
|------|-----|-----|
| **용도** | 스키마 정의 | 데이터 조작 |
| **예시** | CREATE TABLE, ALTER TABLE | INSERT, UPDATE, DELETE |
| **파일** | `schema/schema.sql` | `migrations/dml/*.sql` |
| **도구** | hammer | wrench |
| **명령어** | `make db-apply` | `make migrate-dml` |

## 주의사항

1. **파일명 규칙**: 숫자 접두사 사용 (실행 순서 보장)
   - ✅ `001_sample_users.sql`
   - ✅ `002_categories.sql`
   - ❌ `sample_users.sql` (순서 불명확)

2. **멱등성**: 여러 번 실행해도 안전하게
   ```sql
   -- ✅ 멱등성 있는 INSERT (ON CONFLICT 사용 또는 조건부 삽입)
   INSERT INTO users (id, email, name, created_at, updated_at)
   SELECT 'user-1', 'alice@example.com', 'Alice', CURRENT_TIMESTAMP(), CURRENT_TIMESTAMP()
   WHERE NOT EXISTS (SELECT 1 FROM users WHERE id = 'user-1');
   
   -- ❌ 멱등성 없는 INSERT (중복 실행 시 에러)
   INSERT INTO users (id, email, name, created_at, updated_at)
   VALUES ('user-1', 'alice@example.com', 'Alice', CURRENT_TIMESTAMP(), CURRENT_TIMESTAMP());
   ```

3. **트랜잭션**: 각 파일은 별도 트랜잭션으로 실행됨

## 예제 생성

DML 디렉토리와 샘플 파일 생성:

```bash
# DML 디렉토리 생성
mkdir -p migrations/dml

# 샘플 파일 생성
cat > migrations/dml/001_sample_users.sql << 'EOF'
-- 샘플 사용자 데이터
INSERT INTO users (id, email, name, created_at, updated_at)
SELECT 'sample-user-1', 'alice@example.com', 'Alice', CURRENT_TIMESTAMP(), CURRENT_TIMESTAMP()
WHERE NOT EXISTS (SELECT 1 FROM users WHERE id = 'sample-user-1');

INSERT INTO users (id, email, name, created_at, updated_at)
SELECT 'sample-user-2', 'bob@example.com', 'Bob', CURRENT_TIMESTAMP(), CURRENT_TIMESTAMP()
WHERE NOT EXISTS (SELECT 1 FROM users WHERE id = 'sample-user-2');
EOF

# 실행
make migrate-dml
```

## 참고

- [wrench 문서](https://github.com/cloudspannerecosystem/wrench)
- [Spanner DML 문서](https://cloud.google.com/spanner/docs/dml-tasks)

