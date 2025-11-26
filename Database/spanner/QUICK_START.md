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
- Spanner Database 생성: `test-database`

### 3️⃣ 마이그레이션 실행

```bash
# Wrench 사용 (권장)
make migrate-up-wrench

# 또는 Hammer 사용
make migrate-up-hammer
```

**무엇을 하나요?**:
- `migrations/*.up.sql` 파일 실행
- 테이블 생성 (users, posts)
- 인덱스 생성

### 4️⃣ yo로 코드 생성

```bash
make generate-yo
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

### 5️⃣ 서버 실행

```bash
make run
```

### 6️⃣ API 테스트

```bash
# 다른 터미널에서
make test
```

## 🎯 주요 명령어

### 개발 중

```bash
# 스키마 변경 시
1. migrations/ 에 새 SQL 파일 추가
2. make migrate-up-wrench
3. make generate-yo

# DB 리셋
make reset

# 생성된 코드 확인
ls -lh models/
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
make migrate-up-wrench
make generate-yo
```

### 마이그레이션 파일 작성

```sql
-- migrations/000003_add_column.up.sql
ALTER TABLE users ADD COLUMN age INT64;

-- migrations/000003_add_column.down.sql
ALTER TABLE users DROP COLUMN age;
```

그리고:
```bash
make reset  # 마이그레이션 + 코드 재생성
```

## 🎉 완성!

이제 Spanner + yo를 사용할 준비가 되었습니다! 🚀

더 자세한 내용은 `README.md`를 확인하세요.

