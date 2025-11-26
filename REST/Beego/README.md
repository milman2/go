# Beego REST API 서버

Beego 프레임워크를 사용한 간단한 REST API 서버 예제입니다.

## 특징

- ✅ MVC 아키텍처
- ✅ 컨트롤러 기반 라우팅
- ✅ 네임스페이스 지원
- ✅ ORM 내장 (Bee ORM)
- ✅ 자동 API 문서 생성
- ✅ 캐싱 지원
- ✅ 세션 관리
- ✅ i18n 지원

## Beego의 특별한 점

Beego는 다른 프레임워크와 달리 **풀스택 MVC 프레임워크**입니다:

### vs Gin/Chi/Echo
- MVC 패턴 (Model-View-Controller)
- ORM 내장 (다른 프레임워크는 외부 ORM 필요)
- 더 많은 엔터프라이즈 기능
- 자동 문서 생성 (Swagger)
- 프로젝트 스캐폴딩 도구

### 언제 Beego를 사용할까?
- 대규모 엔터프라이즈 애플리케이션
- 풀스택 웹 애플리케이션
- ORM이 필요한 경우
- 자동 문서화가 필요한 경우
- 중국 시장 타겟팅

## 설치 및 실행

```bash
# 의존성 설치
cd REST/Beego
go mod tidy

# 서버 실행
go run main.go

# 또는 빌드 후 실행
go build -o beego-server
./beego-server
```

서버는 `http://localhost:8080`에서 실행됩니다.

## Beego CLI 도구 (선택사항)

```bash
# Beego CLI 설치
go install github.com/beego/bee/v2@latest

# 새 프로젝트 생성
bee new myproject
bee api myapi

# 서버 실행 (핫 리로드)
bee run

# API 문서 생성
bee generate docs
```

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
    "name": "데스크탑",
    "description": "고성능 게이밍 PC",
    "price": 2000000
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
    "name": "데스크탑 (수정)",
    "description": "최신 RTX 4090 게이밍 PC",
    "price": 3000000
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
    "name": "데스크탑",
    "description": "고성능 게이밍 PC",
    "price": 2000000
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
Beego/
├── main.go          # 메인 애플리케이션 (컨트롤러 포함)
├── go.mod          # Go 모듈 파일
└── README.md       # 문서
```

실제 Beego 프로젝트에서는 다음과 같은 구조를 사용합니다:

```
myapp/
├── main.go
├── controllers/
│   ├── default.go
│   └── item.go
├── models/
│   └── item.go
├── routers/
│   └── router.go
├── conf/
│   └── app.conf
├── views/
│   └── index.tpl
└── static/
    ├── css/
    ├── js/
    └── img/
```

## 주요 기능

### 1. MVC 컨트롤러
Beego는 컨트롤러 기반입니다:

```go
type ItemsController struct {
    web.Controller
}

func (c *ItemsController) Get() {
    // ID 가져오기
    id := c.Ctx.Input.Param(":id")
    
    // JSON 응답
    c.Data["json"] = data
    c.ServeJSON()
}
```

### 2. 네임스페이스 라우팅
```go
ns := web.NewNamespace("/api/v1",
    web.NSRouter("/items", &ItemsController{}, "get:GetAll;post:Post"),
    web.NSRouter("/items/:id", &ItemsController{}, "get:Get;put:Put;delete:Delete"),
)
web.AddNamespace(ns)
```

### 3. ORM 사용 (Bee ORM)
```go
import "github.com/beego/beego/v2/client/orm"

type Item struct {
    Id          int
    Name        string
    Description string
    Price       int
}

func init() {
    orm.RegisterModel(new(Item))
}

// CRUD 작업
o := orm.NewOrm()
item := Item{Name: "Test", Price: 100}
id, err := o.Insert(&item)
```

### 4. 자동 라우팅
```go
// 컨트롤러 메서드 이름으로 자동 라우팅
type UserController struct {
    web.Controller
}

func (c *UserController) Get() {}    // GET /user/:id
func (c *UserController) Post() {}   // POST /user
func (c *UserController) Put() {}    // PUT /user/:id
func (c *UserController) Delete() {} // DELETE /user/:id
```

### 5. 필터 (미들웨어)
```go
// 인증 필터
var FilterAuth = func(ctx *context.Context) {
    token := ctx.Input.Header("Authorization")
    if token == "" {
        ctx.Output.SetStatus(401)
        ctx.Output.JSON(map[string]string{
            "error": "Unauthorized",
        }, false, false)
        return
    }
}

// 필터 적용
web.InsertFilter("/api/v1/*", web.BeforeRouter, FilterAuth)
```

### 6. 캐싱
```go
import "github.com/beego/beego/v2/client/cache"

// 캐시 초기화
bm, err := cache.NewCache("memory", `{"interval":60}`)

// 캐시 사용
bm.Put(ctx, "key", "value", 3600*time.Second)
val := bm.Get(ctx, "key")
```

### 7. 세션 관리
```go
// 세션 활성화 (app.conf)
sessionon = true

// 컨트롤러에서 사용
func (c *UserController) Login() {
    c.SetSession("user", "username")
    user := c.GetSession("user")
    c.DelSession("user")
}
```

### 8. 설정 파일 (app.conf)
```ini
appname = myapp
httpport = 8080
runmode = dev

# 데이터베이스
dbdriver = mysql
dbuser = root
dbpass = password
dbname = mydb
dbhost = localhost
dbport = 3306

# 세션
sessionon = true
sessionprovider = file
sessionpath = ./tmp/sessions

# 로그
loglevel = debug
```

## Beego의 고급 기능

### 1. Swagger 문서 자동 생성
```go
// @Title Get Item
// @Description get item by id
// @Param   id     path    int     true        "Item ID"
// @Success 200 {object} Item
// @Failure 404 not found
// @router /:id [get]
func (c *ItemsController) Get() {
    // ...
}
```

```bash
# 문서 생성
bee generate docs

# Swagger UI 접근
# http://localhost:8080/swagger/
```

### 2. 태스크 스케줄링
```go
import "github.com/beego/beego/v2/task"

func init() {
    tk := task.NewTask("clean_cache", "0 0 * * * *", cleanCache)
    task.AddTask("clean_cache", tk)
    task.StartTask()
}

func cleanCache() error {
    // 캐시 정리 로직
    return nil
}
```

### 3. i18n 다국어 지원
```go
import "github.com/beego/beego/v2/core/logs"

// messages 디렉토리에 언어 파일 생성
// en-US.ini
// ko-KR.ini

func (c *MainController) Get() {
    c.Data["hello"] = c.Tr("hello")
    c.TplName = "index.tpl"
}
```

### 4. 로깅
```go
import "github.com/beego/beego/v2/core/logs"

logs.Debug("Debug message")
logs.Info("Info message")
logs.Warn("Warning message")
logs.Error("Error message")
logs.Critical("Critical message")
```

## 데이터베이스 연동 예제

### MySQL 연동
```go
import (
    "github.com/beego/beego/v2/client/orm"
    _ "github.com/go-sql-driver/mysql"
)

func init() {
    // DB 등록
    orm.RegisterDriver("mysql", orm.DRMySQL)
    orm.RegisterDataBase("default", "mysql", 
        "root:password@tcp(127.0.0.1:3306)/mydb?charset=utf8")
    
    // 모델 등록
    orm.RegisterModel(new(Item))
    
    // 테이블 자동 생성
    orm.RunSyncdb("default", false, true)
}
```

### CRUD 작업
```go
o := orm.NewOrm()

// Create
item := Item{Name: "Test", Price: 100}
id, err := o.Insert(&item)

// Read
item := Item{Id: 1}
err := o.Read(&item)

// Update
item.Price = 200
num, err := o.Update(&item)

// Delete
num, err := o.Delete(&item)

// Query
var items []Item
num, err := o.QueryTable("item").All(&items)
```

## 성능 비교

| 프레임워크 | Requests/sec | 메모리 사용 |
|----------|--------------|-------------|
| Gin | 50,000+ | 낮음 |
| Chi | 50,000+ | 매우 낮음 |
| Echo | 48,000+ | 낮음 |
| Beego | 40,000+ | 보통 |

Beego는 다른 프레임워크보다 약간 느리지만, 더 많은 기능을 제공합니다.

## Gin/Chi/Echo vs Beego

### Gin/Chi/Echo (마이크로 프레임워크)
- ✅ 더 빠름
- ✅ 더 가벼움
- ✅ 유연성
- ❌ ORM 별도 설치
- ❌ 많은 설정 필요

### Beego (풀스택 프레임워크)
- ✅ ORM 내장
- ✅ 모든 기능 포함
- ✅ 자동 문서화
- ✅ 엔터프라이즈 기능
- ❌ 약간 느림
- ❌ 더 무거움

## 언제 Beego를 선택해야 할까?

### Beego를 선택하세요 ✅
- 대규모 엔터프라이즈 애플리케이션
- ORM이 필요한 경우
- 자동 문서화가 필요한 경우
- 풀스택 MVC 패턴을 원할 때
- 캐싱, 세션 등 모든 기능이 필요한 경우
- 중국 시장 (Beego는 중국에서 매우 인기)

### 다른 프레임워크를 선택하세요
- **Gin**: 마이크로서비스, 최고 성능
- **Chi**: 표준 호환, 경량
- **Echo**: 균형잡힌 기능과 성능

## 다음 단계

1. **ORM 사용**: MySQL/PostgreSQL 연동
2. **Swagger 문서**: 자동 API 문서 생성
3. **인증**: JWT, OAuth2
4. **캐싱**: Redis 연동
5. **세션**: 세션 관리
6. **Task**: 스케줄링
7. **i18n**: 다국어 지원
8. **테스트**: httptest 사용

## Bee CLI 명령어

```bash
# 새 프로젝트
bee new myproject
bee api myapi

# 코드 생성
bee generate scaffold item -fields="name:string,price:int"
bee generate model item
bee generate controller item

# 서버 실행
bee run

# 문서 생성
bee generate docs

# 패키징
bee pack
```

## 참고 자료

- [Beego 공식 문서](https://beego.wiki/)
- [Beego GitHub](https://github.com/beego/beego)
- [Bee ORM 문서](https://beego.wiki/docs/mvc/model/overview/)
- [Beego 예제](https://github.com/beego/samples)

