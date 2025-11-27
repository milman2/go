# ✅ Spanner + yo 설치 체크리스트

## 📋 시작 전 확인사항

### 1. Docker 설치 확인
```bash
docker --version
# Docker version 20.10.x 이상
```

### 2. Go 설치 확인
```bash
go version
# go version go1.21 이상
```

### 3. gcloud CLI 설치 확인
```bash
gcloud --version
# Google Cloud SDK 400.0.0 이상
```

## 🔧 도구 설치

### yo 설치
```bash
go install go.mercari.io/yo@latest
which yo
# /home/milman2/go/bin/yo
```

### Hammer 설치
```bash
go install github.com/daichirata/hammer@v0.6.0
which hammer
# /home/milman2/go/bin/hammer
```

### Wrench 설치
```bash
go install github.com/cloudspannerecosystem/wrench@v1.0.4
which wrench
# /home/milman2/go/bin/wrench
```

**또는 한번에:**
```bash
make install-tools
```

## 🐳 Spanner Emulator 설정

### 기존 Spanner 확인
```bash
docker ps | grep spanner
```

**실행 중이면:** ✅ 건너뛰기

**실행 중이 아니면:**
```bash
make docker-up
```

### 연결 테스트
```bash
curl http://localhost:9020
# 응답이 오면 OK
```

## 🗄️ Instance/Database 생성

### gcloud 설정 (Emulator용)
```bash
gcloud config set auth/disable_credentials true
gcloud config set project test-project
gcloud config set api_endpoint_overrides/spanner http://localhost:9020/
```

### Instance 생성
```bash
gcloud spanner instances create test-instance \
  --config=emulator-config \
  --description="Test Instance" \
  --nodes=1
```

### Database 생성
```bash
gcloud spanner databases create test-db \
  --instance=test-instance
```

**또는:**
```bash
make setup-instance
```

### 확인
```bash
gcloud spanner instances list
gcloud spanner databases list --instance=test-instance
```

## 📊 마이그레이션 실행

### 환경 변수 설정
```bash
export SPANNER_EMULATOR_HOST=localhost:9010
```

### Wrench로 마이그레이션
```bash
wrench migrate up \
  --directory migrations \
  --database projects/test-project/instances/test-instance/databases/test-db
```

**또는:**
```bash
make migrate-up-wrench
```

### 확인
```bash
make show-schema
```

출력:
```
Applied Migrations:
  000001_create_users
  000002_create_posts
```

## 🔨 yo로 코드 생성

### yo 실행
```bash
yo test-project test-instance test-db \
  -o models -p models
```

**또는:**
```bash
make generate-yo
```

### 생성 파일 확인
```bash
ls -lh models/
# user.yo.go
# post.yo.go
# yo_db.yo.go
```

## ✅ 최종 체크리스트

- [ ] Docker Spanner emulator 실행 중
- [ ] yo, hammer, wrench 설치됨
- [ ] gcloud CLI emulator 설정 완료
- [ ] Instance `test-instance` 생성
- [ ] Database `test-db` 생성
- [ ] 마이그레이션 2개 적용됨
- [ ] models/ 디렉토리에 3개 파일 생성
- [ ] SPANNER_EMULATOR_HOST 환경 변수 설정

## 🚀 서버 실행

```bash
make run
```

출력:
```
✅ Spanner 연결 성공: projects/test-project/instances/test-instance/databases/test-db

🚀 Spanner + yo 서버 시작
=========================================
📦 Database: Google Cloud Spanner
🔨 Code Generator: yo (go.mercari.io/yo)
🔧 Migration: Hammer + Wrench

🌐 서버 주소: http://localhost:8080
=========================================
```

## 🧪 테스트 실행

다른 터미널:
```bash
make test
```

출력:
```
🧪 Google Cloud Spanner + yo API 테스트
===========================================

1️⃣ Health Check
{"status":"ok","database":"spanner"}

2️⃣ 사용자 생성 (Spanner INSERT)
{
  "id": "...",
  "email": "alice@spanner.com",
  "name": "Alice (Spanner)"
}

...

✅ 테스트 완료
```

## 🎉 완료!

모든 체크리스트를 통과했다면 설치가 완료되었습니다!

### 다음 단계

1. **코드 탐색**: `models/` 디렉토리 확인
2. **API 테스트**: `test.sh` 실행
3. **문서 읽기**:
   - `README.md`: 전체 개요
   - `QUICK_START.md`: 빠른 시작
   - `YO_GUIDE.md`: yo 상세 가이드
   - `USAGE.md`: 사용법
   - `DOCKER_GUIDE.md`: Docker 가이드

### 문제가 있다면?

각 섹션의 "확인" 단계를 다시 실행해보세요.

```bash
# 전체 초기화 (한번에 모든 것을 다시)
make init
```

Happy Coding! 🚀

