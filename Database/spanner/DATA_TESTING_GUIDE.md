# 샘플 데이터 테스팅 가이드

Spanner 데이터베이스에 샘플 데이터를 삽입하고 테스트하는 방법을 설명합니다.

## 🚀 빠른 시작

### 1. 처음부터 시작 (DB + 샘플 데이터)

```bash
# 전체 초기화
make init

# 샘플 데이터 삽입
make seed-data

# 데이터 확인
make test-query
```

### 2. 기존 DB에 샘플 데이터만 추가

```bash
# 샘플 데이터 삽입
make seed-data

# 데이터 확인
make test-query
```

## 📋 주요 명령어

### seed-data - 샘플 데이터 삽입

```bash
make seed-data
```

**무엇을 하나요?**
- 3명의 사용자 생성
- 5개의 게시글 생성 (published 3개, draft 2개)

**샘플 데이터:**

| 사용자 | 이메일 | 게시글 수 |
|--------|--------|-----------|
| John Doe | john.doe@example.com | 2개 (1 published, 1 draft) |
| Jane Smith | jane.smith@example.com | 2개 (2 published) |
| Bob Johnson | bob.johnson@example.com | 1개 (1 draft) |

### test-query - 데이터 확인

```bash
make test-query
```

**출력 예시:**
```
📊 Users count:
user_count
3

📊 Posts count:
post_count
5

📊 Published posts by user:
name        post_count
Jane Smith  2
John Doe    1
```

### clear-data - 모든 데이터 삭제

```bash
make clear-data
```

**주의:** 
- 테이블 구조는 유지됩니다
- 모든 데이터만 삭제됩니다
- 확인 프롬프트가 표시됩니다

## 🔄 워크플로우 예시

### 시나리오 1: 개발 중 데이터 리셋

```bash
# 1. 데이터만 삭제
make clear-data

# 2. 새로운 샘플 데이터 삽입
make seed-data

# 3. 확인
make test-query
```

### 시나리오 2: 완전 리셋 (스키마 + 데이터)

```bash
# 1. DB 전체 리셋
make resetdb

# 2. 샘플 데이터 삽입
make seed-data

# 3. 모델 재생성
make generate-models
```

### 시나리오 3: 커스텀 쿼리 테스트

```bash
# 샘플 데이터 삽입
make seed-data

# 직접 SQL 실행
export SPANNER_EMULATOR_HOST=localhost:9010
gcloud spanner databases execute-sql test-db \
  --instance=test-instance \
  --sql="SELECT u.name, p.title, p.published 
         FROM users u 
         JOIN posts p ON u.id = p.user_id 
         WHERE u.email = 'john.doe@example.com'"
```

## 💡 유용한 쿼리 모음

### 1. 특정 사용자의 모든 게시글

```sql
SELECT p.title, p.published, p.created_at
FROM posts p
JOIN users u ON p.user_id = u.id
WHERE u.email = 'john.doe@example.com'
ORDER BY p.created_at DESC;
```

### 2. 발행된 게시글만 조회

```sql
SELECT u.name, p.title, p.created_at
FROM posts p
JOIN users u ON p.user_id = u.id
WHERE p.published = TRUE
ORDER BY p.created_at DESC;
```

### 3. 사용자별 게시글 통계

```sql
SELECT 
  u.name,
  COUNT(p.id) as total_posts,
  SUM(CASE WHEN p.published THEN 1 ELSE 0 END) as published,
  SUM(CASE WHEN NOT p.published THEN 1 ELSE 0 END) as drafts
FROM users u
LEFT JOIN posts p ON u.id = p.user_id
GROUP BY u.name;
```

### 4. 최근 게시글 5개

```sql
SELECT u.name as author, p.title, p.published, p.created_at
FROM posts p
JOIN users u ON p.user_id = u.id
ORDER BY p.created_at DESC
LIMIT 5;
```

## 🛠️ 고급 사용법

### 커스텀 샘플 데이터 추가

`scripts/seed_data.go` 파일을 수정하여 자신만의 샘플 데이터를 추가할 수 있습니다:

```go
// 새로운 사용자 추가
spanner.Insert("users",
    []string{"id", "email", "name", "created_at", "updated_at"},
    []interface{}{
        "your-uuid-here",
        "new.user@example.com",
        "New User",
        spanner.CommitTimestamp,
        spanner.CommitTimestamp,
    }),
```

### 프로그래밍 방식으로 데이터 삽입

```go
package main

import (
    "context"
    "fmt"
    
    "cloud.google.com/go/spanner"
)

func main() {
    ctx := context.Background()
    client, _ := spanner.NewClient(ctx, "projects/test-project/instances/test-instance/databases/test-db")
    defer client.Close()

    _, err := client.Apply(ctx, []*spanner.Mutation{
        spanner.Insert("users",
            []string{"id", "email", "name", "created_at", "updated_at"},
            []interface{}{"uuid", "test@example.com", "Test User", spanner.CommitTimestamp, spanner.CommitTimestamp}),
    })
    
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    }
}
```

## 🐛 트러블슈팅

### 문제: "Row already exists"

**증상:**
```
Error: Row [uuid] already exists
```

**해결:**
```bash
make clear-data  # 기존 데이터 삭제
make seed-data   # 다시 삽입
```

### 문제: "Foreign key constraint violation"

**증상:**
```
Error: Foreign key constraint violated
```

**원인:** posts를 users보다 먼저 삽입하려고 함

**해결:** `seed_data.go`에서 users를 먼저 삽입하는지 확인

### 문제: 데이터가 보이지 않음

**확인 사항:**
1. Emulator가 실행 중인지 확인: `docker ps | grep spanner`
2. 올바른 데이터베이스에 연결했는지 확인
3. 트랜잭션이 커밋되었는지 확인

```bash
# 데이터 확인
make test-query
```

## 📚 추가 리소스

- **샘플 데이터 상세:** `migrations/dml/README.md`
- **스크립트 코드:** `scripts/seed_data.go`
- **전체 명령어:** `make help`

## 🎯 CI/CD에서 사용

```yaml
# .github/workflows/test.yml
- name: Setup Database
  run: |
    make docker-up
    make setup-instance
    make createdb
    make seed-data

- name: Run Tests
  run: |
    make test-query
    go test ./...
```

## 📊 성능 테스트용 대량 데이터

대량의 테스트 데이터가 필요한 경우:

```bash
# scripts/seed_data.go를 수정하여 루프 추가
for i := 0; i < 1000; i++ {
    // 1000명의 사용자 생성
}
```

## 🔐 보안 주의사항

⚠️ **중요:**
- 샘플 데이터는 **개발/테스트 환경**에만 사용
- 운영 환경에서는 절대 사용하지 마세요
- 실제 이메일 주소 사용 금지
- 민감한 정보 포함 금지

## ✅ 체크리스트

샘플 데이터 테스트 전:
- [ ] Spanner emulator 실행 중
- [ ] Instance 생성 완료
- [ ] Database 생성 완료
- [ ] 스키마 적용 완료

샘플 데이터 삽입 후:
- [ ] `make test-query`로 데이터 확인
- [ ] 사용자 3명 존재
- [ ] 게시글 5개 존재
- [ ] 조인 쿼리 정상 작동

