# Gin vs Chi 비교

두 프레임워크의 실제 코드를 비교해봅시다.

## 1. 프로젝트 초기화

### Gin
```go
r := gin.Default() // Logger + Recovery 미들웨어 포함
```

### Chi
```go
r := chi.NewRouter()
r.Use(middleware.Logger)
r.Use(middleware.Recoverer)
```

**차이점**: Gin은 기본 미들웨어를 자동으로 포함하지만, Chi는 명시적으로 추가해야 합니다.

---

## 2. 라우팅

### Gin - 그룹화
```go
v1 := r.Group("/api/v1")
{
    items := v1.Group("/items")
    {
        items.GET("", getItems)
        items.GET("/:id", getItem)
    }
}
```

### Chi - 중첩 라우팅
```go
r.Route("/api/v1", func(r chi.Router) {
    r.Route("/items", func(r chi.Router) {
        r.Get("/", getItems)
        r.Get("/{id}", getItem)
    })
})
```

**차이점**: Chi의 중첩 라우팅이 더 명확하고 함수형 스타일입니다.

---

## 3. URL 파라미터

### Gin
```go
id := c.Param("id")
```

### Chi
```go
id := chi.URLParam(r, "id")
```

**차이점**: Gin은 Context에서, Chi는 Request에서 파라미터를 추출합니다.

---

## 4. JSON 처리

### Gin - 자동 바인딩
```go
var item Item
if err := c.ShouldBindJSON(&item); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
        "error": err.Error(),
    })
    return
}

// 응답
c.JSON(http.StatusOK, gin.H{
    "data": item,
})
```

### Chi - 수동 처리
```go
var item Item
if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
    respondJSON(w, http.StatusBadRequest, Response{
        Error: err.Error(),
    })
    return
}

// 응답 (헬퍼 함수 필요)
respondJSON(w, http.StatusOK, Response{
    Data: item,
})
```

**차이점**: Gin은 자동 바인딩과 유효성 검증을 제공하지만, Chi는 표준 라이브러리를 사용합니다.

---

## 5. 유효성 검증

### Gin - 구조체 태그
```go
type Item struct {
    Name  string `json:"name" binding:"required"`
    Price int    `json:"price" binding:"required,min=0"`
}
```

### Chi - 수동 검증
```go
type Item struct {
    Name  string `json:"name"`
    Price int    `json:"price"`
}

// 수동 검증 필요
if item.Name == "" {
    respondJSON(w, http.StatusBadRequest, Response{
        Error: "이름은 필수 항목입니다",
    })
    return
}
```

**차이점**: Gin은 선언적 유효성 검증을 지원하지만, Chi는 직접 구현해야 합니다.

---

## 6. 미들웨어

### Gin
```go
r.Use(gin.Logger())
r.Use(gin.Recovery())
r.Use(CustomMiddleware())
```

### Chi
```go
r.Use(middleware.Logger)
r.Use(middleware.Recoverer)
r.Use(CustomMiddleware)
```

**차이점**: Chi의 미들웨어는 표준 `http.Handler` 인터페이스를 따릅니다.

---

## 7. 핸들러 시그니처

### Gin
```go
func handler(c *gin.Context) {
    // gin.Context는 모든 것을 포함
}
```

### Chi
```go
func handler(w http.ResponseWriter, r *http.Request) {
    // 표준 net/http 시그니처
}
```

**차이점**: Chi는 표준 HTTP 핸들러를 사용하여 다른 라이브러리와 쉽게 통합됩니다.

---

## 성능 비교

### Gin
- 더 빠른 라우팅 (httprouter 기반)
- 벤치마크: ~164 ns/op (파라미터 1개)

### Chi
- Radix tree 기반
- 벤치마크: ~164 ns/op (파라미터 1개)
- 거의 동일한 성능

**결론**: 실제 성능 차이는 미미합니다.

---

## 코드 크기

### Gin 예제
- 라인 수: ~198줄
- 의존성: 28개

### Chi 예제
- 라인 수: ~236줄
- 의존성: 1개

**차이점**: Chi는 의존성이 적지만, 보일러플레이트 코드가 더 많습니다.

---

## 사용 사례

### Gin을 선택하세요
- ✅ 빠른 프로토타이핑
- ✅ 자동 바인딩/검증이 필요한 경우
- ✅ 풀스택 웹 애플리케이션
- ✅ 학습 곡선이 중요한 경우

### Chi를 선택하세요
- ✅ 마이크로서비스
- ✅ 표준 라이브러리와의 호환성이 중요
- ✅ 미들웨어 체인이 복잡한 경우
- ✅ 코드 명확성이 우선
- ✅ 기존 net/http 코드와 통합

---

## 실제 벤치마크

```bash
# Gin 서버 벤치마크
wrk -t12 -c400 -d30s http://localhost:8080/api/v1/items

# Chi 서버 벤치마크
wrk -t12 -c400 -d30s http://localhost:8080/api/v1/items
```

일반적으로 두 프레임워크 모두:
- 수만 req/s 처리 가능
- 낮은 메모리 사용량
- 프로덕션 준비 완료

---

## 커뮤니티 & 생태계

### Gin
- ⭐ GitHub Stars: ~77k
- 📦 더 많은 플러그인/미들웨어
- 📚 더 많은 튜토리얼
- 🇨🇳 중국에서 매우 인기

### Chi
- ⭐ GitHub Stars: ~18k
- 📦 핵심 기능에 집중
- 📚 깔끔한 문서
- 🌍 서양에서 인기

---

## 마이그레이션

### Gin → Chi
```go
// Gin
c.JSON(200, gin.H{"data": item})
c.Param("id")

// Chi
respondJSON(w, 200, Response{Data: item})
chi.URLParam(r, "id")
```

### Chi → Gin
```go
// Chi
respondJSON(w, 200, Response{Data: item})
chi.URLParam(r, "id")

// Gin
c.JSON(200, gin.H{"data": item})
c.Param("id")
```

---

## 최종 추천

### 초보자 / 빠른 개발
→ **Gin** (더 많은 기능, 쉬운 시작)

### 경험자 / 마이크로서비스
→ **Chi** (더 명확한 코드, 표준 호환)

### 어느 쪽이든
→ 둘 다 훌륭한 선택입니다! 🎉

