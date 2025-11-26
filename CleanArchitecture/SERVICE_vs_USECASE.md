# Service vs Use Case: 빠른 비교

## 🎯 한 줄 요약

**Service = Use Case** (같은 역할, 다른 이름)

## 📊 시각적 비교

```
전통적 아키�ecture          Clean Architecture
──────────────────          ──────────────────
┌────────────────┐          ┌────────────────┐
│  Controller    │          │   Handler      │
└────────┬───────┘          └────────┬───────┘
         │                           │
         ↓                           ↓
┌────────────────┐          ┌────────────────┐
│  Service ⭐    │    =     │ Use Case ⭐    │  ← 같은 역할!
│  (비즈니스)    │          │  (비즈니스)    │
└────────┬───────┘          └────────┬───────┘
         │                           │
         ↓                           ↓
┌────────────────┐          ┌────────────────┐
│  Repository    │          │  Repository    │
└────────────────┘          └────────────────┘
```

## 🔄 실제 코드 비교

### 전통적인 방식 (Service)
```go
// service/user_service.go
type UserService struct {
    userRepo UserRepository
}

func (s *UserService) CreateUser(email, name string) (*User, error) {
    // 비즈니스 로직
    user := NewUser(email, name)
    return s.userRepo.Create(user)
}
```

### Clean Architecture (Use Case)
```go
// usecase/user_usecase.go
type UserUseCase struct {
    userRepo UserRepository
}

func (uc *UserUseCase) CreateUser(email, name string) (*User, error) {
    // 비즈니스 로직 (동일!)
    user := NewUser(email, name)
    return uc.userRepo.Create(user)
}
```

**차이점**: 이름뿐! 내용은 100% 동일

## 🚀 현재 프로젝트 실행 방법

### Option 1: Use Case 용어 (Clean Architecture 표준)
```bash
go run cmd/api/main.go
```

### Option 2: Service 용어 (익숙한 용어)
```bash
go run cmd/api/main_with_service.go
```

**둘 다 같은 API를 제공합니다!**

## 📁 프로젝트 구조

```
internal/
├── domain/           # 도메인 엔티티
├── usecase/          # ⭐ Use Case (비즈니스 로직)
├── service/          # ⭐ Service (비즈니스 로직, usecase와 동일)
├── repository/       # 데이터 접근
└── delivery/         # HTTP 처리
```

## 💡 어떤 용어를 사용할까?

### Use Case 추천 ✅ (Clean Architecture 표준)
- 명확한 의도: "이 시스템이 무엇을 하는가"
- Uncle Bob의 공식 용어
- 비즈니스 시나리오 중심

**사용 예**:
```go
userUseCase := usecase.NewUserUseCase(userRepo)
user, err := userUseCase.CreateUser(email, name)
```

### Service도 가능 ✅ (익숙한 용어)
- 개발자들에게 익숙함
- 전통적 아키텍처와 일관성
- 팀의 선호도에 따라

**사용 예**:
```go
userService := service.NewUserService(userRepo)
user, err := userService.CreateUser(email, name)
```

## 🎓 학습 가이드

1. **기본 학습**: `internal/usecase/` 디렉토리 (Use Case)
2. **비교 학습**: `internal/service/` 디렉토리 (Service)
3. **상세 비교**: `LAYERS_COMPARISON.md` 문서

## 📝 테스트

```bash
# Use Case 버전 실행
go run cmd/api/main.go

# 다른 터미널에서
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","name":"Test User"}'

# 서버 종료 후 Service 버전 실행
go run cmd/api/main_with_service.go

# 같은 API 호출 (동일하게 작동!)
curl http://localhost:8080/api/v1/users
```

## 🎯 최종 결론

### 중요한 것은 이름이 아닙니다!

**중요한 것**:
- ✅ 비즈니스 로직을 분리
- ✅ 의존성을 올바른 방향으로
- ✅ 테스트 가능한 구조
- ✅ 단일 책임 원칙

**덜 중요한 것**:
- ❌ Service vs Use Case (용어 선택)

### 팀에서 편한 용어를 선택하세요!

둘 다 완벽하게 유효한 선택입니다. 🎉

---

## 🔗 관련 문서

- **LAYERS_COMPARISON.md** - 레이어 용어 상세 비교
- **ARCHITECTURE.md** - 아키텍처 상세 다이어그램
- **README.md** - 전체 프로젝트 개요
- **internal/service/README.md** - Service 레이어 설명

