# 🚀 빠른 시작 가이드

## ⚡ 30초 시작

```bash
# 1. 전체 초기화 (한번에!)
make init

# 2. 서버 실행
make run

# 3. 테스트 (다른 터미널)
make test
```

끝! 🎉

## 📋 단계별 설명

### 1️⃣ Spanner Emulator 시작

```bash
# 기존 Spanner가 실행 중이라면 건너뛰기
docker ps | grep spanner

# 새로 시작하려면
make docker-up
```

**확인**:
```bash
curl http://localhost:9020
# Spanner emulator 응답이 오면 OK
```

### 2️⃣ Instance & Database 생성

```bash
make setup-instance
```

**무엇을 하나요?**:
- Spanner Instance 생성: `test-instance`
- Spanner Database 생성: `test-db`

### 3️⃣ 데이터베이스 생성 (hammer)

```bash
make createdb
```

**무엇을 하나요?**:
- `schema/schema.sql` 파일 읽기
- 데이터베이스 생성 및 스키마 적용
- 테이블 생성 (users, posts)
- 인덱스 생성

### 4️⃣ 샘플 데이터 삽입 (선택사항)

```bash
make seed-data
```

**무엇을 하나요?**:
- 3명의 샘플 사용자 생성
- 5개의 샘플 게시글 생성
- 개발/테스트에 유용

**확인**:
```bash
make test-query
# 사용자 3명, 게시글 5개 확인
```

### 5️⃣ yo로 코드 생성

```bash
make generate-models
```

**무엇을 하나요?**:
- Spanner 스키마를 읽음
- `models/` 디렉토리에 Go 코드 생성
- User, Post 구조체 및 CRUD 메서드 생성

**생성되는 파일**:
```
models/
├── user.yo.go       # User 모델 + CRUD
├── post.yo.go       # Post 모델 + CRUD
└── yo_db.yo.go      # DB 헬퍼 함수
```

### 6️⃣ 서버 실행

```bash
make run
```

### 7️⃣ API 테스트

```bash
# 다른 터미널에서
make test
```

## 🎯 주요 명령어

### 개발 중

```bash
# 스키마 변경 시
1. schema/schema.sql 파일 수정
2. make db-diff          # 변경사항 확인
3. make db-apply         # 변경사항 적용
4. make generate-models  # 코드 재생성

# DB 리셋
make resetdb             # DB 전체 리셋
make seed-data           # 샘플 데이터 다시 삽입

# 데이터만 리셋
make clear-data          # 데이터 삭제
make seed-data           # 샘플 데이터 삽입

# 생성된 코드 확인
ls -lh models/*.yo.go
```

### 디버깅

```bash
# Docker 상태
make docker-ps

# Spanner 설정 확인
make info

# 스키마 상태
make show-schema

# Spanner CLI 접속
make spanner-cli
```

## 🧹 정리

```bash
# 생성된 파일만 삭제
make clean

# Docker도 중지
make docker-down
```

## 💡 팁

### Spanner Emulator 이미 실행 중?

```bash
# 확인
docker ps | grep spanner

# 사용
export SPANNER_EMULATOR_HOST=localhost:9010
make setup-instance
make createdb
make generate-models
```

### 스키마 파일 수정

```sql
-- schema/schema.sql
CREATE TABLE users (
  id STRING(36) NOT NULL,
  email STRING(255) NOT NULL,
  name STRING(100) NOT NULL,
  age INT64,  -- ✨ 새 컬럼 추가
  created_at TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true),
  updated_at TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true),
) PRIMARY KEY (id);
```

그리고:
```bash
make db-diff   # 차이 확인
make db-apply  # 적용
make generate-models  # 코드 재생성
```

## 🎉 완성!

이제 Spanner + yo를 사용할 준비가 되었습니다! 🚀

더 자세한 내용은 `README.md`를 확인하세요.

