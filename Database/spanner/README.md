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
├── migrations/                  # 마이그레이션 파일
│   ├── 000001_create_users.up.sql
│   ├── 000001_create_users.down.sql
│   ├── 000002_create_posts.up.sql
│   └── 000002_create_posts.down.sql
│
├── models/                      # yo가 생성하는 코드
│   ├── user.yo.go              # User 모델 (자동 생성)
│   ├── post.yo.go              # Post 모델 (자동 생성)
│   └── yo_db.yo.go             # DB 헬퍼 (자동 생성)
│
├── cmd/api/
│   └── main.go                  # 서버 진입점
│
├── docker-compose.yml           # Spanner emulator
├── Makefile                     # 자동화 스크립트
├── go.mod
└── README.md
```

## 🚀 빠른 시작

### 1. 필수 도구 설치

```bash
make install-tools
```

설치되는 도구:
- **yo**: Spanner 코드 생성기
- **hammer**: 마이그레이션 도구 #1
- **wrench**: 마이그레이션 도구 #2

### 2. 전체 초기화 (한번에)

```bash
make init
```

이 명령어는:
1. ✅ Docker Spanner emulator 시작
2. ✅ Instance & Database 생성
3. ✅ 마이그레이션 실행
4. ✅ yo로 코드 생성

### 3. 서버 실행

```bash
make run
```

### 4. API 테스트

```bash
make test
```

## 📋 Makefile 명령어 전체 목록

```bash
make help                # 모든 명령어 보기
make docker-up           # Docker 시작
make docker-down         # Docker 중지
make docker-ps           # Docker 상태 확인
make setup-instance      # Instance/Database 생성
make migrate-up-hammer   # Hammer로 마이그레이션 UP
make migrate-down-hammer # Hammer로 마이그레이션 DOWN
make migrate-up-wrench   # Wrench로 마이그레이션 UP
make migrate-down-wrench # Wrench로 마이그레이션 DOWN
make generate-yo         # yo로 코드 생성
make clean               # 생성된 파일 삭제
make reset               # DB 리셋 & 코드 재생성
make run                 # 서버 실행
make test                # API 테스트
make spanner-cli         # Spanner CLI 접속
make show-schema         # 스키마 확인
make info                # 설정 정보 보기
```

## 🔧 마이그레이션 도구 비교

### Hammer vs Wrench

| 특징 | Hammer | Wrench |
|------|--------|--------|
| **개발사** | daichirata | Google Cloud Spanner Ecosystem |
| **방식** | 파일 기반 | 파일 기반 |
| **설정** | CLI 플래그 | CLI 플래그 |
| **복잡도** | 간단 | 간단 |
| **상태 추적** | ✅ | ✅ |

**결론**: 둘 다 유사, 취향에 따라 선택

### Hammer 사용법

```bash
# UP
SPANNER_EMULATOR_HOST=localhost:9010 \
hammer -p test-project -i test-instance -d test-database \
  -m migrations up

# DOWN
hammer -p test-project -i test-instance -d test-database \
  -m migrations down
```

### Wrench 사용법

```bash
# UP
SPANNER_EMULATOR_HOST=localhost:9010 \
wrench migrate up \
  --directory migrations \
  --database projects/test-project/instances/test-instance/databases/test-database

# DOWN
wrench migrate down \
  --directory migrations \
  --database projects/test-project/instances/test-instance/databases/test-database
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
        "projects/test-project/instances/test-instance/databases/test-database")
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

### 일반적인 개발 흐름

```
1. 마이그레이션 파일 작성
   migrations/000003_add_column.up.sql
   
2. 마이그레이션 실행
   make migrate-up-wrench
   
3. yo로 코드 재생성
   make generate-yo
   
4. 생성된 모델 사용
   import "project/models"
   user := &models.User{...}
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

## 🔄 마이그레이션 워크플로우

### Hammer 사용

```bash
# 마이그레이션 실행
make migrate-up-hammer

# 롤백
make migrate-down-hammer
```

### Wrench 사용 (권장)

```bash
# 마이그레이션 실행
make migrate-up-wrench

# 롤백
make migrate-down-wrench

# 상태 확인
make show-schema
```

## 📚 추가 학습 자료

- [yo GitHub](https://github.com/mercari/yo)
- [Cloud Spanner 문서](https://cloud.google.com/spanner/docs)
- [Hammer GitHub](https://github.com/daichirata/hammer)
- [Wrench GitHub](https://github.com/cloudspannerecosystem/wrench)

## 🎓 다음 단계

1. **마이그레이션 추가**: `migrations/` 디렉토리에 새 파일 추가
2. **코드 재생성**: `make generate-yo`
3. **Clean Architecture 적용**: Repository 레이어에서 yo 모델 사용
4. **관계 추가**: Foreign Key 및 인덱스 활용
5. **트랜잭션**: Spanner의 강력한 트랜잭션 기능 활용

## 🎉 결론

**yo + Spanner = 타입 안전 + 자동 코드 생성**

Clean Architecture와 결합하면:
- ✅ Domain은 순수하게
- ✅ yo 모델은 Repository 레이어에
- ✅ 쉽게 테스트 가능
- ✅ 비즈니스 로직 보호

Happy Coding with Spanner! 🚀

