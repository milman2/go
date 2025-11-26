# Clean Architecture 예제

Robert C. Martin(Uncle Bob)의 Clean Architecture 원칙을 적용한 Go 애플리케이션 예제입니다.

## 🏛️ Clean Architecture란?

Clean Architecture는 소프트웨어를 레이어로 분리하여 **비즈니스 로직을 외부 의존성으로부터 독립**시키는 아키텍처 패턴입니다.

### 핵심 원칙

1. **독립성**: 프레임워크, UI, 데이터베이스, 외부 라이브러리에 독립적
2. **테스트 가능성**: 비즈니스 규칙을 외부 요소 없이 테스트 가능
3. **의존성 규칙**: 소스 코드 의존성은 항상 안쪽(고수준)을 향해야 함

## 📐 레이어 구조

```
┌─────────────────────────────────────────┐
│  Frameworks & Drivers (가장 바깥)       │  ← 외부 세계
│  - Web Framework (Chi)                  │
│  - Database Driver                      │
│  - External APIs                        │
├─────────────────────────────────────────┤
│  Interface Adapters                     │
│  - HTTP Handlers (Controllers)          │
│  - Repositories (Gateways)              │
│  - Presenters                           │
├─────────────────────────────────────────┤
│  Use Cases (Application Business Rules) │  ← 의존성 방향
│  - User Use Case                        │     (안쪽으로)
│  - Business Logic Orchestration         │
├─────────────────────────────────────────┤
│  Entities (Enterprise Business Rules)   │
│  - Domain Models                        │  ← 가장 안쪽
│  - Business Logic                       │     (핵심)
└─────────────────────────────────────────┘
```

## 🗂️ 프로젝트 구조

```
CleanArchitecture/
├── cmd/
│   └── api/
│       └── main.go                 # 애플리케이션 진입점 (의존성 주입)
│
├── internal/
│   ├── domain/                     # 🔵 Entities (가장 안쪽)
│   │   ├── user.go                 # 도메인 엔티티
│   │   └── errors.go               # 도메인 에러
│   │
│   ├── usecase/                    # 🟢 Use Cases
│   │   ├── user_usecase.go         # 비즈니스 로직
│   │   └── interfaces.go           # 포트 (인터페이스)
│   │
│   ├── repository/                 # 🟡 Interface Adapters
│   │   └── memory/
│   │       └── user_repository.go  # 리포지토리 구현 (어댑터)
│   │
│   └── delivery/                   # 🔴 Frameworks & Drivers
│       └── http/
│           ├── handler.go          # HTTP 핸들러
│           └── router.go           # 라우터 설정
│
├── go.mod
└── README.md
```

## 🎯 각 레이어 설명

### 1️⃣ Domain (Entities) - 가장 안쪽 레이어

**위치**: `internal/domain/`

**책임**: 
- 비즈니스 규칙의 핵심
- 외부 의존성 전혀 없음
- 순수한 비즈니스 로직

**예시**:
```go
// User 엔티티
type User struct {
    ID        string
    Email     string
    Name      string
    CreatedAt time.Time
}

// 비즈니스 규칙
func (u *User) UpdateName(name string) error {
    if name == "" {
        return ErrInvalidName
    }
    u.Name = name
    return nil
}
```

**특징**:
- ✅ 외부 패키지 import 없음 (표준 라이브러리만)
- ✅ 프레임워크 독립적
- ✅ 데이터베이스 독립적

---

### 2️⃣ Use Cases - 애플리케이션 비즈니스 규칙

**위치**: `internal/usecase/`

**책임**:
- 애플리케이션의 비즈니스 흐름 조정
- 엔티티 조작
- 외부 레이어와의 인터페이스 정의 (포트)

**예시**:
```go
type UserUseCase struct {
    userRepo UserRepository  // 인터페이스 (포트)
}

func (uc *UserUseCase) CreateUser(email, name string) (*User, error) {
    // 1. 도메인 엔티티 생성
    user, err := domain.NewUser(email, name)
    
    // 2. 중복 체크
    existing, _ := uc.userRepo.GetByEmail(email)
    if existing != nil {
        return nil, ErrUserExists
    }
    
    // 3. 저장
    return user, uc.userRepo.Create(user)
}
```

**특징**:
- ✅ 도메인 레이어만 의존
- ✅ 인터페이스로 외부 레이어와 통신
- ✅ 의존성 역전 원칙 (DIP) 적용

---

### 3️⃣ Interface Adapters - 어댑터 레이어

**위치**: `internal/repository/`, `internal/delivery/http/`

**책임**:
- Use Case와 외부 세계를 연결
- 데이터 형식 변환
- 인터페이스 구현

**예시 - Repository (어댑터)**:
```go
type UserRepository struct {
    users map[string]*domain.User
}

// UserRepository 인터페이스 구현
func (r *UserRepository) Create(user *domain.User) error {
    r.users[user.ID] = user
    return nil
}
```

**예시 - HTTP Handler**:
```go
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    // 1. HTTP 요청 → DTO
    var req CreateUserRequest
    json.NewDecoder(r.Body).Decode(&req)
    
    // 2. Use Case 호출
    user, err := h.userUseCase.CreateUser(req.Email, req.Name)
    
    // 3. 도메인 → HTTP 응답
    respondJSON(w, http.StatusCreated, toUserResponse(user))
}
```

---

### 4️⃣ Frameworks & Drivers - 외부 레이어

**위치**: `cmd/api/main.go`

**책임**:
- 프레임워크 설정
- 의존성 주입
- 서버 시작

**예시**:
```go
func main() {
    // 의존성 주입 (바깥→안쪽)
    userRepo := memory.NewUserRepository()
    userUseCase := usecase.NewUserUseCase(userRepo)
    userHandler := http.NewUserHandler(userUseCase)
    router := http.NewRouter(userHandler)
    
    http.ListenAndServe(":8080", router)
}
```

---

## 🔄 의존성 흐름

```
[HTTP Request] 
    ↓
[HTTP Handler] ────→ [User Use Case] ────→ [User Entity]
    ↑                      ↑
    │                      │
[Repository Interface] ←───┘
    ↑
    │
[Memory Repository]
```

**의존성 방향**: HTTP → Use Case → Domain
**제어 흐름**: HTTP → Use Case → Repository → Use Case → HTTP

## 🚀 실행 방법

```bash
# 의존성 설치
cd CleanArchitecture
go mod tidy

# 옵션 1: Use Case 용어 사용 (기본)
go run cmd/api/main.go

# 옵션 2: Service 용어 사용
go run cmd/api/main_with_service.go

# 또는 빌드 후 실행
go build -o app cmd/api/main.go
./app
```

**참고**: `main.go`와 `main_with_service.go`는 용어만 다르고 기능은 동일합니다!
- `main.go` → **Use Case** 레이어 사용
- `main_with_service.go` → **Service** 레이어 사용

## 📝 API 사용 예제

### 사용자 생성
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "name": "John Doe"
  }'
```

### 모든 사용자 조회
```bash
curl http://localhost:8080/api/v1/users
```

### 특정 사용자 조회
```bash
curl http://localhost:8080/api/v1/users/{user-id}
```

### 사용자 수정
```bash
curl -X PUT http://localhost:8080/api/v1/users/{user-id} \
  -H "Content-Type: application/json" \
  -d '{"name": "Jane Doe"}'
```

### 사용자 삭제
```bash
curl -X DELETE http://localhost:8080/api/v1/users/{user-id}
```

## ✨ Clean Architecture의 장점

### 1. 테스트 용이성
```go
// Use Case를 독립적으로 테스트 가능
func TestCreateUser(t *testing.T) {
    mockRepo := &MockUserRepository{}
    useCase := usecase.NewUserUseCase(mockRepo)
    
    user, err := useCase.CreateUser("test@example.com", "Test User")
    
    assert.NoError(t, err)
    assert.NotNil(t, user)
}
```

### 2. 데이터베이스 교체 가능
```go
// 메모리 → PostgreSQL로 쉽게 교체
// userRepo := memory.NewUserRepository()
userRepo := postgres.NewUserRepository(db)

// Use Case는 변경 불필요!
userUseCase := usecase.NewUserUseCase(userRepo)
```

### 3. 프레임워크 독립성
```go
// Chi → Gin으로 교체해도 비즈니스 로직은 그대로
// router := chi.NewRouter()
router := gin.Default()

// 핸들러만 어댑터 교체
handler := gin_adapter.NewUserHandler(userUseCase)
```

## 🎯 핵심 패턴

### 1. 의존성 역전 원칙 (DIP)
```go
// ❌ 나쁜 예: Use Case가 구체적 구현에 의존
type UserUseCase struct {
    repo *PostgresRepository  // 구체적 타입
}

// ✅ 좋은 예: Use Case가 인터페이스에 의존
type UserUseCase struct {
    repo UserRepository  // 인터페이스
}
```

### 2. 포트와 어댑터 (Ports & Adapters)
```go
// Port (인터페이스)
type UserRepository interface {
    Create(user *User) error
}

// Adapter (구현)
type MemoryUserRepository struct { ... }
type PostgresUserRepository struct { ... }
type MongoUserRepository struct { ... }
```

### 3. DTO (Data Transfer Object)
```go
// Domain Entity
type User struct {
    ID    string
    Email string
    ...
}

// HTTP DTO
type UserResponse struct {
    ID    string `json:"id"`
    Email string `json:"email"`
}

// DTO 변환
func toUserResponse(user *User) UserResponse { ... }
```

## 📚 비교: 일반 구조 vs Clean Architecture

### 일반적인 MVC 구조
```
controllers/ ─┐
              ├─→ models/ ─→ database/
handlers/  ───┘
```
**문제**: 모든 것이 데이터베이스에 의존

### Clean Architecture
```
domain/ (독립)
    ↑
usecase/ (domain에만 의존)
    ↑
repository/, handler/ (usecase에 의존)
    ↑
main.go (모든 것을 조립)
```
**장점**: 비즈니스 로직이 완전히 독립적

## 🔧 확장 방법

### 1. 새 리포지토리 추가 (예: PostgreSQL)
```go
// internal/repository/postgres/user_repository.go
type UserRepository struct {
    db *sql.DB
}

func (r *UserRepository) Create(user *domain.User) error {
    // PostgreSQL 구현
}
```

### 2. 새 Use Case 추가
```go
// internal/usecase/auth_usecase.go
type AuthUseCase struct {
    userRepo UserRepository
}

func (uc *AuthUseCase) Login(email, password string) (*User, error) {
    // 로그인 로직
}
```

### 3. 새 Delivery 추가 (예: gRPC)
```go
// internal/delivery/grpc/user_service.go
type UserService struct {
    userUseCase *usecase.UserUseCase
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) {
    // gRPC 구현
}
```

## 📖 학습 자료

- [The Clean Architecture - Uncle Bob](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Clean Architecture in Go](https://github.com/bxcodec/go-clean-arch)
- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)

## 🎓 Clean Architecture를 언제 사용할까?

### ✅ 사용하면 좋을 때
- 장기간 유지보수할 프로젝트
- 복잡한 비즈니스 로직
- 테스트가 중요한 경우
- 여러 팀이 협업하는 프로젝트

### ❌ 과할 수 있는 경우
- 간단한 CRUD 애플리케이션
- 빠른 프로토타입
- 소규모 프로젝트

## 🎉 결론

Clean Architecture는 **비즈니스 로직을 보호**하고 **변경에 유연**한 구조를 제공합니다.

초기 설정은 복잡할 수 있지만, 장기적으로 **유지보수성**과 **테스트 용이성**이 크게 향상됩니다!

