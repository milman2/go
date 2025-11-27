# 🔍 Spanner Emulator 연결 및 테스트 가이드

## 📋 목차

1. [현재 실행 중인 Spanner 확인](#현재-실행-중인-spanner-확인)
2. [gcloud CLI로 연결](#gcloud-cli로-연결)
3. [Spanner CLI로 연결](#spanner-cli로-연결)
4. [Go 코드로 직접 연결](#go-코드로-직접-연결)
5. [데이터베이스 정보 조회](#데이터베이스-정보-조회)
6. [CRUD 작업 테스트](#crud-작업-테스트)
7. [트랜잭션 테스트](#트랜잭션-테스트)
8. [인덱스 활용 테스트](#인덱스-활용-테스트)
9. [고급 기능 테스트](#고급-기능-테스트)

---

## 현재 실행 중인 Spanner 확인

### Docker 컨테이너 상태

```bash
# Spanner 컨테이너 확인
docker ps | grep spanner
```

**출력:**
```
school-live-api-spanner-1   gcr.io/cloud-spanner-emulator/emulator:1.5.33
Up 2 months   0.0.0.0:9010->9010/tcp, 0.0.0.0:9020->9020/tcp
```

### 포트 확인

```bash
# gRPC 포트 (9010)
lsof -i :9010

# HTTP 포트 (9020)
lsof -i :9020
```

### 연결 테스트

```bash
# HTTP 엔드포인트 테스트
curl http://localhost:9020

# 정상 응답이 오면 OK

# 러브 라이브 local spanner emulator의 경우
curl http://localhost:9020/v1/projects/school-live-local/instances/school-app-instance/databases/school-app/sessions
```

---

## gcloud CLI로 연결

### 1. gcloud 설정 (Emulator용)

```bash
# 현재 설정 백업 (선택사항)
gcloud config configurations create spanner-emulator

# Emulator 설정
gcloud config set auth/disable_credentials true
gcloud config set project test-project
gcloud config set api_endpoint_overrides/spanner http://localhost:9020/

# 환경 변수 설정
export SPANNER_EMULATOR_HOST=localhost:9020
```

### 2. Instance 목록 조회

```bash
# 모든 Instance 조회
gcloud spanner instances list

# 출력 예:
# NAME           DISPLAY_NAME    CONFIG              NODE_COUNT  STATE
# test-instance  Test Instance   emulator-config     1           READY
```

### 3. Instance 상세 정보

```bash
# Instance 정보 조회
gcloud spanner instances describe test-instance

# JSON 형식으로 출력
gcloud spanner instances describe test-instance --format=json
```

### 4. Database 목록 조회

```bash
# 특정 Instance의 Database 목록
gcloud spanner databases list --instance=test-instance

# 출력 예:
# NAME          STATE
# test-db READY
```

### 5. Database 상세 정보

```bash
# Database DDL 조회
gcloud spanner databases ddl describe test-db \
  --instance=test-instance

# 출력: 모든 CREATE TABLE, CREATE INDEX 문
```

### 6. 새 Instance/Database 생성

```bash
# Instance 생성
gcloud spanner instances create my-instance \
  --config=emulator-config \
  --description="My Test Instance" \
  --nodes=1

# Database 생성
gcloud spanner databases create my-database \
  --instance=my-instance

# Database에 DDL 적용
gcloud spanner databases ddl update my-database \
  --instance=my-instance \
  --ddl='CREATE TABLE users (id STRING(36) NOT NULL, name STRING(100)) PRIMARY KEY (id)'
```

---

## Spanner CLI로 연결

### 1. Docker로 Spanner CLI 실행

```bash
# 기존 docker-compose 사용
docker-compose up -d spanner-cli

# CLI 접속
docker-compose exec spanner-cli spanner-cli \
  -p test-project \
  -i test-instance \
  -d test-db
```

**또는 Makefile 사용:**
```bash
make spanner-cli
```

### 2. 스키마 조회

```sql
-- 모든 테이블 보기
SHOW TABLES;

-- 테이블 정의 보기
SHOW CREATE TABLE users;

-- 인덱스 보기
SHOW INDEXES FROM users;
```

### 3. 데이터 조회

```sql
-- 전체 조회
SELECT * FROM users;

-- 조건부 조회
SELECT * FROM users WHERE email = 'test@example.com';

-- COUNT
SELECT COUNT(*) as total FROM users;

-- JOIN
SELECT u.*, COUNT(p.id) as post_count
FROM users u
LEFT JOIN posts p ON p.user_id = u.id
GROUP BY u.id, u.email, u.name, u.created_at, u.updated_at;
```

### 4. 데이터 삽입

```sql
-- 단일 INSERT (Spanner는 GENERATE_UUID 대신 UUID() 사용)
INSERT INTO users (id, email, name, created_at, updated_at)
VALUES (
  'user-001',
  'alice@example.com',
  'Alice',
  CURRENT_TIMESTAMP(),
  CURRENT_TIMESTAMP()
);

-- Posts 테이블 (DEFAULT 값 활용)
INSERT INTO posts (id, user_id, title, content, created_at, updated_at)
VALUES (
  'post-001',
  'user-001',
  'My First Post',
  'Hello World',
  -- published는 DEFAULT (false)로 자동 설정됨
  CURRENT_TIMESTAMP(),
  CURRENT_TIMESTAMP()
);

-- 여러 행 INSERT
INSERT INTO users (id, email, name, created_at, updated_at)
VALUES 
  ('user-002', 'bob@example.com', 'Bob', CURRENT_TIMESTAMP(), CURRENT_TIMESTAMP()),
  ('user-003', 'charlie@example.com', 'Charlie', CURRENT_TIMESTAMP(), CURRENT_TIMESTAMP());
```

### 5. 데이터 수정

```sql
-- UPDATE
UPDATE users
SET name = 'Alice Updated', updated_at = CURRENT_TIMESTAMP()
WHERE id = 'user-001';

-- 조건부 UPDATE
UPDATE users
SET name = CONCAT(name, ' (Verified)')
WHERE email LIKE '%@example.com';
```

### 6. 데이터 삭제

```sql
-- 단일 DELETE
DELETE FROM users WHERE id = 'user-001';

-- 조건부 DELETE
DELETE FROM users WHERE created_at < TIMESTAMP('2024-01-01');
```

---

## Go 코드로 직접 연결

### 1. 기본 연결 테스트

```bash
cd /home/milman2/go-api/go/Database/Spanner
```

**test_connection.go** 생성:
```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"
)

func main() {
	ctx := context.Background()
	
	// 환경 변수 설정
	os.Setenv("SPANNER_EMULATOR_HOST", "localhost:9010")
	
	// Spanner 클라이언트 생성
	database := "projects/test-project/instances/test-instance/databases/test-db"
	client, err := spanner.NewClient(ctx, database)
	if err != nil {
		log.Fatalf("연결 실패: %v", err)
	}
	defer client.Close()
	
	fmt.Println("✅ Spanner 연결 성공!")
	fmt.Println("Database:", database)
	
	// 간단한 쿼리 실행
	stmt := spanner.Statement{SQL: `SELECT 1 as test`}
	iter := client.Single().Query(ctx, stmt)
	defer iter.Stop()
	
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		
		var test int64
		if err := row.Columns(&test); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("쿼리 결과: %d\n", test)
	}
}
```

**실행:**
```bash
SPANNER_EMULATOR_HOST=localhost:9010 go run test_connection.go
```

### 2. 테이블 정보 조회

**test_tables.go** 생성:
```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"
)

func main() {
	ctx := context.Background()
	os.Setenv("SPANNER_EMULATOR_HOST", "localhost:9010")
	
	database := "projects/test-project/instances/test-instance/databases/test-db"
	client, err := spanner.NewClient(ctx, database)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	
	// INFORMATION_SCHEMA에서 테이블 목록 조회
	stmt := spanner.Statement{
		SQL: `SELECT table_name, parent_table_name
		      FROM INFORMATION_SCHEMA.TABLES
		      WHERE table_catalog = '' AND table_schema = ''
		      ORDER BY table_name`,
	}
	
	iter := client.Single().Query(ctx, stmt)
	defer iter.Stop()
	
	fmt.Println("📊 테이블 목록:")
	fmt.Println("================")
	
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		
		var tableName, parentTable spanner.NullString
		if err := row.Columns(&tableName, &parentTable); err != nil {
			log.Fatal(err)
		}
		
		fmt.Printf("- %s", tableName.StringVal)
		if parentTable.Valid {
			fmt.Printf(" (parent: %s)", parentTable.StringVal)
		}
		fmt.Println()
	}
}
```

### 3. 컬럼 정보 조회

**test_columns.go** 생성:
```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"
)

func main() {
	ctx := context.Background()
	os.Setenv("SPANNER_EMULATOR_HOST", "localhost:9010")
	
	database := "projects/test-project/instances/test-instance/databases/test-db"
	client, err := spanner.NewClient(ctx, database)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	
	// 테이블 이름 (인자로 받거나 하드코딩)
	tableName := "users"
	
	stmt := spanner.Statement{
		SQL: `SELECT column_name, spanner_type, is_nullable
		      FROM INFORMATION_SCHEMA.COLUMNS
		      WHERE table_name = @tableName
		      ORDER BY ordinal_position`,
		Params: map[string]interface{}{
			"tableName": tableName,
		},
	}
	
	iter := client.Single().Query(ctx, stmt)
	defer iter.Stop()
	
	fmt.Printf("📋 테이블 '%s'의 컬럼:\n", tableName)
	fmt.Println("=====================================")
	
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		
		var columnName, spannerType, isNullable string
		if err := row.Columns(&columnName, &spannerType, &isNullable); err != nil {
			log.Fatal(err)
		}
		
		nullable := ""
		if isNullable == "YES" {
			nullable = " (nullable)"
		}
		fmt.Printf("- %-20s %s%s\n", columnName, spannerType, nullable)
	}
}
```

---

## 데이터베이스 정보 조회

### 1. 모든 테이블 조회

```bash
# gcloud 사용
gcloud spanner databases ddl describe test-db \
  --instance=test-instance
```

### 2. 인덱스 정보 조회

**SQL (Spanner CLI):**
```sql
SELECT 
  index_name,
  table_name,
  index_type,
  is_unique,
  is_null_filtered
FROM INFORMATION_SCHEMA.INDEXES
WHERE table_catalog = '' AND table_schema = ''
ORDER BY table_name, index_name;
```

### 3. Primary Key 조회

```sql
SELECT 
  table_name,
  column_name,
  ordinal_position
FROM INFORMATION_SCHEMA.INDEX_COLUMNS
WHERE index_name = 'PRIMARY_KEY'
ORDER BY table_name, ordinal_position;
```

### 4. Foreign Key 조회

```sql
SELECT
  constraint_name,
  table_name,
  column_name,
  CONCAT(referenced_table_name, '.', referenced_column_name) as references
FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE
WHERE constraint_name LIKE 'FK_%'
ORDER BY table_name;
```

### 5. 통계 정보

```sql
-- 각 테이블의 레코드 수
SELECT 
  'users' as table_name,
  COUNT(*) as row_count
FROM users
UNION ALL
SELECT 
  'posts' as table_name,
  COUNT(*) as row_count
FROM posts;
```

---

## CRUD 작업 테스트

### 1. Create (생성)

#### SQL
```sql
-- 단일 INSERT
INSERT INTO users (id, email, name, created_at, updated_at)
VALUES (
  GENERATE_UUID(),
  'test@example.com',
  'Test User',
  CURRENT_TIMESTAMP(),
  CURRENT_TIMESTAMP()
);
```

#### Go (Mutation API)
```go
func testCreate(client *spanner.Client) {
	ctx := context.Background()
	
	// Mutation 생성
	m := spanner.InsertMap("users", map[string]interface{}{
		"id":         uuid.New().String(),
		"email":      "go@example.com",
		"name":       "Go User",
		"created_at": spanner.CommitTimestamp,
		"updated_at": spanner.CommitTimestamp,
	})
	
	// Apply
	_, err := client.Apply(ctx, []*spanner.Mutation{m})
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Println("✅ 사용자 생성 완료")
}
```

### 2. Read (조회)

#### SQL - 단일 조회
```sql
SELECT * FROM users WHERE id = 'user-001';
```

#### SQL - 전체 조회
```sql
SELECT * FROM users ORDER BY created_at DESC LIMIT 10;
```

#### Go - Query
```go
func testRead(client *spanner.Client) {
	ctx := context.Background()
	
	stmt := spanner.Statement{
		SQL: `SELECT id, email, name FROM users 
		      WHERE email = @email`,
		Params: map[string]interface{}{
			"email": "test@example.com",
		},
	}
	
	iter := client.Single().Query(ctx, stmt)
	defer iter.Stop()
	
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		
		var id, email, name string
		row.Columns(&id, &email, &name)
		fmt.Printf("User: %s (%s)\n", name, email)
	}
}
```

#### Go - Read (Key-based)
```go
func testReadByKey(client *spanner.Client, userID string) {
	ctx := context.Background()
	
	row, err := client.Single().ReadRow(ctx, "users",
		spanner.Key{userID},
		[]string{"id", "email", "name"})
	if err != nil {
		log.Fatal(err)
	}
	
	var id, email, name string
	row.Columns(&id, &email, &name)
	fmt.Printf("User: %s (%s)\n", name, email)
}
```

### 3. Update (수정)

#### SQL
```sql
UPDATE users
SET name = 'Updated Name', updated_at = CURRENT_TIMESTAMP()
WHERE id = 'user-001';
```

#### Go
```go
func testUpdate(client *spanner.Client, userID string) {
	ctx := context.Background()
	
	m := spanner.UpdateMap("users", map[string]interface{}{
		"id":         userID,
		"name":       "Updated via Go",
		"updated_at": spanner.CommitTimestamp,
	})
	
	_, err := client.Apply(ctx, []*spanner.Mutation{m})
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Println("✅ 사용자 수정 완료")
}
```

### 4. Delete (삭제)

#### SQL
```sql
DELETE FROM users WHERE id = 'user-001';
```

#### Go
```go
func testDelete(client *spanner.Client, userID string) {
	ctx := context.Background()
	
	m := spanner.Delete("users", spanner.Key{userID})
	
	_, err := client.Apply(ctx, []*spanner.Mutation{m})
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Println("✅ 사용자 삭제 완료")
}
```

---

## 트랜잭션 테스트

### 1. Read-Write 트랜잭션

```go
func testTransaction(client *spanner.Client) {
	ctx := context.Background()
	
	_, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		// 1. 읽기
		row, err := txn.ReadRow(ctx, "users",
			spanner.Key{"user-001"},
			[]string{"id", "name"})
		if err != nil {
			return err
		}
		
		var id, name string
		row.Columns(&id, &name)
		
		// 2. 수정
		newName := name + " (Updated)"
		m := spanner.UpdateMap("users", map[string]interface{}{
			"id":         id,
			"name":       newName,
			"updated_at": spanner.CommitTimestamp,
		})
		
		// 3. 버퍼에 추가
		return txn.BufferWrite([]*spanner.Mutation{m})
	})
	
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Println("✅ 트랜잭션 완료")
}
```

### 2. Batch Write (여러 Mutation)

```go
func testBatchWrite(client *spanner.Client) {
	ctx := context.Background()
	
	mutations := []*spanner.Mutation{
		spanner.InsertMap("users", map[string]interface{}{
			"id":         uuid.New().String(),
			"email":      "batch1@example.com",
			"name":       "Batch User 1",
			"created_at": spanner.CommitTimestamp,
			"updated_at": spanner.CommitTimestamp,
		}),
		spanner.InsertMap("users", map[string]interface{}{
			"id":         uuid.New().String(),
			"email":      "batch2@example.com",
			"name":       "Batch User 2",
			"created_at": spanner.CommitTimestamp,
			"updated_at": spanner.CommitTimestamp,
		}),
	}
	
	_, err := client.Apply(ctx, mutations)
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Println("✅ Batch write 완료")
}
```

### 3. 트랜잭션 롤백 테스트

```go
func testRollback(client *spanner.Client) {
	ctx := context.Background()
	
	_, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		// 여러 작업 수행
		m1 := spanner.InsertMap("users", map[string]interface{}{
			"id":         uuid.New().String(),
			"email":      "rollback@example.com",
			"name":       "Rollback Test",
			"created_at": spanner.CommitTimestamp,
			"updated_at": spanner.CommitTimestamp,
		})
		txn.BufferWrite([]*spanner.Mutation{m1})
		
		// 에러 발생 시 자동 롤백
		return fmt.Errorf("의도적인 에러 - 롤백됨")
	})
	
	if err != nil {
		fmt.Println("❌ 트랜잭션 롤백:", err)
	}
}
```

---

## 인덱스 활용 테스트

### 1. 인덱스 스캔 vs Full Scan

```sql
-- 인덱스 사용 (users_email_idx)
SELECT * FROM users WHERE email = 'test@example.com';

-- Full Table Scan
SELECT * FROM users WHERE name LIKE '%test%';
```

### 2. FORCE_INDEX 힌트

```sql
-- 특정 인덱스 강제 사용
SELECT * FROM users@{FORCE_INDEX=users_email_idx}
WHERE email = 'test@example.com';
```

### 3. 인덱스 효율성 테스트

```go
func testIndexPerformance(client *spanner.Client) {
	ctx := context.Background()
	
	// 인덱스 사용
	start := time.Now()
	stmt1 := spanner.Statement{
		SQL: `SELECT * FROM users WHERE email = @email`,
		Params: map[string]interface{}{
			"email": "test@example.com",
		},
	}
	iter1 := client.Single().Query(ctx, stmt1)
	iter1.Stop()
	fmt.Printf("인덱스 사용: %v\n", time.Since(start))
	
	// Full Scan
	start = time.Now()
	stmt2 := spanner.Statement{
		SQL: `SELECT * FROM users WHERE name LIKE '%test%'`,
	}
	iter2 := client.Single().Query(ctx, stmt2)
	iter2.Stop()
	fmt.Printf("Full Scan: %v\n", time.Since(start))
}
```

---

## 고급 기능 테스트

### 0. Spanner 주요 기능

#### DEFAULT 값 (괄호 필수!)

```sql
-- ✅ 올바른 방법
CREATE TABLE posts (
  id STRING(36) NOT NULL,
  published BOOL NOT NULL DEFAULT (false),  -- 괄호 필수!
  view_count INT64 DEFAULT (0),
) PRIMARY KEY (id);

-- ❌ 잘못된 방법
published BOOL NOT NULL DEFAULT false  -- 에러!
```

#### FOREIGN KEY (CASCADE 미지원)

```sql
-- ✅ FOREIGN KEY 지원
CREATE TABLE posts (
  id STRING(36) NOT NULL,
  user_id STRING(36) NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users (id)
) PRIMARY KEY (id);

-- ❌ CASCADE는 미지원
-- FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
-- 대신 INTERLEAVE 사용
```

#### INTERLEAVE (부모-자식 + CASCADE DELETE)

```sql
-- 부모 테이블
CREATE TABLE users (
  id STRING(36) NOT NULL,
  name STRING(100),
) PRIMARY KEY (id);

-- 자식 테이블 (INTERLEAVE + CASCADE DELETE)
CREATE TABLE user_profiles (
  user_id STRING(36) NOT NULL,
  profile_id STRING(36) NOT NULL,
  bio STRING(MAX),
) PRIMARY KEY (user_id, profile_id),
  INTERLEAVE IN PARENT users ON DELETE CASCADE;

-- 장점:
-- 1. 부모 삭제 시 자식도 자동 삭제 (CASCADE DELETE)
-- 2. 같은 물리적 위치에 저장되어 성능 최적화
-- 3. 부모-자식 조인 쿼리가 매우 빠름
```

### 1. PENDING_COMMIT_TIMESTAMP

```sql
-- 커밋 타임스탬프 자동 설정
INSERT INTO users (id, email, name, created_at, updated_at)
VALUES (
  'user-uuid',
  'timestamp@example.com',
  'Timestamp Test',
  PENDING_COMMIT_TIMESTAMP(),
  PENDING_COMMIT_TIMESTAMP()
);
```

### 2. ARRAY 타입

```sql
-- ARRAY 타입 테이블 생성 (마이그레이션)
CREATE TABLE tags (
  id STRING(36) NOT NULL,
  names ARRAY<STRING(50)>,
) PRIMARY KEY (id);

-- INSERT
INSERT INTO tags (id, names)
VALUES ('tag-001', ['go', 'spanner', 'database']);

-- 조회
SELECT id, names FROM tags WHERE 'go' IN UNNEST(names);
```

### 3. JSON 타입 (Spanner 지원 시)

```sql
-- JSON 컬럼
CREATE TABLE settings (
  id STRING(36) NOT NULL,
  config JSON,
) PRIMARY KEY (id);

-- INSERT
INSERT INTO settings (id, config)
VALUES ('setting-001', JSON '{"theme": "dark", "lang": "ko"}');

-- JSON 쿼리
SELECT id, JSON_VALUE(config, '$.theme') as theme
FROM settings;
```

### 4. Interleaved Tables (부모-자식) - 상세 예제

```sql
-- 부모 테이블
CREATE TABLE authors (
  author_id STRING(36) NOT NULL,
  name STRING(100),
) PRIMARY KEY (author_id);

-- 자식 테이블 (Interleaved)
CREATE TABLE books (
  author_id STRING(36) NOT NULL,
  book_id STRING(36) NOT NULL,
  title STRING(200),
) PRIMARY KEY (author_id, book_id),
  INTERLEAVE IN PARENT authors ON DELETE CASCADE;

-- 테스트
INSERT INTO authors (author_id, name) VALUES ('author-1', 'Alice');
INSERT INTO books (author_id, book_id, title) VALUES ('author-1', 'book-1', 'Book 1');
INSERT INTO books (author_id, book_id, title) VALUES ('author-1', 'book-2', 'Book 2');

-- 부모 삭제 시 자식도 함께 삭제됨
DELETE FROM authors WHERE author_id = 'author-1';
-- books의 두 레코드도 자동 삭제!

-- INTERLEAVE vs FOREIGN KEY:
-- FOREIGN KEY: 참조 무결성만 보장, CASCADE 미지원
-- INTERLEAVE: CASCADE DELETE 지원 + 성능 최적화 (같은 물리 저장소)
```

### 5. Partitioned DML

```go
func testPartitionedDML(client *spanner.Client) {
	ctx := context.Background()
	
	// 대량 업데이트
	stmt := spanner.Statement{
		SQL: `UPDATE users SET name = CONCAT(name, ' (Bulk)')
		      WHERE created_at < @cutoff`,
		Params: map[string]interface{}{
			"cutoff": time.Now().Add(-24 * time.Hour),
		},
	}
	
	rowCount, err := client.PartitionedUpdate(ctx, stmt)
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("✅ %d rows updated\n", rowCount)
}
```

### 6. Stale Read (오래된 데이터 읽기)

```go
func testStaleRead(client *spanner.Client) {
	ctx := context.Background()
	
	// 15초 전 데이터 읽기
	ro := client.Single().WithTimestampBound(
		spanner.ExactStaleness(15 * time.Second))
	
	stmt := spanner.Statement{SQL: `SELECT * FROM users`}
	iter := ro.Query(ctx, stmt)
	defer iter.Stop()
	
	// 결과 처리...
}
```

---

## 종합 테스트 스크립트

**test_all.sh** 생성:

```bash
#!/bin/bash

echo "🧪 Spanner Emulator 종합 테스트"
echo "=================================="

# 환경 변수
export SPANNER_EMULATOR_HOST=localhost:9010
export PROJECT=test-project
export INSTANCE=test-instance
export DATABASE=test-db

echo ""
echo "1️⃣ 연결 테스트"
curl -s http://localhost:9020 > /dev/null && echo "✅ HTTP OK" || echo "❌ HTTP 실패"

echo ""
echo "2️⃣ Instance 목록"
gcloud spanner instances list

echo ""
echo "3️⃣ Database 목록"
gcloud spanner databases list --instance=$INSTANCE

echo ""
echo "4️⃣ 테이블 정보"
gcloud spanner databases ddl describe $DATABASE --instance=$INSTANCE

echo ""
echo "5️⃣ 데이터 조회 (SQL)"
gcloud spanner databases execute-sql $DATABASE \
  --instance=$INSTANCE \
  --sql="SELECT COUNT(*) as total FROM users"

echo ""
echo "6️⃣ 데이터 삽입 테스트"
gcloud spanner databases execute-sql $DATABASE \
  --instance=$INSTANCE \
  --sql="INSERT INTO users (id, email, name, created_at, updated_at) 
         VALUES ('test-$(date +%s)', 'test@test.com', 'Test', CURRENT_TIMESTAMP(), CURRENT_TIMESTAMP())"

echo ""
echo "7️⃣ Go 코드 테스트"
if [ -f test_connection.go ]; then
  go run test_connection.go
fi

echo ""
echo "=================================="
echo "✅ 테스트 완료"
```

**실행:**
```bash
chmod +x test_all.sh
./test_all.sh
```

---

## 유용한 명령어 모음

### gcloud 명령어

```bash
# Instance 생성
gcloud spanner instances create INSTANCE_NAME \
  --config=emulator-config --nodes=1

# Database 생성
gcloud spanner databases create DATABASE_NAME \
  --instance=INSTANCE_NAME

# DDL 실행
gcloud spanner databases ddl update DATABASE_NAME \
  --instance=INSTANCE_NAME \
  --ddl='CREATE TABLE ...'

# SQL 실행
gcloud spanner databases execute-sql DATABASE_NAME \
  --instance=INSTANCE_NAME \
  --sql='SELECT * FROM users'

# Database 삭제
gcloud spanner databases delete DATABASE_NAME \
  --instance=INSTANCE_NAME

# Instance 삭제
gcloud spanner instances delete INSTANCE_NAME
```

### Spanner CLI 명령어

```sql
-- 메타 명령어
\h                    -- 도움말
\q                    -- 종료
\d                    -- 테이블 목록
\d TABLE_NAME         -- 테이블 정의

-- 트랜잭션
BEGIN;                -- 트랜잭션 시작
COMMIT;               -- 커밋
ROLLBACK;             -- 롤백
```

---

## 문제 해결

### "database not found" 에러

```bash
# Database 목록 확인
gcloud spanner databases list --instance=test-instance

# 없다면 생성
make setup-instance
```

### "permission denied" 에러

```bash
# Emulator 설정 확인
gcloud config get-value auth/disable_credentials
# true 여야 함

# 재설정
gcloud config set auth/disable_credentials true
```

### 연결 타임아웃

```bash
# Docker 상태 확인
docker ps | grep spanner

# 재시작
docker restart CONTAINER_ID
```

---

## Makefile 통합

Makefile에 추가할 명령어:

```makefile
.PHONY: test-connection
test-connection: ## Spanner 연결 테스트
	@SPANNER_EMULATOR_HOST=localhost:9010 go run test_connection.go

.PHONY: test-tables
test-tables: ## 테이블 정보 조회
	@SPANNER_EMULATOR_HOST=localhost:9010 go run test_tables.go

.PHONY: sql
sql: ## SQL 실행 (SQL=<query>)
	@SPANNER_EMULATOR_HOST=localhost:9010 \
		gcloud spanner databases execute-sql $(SPANNER_DATABASE_ID) \
		--instance=$(SPANNER_INSTANCE_ID) \
		--sql='$(SQL)'

.PHONY: test-all
test-all: ## 종합 테스트
	@./test_all.sh
```

**사용:**
```bash
make test-connection
make test-tables
make sql SQL="SELECT * FROM users"
make test-all
```

---

## 다음 단계

1. **실제 데이터로 테스트**: 마이그레이션 실행 후 CRUD 테스트
2. **성능 측정**: 인덱스 유무에 따른 성능 비교
3. **트랜잭션 연습**: 복잡한 비즈니스 로직 구현
4. **yo 코드 활용**: 생성된 모델로 타입 안전 CRUD

Happy Testing! 🚀

