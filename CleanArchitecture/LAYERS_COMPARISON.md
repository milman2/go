# 레이어 용어 비교: Clean Architecture vs 전통적인 아키텍처

## 🎯 Service 레이어는 어디에?

**답**: Clean Architecture에서 **Use Case 레이어 = Service 레이어**입니다!

## 📊 아키텍처 용어 비교

### 1. 전통적인 3-Layer (MVC) Architecture

```
┌─────────────────────────────────────┐
│  Presentation Layer                 │
│  (Controllers / Handlers)           │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│  Service Layer (Business Logic) ⭐  │  ← 여기가 Service!
│  - UserService                      │
│  - OrderService                     │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│  Data Access Layer                  │
│  (Repository / DAO)                 │
└─────────────────────────────────────┘
```

### 2. Clean Architecture

```
┌─────────────────────────────────────┐
│  Frameworks & Drivers               │
│  (HTTP Handlers, Database)          │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│  Interface Adapters                 │
│  (Controllers, Repositories)        │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│  Use Cases (Business Logic) ⭐      │  ← Service와 동일!
│  - UserUseCase                      │
│  - OrderUseCase                     │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│  Entities (Domain)                  │
│  - User, Order                      │
└─────────────────────────────────────┘
```

## 🔄 용어 매핑표

| 전통적인 아키텍처 | Clean Architecture | 역할 |
|------------------|-------------------|------|
| **Controller** | HTTP Handler | HTTP 요청/응답 처리 |
| **Service** ⭐ | **Use Case** ⭐ | 비즈니스 로직 |
| **Repository** | Repository (Interface Adapter) | 데이터 접근 |
| **Entity/Model** | Entity (Domain) | 도메인 모델 |

## 💡 왜 "Service"가 아니라 "Use Case"인가?

### Service라는 이름의 문제점
```go
// ❌ Service라는 이름은 너무 포괄적
type UserService struct {
    // 무엇을 하는 서비스인지 불명확
    // HTTP 서비스? 비즈니스 서비스? 데이터 서비스?
}
```

### Use Case라는 이름의 장점
```go
// ✅ Use Case는 명확한 의도를 표현
type UserUseCase struct {
    // "사용자 관련 유스케이스를 처리한다"
    // 비즈니스 시나리오를 구현한다는 의미
}

// 각 메서드가 하나의 유스케이스
func (uc *UserUseCase) CreateUser(...)    // 유스케이스: 사용자 생성
func (uc *UserUseCase) UpdateUser(...)    // 유스케이스: 사용자 수정
func (uc *UserUseCase) DeleteUser(...)    // 유스케이스: 사용자 삭제
```

**Use Case = 시스템이 수행하는 구체적인 비즈니스 시나리오**

## 📝 현재 프로젝트의 "Service" 레이어

우리 프로젝트에서:

```
internal/usecase/
├── user_usecase.go      ← 이것이 UserService와 동일!
└── interfaces.go        ← Repository 인터페이스 (포트)
```

### UserUseCase = UserService

```go
// internal/usecase/user_usecase.go
type UserUseCase struct {
    userRepo UserRepository
}

// 이것들이 전통적인 Service 메서드와 동일
func (uc *UserUseCase) CreateUser(ctx context.Context, email, name string) (*domain.User, error)
func (uc *UserUseCase) GetUser(ctx context.Context, id string) (*domain.User, error)
func (uc *UserUseCase) UpdateUser(ctx context.Context, id, name string) (*domain.User, error)
func (uc *UserUseCase) DeleteUser(ctx context.Context, id string) error
```

## 🔀 원한다면 Service라는 이름도 사용 가능

### 옵션 1: type alias 사용
```go
// internal/usecase/user_usecase.go
type UserUseCase struct {
    userRepo UserRepository
}

// Service는 UseCase의 별칭
type UserService = UserUseCase

// 둘 다 사용 가능
var _ UserService = (*UserUseCase)(nil)
```

### 옵션 2: Service라는 이름으로 변경
```go
// internal/service/user_service.go
package service

type UserService struct {
    userRepo UserRepository
}

func NewUserService(userRepo UserRepository) *UserService {
    return &UserService{userRepo: userRepo}
}

func (s *UserService) CreateUser(...) (*domain.User, error) {
    // 비즈니스 로직
}
```

## 🆚 전체 비교: 전통적 vs Clean Architecture

### 전통적인 3-Layer Architecture

```go
// controller/user_controller.go
type UserController struct {
    userService *service.UserService
}

func (c *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
    user, err := c.userService.CreateUser(email, name)
    json.NewEncoder(w).Encode(user)
}

// service/user_service.go
type UserService struct {
    userRepo *repository.UserRepository
}

func (s *UserService) CreateUser(email, name string) (*model.User, error) {
    // 비즈니스 로직
    return s.userRepo.Create(user)
}

// repository/user_repository.go
type UserRepository struct {
    db *sql.DB
}

func (r *UserRepository) Create(user *model.User) error {
    // DB 저장
}
```

### Clean Architecture (현재 프로젝트)

```go
// internal/delivery/http/handler.go
type UserHandler struct {
    userUseCase *usecase.UserUseCase  // ← Service와 동일!
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    user, err := h.userUseCase.CreateUser(ctx, email, name)
    respondJSON(w, http.StatusCreated, toUserResponse(user))
}

// internal/usecase/user_usecase.go
type UserUseCase struct {  // ← 이것이 Service!
    userRepo UserRepository  // 인터페이스
}

func (uc *UserUseCase) CreateUser(ctx context.Context, email, name string) (*domain.User, error) {
    // 비즈니스 로직
    return user, uc.userRepo.Create(ctx, user)
}

// internal/repository/memory/user_repository.go
type UserRepository struct {
    users map[string]*domain.User
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
    // 저장
}
```

## 📐 레이어 매핑 상세

### 전통적인 아키텍처
```
Controller
    ↓ calls
Service ⭐ (비즈니스 로직)
    ↓ calls
Repository (데이터 접근)
    ↓ uses
Model/Entity
```

### Clean Architecture
```
Handler (Controller와 동일)
    ↓ calls
Use Case ⭐ (Service와 동일)
    ↓ calls via interface
Repository Adapter
    ↓ uses
Domain Entity
```

## 🎯 결론

### Service = Use Case

| 측면 | Service | Use Case |
|------|---------|----------|
| **역할** | 비즈니스 로직 | 비즈니스 로직 |
| **위치** | 중간 레이어 | 중간 레이어 |
| **의존** | Repository | Repository (인터페이스) |
| **호출자** | Controller | Handler |

**차이점**:
1. **용어**: Service vs Use Case (의미 강조 차이)
2. **의존성**: Service는 구체 타입, Use Case는 인터페이스
3. **철학**: Service는 레이어, Use Case는 시나리오

## 💼 실무에서는?

### 많은 프로젝트가 혼용
```go
// 이것도 OK
type UserService struct { ... }

// 이것도 OK
type UserUseCase struct { ... }

// 심지어 이것도 OK (DDD 스타일)
type UserApplicationService struct { ... }
```

**중요한 것**: 
- ✅ 비즈니스 로직을 분리
- ✅ 의존성을 올바른 방향으로
- ✅ 테스트 가능한 구조
- ❌ 이름에 집착하지 말 것

## 🔄 원한다면 Service로 이름 변경하기

현재 프로젝트를 Service 용어로 변경하고 싶다면:

```bash
# 1. 디렉토리 이름 변경
mv internal/usecase internal/service

# 2. 파일 이름 변경
mv internal/service/user_usecase.go internal/service/user_service.go

# 3. 타입 이름 변경
# UserUseCase → UserService

# 4. 패키지 이름 변경
# package usecase → package service
```

하지만 **Clean Architecture에서는 Use Case 용어를 권장**합니다!

## 📚 참고: 다양한 아키텍처 용어

| 아키텍처 스타일 | 비즈니스 로직 레이어 이름 |
|----------------|------------------------|
| 3-Layer | Service Layer |
| Clean Architecture | Use Case Layer |
| Hexagonal (Ports & Adapters) | Application Services |
| DDD | Application Services |
| Onion Architecture | Application Services |
| MVC | Service Layer / Model |

**모두 같은 역할**: 비즈니스 로직 처리! 🎯

