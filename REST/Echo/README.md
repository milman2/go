# Echo REST API 서버

Echo 프레임워크를 사용한 간단한 REST API 서버 예제입니다.

## 특징

- ✅ 고성능 HTTP 라우터
- ✅ 풍부한 내장 미들웨어
- ✅ 자동 JSON 바인딩
- ✅ 에러 핸들링
- ✅ 데이터 검증 지원
- ✅ HTTP/2, WebSocket 지원

## Echo의 강점

Echo는 Gin과 Chi의 장점을 결합한 프레임워크입니다:

### vs Gin
- 더 많은 내장 미들웨어 (CORS, JWT, Secure 등)
- 더 나은 문서화
- HTTP/2 기본 지원
- 더 유연한 라우팅

### vs Chi
- 자동 바인딩 및 검증
- 더 많은 편의 기능
- Gin보다 가벼우면서도 기능이 풍부

## 설치 및 실행

```bash
# 의존성 설치
cd REST/Echo
go mod tidy

# 서버 실행
go run main.go

# 또는 빌드 후 실행
go build -o echo-server
./echo-server
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
    "name": "헤드셋",
    "description": "게이밍 헤드셋",
    "price": 80000
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
    "name": "헤드셋 (수정)",
    "description": "7.1 채널 게이밍 헤드셋",
    "price": 120000
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
    "name": "헤드셋",
    "description": "게이밍 헤드셋",
    "price": 80000
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
Echo/
├── main.go          # 메인 애플리케이션
├── go.mod          # Go 모듈 파일
└── README.md       # 문서
```

## 주요 기능

### 1. Echo Context
Echo는 강력한 Context를 제공합니다:

```go
func handler(c echo.Context) error {
    // 요청 데이터 바인딩
    var data MyData
    c.Bind(&data)
    
    // URL 파라미터
    id := c.Param("id")
    
    // 쿼리 파라미터
    page := c.QueryParam("page")
    
    // JSON 응답
    return c.JSON(http.StatusOK, data)
}
```

### 2. 라우트 그룹화
```go
v1 := e.Group("/api/v1")
{
    items := v1.Group("/items")
    {
        items.GET("", getItems)
        items.POST("", createItem)
    }
}
```

### 3. 내장 미들웨어
Echo는 풍부한 미들웨어를 제공합니다:

```go
// 로깅
e.Use(middleware.Logger())

// 패닉 복구
e.Use(middleware.Recover())

// CORS
e.Use(middleware.CORS())

// Gzip 압축
e.Use(middleware.Gzip())

// JWT 인증
e.Use(middleware.JWT([]byte("secret")))

// Rate Limiting
e.Use(middleware.RateLimiter(
    middleware.NewRateLimiterMemoryStore(20),
))

// Secure 헤더
e.Use(middleware.Secure())

// Request ID
e.Use(middleware.RequestID())
```

### 4. 에러 핸들링
Echo는 일관된 에러 핸들링을 제공합니다:

```go
// 커스텀 에러 핸들러
e.HTTPErrorHandler = func(err error, c echo.Context) {
    code := http.StatusInternalServerError
    message := "Internal Server Error"
    
    if he, ok := err.(*echo.HTTPError); ok {
        code = he.Code
        message = he.Message.(string)
    }
    
    c.JSON(code, map[string]string{
        "error": message,
    })
}
```

### 5. 데이터 바인딩 및 검증
```go
type User struct {
    Name  string `json:"name" validate:"required"`
    Email string `json:"email" validate:"required,email"`
}

func createUser(c echo.Context) error {
    u := new(User)
    if err := c.Bind(u); err != nil {
        return err
    }
    if err := c.Validate(u); err != nil {
        return err
    }
    return c.JSON(http.StatusOK, u)
}
```

## Echo의 특별한 기능

### 1. HTTP/2 지원
```go
e.StartTLS(":443", "cert.pem", "key.pem")
```

### 2. WebSocket
```go
e.GET("/ws", func(c echo.Context) error {
    websocket.Handler(func(ws *websocket.Conn) {
        // WebSocket 로직
    }).ServeHTTP(c.Response(), c.Request())
    return nil
})
```

### 3. 서브도메인 라우팅
```go
// api.example.com
api := e.Host("api.example.com")
api.GET("/users", getUsers)

// admin.example.com
admin := e.Host("admin.example.com")
admin.GET("/dashboard", getDashboard)
```

### 4. 파일 업로드
```go
func upload(c echo.Context) error {
    file, err := c.FormFile("file")
    if err != nil {
        return err
    }
    
    src, err := file.Open()
    if err != nil {
        return err
    }
    defer src.Close()
    
    // 파일 저장 로직
    return c.JSON(http.StatusOK, map[string]string{
        "message": "uploaded",
    })
}
```

## 성능

Echo는 매우 빠른 성능을 제공합니다:

```
BenchmarkEcho_Param           20000000    100 ns/op    0 B/op   0 allocs/op
BenchmarkEcho_Param5          10000000    191 ns/op    0 B/op   0 allocs/op
BenchmarkEcho_ParamWrite      10000000    203 ns/op    0 B/op   0 allocs/op
```

## Gin vs Chi vs Echo 비교

| 특징 | Gin | Chi | Echo |
|------|-----|-----|------|
| 성능 | ⚡⚡⚡ | ⚡⚡⚡ | ⚡⚡⚡ |
| 학습 곡선 | 쉬움 | 보통 | 쉬움 |
| 미들웨어 | 많음 | 보통 | 매우 많음 |
| 자동 바인딩 | ✅ | ❌ | ✅ |
| 표준 호환 | ❌ | ✅ | ❌ |
| 문서화 | 좋음 | 좋음 | 매우 좋음 |
| 커뮤니티 | 큼 | 중간 | 큼 |
| HTTP/2 | ✅ | ✅ | ✅ |
| WebSocket | 수동 | 수동 | 내장 |

## 언제 Echo를 선택해야 할까?

### Echo를 선택하세요 ✅
- 풍부한 미들웨어가 필요한 경우
- HTTP/2, WebSocket 등 고급 기능이 필요한 경우
- 좋은 문서와 큰 커뮤니티가 중요한 경우
- Gin보다 더 많은 기능을 원하지만 무겁지 않은 프레임워크를 원할 때
- 실시간 애플리케이션 (WebSocket)

### 다른 프레임워크를 선택하세요
- **Gin**: 가장 빠른 개발 속도와 큰 커뮤니티
- **Chi**: 표준 호환성과 마이크로서비스

## 다음 단계

이 예제를 확장하려면:

1. **데이터베이스 연동**: PostgreSQL, MySQL, MongoDB 등
2. **JWT 인증**: `middleware.JWT()`
3. **CORS 설정**: `middleware.CORSWithConfig()`
4. **Rate Limiting**: `middleware.RateLimiter()`
5. **로깅**: 구조화된 로깅 (zap, logrus)
6. **설정 관리**: 환경 변수, 설정 파일
7. **테스트**: `httptest` 패키지 사용
8. **문서화**: Swagger/OpenAPI
9. **모니터링**: Prometheus 메트릭
10. **WebSocket**: 실시간 통신

## 추가 Echo 패키지

```bash
# JWT 미들웨어
go get github.com/labstack/echo-jwt/v4

# Contrib (추가 미들웨어)
go get github.com/labstack/echo-contrib
```

## 예제: JWT 인증

```go
import (
    echojwt "github.com/labstack/echo-jwt/v4"
)

func main() {
    e := echo.New()
    
    // 로그인 (JWT 토큰 생성)
    e.POST("/login", login)
    
    // 보호된 라우트
    r := e.Group("/api")
    r.Use(echojwt.JWT([]byte("secret")))
    r.GET("/users", getUsers)
    
    e.Start(":8080")
}
```

## 참고 자료

- [Echo 공식 문서](https://echo.labstack.com/docs)
- [Echo GitHub](https://github.com/labstack/echo)
- [Echo 예제 모음](https://echo.labstack.com/cookbook)
- [Echo 미들웨어 가이드](https://echo.labstack.com/middleware)

