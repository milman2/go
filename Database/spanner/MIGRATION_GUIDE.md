# 🔄 마이그레이션 가이드

이 문서는 최근 변경 사항과 더 이상 사용하지 않는 명령어/개념을 정리합니다.

## 📊 주요 변경 사항 요약

### ✅ 변경된 것

| 구분 | 이전 | 현재 |
|------|------|------|
| **데이터베이스 이름** | `test-database` | `test-db` |
| **코드 생성 명령어** | `make generate-yo` | `make generate-models` |
| **스키마 관리** | migrations/*.sql | `schema/schema.sql` |
| **도구 설치** | `make install-tools` | `make build/ext` |

### ❌ 삭제된 명령어

다음 명령어들은 **더 이상 사용하지 않습니다**:

```bash
# ❌ 삭제됨
make migrate-up-wrench
make migrate-down-wrench
make migrate-up-hammer
make migrate-down-hammer
make install-tools
make generate-yo
```

**이유**: 
- wrench와 hammer는 up/down을 지원하지 않음
- schema.sql 기반으로 통합됨

### ✨ 새로운 명령어

```bash
# 데이터베이스 관리
make createdb          # hammer create로 DB 생성
make dropdb            # wrench drop으로 DB 삭제
make resetdb           # DB 리셋 (삭제 후 재생성)

# 스키마 관리
make db-apply          # hammer apply로 스키마 변경 적용
make db-diff           # hammer diff로 차이 확인
make db-export         # hammer export로 스키마 내보내기

# 샘플 데이터
make seed-data         # 샘플 데이터 삽입
make clear-data        # 모든 데이터 삭제
make test-query        # 샘플 쿼리 테스트

# 도구
make build/ext         # 외부 도구 빌드
```

## 🎯 도구별 역할 재정의

### Hammer (DDL 관리)
```bash
# ✅ 사용
hammer create   # DB 생성 + 스키마 적용
hammer apply    # 스키마 변경 적용
hammer diff     # 스키마 차이 확인
hammer export   # 스키마 내보내기

# ❌ 지원 안 함
hammer up       # 지원 안 함
hammer down     # 지원 안 함
```

### Wrench (DB 관리)
```bash
# ✅ 사용
wrench drop     # 데이터베이스 삭제

# ❌ 지원 안 함
wrench migrate up      # 지원 안 함
wrench migrate down    # 지원 안 함
```

### Yo (코드 생성)
```bash
# ✅ 사용
yo project instance database -o models/

# Makefile 명령어
make generate-models   # 권장
```

## 📝 마이그레이션 체크리스트

기존 프로젝트를 업데이트하는 경우:

### 1. 환경변수 확인
```bash
# 이전
export SPANNER_DATABASE_ID=test-database

# 현재
export SPANNER_DATABASE_ID=test-db
```

### 2. Makefile 명령어 변경

```bash
# 이전
make generate-yo

# 현재
make generate-models
```

### 3. 스키마 관리 방식 변경

**이전 방식** (migrations 기반):
```bash
# migrations/001_create_users.up.sql
# migrations/001_create_users.down.sql
make migrate-up-wrench
```

**현재 방식** (schema.sql 기반):
```bash
# schema/schema.sql (모든 DDL)
make createdb         # 처음 생성
make db-apply         # 변경 적용
```

### 4. 샘플 데이터 사용

```bash
# 새로운 기능!
make seed-data        # 샘플 데이터 삽입
make test-query       # 데이터 확인
make clear-data       # 데이터 삭제
```

## 🔧 트러블슈팅

### 문제: "명령어를 찾을 수 없습니다"

**증상:**
```bash
make migrate-up-wrench
make: *** No rule to make target 'migrate-up-wrench'
```

**해결:**
```bash
# 새 명령어 사용
make createdb  # 또는
make db-apply
```

### 문제: 데이터베이스 이름 불일치

**증상:**
```
Error: Database not found: test-database
```

**해결:**
```bash
# 1. Makefile 확인
grep SPANNER_DATABASE_ID Makefile
# SPANNER_DATABASE_ID=test-db 인지 확인

# 2. 환경변수 설정
export SPANNER_DATABASE_ID=test-db

# 3. 다시 생성
make resetdb
```

### 문제: 마이그레이션 down이 필요한 경우

**해결:**
```bash
# down 대신 resetdb 사용
make resetdb           # DB 전체 리셋
make seed-data         # 샘플 데이터 재삽입
```

## 📚 관련 문서

- **빠른 시작**: `QUICK_START.md`
- **전체 가이드**: `README.md`
- **데이터 테스팅**: `DATA_TESTING_GUIDE.md`
- **스키마 관리**: `schema/README.md`
- **변경 이력**: `CHANGELOG.md`

## ⚡ 주요 사실 (Facts)

### Spanner의 특징
1. ✅ **DEFAULT 값**: `DEFAULT (false)` (괄호 필수!)
2. ✅ **INTERLEAVE**: CASCADE DELETE 지원
3. ❌ **CASCADE**: FOREIGN KEY의 CASCADE는 미지원
4. ❌ **AUTO_INCREMENT**: 미지원 (UUID 사용)

### 도구의 한계
1. ❌ **wrench**: migrate down 미지원
2. ❌ **hammer**: up/down 미지원
3. ✅ **대안**: `schema.sql` + `resetdb`로 해결

## 🎉 정리

이제 다음과 같은 깔끔한 워크플로우를 사용할 수 있습니다:

```bash
# 1. 초기 설정
make init

# 2. 샘플 데이터
make seed-data

# 3. 개발
vi schema/schema.sql      # 스키마 수정
make db-diff              # 확인
make db-apply             # 적용
make generate-models      # 코드 재생성

# 4. 리셋 (필요 시)
make resetdb
make seed-data
```

더 이상 복잡한 마이그레이션 파일 관리가 필요 없습니다! 🚀

