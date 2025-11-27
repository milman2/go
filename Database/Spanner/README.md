# Google Cloud Spanner + yo 예제

Google Cloud Spanner Emulator와 yo를 사용한 Clean Architecture 예제입니다.

## 🎯 yo란?

**yo** = Code generator for Google Cloud Spanner

Mercari에서 만든 Cloud Spanner용 코드 생성 도구로, 데이터베이스 스키마에서 Go 코드를 자동 생성합니다.

- GitHub: https://github.com/mercari/yo
- pkg.go.dev: https://pkg.go.dev/go.mercari.io/yo

## 🏗️ 프로젝트 구조

```
spanner/
├── schema/                      # DDL 스키마 정의 (단일 진실 공급원)
│   └── schema.sql              # 전체 스키마 (테이블, 인덱스)
│
├── models/                      # yo가 생성하는 코드
│   ├── user.yo.go              # User 모델 (자동 생성)
│   ├── post.yo.go              # Post 모델 (자동 생성)
│   └── yo_db.yo.go             # DB 헬퍼 (자동 생성)
│
├── ext/                         # 외부 도구 빌드 설정
│   ├── hammer.go               # hammer 빌드
│   ├── wrench.go               # wrench 빌드
│   └── yo.go                   # yo 빌드
│
├── bin/ext/                     # 빌드된 도구들
│   ├── hammer
│   ├── wrench
│   └── yo
│
├── cmd/api/
│   └── main.go                  # 서버 진입점
│
├── migrations/                  # DML 마이그레이션 (선택사항)
│   └── dml/                    # 샘플 데이터, 마스터 데이터
│
├── docker-compose.yml           # Spanner emulator
├── Makefile                     # 자동화 스크립트
├── go.mod
└── README.md
```

## 🚀 빠른 시작

### 1. 전체 초기화 (한번에)

```bash
make init
```

이 명령어는:
1. ✅ Docker Spanner emulator 시작
2. ✅ 외부 도구 빌드 (hammer, wrench, yo)
3. ✅ Instance 생성
4. ✅ Database 및 스키마 생성 (hammer)
5. ✅ yo로 코드 생성

### 2. Spanner 주요 기능

**DEFAULT 값**:
```sql
published BOOL NOT NULL DEFAULT (false)  -- 괄호 필수!
```

**FOREIGN KEY**:
```sql
-- 기본 지원, CASCADE는 미지원
FOREIGN KEY (user_id) REFERENCES users (id)
```

**INTERLEAVE** (부모-자식 관계):
```sql
-- CASCADE DELETE 지원 + 성능 최적화
CREATE TABLE comments (
  user_id STRING(36) NOT NULL,
  comment_id STRING(36) NOT NULL,
  ...
) PRIMARY KEY (user_id, comment_id),
  INTERLEAVE IN PARENT users ON DELETE CASCADE;
```

### 3. 서버 실행

```bash
make run
```

### 4. API 테스트

```bash
make test
```

## 📋 Makefile 명령어 전체 목록

### 기본 명령어
```bash
make help                # 모든 명령어 보기
make init                # 전체 초기화 (Docker + 도구 빌드 + DB 생성)
make run                 # 서버 실행
make test                # API 테스트
make info                # 설정 정보 보기
```

### Docker 관리
```bash
make docker-up           # Docker 시작
make docker-down         # Docker 중지
make docker-ps           # Docker 상태 확인
```

### 데이터베이스 관리
```bash
make setup-instance      # Instance 생성
make createdb            # Database 생성 (hammer create)
make dropdb              # Database 삭제 (wrench drop)
make resetdb             # Database 리셋 (삭제 후 재생성)
make reset               # DB 리셋 + 코드 재생성
```

### 스키마 관리 (hammer)
```bash
make db-apply            # 스키마 변경사항 적용
make db-diff             # 현재 DB와 스키마 파일 차이 확인
make db-export           # 현재 스키마를 파일로 내보내기
make show-schema         # 현재 스키마 확인
```

### DML 관리 (wrench)
```bash
make migrate-dml         # DML 마이그레이션 실행
```

### 코드 생성 (yo)
```bash
make generate-models     # Go 코드 생성
make build/ext           # 외부 도구 빌드
make clean               # 생성된 파일 삭제
```

### 개발 도구
```bash
make spanner-cli         # Spanner CLI 접속
make test-connection     # 연결 테스트
make test-tables         # 테이블 정보 조회
make test-crud           # CRUD 테스트
make test-all            # 종합 테스트
make sql SQL="..."       # SQL 직접 실행
```

## 🔧 도구 역할 구분

### Hammer - DDL 스키마 관리

**용도**: 데이터베이스 스키마(테이블, 인덱스) 관리

| 명령어 | 설명 |
|--------|------|
| `hammer create` | 스키마 파일로부터 데이터베이스 생성 |
| `hammer apply` | 스키마 변경사항을 기존 DB에 적용 |
| `hammer diff` | 현재 DB 스키마와 파일의 차이점 확인 |
| `hammer export` | 현재 DB 스키마를 SQL 파일로 내보내기 |

**사용 예시**:
```bash
# 데이터베이스 생성
hammer create spanner://projects/test-project/instances/test-instance/databases/test-db schema/schema.sql

# 스키마 변경사항 적용
hammer apply spanner://projects/test-project/instances/test-instance/databases/test-db schema/schema.sql

# 스키마 차이 확인
hammer diff spanner://projects/test-project/instances/test-instance/databases/test-db schema/schema.sql
```

### Wrench - 데이터베이스 관리

**용도**: 데이터베이스 삭제 및 DML(데이터 조작어) 실행

| 명령어 | 설명 |
|--------|------|
| `wrench drop` | 데이터베이스 완전 삭제 |
| `wrench apply --dml` | DML 파일 실행 (INSERT, UPDATE, DELETE) |

**사용 예시**:
```bash
# 데이터베이스 삭제
wrench drop --project test-project --instance test-instance --database test-db

# DML 실행
wrench apply --dml migrations/dml/sample_data.sql
```

### yo - Go 코드 생성

**용도**: Spanner 데이터베이스 스키마로부터 Go 코드 자동 생성

**사용 예시**:
```bash
# 모델 코드 생성
yo test-project test-instance test-db -o models -p models --ignore-tables SchemaMigrations
```

## 🔨 yo 코드 생성

### 기본 사용법

```bash
SPANNER_EMULATOR_HOST=localhost:9010 \
yo PROJECT_NAME INSTANCE_NAME DATABASE_NAME -o models -p models
```

### 생성되는 코드

#### 1. Struct (모델)

```go
// models/user.yo.go
type User struct {
    ID        string    `spanner:"id" json:"id"`
    Email     string    `spanner:"email" json:"email"`
    Name      string    `spanner:"name" json:"name"`
    CreatedAt time.Time `spanner:"created_at" json:"created_at"`
    UpdatedAt time.Time `spanner:"updated_at" json:"updated_at"`
}
```

#### 2. Mutation Methods (INSERT/UPDATE/DELETE)

```go
// Insert
func (u *User) Insert(ctx context.Context, db YODB) error

// Update
func (u *User) Update(ctx context.Context, db YODB) error

// InsertOrUpdate
func (u *User) InsertOrUpdate(ctx context.Context, db YODB) error

// UpdateColumns (특정 컬럼만)
func (u *User) UpdateColumns(ctx context.Context, db YODB, columns ...string) error

// Delete
func (u *User) Delete(ctx context.Context, db YODB) error
```

#### 3. Read Functions (인덱스 기반)

```go
// Primary Key로 조회
func FindUserByID(ctx context.Context, db YODB, id string) (*User, error)

// Unique Index로 조회
func FindUserByEmail(ctx context.Context, db YODB, email string) (*User, error)

// 전체 조회
func FindAllUsers(ctx context.Context, db YODB) ([]*User, error)
```

## 📝 사용 예제

### 생성된 코드 사용하기

```go
import (
    "context"
    "github.com/milman2/go-api/spanner-yo/models"
    "cloud.google.com/go/spanner"
)

func main() {
    ctx := context.Background()
    
    // Spanner 클라이언트 생성
    client, _ := spanner.NewClient(ctx, 
        "projects/test-project/instances/test-instance/databases/test-db")
    defer client.Close()
    
    // 사용자 생성
    user := &models.User{
        ID:    uuid.New().String(),
        Email: "test@example.com",
        Name:  "Test User",
    }
    
    // INSERT
    _, err := client.Apply(ctx, []*spanner.Mutation{
        user.Insert(ctx),
    })
    
    // 조회
    user, err := models.FindUserByEmail(ctx, client, "test@example.com")
    
    // 수정
    user.Name = "Updated Name"
    _, err = client.Apply(ctx, []*spanner.Mutation{
        user.Update(ctx),
    })
    
    // 삭제
    _, err = client.Apply(ctx, []*spanner.Mutation{
        user.Delete(ctx),
    })
}
```

## 🐳 Docker Spanner Emulator

### 현재 실행 중인 Spanner

사용자 시스템에 이미 실행 중입니다:
```
Container: school-live-api-spanner-1
Ports: 
  - 9010 (gRPC)
  - 9020 (HTTP)
```

### 새로 실행하려면

```bash
# docker-compose.yml 사용
make docker-up

# 또는 직접
docker run -d -p 9010:9010 -p 9020:9020 \
  gcr.io/cloud-spanner-emulator/emulator:1.5.33
```

## 🎯 워크플로우

### 스키마 변경 워크플로우

```
1. 스키마 파일 수정
   schema/schema.sql 편집
   
2. 변경사항 확인
   make db-diff
   
3. 변경사항 적용
   make db-apply
   
4. Go 코드 재생성
   make generate-models
   
5. 생성된 모델 사용
   import "project/models"
   user := &models.User{...}
```

### 전체 리셋

```bash
# DB 완전 리셋 + 코드 재생성
make reset
```

### 샘플 데이터 추가

```bash
# 1. migrations/dml/ 디렉토리에 SQL 파일 추가
# 2. DML 마이그레이션 실행
make migrate-dml
```

## ✨ yo의 장점

### 1. 타입 안전
- 스키마에서 직접 생성 → 타입 불일치 없음
- 컴파일 타임 체크

### 2. 보일러플레이트 제거
- CRUD 메서드 자동 생성
- 인덱스 기반 조회 함수 자동 생성

### 3. Spanner 최적화
- Mutation API 활용
- 인덱스 기반 효율적 조회

### 4. 일관성
- 모든 테이블에 동일한 패턴
- 유지보수 용이

## 📝 스키마 관리 Best Practices

### 단일 진실 공급원 (Single Source of Truth)

`schema/schema.sql` 파일이 전체 스키마의 유일한 소스입니다.

```sql
-- schema/schema.sql
-- Spanner 주요 기능:
-- 1. DEFAULT 값: DEFAULT (값) 형식으로 괄호 필수
-- 2. FOREIGN KEY: 기본 지원 (CASCADE 미지원)
-- 3. INTERLEAVE: 부모-자식 관계 + CASCADE DELETE 지원 + 성능 최적화

CREATE TABLE users (
  id STRING(36) NOT NULL,
  email STRING(255) NOT NULL,
  name STRING(100) NOT NULL,
  created_at TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true),
  updated_at TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true),
) PRIMARY KEY (id);

CREATE UNIQUE INDEX users_email_idx ON users(email);
```

### 스키마 변경 절차

1. **로컬에서 변경사항 확인**
   ```bash
   make db-diff
   ```

2. **변경사항 적용**
   ```bash
   make db-apply
   ```

3. **코드 재생성**
   ```bash
   make generate-models
   ```

## 📚 추가 학습 자료

- [yo GitHub](https://github.com/mercari/yo)
- [Cloud Spanner 문서](https://cloud.google.com/spanner/docs)
- [Hammer GitHub](https://github.com/daichirata/hammer)
- [Wrench GitHub](https://github.com/cloudspannerecosystem/wrench)

## 🎓 다음 단계

1. **스키마 확장**: `schema/schema.sql`에 새 테이블 추가
2. **코드 재생성**: `make generate-models`
3. **Clean Architecture 적용**: Repository 레이어에서 yo 모델 사용
4. **관계 활용**: Foreign Key, Index, INTERLEAVE 활용
5. **트랜잭션**: Spanner의 강력한 트랜잭션 기능 활용
6. **샘플 데이터**: `migrations/dml/` 디렉토리에 샘플 데이터 추가

## 🧪 Spanner 테스트

### 빠른 테스트

```bash
# 연결 테스트
make test-connection

# 테이블 정보 조회
make test-tables

# CRUD 테스트
make test-crud

# 종합 테스트 (모두 실행)
make test-all
```

### SQL 직접 실행

```bash
# 단일 SQL 실행
make sql SQL="SELECT * FROM users"

# 카운트 조회
make sql SQL="SELECT COUNT(*) FROM users"

# 조건부 조회
make sql SQL="SELECT * FROM users WHERE email LIKE '%@example.com'"
```

### Spanner CLI 사용

```bash
# CLI 접속
make spanner-cli

# CLI에서 사용
spanner> SELECT * FROM users;
spanner> \d users  # 테이블 정의
spanner> \q        # 종료
```

### 상세 가이드

테스트에 대한 자세한 내용은 다음 문서를 참고하세요:
- **SPANNER.md**: 완벽한 테스트 가이드
- **USAGE.md**: 사용법 및 예제
- **test_*.go**: Go 테스트 스크립트
- **test_all.sh**: 종합 테스트 스크립트

## 🎉 결론

**yo + Spanner = 타입 안전 + 자동 코드 생성**

Clean Architecture와 결합하면:
- ✅ Domain은 순수하게
- ✅ yo 모델은 Repository 레이어에
- ✅ 쉽게 테스트 가능
- ✅ 비즈니스 로직 보호

Happy Coding with Spanner! 🚀

