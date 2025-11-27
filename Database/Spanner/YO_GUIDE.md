# yo 완전 가이드 (Google Cloud Spanner Code Generator)

## 🎯 yo란?

**yo** = **Y**et another **O**RM? 아니요!

**yo** = Code generator for **Google Cloud Spanner**

Mercari에서 개발한 Spanner용 코드 생성 도구입니다.

- 공식 문서: https://pkg.go.dev/go.mercari.io/yo
- GitHub: https://github.com/mercari/yo

## 🔑 핵심 개념

### yo는 ORM이 아닙니다!

```
ORM (GORM, Hibernate)        yo
─────────────────────        ───────────────────────
런타임에 SQL 생성            빌드 타임에 코드 생성
동적 쿼리 작성               스키마 기반 타입 안전
성능 오버헤드 있음           오버헤드 없음 (네이티브)
```

### yo는 코드 생성기입니다!

```
Spanner Schema (INFORMATION_SCHEMA)
         ↓
      yo 실행
         ↓
  Go Code (models/*.yo.go)
```

## 📦 설치

```bash
# Go 1.16+
go install go.mercari.io/yo@latest

# 또는 Makefile
make install-tools
```

## 🚀 기본 사용법

### 1. 스키마 준비

```sql
CREATE TABLE users (
  id STRING(36) NOT NULL,
  email STRING(255) NOT NULL,
  name STRING(100) NOT NULL,
  created_at TIMESTAMP NOT NULL,
) PRIMARY KEY (id);

CREATE UNIQUE INDEX users_email_idx ON users(email);
```

### 2. yo 실행

```bash
SPANNER_EMULATOR_HOST=localhost:9010 \
yo PROJECT_ID INSTANCE_ID DATABASE_ID \
  -o models \
  -p models
```

### 3. 생성되는 코드

#### 구조체

```go
// models/user.yo.go
type User struct {
    ID        string    `spanner:"id" json:"id"`
    Email     string    `spanner:"email" json:"email"`
    Name      string    `spanner:"name" json:"name"`
    CreatedAt time.Time `spanner:"created_at" json:"created_at"`
}
```

#### Mutation Methods (INSERT/UPDATE/DELETE)

```go
// Insert - 새 레코드 삽입
func (u *User) Insert(ctx context.Context) *spanner.Mutation

// Update - 전체 컬럼 업데이트
func (u *User) Update(ctx context.Context) *spanner.Mutation

// InsertOrUpdate - Upsert
func (u *User) InsertOrUpdate(ctx context.Context) *spanner.Mutation

// UpdateColumns - 특정 컬럼만 업데이트
func (u *User) UpdateColumns(ctx context.Context, columns ...string) *spanner.Mutation

// Delete - 레코드 삭제
func (u *User) Delete(ctx context.Context) *spanner.Mutation
```

#### Read Functions (인덱스 기반)

```go
// Primary Key로 조회
func FindUserByID(ctx context.Context, db YODB, id string) (*User, error)

// Unique Index로 조회
func FindUserByEmail(ctx context.Context, db YODB, email string) (*User, error)

// 전체 조회
func FindAllUsers(ctx context.Context, db YODB) ([]*User, error)
```

## 💻 실제 사용 예제

### 생성 (Insert)

```go
import (
    "cloud.google.com/go/spanner"
    "github.com/milman2/go-api/spanner-yo/models"
)

func createUser(client *spanner.Client) error {
    ctx := context.Background()
    
    // yo 생성 모델 사용
    user := &models.User{
        ID:    uuid.New().String(),
        Email: "test@example.com",
        Name:  "Test User",
    }
    
    // Mutation 생성 (yo가 생성한 메서드)
    mutation := user.Insert(ctx)
    
    // Spanner에 적용
    _, err := client.Apply(ctx, []*spanner.Mutation{mutation})
    return err
}
```

### 조회 (Read)

```go
func getUser(client *spanner.Client, id string) (*models.User, error) {
    ctx := context.Background()
    
    // yo가 생성한 Find 함수 사용
    user, err := models.FindUserByID(ctx, client, id)
    if err != nil {
        return nil, err
    }
    
    return user, nil
}

func getUserByEmail(client *spanner.Client, email string) (*models.User, error) {
    ctx := context.Background()
    
    // Unique Index 기반 조회
    user, err := models.FindUserByEmail(ctx, client, email)
    return user, err
}

func getAllUsers(client *spanner.Client) ([]*models.User, error) {
    ctx := context.Background()
    
    // 전체 조회
    users, err := models.FindAllUsers(ctx, client)
    return users, err
}
```

### 수정 (Update)

```go
func updateUser(client *spanner.Client, id, newName string) error {
    ctx := context.Background()
    
    // 1. 조회
    user, err := models.FindUserByID(ctx, client, id)
    if err != nil {
        return err
    }
    
    // 2. 값 변경
    user.Name = newName
    
    // 3. 업데이트 (전체 컬럼)
    _, err = client.Apply(ctx, []*spanner.Mutation{
        user.Update(ctx),
    })
    
    return err
}

func updateUserColumns(client *spanner.Client, id, newName string) error {
    ctx := context.Background()
    
    user, _ := models.FindUserByID(ctx, client, id)
    user.Name = newName
    
    // 특정 컬럼만 업데이트
    _, err := client.Apply(ctx, []*spanner.Mutation{
        user.UpdateColumns(ctx, "name", "updated_at"),
    })
    
    return err
}
```

### 삭제 (Delete)

```go
func deleteUser(client *spanner.Client, id string) error {
    ctx := context.Background()
    
    // 1. 조회
    user, err := models.FindUserByID(ctx, client, id)
    if err != nil {
        return err
    }
    
    // 2. 삭제
    _, err = client.Apply(ctx, []*spanner.Mutation{
        user.Delete(ctx),
    })
    
    return err
}
```

## 🎨 yo 고급 기능

### 1. 커스텀 타입

```yaml
# custom_column_types.yml
tables:
  users:
    columns:
      status:
        type: UserStatus  # 커스텀 타입 사용
```

```bash
yo PROJECT INSTANCE DATABASE \
  --custom-types-file custom_column_types.yml \
  --custom-type-package ./types \
  -o models
```

### 2. 필드/테이블 제외

```bash
yo PROJECT INSTANCE DATABASE \
  --ignore-tables "migrations,schema_history" \
  --ignore-fields "internal_field" \
  -o models
```

### 3. 단일 파일 생성

```bash
yo PROJECT INSTANCE DATABASE \
  --single-file \
  -o models/all.yo.go
```

### 4. 커스텀 템플릿

```bash
# 템플릿 복사
cp -r $GOPATH/src/github.com/mercari/yo/templates ./templates

# 템플릿 수정
vi templates/type.go.tpl

# 커스텀 템플릿 사용
yo PROJECT INSTANCE DATABASE \
  --template-path ./templates \
  -o models
```

## 🔄 워크플로우

### 일반적인 개발 흐름

```
1. 마이그레이션 파일 작성
   └─ migrations/000003_add_table.up.sql

2. 마이그레이션 실행
   └─ make migrate-up-wrench

3. yo로 코드 생성
   └─ make generate-yo

4. 생성된 모델 사용
   └─ import "project/models"
      user := &models.User{...}
      user.Insert(ctx)
```

### 스키마 변경 시

```bash
# 1. 마이그레이션 파일 작성
vim migrations/000003_add_age.up.sql

# 2. 마이그레이션 + 코드 재생성
make reset
```

## 💡 yo vs ORM

| 특징 | yo | GORM/ORM |
|------|-----|-----------|
| **타입 안전** | ✅✅✅ 완벽 | ✅ 좋음 |
| **성능** | ✅✅✅ 네이티브 | ✅ 오버헤드 있음 |
| **보일러플레이트** | ✅ 자동 생성 | ✅ 적음 |
| **복잡한 쿼리** | ❌ Raw SQL 필요 | ✅ 체이닝 가능 |
| **러닝 커브** | 낮음 | 보통 |
| **스키마 변경** | 재생성 필요 | 자동 반영 |

## 🎯 yo의 장점

### 1. 완벽한 타입 안전성
```go
// 컴파일 타임에 모든 것을 체크
user := &models.User{
    ID:    "abc",  // string
    Email: "...",  // string
    Name:  123,    // ❌ 컴파일 에러!
}
```

### 2. 제로 오버헤드
```go
// ORM: 런타임에 SQL 생성
db.Where("email = ?", email).First(&user)

// yo: 미리 생성된 코드 실행
models.FindUserByEmail(ctx, client, email)
```

### 3. Spanner 네이티브
```go
// Spanner의 Mutation API를 직접 사용
user.Insert(ctx)  // spanner.Mutation 반환
```

## 🚨 yo의 단점

### 1. 스키마 변경 시 재생성 필요

```bash
# 스키마 변경 후
make generate-yo  # 매번 실행 필요
```

### 2. 복잡한 쿼리는 직접 작성

```go
// yo는 간단한 CRUD만
user, _ := models.FindUserByID(ctx, client, id)

// 복잡한 쿼리는 Raw SQL
stmt := spanner.Statement{
    SQL: `SELECT u.*, COUNT(p.id) as post_count
          FROM users u
          LEFT JOIN posts p ON p.user_id = u.id
          GROUP BY u.id`,
}
```

## 📊 생성 코드 구조

### models/user.yo.go

```go
package models

import (
    "cloud.google.com/go/spanner"
    "context"
    "time"
)

// User - yo가 생성한 구조체
type User struct {
    ID        string    `spanner:"id" json:"id"`
    Email     string    `spanner:"email" json:"email"`
    Name      string    `spanner:"name" json:"name"`
    CreatedAt time.Time `spanner:"created_at" json:"created_at"`
    UpdatedAt time.Time `spanner:"updated_at" json:"updated_at"`
}

// Insert - INSERT Mutation 생성
func (u *User) Insert(ctx context.Context) *spanner.Mutation {
    return spanner.Insert("users", /* ... */)
}

// FindUserByID - Primary Key로 조회
func FindUserByID(ctx context.Context, db YODB, id string) (*User, error) {
    /* ... */
}

// FindUserByEmail - Unique Index로 조회
func FindUserByEmail(ctx context.Context, db YODB, email string) (*User, error) {
    /* ... */
}
```

## 🔧 Makefile 통합

현재 프로젝트의 Makefile:

```makefile
# 코드 생성
make generate-yo

# 마이그레이션 + 생성
make reset

# 전체 초기화
make init
```

## 🎓 모범 사례

### 1. 마이그레이션 파일 명명

```
000001_create_users.up.sql      # UP
000001_create_users.down.sql    # DOWN
000002_create_posts.up.sql
000002_create_posts.down.sql
```

### 2. yo 재생성은 자주

```bash
# 스키마 변경 후 즉시
make generate-yo
```

### 3. 생성 파일은 Git에 포함

```gitignore
# ❌ 제외하지 마세요
# models/*.yo.go

# ✅ 포함하세요 (검토 가능)
models/*.yo.go
```

### 4. Clean Architecture에서 사용

```go
// Repository 레이어에서만 yo 모델 사용
package repository

import "project/models"

type UserRepository struct {
    client *spanner.Client
}

func (r *UserRepository) Create(user *domain.User) error {
    // domain.User → models.User 변환
    model := &models.User{
        ID:    user.ID,
        Email: user.Email,
        Name:  user.Name,
    }
    
    // yo 생성 메서드 사용
    _, err := r.client.Apply(ctx, []*spanner.Mutation{
        model.Insert(ctx),
    })
    
    return err
}
```

## 📚 참고 자료

- [yo 공식 문서](https://pkg.go.dev/go.mercari.io/yo)
- [yo GitHub](https://github.com/mercari/yo)
- [Cloud Spanner](https://cloud.google.com/spanner)
- [Spanner Go Client](https://pkg.go.dev/cloud.google.com/go/spanner)

## 🎉 결론

**yo는 Spanner를 위한 최고의 코드 생성 도구입니다!**

장점:
- ✅ 타입 안전
- ✅ 제로 오버헤드
- ✅ Spanner 네이티브
- ✅ 보일러플레이트 제거

단점:
- ❌ Spanner 전용
- ❌ 재생성 필요

**Spanner를 사용한다면 yo는 필수!** 🚀

