# Gin REST API 서버

Gin 프레임워크를 사용한 간단한 REST API 서버 예제입니다.

## 특징

- ✅ 고성능 HTTP 라우터
- ✅ JSON 바인딩 및 유효성 검증
- ✅ 미들웨어 지원
- ✅ 에러 핸들링
- ✅ RESTful 라우팅

## 설치 및 실행

```bash
# 의존성 설치
cd REST/Gin
go mod tidy

# 서버 실행
go run main.go

# 또는 빌드 후 실행
go build -o gin-server
./gin-server
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
    "name": "노트북",
    "description": "고성능 노트북",
    "price": 1500000
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
    "name": "노트북 (수정)",
    "description": "최신 고성능 노트북",
    "price": 1800000
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
    "name": "노트북",
    "description": "고성능 노트북",
    "price": 1500000
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
Gin/
├── main.go          # 메인 애플리케이션
├── go.mod          # Go 모듈 파일
└── README.md       # 문서
```

## 주요 기능

### 1. 자동 JSON 바인딩
```go
type Item struct {
    Name  string `json:"name" binding:"required"`
    Price int    `json:"price" binding:"required,min=0"`
}
```

### 2. 라우트 그룹화
```go
v1 := r.Group("/api/v1")
{
    items := v1.Group("/items")
    {
        items.GET("", getItems)
        items.POST("", createItem)
    }
}
```

### 3. 미들웨어
Gin은 기본적으로 다음 미들웨어를 포함합니다:
- Logger: 요청 로깅
- Recovery: 패닉 복구

### 4. 에러 핸들링
```go
if err := c.ShouldBindJSON(&item); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
        "error": err.Error(),
    })
    return
}
```

## 성능

Gin은 다음과 같은 이유로 고성능입니다:
- httprouter 기반의 빠른 라우팅
- 낮은 메모리 사용량
- 효율적인 미들웨어 체인

## 다음 단계

이 예제를 확장하려면:

1. **데이터베이스 연동**: PostgreSQL, MySQL, MongoDB 등
2. **인증/인가**: JWT, OAuth2
3. **미들웨어 추가**: CORS, Rate Limiting
4. **로깅**: 구조화된 로깅 (zap, logrus)
5. **설정 관리**: 환경 변수, 설정 파일
6. **테스트**: 유닛 테스트, 통합 테스트
7. **문서화**: Swagger/OpenAPI

## 참고 자료

- [Gin 공식 문서](https://gin-gonic.com/docs/)
- [Gin GitHub](https://github.com/gin-gonic/gin)

