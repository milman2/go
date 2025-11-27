# DML (Data Manipulation Language) Files

이 디렉토리는 샘플 데이터를 위한 DML(INSERT, UPDATE, DELETE) SQL 파일을 포함합니다.

## 📁 파일 구조

```
dml/
├── README.md
├── 001_seed_users.sql    # 샘플 사용자 데이터
└── 002_seed_posts.sql    # 샘플 게시글 데이터
```

## 🎯 사용 방법

### 1. 샘플 데이터 삽입

```bash
make seed-data
```

이 명령어는 `dml/*.sql` 파일을 순서대로 실행하여 샘플 데이터를 삽입합니다.

### 2. 데이터 확인

```bash
make test-query
```

삽입된 데이터를 확인하는 테스트 쿼리를 실행합니다:
- 사용자 수
- 게시글 수
- 사용자별 발행 게시글 수

### 3. 모든 데이터 삭제

```bash
make clear-data
```

테이블 구조는 유지하고 모든 데이터만 삭제합니다.

### 4. 특정 쿼리 실행

```bash
# 직접 SQL 실행
export SPANNER_EMULATOR_HOST=localhost:9010
gcloud spanner databases execute-sql test-db \
  --instance=test-instance \
  --sql="SELECT * FROM users LIMIT 10"
```

## 📊 샘플 데이터 내용

### Users (3명)

| ID | Email | Name |
|----|-------|------|
| ...0001 | john.doe@example.com | John Doe |
| ...0002 | jane.smith@example.com | Jane Smith |
| ...0003 | bob.johnson@example.com | Bob Johnson |

### Posts (5개)

| User | Title | Published |
|------|-------|-----------|
| John Doe | Getting Started with Cloud Spanner | ✅ |
| John Doe | Advanced Spanner Features | ❌ (draft) |
| Jane Smith | Building Scalable Applications | ✅ |
| Jane Smith | Database Design Best Practices | ✅ |
| Bob Johnson | Work in Progress | ❌ (draft) |

## ✍️ 새로운 DML 파일 추가

### 파일 명명 규칙

```
<순서>_<설명>.sql

예시:
003_seed_comments.sql
004_update_user_status.sql
```

### 파일 템플릿

```sql
-- 설명: 이 파일이 하는 일
-- 사용법: make seed-data

INSERT INTO table_name (column1, column2, ...)
VALUES (
  'value1',
  'value2',
  ...
);
```

## 🔄 워크플로우 예시

### 시나리오 1: 처음부터 시작

```bash
# 1. 데이터베이스 생성
make createdb

# 2. 샘플 데이터 삽입
make seed-data

# 3. 데이터 확인
make test-query

# 4. 모델 생성
make generate-models
```

### 시나리오 2: 데이터 리셋

```bash
# 1. 데이터만 삭제
make clear-data

# 2. 다시 삽입
make seed-data
```

### 시나리오 3: 완전 리셋

```bash
# DB 전체 재생성 (스키마 + 데이터)
make resetdb
make seed-data
```

## 💡 유용한 쿼리 예시

### 1. 특정 사용자의 모든 게시글

```sql
SELECT p.* 
FROM posts p
JOIN users u ON p.user_id = u.id
WHERE u.email = 'john.doe@example.com'
ORDER BY p.created_at DESC;
```

### 2. 발행된 게시글만 조회

```sql
SELECT u.name, p.title, p.created_at
FROM posts p
JOIN users u ON p.user_id = u.id
WHERE p.published = TRUE
ORDER BY p.created_at DESC;
```

### 3. 사용자별 게시글 통계

```sql
SELECT 
  u.name,
  COUNT(p.id) as total_posts,
  SUM(CASE WHEN p.published THEN 1 ELSE 0 END) as published_posts,
  SUM(CASE WHEN NOT p.published THEN 1 ELSE 0 END) as draft_posts
FROM users u
LEFT JOIN posts p ON u.id = p.user_id
GROUP BY u.name;
```

## 🛠️ 트러블슈팅

### 에러: "Row already exists"

이미 데이터가 있는 경우 발생합니다.

**해결:**
```bash
make clear-data  # 먼저 데이터 삭제
make seed-data   # 다시 삽입
```

### 에러: "Foreign key constraint violation"

users를 먼저 삽입하지 않고 posts를 삽입한 경우.

**해결:** DML 파일이 올바른 순서로 실행되는지 확인 (001_users → 002_posts)

## 📝 참고

- DML은 DDL(테이블 생성)과 달리 **데이터**를 다룹니다
- 개발/테스트 환경에서만 사용하세요
- 운영 환경에서는 신중하게 사용해야 합니다

