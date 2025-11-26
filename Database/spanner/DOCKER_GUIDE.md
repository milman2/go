# 🐳 Docker Spanner Emulator 가이드

## 현재 실행 중인 Spanner

현재 시스템에서 실행 중인 Spanner 컨테이너:

```
Container: school-live-api-spanner-1
Image: gcr.io/cloud-spanner-emulator/emulator:1.5.33
Ports:
  - 9010 (gRPC)
  - 9020 (HTTP)
Status: Up 2 months
```

## 🎯 기존 Spanner 사용하기

### 1. 환경 변수 설정

```bash
export SPANNER_EMULATOR_HOST=localhost:9010
```

### 2. 새 Instance/Database 생성

```bash
# gcloud 설정
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

또는 Makefile 사용:

```bash
make setup-instance
```

### 3. 마이그레이션 실행

```bash
# Wrench 사용 (권장)
make migrate-up-wrench

# 또는 Hammer 사용
make migrate-up-hammer
```

### 4. yo로 코드 생성

```bash
make generate-yo
```

## 🆕 새 Spanner 띄우기

### docker-compose 사용

```bash
# 시작
docker-compose up -d

# 중지
docker-compose down

# 상태 확인
docker-compose ps
```

### 직접 실행

```bash
docker run -d \
  --name spanner-emulator \
  -p 9010:9010 \
  -p 9020:9020 \
  gcr.io/cloud-spanner-emulator/emulator:1.5.33
```

## 🔍 Spanner 상태 확인

### Docker 컨테이너 확인

```bash
docker ps | grep spanner
```

출력:
```
ffc20c6b3aac   gcr.io/cloud-spanner-emulator/emulator:1.5.33
  Up 2 months   0.0.0.0:9010->9010/tcp, 0.0.0.0:9020->9020/tcp
```

### HTTP 엔드포인트 테스트

```bash
curl http://localhost:9020
```

정상 응답이 오면 OK!

### gRPC 연결 테스트

```bash
# Go 코드로 테스트
SPANNER_EMULATOR_HOST=localhost:9010 go run cmd/api/main.go
```

## 📋 Instance/Database 목록 확인

```bash
# Instance 목록
gcloud spanner instances list

# Database 목록
gcloud spanner databases list --instance=test-instance
```

## 🔧 gcloud CLI 설정

### emulator용 설정

```bash
# 인증 비활성화
gcloud config set auth/disable_credentials true

# 프로젝트 설정
gcloud config set project YOUR_PROJECT_ID

# Spanner API 엔드포인트 변경
gcloud config set api_endpoint_overrides/spanner http://localhost:9020/
```

### 원래대로 복원

```bash
# 인증 활성화
gcloud config set auth/disable_credentials false

# API 엔드포인트 초기화
gcloud config unset api_endpoint_overrides/spanner
```

## 💻 Spanner CLI 사용

### docker-compose로 CLI 컨테이너 실행

```bash
# CLI 컨테이너 시작
docker-compose up -d spanner-cli

# CLI 접속
docker-compose exec spanner-cli spanner-cli \
  -p test-project \
  -i test-instance \
  -d test-database
```

또는 Makefile:

```bash
make spanner-cli
```

### CLI에서 쿼리 실행

```sql
spanner> SELECT * FROM users;
spanner> INSERT INTO users (id, email, name, created_at, updated_at)
         VALUES ('test-id', 'test@example.com', 'Test', PENDING_COMMIT_TIMESTAMP(), PENDING_COMMIT_TIMESTAMP());
spanner> SHOW TABLES;
```

## 🐛 문제 해결

### Spanner가 실행되지 않음

```bash
# 상태 확인
docker ps -a | grep spanner

# 재시작
docker restart CONTAINER_ID

# 또는 새로 시작
make docker-up
```

### 포트 충돌

```bash
# 9010 포트 사용 확인
lsof -i :9010

# 9020 포트 사용 확인
lsof -i :9020

# 다른 포트 사용
docker run -d -p 19010:9010 -p 19020:9020 \
  gcr.io/cloud-spanner-emulator/emulator:1.5.33

# 환경 변수도 변경
export SPANNER_EMULATOR_HOST=localhost:19010
```

### gRPC 연결 실패

```bash
# 환경 변수 확인
echo $SPANNER_EMULATOR_HOST
# 출력: localhost:9010

# 설정되지 않았다면
export SPANNER_EMULATOR_HOST=localhost:9010

# Makefile에서는 자동 설정됨
make run
```

### Database not found 에러

```bash
# Instance/Database 생성 확인
gcloud spanner instances list
gcloud spanner databases list --instance=test-instance

# 없다면 생성
make setup-instance
```

## 📊 Spanner Emulator vs 실제 Spanner

| 특징 | Emulator | 실제 Spanner |
|------|----------|--------------|
| **비용** | ✅ 무료 | 💰 유료 |
| **성능** | 개발용 | 프로덕션 |
| **저장** | 메모리 | 디스크 |
| **재시작** | 데이터 삭제 | 데이터 유지 |
| **기능** | 대부분 지원 | 전체 지원 |
| **네트워크** | localhost | 글로벌 |

## 🎯 Emulator 제한사항

### 지원하지 않는 기능

- **IAM**: 인증/권한 없음
- **Backup/Restore**: 백업 불가
- **Multi-Region**: 단일 노드
- **Query Optimizer**: 최적화 제한
- **Monitoring**: 모니터링 불가

### 데이터 영속성

```bash
# Emulator는 메모리 사용 → 재시작 시 데이터 삭제

# 데이터 영속성 필요하면 실제 Spanner 사용
gcloud spanner instances create real-instance \
  --config=regional-us-central1 \
  --nodes=1
```

## 🔄 여러 Database 사용

```bash
# Database 1
gcloud spanner databases create db1 --instance=test-instance

# Database 2
gcloud spanner databases create db2 --instance=test-instance

# yo로 각각 생성
yo test-project test-instance db1 -o models/db1
yo test-project test-instance db2 -o models/db2
```

## 📝 환경별 설정

### 개발 (Development)

```bash
export SPANNER_EMULATOR_HOST=localhost:9010
export SPANNER_PROJECT_ID=dev-project
export SPANNER_INSTANCE_ID=dev-instance
export SPANNER_DATABASE_ID=dev-database
```

### 테스트 (Test)

```bash
export SPANNER_EMULATOR_HOST=localhost:9010
export SPANNER_PROJECT_ID=test-project
export SPANNER_INSTANCE_ID=test-instance
export SPANNER_DATABASE_ID=test-database
```

### 프로덕션 (Production)

```bash
unset SPANNER_EMULATOR_HOST  # 실제 Spanner 사용
export SPANNER_PROJECT_ID=prod-project
export SPANNER_INSTANCE_ID=prod-instance
export SPANNER_DATABASE_ID=prod-database
```

## 🎉 정리

### Emulator 사용 시 체크리스트

- [x] Docker 컨테이너 실행 중
- [x] 포트 9010, 9020 열림
- [x] SPANNER_EMULATOR_HOST 설정
- [x] Instance/Database 생성
- [x] 마이그레이션 실행
- [x] yo 코드 생성

### 유용한 명령어

```bash
# 상태 확인
make info

# Docker 확인
make docker-ps

# 스키마 확인
make show-schema

# CLI 접속
make spanner-cli
```

Happy Coding with Spanner Emulator! 🐳🚀

