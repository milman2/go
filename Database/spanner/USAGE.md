# 📘 Spanner + yo 사용 가이드

## 🎯 전체 워크플로우

```
1. Spanner 준비 (Docker Emulator)
   ↓
2. Instance/Database 생성
   ↓
3. 마이그레이션 실행 (Hammer/Wrench)
   ↓
4. yo로 코드 생성
   ↓
5. 생성된 모델 사용
   ↓
6. 서버 실행 & 테스트
```

## 🚀 빠른 시작 (3단계)

### 1단계: 전체 초기화

```bash
cd /home/milman2/go-api/go/Database/spanner
make init
```

**실행 내용**:
- ✅ Docker Spanner emulator 시작
- ✅ 도구 설치 (yo, hammer, wrench)
- ✅ Instance/Database 생성
- ✅ 마이그레이션 실행
- ✅ yo 코드 생성

### 2단계: 서버 실행

```bash
make run
```

### 3단계: 테스트

```bash
# 다른 터미널
make test
```

끝! 🎉

## 📝 단계별 설명

### Step 1: 기존 Spanner 사용

현재 시스템에 Spanner가 실행 중:

```bash
# 확인
docker ps | grep spanner

# 출력:
# school-live-api-spanner-1 ... Up 2 months ... 0.0.0.0:9010->9010/tcp
```

**이미 실행 중이므로 건너뛰고 다음 단계로!**

### Step 2: Instance/Database 생성

```bash
# gcloud 설정 (Emulator용)
gcloud config set auth/disable_credentials true
gcloud config set project test-project
gcloud config set api_endpoint_overrides/spanner http://localhost:9020/

# Instance 생성
gcloud spanner instances create test-instance \
  --config=emulator-config \
  --description="Test Instance" \
  --nodes=1

# Database 생성
gcloud spanner databases create test-database \
  --instance=test-instance
```

**또는 Makefile:**

```bash
make setup-instance
```

### Step 3: 마이그레이션 실행

#### Wrench 사용 (권장)

```bash
# UP
SPANNER_EMULATOR_HOST=localhost:9010 \
wrench migrate up \
  --directory migrations \
  --database projects/test-project/instances/test-instance/databases/test-database
```

**또는:**

```bash
make migrate-up-wrench
```

#### Hammer 사용

```bash
# UP
SPANNER_EMULATOR_HOST=localhost:9010 \
hammer -p test-project -i test-instance -d test-database \
  -m migrations up
```

**또는:**

```bash
make migrate-up-hammer
```

### Step 4: yo로 코드 생성

```bash
# yo 실행
SPANNER_EMULATOR_HOST=localhost:9010 \
yo test-project test-instance test-database \
  -o models -p models
```

**또는:**

```bash
make generate-yo
```

**생성되는 파일:**

```
models/
├── user.yo.go       # User 모델
├── post.yo.go       # Post 모델
└── yo_db.yo.go      # DB 헬퍼
```

### Step 5: 생성된 코드 확인

```bash
# 파일 목록
ls -lh models/

# 내용 확인
cat models/user.yo.go
```

### Step 6: 서버 실행

```bash
# 환경 변수와 함께 실행
SPANNER_EMULATOR_HOST=localhost:9010 \
SPANNER_PROJECT_ID=test-project \
SPANNER_INSTANCE_ID=test-instance \
SPANNER_DATABASE_ID=test-database \
go run cmd/api/main.go
```

**또는:**

```bash
make run
```

### Step 7: API 테스트

```bash
# Health Check
curl http://localhost:8080/health

# 사용자 생성
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@spanner.com",
    "name": "Test User"
  }'

# 사용자 목록
curl http://localhost:8080/api/v1/users

# 또는 테스트 스크립트
./test.sh
# 또는
make test
```

## 🔄 개발 워크플로우

### 새 테이블 추가

#### 1. 마이그레이션 파일 작성

```bash
# UP 파일
vim migrations/000003_create_comments.up.sql
```

```sql
CREATE TABLE comments (
  id STRING(36) NOT NULL,
  post_id STRING(36) NOT NULL,
  user_id STRING(36) NOT NULL,
  content STRING(MAX),
  created_at TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true),
  CONSTRAINT fk_comments_post FOREIGN KEY (post_id) REFERENCES posts (id),
  CONSTRAINT fk_comments_user FOREIGN KEY (user_id) REFERENCES users (id),
) PRIMARY KEY (id);

CREATE INDEX comments_post_id_idx ON comments(post_id);
```

```bash
# DOWN 파일
vim migrations/000003_create_comments.down.sql
```

```sql
DROP INDEX comments_post_id_idx;
DROP TABLE comments;
```

#### 2. 마이그레이션 & 코드 생성

```bash
make reset
```

이 명령어는:
- ✅ 마이그레이션 실행
- ✅ yo 코드 재생성

#### 3. 생성된 모델 확인

```bash
ls models/comment.yo.go
```

### 스키마 변경 (컬럼 추가)

#### 1. 마이그레이션 파일

```sql
-- migrations/000004_add_user_age.up.sql
ALTER TABLE users ADD COLUMN age INT64;

-- migrations/000004_add_user_age.down.sql
ALTER TABLE users DROP COLUMN age;
```

#### 2. 적용 & 재생성

```bash
make reset
```

#### 3. User 구조체 확인

```go
// models/user.yo.go
type User struct {
    ID    string `spanner:"id"`
    Email string `spanner:"email"`
    Name  string `spanner:"name"`
    Age   int64  `spanner:"age"`  // ← 추가됨!
}
```

## 💻 생성된 코드 사용법

### 1. Insert (생성)

```go
import (
    "github.com/milman2/go-api/spanner-yo/models"
    "github.com/google/uuid"
)

func createUser(client *spanner.Client) error {
    ctx := context.Background()
    
    user := &models.User{
        ID:    uuid.New().String(),
        Email: "alice@example.com",
        Name:  "Alice",
    }
    
    // yo 생성 메서드 사용
    mutation := user.Insert(ctx)
    
    _, err := client.Apply(ctx, []*spanner.Mutation{mutation})
    return err
}
```

### 2. Read (조회)

```go
func getUser(client *spanner.Client, id string) (*models.User, error) {
    ctx := context.Background()
    
    // yo가 생성한 Find 함수
    user, err := models.FindUserByID(ctx, client, id)
    return user, err
}

func getUserByEmail(client *spanner.Client, email string) (*models.User, error) {
    ctx := context.Background()
    
    // Unique Index 조회
    user, err := models.FindUserByEmail(ctx, client, email)
    return user, err
}
```

### 3. Update (수정)

```go
func updateUser(client *spanner.Client, id, newName string) error {
    ctx := context.Background()
    
    // 조회
    user, err := models.FindUserByID(ctx, client, id)
    if err != nil {
        return err
    }
    
    // 수정
    user.Name = newName
    
    // 업데이트
    mutation := user.Update(ctx)
    _, err = client.Apply(ctx, []*spanner.Mutation{mutation})
    
    return err
}
```

### 4. Delete (삭제)

```go
func deleteUser(client *spanner.Client, id string) error {
    ctx := context.Background()
    
    user, err := models.FindUserByID(ctx, client, id)
    if err != nil {
        return err
    }
    
    mutation := user.Delete(ctx)
    _, err = client.Apply(ctx, []*spanner.Mutation{mutation})
    
    return err
}
```

## 🔧 Makefile 명령어

```bash
# 도움말
make help

# Docker 관련
make docker-up           # Spanner 시작
make docker-down         # Spanner 중지
make docker-ps           # 상태 확인

# 초기 설정
make setup-instance      # Instance/Database 생성
make install-tools       # yo, hammer, wrench 설치

# 마이그레이션
make migrate-up-wrench   # Wrench UP
make migrate-down-wrench # Wrench DOWN
make migrate-up-hammer   # Hammer UP
make migrate-down-hammer # Hammer DOWN

# 코드 생성
make generate-yo         # yo 실행

# 통합
make init                # 전체 초기화
make reset               # 마이그레이션 + 코드 재생성
make clean               # 생성 파일 삭제

# 실행/테스트
make run                 # 서버 실행
make test                # API 테스트

# 디버깅
make spanner-cli         # Spanner CLI 접속
make show-schema         # 스키마 확인
make info                # 설정 정보
```

## 🐛 문제 해결

### yo 실행 시 "command not found"

```bash
# yo 설치 확인
which yo

# 없다면 설치
go install go.mercari.io/yo@latest

# 또는
make install-tools
```

### "database not found" 에러

```bash
# Instance/Database 확인
gcloud spanner instances list
gcloud spanner databases list --instance=test-instance

# 없다면 생성
make setup-instance
```

### 마이그레이션 실패

```bash
# 상태 확인
make show-schema

# 롤백 후 재시도
make migrate-down-wrench
make migrate-up-wrench
```

### Spanner 연결 실패

```bash
# 환경 변수 확인
echo $SPANNER_EMULATOR_HOST
# 출력: localhost:9010

# 없다면 설정
export SPANNER_EMULATOR_HOST=localhost:9010

# Docker 확인
docker ps | grep spanner
```

## 📊 프로젝트 구조

```
spanner/
├── cmd/
│   └── api/
│       └── main.go              # 서버 진입점
│
├── migrations/                  # 마이그레이션 파일
│   ├── 000001_create_users.up.sql
│   ├── 000001_create_users.down.sql
│   ├── 000002_create_posts.up.sql
│   └── 000002_create_posts.down.sql
│
├── models/                      # yo 생성 코드
│   ├── user.yo.go              # ← yo가 생성
│   ├── post.yo.go              # ← yo가 생성
│   └── yo_db.yo.go             # ← yo가 생성
│
├── docker-compose.yml           # Spanner emulator
├── Makefile                     # 자동화
├── go.mod
├── README.md                    # 전체 문서
├── QUICK_START.md               # 빠른 시작
├── YO_GUIDE.md                  # yo 상세 가이드
├── DOCKER_GUIDE.md              # Docker 가이드
├── USAGE.md                     # 이 파일
└── test.sh                      # 테스트 스크립트
```

## 📚 더 알아보기

### 문서

- `README.md`: 전체 개요 및 소개
- `QUICK_START.md`: 30초 빠른 시작
- `YO_GUIDE.md`: yo 완전 가이드
- `DOCKER_GUIDE.md`: Docker Spanner 가이드
- `USAGE.md`: 이 파일 (상세 사용법)

### 외부 링크

- [yo 공식 문서](https://pkg.go.dev/go.mercari.io/yo)
- [Cloud Spanner 문서](https://cloud.google.com/spanner/docs)
- [Hammer](https://github.com/daichirata/hammer)
- [Wrench](https://github.com/cloudspannerecosystem/wrench)

## 🎉 다음 단계

1. **Clean Architecture 적용**
   - Repository 레이어에서 yo 모델 사용
   - Domain은 순수하게 유지

2. **복잡한 쿼리 추가**
   - Raw SQL + yo 모델 결합

3. **트랜잭션 처리**
   - Spanner의 강력한 트랜잭션 활용

4. **테스트 작성**
   - Emulator로 통합 테스트

5. **프로덕션 배포**
   - 실제 Cloud Spanner로 전환

Happy Coding! 🚀

