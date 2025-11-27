# 🧪 Spanner 테스트 빠른 가이드

## ⚡ 빠른 시작

```bash
# 1. 연결 테스트
make test-connection

# 2. 테이블 정보
make test-tables

# 3. CRUD 테스트
make test-crud

# 4. 종합 테스트 (모두 실행)
make test-all
```

## 📋 테스트 명령어 모음

### Go 스크립트 테스트

| 명령어 | 설명 | 파일 |
|--------|------|------|
| `make test-connection` | Spanner 연결 테스트 | test_connection.go |
| `make test-tables` | 테이블/컬럼/인덱스 정보 조회 | test_tables.go |
| `make test-crud` | CREATE/READ/UPDATE/DELETE 테스트 | test_crud.go |
| `make test-all` | 위 모든 테스트 + gcloud 테스트 | test_all.sh |

### SQL 직접 실행

```bash
# 기본 조회
make sql SQL="SELECT * FROM users"

# COUNT
make sql SQL="SELECT COUNT(*) as total FROM users"

# WHERE 조건
make sql SQL="SELECT * FROM users WHERE email LIKE '%@example.com'"

# JOIN
make sql SQL="SELECT u.name, COUNT(p.id) as posts FROM users u LEFT JOIN posts p ON u.id = p.user_id GROUP BY u.id, u.name"
```

### Spanner CLI

```bash
# CLI 접속
make spanner-cli

# CLI 명령어
spanner> SHOW TABLES;
spanner> SELECT * FROM users;
spanner> \d users          # 테이블 정의
spanner> \h                # 도움말
spanner> \q                # 종료
```

### gcloud CLI

```bash
# Instance 목록
gcloud spanner instances list

# Database 목록
gcloud spanner databases list --instance=test-instance

# DDL 조회
gcloud spanner databases ddl describe test-db \
  --instance=test-instance

# SQL 실행
gcloud spanner databases execute-sql test-db \
  --instance=test-instance \
  --sql="SELECT * FROM users LIMIT 5"
```

## 🎯 테스트 시나리오

### 시나리오 1: 처음 설치 후

```bash
# 1. 전체 초기화
make init

# 2. 연결 확인
make test-connection

# 3. 테이블 확인
make test-tables

# 4. API 서버 테스트
make run  # 다른 터미널
make test # 이 터미널
```

### 시나리오 2: 마이그레이션 후

```bash
# 1. 마이그레이션
make migrate-up-wrench

# 2. yo 코드 생성
make generate-yo

# 3. 테이블 정보 확인
make test-tables

# 4. 수동 데이터 삽입 테스트
make sql SQL="INSERT INTO users (id, email, name, created_at, updated_at) VALUES ('test-1', 'test@example.com', 'Test', CURRENT_TIMESTAMP(), CURRENT_TIMESTAMP())"

# 5. 조회 확인
make sql SQL="SELECT * FROM users WHERE id='test-1'"
```

### 시나리오 3: 개발 중 디버깅

```bash
# 1. 스키마 확인
make show-schema

# 2. 데이터 확인
make test-tables
make sql SQL="SELECT * FROM users"

# 3. 특정 테이블 조회
make spanner-cli
# spanner> SELECT * FROM users WHERE ...

# 4. CRUD 동작 확인
make test-crud
```

## 📊 출력 예제

### test-connection 출력

```
✅ Spanner 연결 성공!
━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Database: projects/test-project/instances/test-instance/databases/test-db
━━━━━━━━━━━━━━━━━━━━━━━━━━━━

🔍 연결 테스트 쿼리 실행...
✅ 쿼리 결과: test=1, message='Hello Spanner'

🎉 테스트 완료!
```

### test-tables 출력

```
📊 Spanner 데이터베이스 정보
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

1️⃣ 테이블 목록:
   ├─ posts
   ├─ users
   └─ 총 2개의 테이블

2️⃣ 인덱스 목록:
   ├─ posts.posts_published_idx (INDEX)
   ├─ posts.posts_user_id_idx (INDEX)
   ├─ users.users_email_idx (INDEX) [UNIQUE]
   └─ 총 3개의 인덱스

3️⃣ 테이블 상세 정보:
   📋 테이블: posts
      ├─ id                      STRING(36)
      ├─ user_id                 STRING(36)
      ├─ title                   STRING(200)
      ├─ content                 STRING(MAX) (nullable)
      ├─ published               BOOL
      ├─ created_at              TIMESTAMP
      ├─ updated_at              TIMESTAMP
      └─ 7개 컬럼

   📋 테이블: users
      ├─ id                      STRING(36)
      ├─ email                   STRING(255)
      ├─ name                    STRING(100)
      ├─ created_at              TIMESTAMP
      ├─ updated_at              TIMESTAMP
      └─ 5개 컬럼

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🎉 조회 완료!
```

### test-crud 출력

```
🧪 Spanner CRUD 테스트
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

1️⃣ CREATE 테스트
✅ 사용자 생성 성공: ID=abc-123, Email=test-1234@example.com

2️⃣ READ 테스트 (Key 기반)
✅ 조회 성공:
   - ID: abc-123
   - Email: test-1234@example.com
   - Name: 테스트 사용자
   - Created: 2024-11-27T10:53:00Z

3️⃣ READ 테스트 (Query)
✅ Query 결과: 테스트 사용자 (test-1234@example.com)
   총 1건 조회

4️⃣ UPDATE 테스트
✅ 사용자 수정 성공: 새 이름='수정된 테스트 사용자'
✅ 수정 확인: 수정된 테스트 사용자

5️⃣ DELETE 테스트
✅ 사용자 삭제 성공: ID=abc-123
✅ 삭제 확인: 사용자가 존재하지 않음 (정상)

6️⃣ BATCH CREATE 테스트
✅ 3명의 사용자 일괄 생성 성공

7️⃣ 전체 조회 테스트
   최근 사용자 5명:
   1. 배치 사용자 3 (batch-2@example.com)
   2. 배치 사용자 2 (batch-1@example.com)
   3. 배치 사용자 1 (batch-0@example.com)

8️⃣ 정리 (테스트 데이터 삭제)
✅ 테스트 데이터 정리 완료 (3건)

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🎉 CRUD 테스트 완료!
```

## 🔍 문제 해결

### "database not found" 에러

```bash
# 먼저 초기화 필요
make init
```

### "connection refused" 에러

```bash
# Docker 확인
docker ps | grep spanner

# 없다면 시작
make docker-up
```

### 테스트 파일이 없다는 에러

```bash
# 현재 디렉토리 확인
pwd
# /home/milman2/go-api/go/Database/spanner 이어야 함

# 파일 확인
ls test_*.go
```

## 📚 추가 문서

- **SPANNER.md**: 완벽한 테스트 가이드 (모든 기능)
- **USAGE.md**: 사용법 및 워크플로우
- **YO_GUIDE.md**: yo 사용법
- **DOCKER_GUIDE.md**: Docker 설정

## 🎉 정리

```bash
# 가장 많이 사용하는 명령어
make test-connection  # 연결 확인
make test-tables      # 스키마 확인
make test-crud        # 동작 확인
make test-all         # 종합 확인
```

Happy Testing! 🧪🚀

