# Chi REST API 서버

Chi 라우터를 사용한 간단한 REST API 서버 예제입니다.

## 특징

- ✅ 경량 및 고성능 라우터
- ✅ 표준 net/http 호환
- ✅ 강력한 미들웨어 지원
- ✅ 중첩 라우트 (nested routes)
- ✅ URL 파라미터 추출
- ✅ 마이크로서비스 친화적

## Chi vs Gin

### Chi 장점
- 표준 `net/http` 인터페이스 사용 (상호 운용성)
- 더 작은 코드베이스
- 중첩 라우팅이 더 명확
- Context 기반의 미들웨어
- 마이크로서비스에 최적화

### Gin 장점
- 더 빠른 성능 (약간)
- 자동 JSON 바인딩 및 유효성 검증
- 더 큰 커뮤니티
- 더 많은 내장 기능

## 설치 및 실행

```bash
# 의존성 설치
cd REST/Chi
go mod tidy

# 서버 실행
go run main.go

# 또는 빌드 후 실행
go build -o chi-server
./chi-server
```

서버는 `http://localhost:8080`에서 실행됩니다.

## API 엔드포인트

### Health Check
```bash
curl http://localhost:8080/health
```

### 아이템 목록 조회
```bash
curl http://localhost:8080/api/v1/items
```

### 아이템 생성
```bash
curl -X POST http://localhost:8080/api/v1/items \
  -H "Content-Type: application/json" \
  -d '{
    "name": "키보드",
    "description": "기계식 키보드",
    "price": 150000
  }'
```

### 특정 아이템 조회
```bash
curl http://localhost:8080/api/v1/items/1
```

### 아이템 수정
```bash
curl -X PUT http://localhost:8080/api/v1/items/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "키보드 (수정)",
    "description": "RGB 기계식 키보드",
    "price": 180000
  }'
```

### 아이템 삭제
```bash
curl -X DELETE http://localhost:8080/api/v1/items/1
```

## 응답 예제

### 성공 응답
```json
{
  "message": "아이템이 생성되었습니다",
  "data": {
    "id": 1,
    "name": "키보드",
    "description": "기계식 키보드",
    "price": 150000
  }
}
```

### 에러 응답
```json
{
  "error": "아이템을 찾을 수 없습니다"
}
```

## 프로젝트 구조

```
Chi/
├── main.go          # 메인 애플리케이션
├── go.mod          # Go 모듈 파일
└── README.md       # 문서
```

## 주요 기능

### 1. 중첩 라우팅 (Nested Routing)
```go
r.Route("/api/v1", func(r chi.Router) {
    r.Route("/items", func(r chi.Router) {
        r.Get("/", getItems)
        r.Post("/", createItem)
        
        r.Route("/{id}", func(r chi.Router) {
            r.Get("/", getItem)
            r.Put("/", updateItem)
            r.Delete("/", deleteItem)
        })
    })
})
```

### 2. URL 파라미터
```go
id := chi.URLParam(r, "id")
```

### 3. 미들웨어
Chi는 다양한 내장 미들웨어를 제공합니다:
- `middleware.Logger`: HTTP 요청 로깅
- `middleware.Recoverer`: 패닉 복구
- `middleware.RequestID`: 요청 ID 생성
- `middleware.RealIP`: 실제 클라이언트 IP 추출
- `middleware.Timeout`: 요청 타임아웃
- `middleware.Throttle`: Rate limiting

### 4. 표준 net/http 호환
```go
// 표준 http.HandlerFunc 사용
func handler(w http.ResponseWriter, r *http.Request) {
    // ...
}
```

## 고급 기능

### 커스텀 미들웨어 만들기
```go
func MyMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 전처리
        log.Println("Before request")
        
        next.ServeHTTP(w, r)
        
        // 후처리
        log.Println("After request")
    })
}

// 사용
r.Use(MyMiddleware)
```

### 서브라우터
```go
// 관리자 전용 라우트
adminRouter := chi.NewRouter()
adminRouter.Use(AdminOnly) // 관리자 전용 미들웨어
adminRouter.Get("/users", listUsers)

r.Mount("/admin", adminRouter)
```

### 타임아웃 설정
```go
r.Use(middleware.Timeout(60 * time.Second))
```

## 성능

Chi는 다음과 같은 특징으로 우수한 성능을 제공합니다:
- Radix tree 기반 라우팅
- 제로 힙 할당 (zero heap allocations)
- 최소한의 메모리 오버헤드
- 컴파일 타임 최적화

## 벤치마크 (참고)

```
BenchmarkChi_Param          10000000    164 ns/op    0 B/op    0 allocs/op
BenchmarkChi_Param5         5000000     287 ns/op    0 B/op    0 allocs/op
BenchmarkChi_ParamWrite     10000000    203 ns/op    0 B/op    0 allocs/op
```

## Gin vs Chi 선택 가이드

### Chi를 선택하는 경우
- 마이크로서비스 아키텍처
- 표준 라이브러리와의 호환성 중요
- 미들웨어 체인이 복잡한 경우
- 코드 간결성과 명확성 우선

### Gin을 선택하는 경우
- 빠른 개발이 필요한 경우
- 자동 바인딩/검증이 필요한 경우
- 커뮤니티 지원이 중요한 경우
- 프로토타이핑

## 다음 단계

이 예제를 확장하려면:

1. **데이터베이스 연동**: PostgreSQL, MySQL, MongoDB 등
2. **인증/인가**: JWT, OAuth2
3. **CORS 설정**: `chi-cors` 미들웨어
4. **로깅**: 구조화된 로깅 (zap, logrus)
5. **설정 관리**: 환경 변수, 설정 파일
6. **테스트**: 유닛 테스트, 통합 테스트
7. **문서화**: Swagger/OpenAPI
8. **Rate Limiting**: `chi/middleware.Throttle`
9. **모니터링**: Prometheus 메트릭

## 추가 Chi 미들웨어

```bash
# CORS
go get github.com/go-chi/cors

# JWT 인증
go get github.com/go-chi/jwtauth/v5

# Rate limiting (내장)
import "github.com/go-chi/chi/v5/middleware"
```

## 참고 자료

- [Chi 공식 문서](https://go-chi.io/)
- [Chi GitHub](https://github.com/go-chi/chi)
- [Chi 예제 모음](https://github.com/go-chi/chi/tree/master/_examples)

