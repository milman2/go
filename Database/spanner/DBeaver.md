# 🔌 DBeaver로 Spanner Emulator 연결하기

## ⚠️ 중요 사항

**Spanner Emulator는 표준 JDBC/ODBC를 지원하지 않습니다.**

DBeaver는 일반적으로 JDBC 드라이버를 사용하여 데이터베이스에 연결하지만, Spanner Emulator는:
- ✅ gRPC API 지원 (포트 9010)
- ✅ HTTP REST API 지원 (포트 9020)
- ❌ JDBC 드라이버 미지원 (Emulator)

## 🎯 대안

### 1. Spanner CLI (권장)

**가장 쉬운 방법:**
```bash
# Docker Compose 사용
make spanner-cli

# 또는 직접
docker-compose exec spanner-cli spanner-cli \
  -p test-project \
  -i test-instance \
  -d test-database
```

**사용 예:**
```sql
spanner> SHOW TABLES;
spanner> SELECT * FROM users;
spanner> \d users  -- 테이블 정의 보기
```

### 2. gcloud CLI

```bash
# SQL 실행
gcloud spanner databases execute-sql test-database \
  --instance=test-instance \
  --sql="SELECT * FROM users"

# DDL 확인
gcloud spanner databases ddl describe test-database \
  --instance=test-instance
```

### 3. Go 스크립트 (이 프로젝트)

```bash
# 연결 테스트
go run test_connection.go

# 테이블 정보
go run test_tables.go

# CRUD 테스트
go run test_crud.go

# 종합 테스트
./test_all.sh
```

### 4. 웹 UI (비공식)

Spanner Emulator는 기본적으로 웹 UI를 제공하지 않지만, 커스텀 UI를 만들 수 있습니다:

```bash
# HTTP API 사용 예
curl http://localhost:9020/v1/projects/test-project/instances
```

## 🔧 실제 Cloud Spanner + DBeaver

**실제 Cloud Spanner (프로덕션)를 사용한다면:**

### 1. JDBC 드라이버 설치

1. DBeaver 설치
2. `Database` → `Driver Manager`
3. `New` 클릭
4. Driver 정보 입력:
   - Driver Name: `Cloud Spanner`
   - Class Name: `com.google.cloud.spanner.jdbc.JdbcDriver`
   - URL Template: `jdbc:cloudspanner:/projects/{project}/instances/{instance}/databases/{database}`

### 2. Maven Dependency

```xml
<dependency>
    <groupId>com.google.cloud</groupId>
    <artifactId>google-cloud-spanner-jdbc</artifactId>
    <version>2.14.0</version>
</dependency>
```

### 3. 연결

- URL: `jdbc:cloudspanner:/projects/my-project/instances/my-instance/databases/my-database`
- 인증: Service Account JSON 키 사용

## 📊 비교표

| 도구 | Emulator | 실제 Spanner | 설명 |
|------|----------|--------------|------|
| **Spanner CLI** | ✅ | ✅ | 공식 CLI 도구 |
| **gcloud** | ✅ | ✅ | Google Cloud CLI |
| **DBeaver (JDBC)** | ❌ | ✅ | GUI 도구 (실제 Spanner만) |
| **Go Client** | ✅ | ✅ | 프로그래밍 방식 |
| **HTTP API** | ✅ | ✅ | REST API |

## 🎯 결론

### Emulator 사용 시 (개발):
```bash
# 가장 편한 방법
make spanner-cli

# 또는 Go 스크립트
go run test_tables.go
go run test_crud.go
```

### 실제 Spanner 사용 시 (프로덕션):
- ✅ DBeaver 사용 가능 (JDBC 드라이버)
- ✅ DataGrip 사용 가능
- ✅ GUI 도구 모두 사용 가능

## 🔍 Spanner CLI 사용 팁

```bash
# 연결
make spanner-cli

# 메타 명령어
spanner> \h        # 도움말
spanner> \d        # 테이블 목록
spanner> \d users  # 테이블 정의
spanner> \q        # 종료

# SQL 실행
spanner> SELECT * FROM users;
spanner> INSERT INTO users (id, email, name, created_at, updated_at)
         VALUES ('test', 'test@example.com', 'Test', 
                 CURRENT_TIMESTAMP(), CURRENT_TIMESTAMP());

# 트랜잭션
spanner> BEGIN;
spanner> INSERT INTO users ...;
spanner> COMMIT;
```

## 💡 추천 워크플로우

```
개발 중:
  1. Spanner CLI로 빠른 쿼리
  2. Go 스크립트로 테스트
  3. gcloud로 DDL 관리

프로덕션:
  1. DBeaver로 데이터 조회
  2. Go 앱으로 프로덕션 작업
  3. 모니터링 도구로 성능 확인
```

Happy Querying! 🚀
