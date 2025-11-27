# 변경 이력

## 2025-11-27 (Latest) - 문서 정리 및 데이터베이스 이름 통일

### 🔄 변경사항

#### 1. 데이터베이스 이름 통일
- ✅ 모든 문서에서 `test-database` → `test-db`로 변경
- ✅ 일관성 있는 데이터베이스 이름 사용

#### 2. 더 이상 유효하지 않은 명령어 정리
- ❌ `make migrate-up-wrench` (삭제됨)
- ❌ `make migrate-down-wrench` (삭제됨)
- ❌ `make migrate-up-hammer` (삭제됨)
- ❌ `make migrate-down-hammer` (삭제됨)
- ❌ `make generate-yo` (이름 변경: `make generate-models`)

#### 3. 새로운 명령어 (샘플 데이터)
- ✨ `make seed-data`: 샘플 데이터 삽입
- ✨ `make clear-data`: 모든 데이터 삭제
- ✨ `make test-query`: 샘플 쿼리 테스트

#### 4. 새 문서 추가
- 📄 `DATA_TESTING_GUIDE.md`: 샘플 데이터 테스팅 가이드
- 📄 `migrations/dml/README.md`: DML 사용 가이드
- 📄 `schema/README.md`: 스키마 관리 가이드
- 🔧 `scripts/seed_data.go`: Go 기반 샘플 데이터 삽입 스크립트

#### 5. 수정된 문서
- 📝 `README.md`: 데이터베이스 이름 통일
- 📝 `USAGE.md`: 더 이상 사용하지 않는 명령어 제거
- 📝 `QUICK_START.md`: 데이터베이스 이름 통일
- 📝 `SPANNER.md`: 데이터베이스 이름 통일
- 📝 `DBeaver.md`: 데이터베이스 이름 통일
- 📝 `SETUP_CHECKLIST.md`: 데이터베이스 이름 통일
- 📝 `TESTING.md`: 데이터베이스 이름 통일
- 📝 `DOCKER_GUIDE.md`: 데이터베이스 이름 통일

#### 6. Spanner 주요 사실 정리
- ✅ **DEFAULT 값**: `DEFAULT (false)` 형식 (괄호 필수!)
- ✅ **INTERLEAVE**: CASCADE DELETE 지원
- ❌ wrench는 `migrate down` 미지원
- ❌ hammer는 `up/down` 미지원 (create/apply/diff/export만 지원)

### 📋 현재 상태

#### 사용 가능한 Makefile 명령어
```bash
# 초기화
make init              # 전체 초기화 (Docker + DB + 스키마 + 모델)

# 데이터베이스 관리
make createdb          # 데이터베이스 생성 (hammer)
make dropdb            # 데이터베이스 삭제 (wrench)
make resetdb           # DB 리셋 (삭제 후 재생성)

# 스키마 관리
make db-apply          # 스키마 변경 적용 (hammer)
make db-diff           # 스키마 차이 확인 (hammer)
make db-export         # 스키마 내보내기 (hammer)
make show-schema       # 현재 스키마 확인

# 샘플 데이터
make seed-data         # 샘플 데이터 삽입
make clear-data        # 모든 데이터 삭제
make test-query        # 샘플 쿼리 테스트

# 코드 생성
make generate-models   # Go 모델 생성 (yo)

# 도구
make build/ext         # 외부 도구 빌드
```

---

## 2025-11-27 - 주요 리팩토링

### 🎯 목표
- hammer와 wrench의 용도를 명확히 구분
- schema/schema.sql을 단일 진실 공급원(Single Source of Truth)으로 사용
- 참고 프로젝트(school-live-api)의 베스트 프랙티스 적용

### ✅ 완료된 작업

#### 1. 디렉토리 구조 변경
- ✨ `schema/` 폴더 추가
  - `schema/schema.sql`: 전체 DDL 스키마 정의
- ✨ `ext/` 폴더 추가
  - `ext/hammer.go`: hammer 빌드 설정
  - `ext/wrench.go`: wrench 빌드 설정
  - `ext/yo.go`: yo 빌드 설정
- 🗑️ `migrations/*.up.sql`, `migrations/*.down.sql` 삭제
  - 더 이상 파일 기반 마이그레이션 사용하지 않음
  - `schema/schema.sql`이 유일한 소스

#### 2. Makefile 개선
- **도구 구분 명확화**:
  - `hammer`: DDL 스키마 관리 (create, apply, diff, export)
  - `wrench`: 데이터베이스 삭제, DML 실행
  - `yo`: Go 코드 생성

- **새 명령어**:
  - `make createdb`: DB 생성 (hammer create)
  - `make db-apply`: 스키마 변경 적용 (hammer apply)
  - `make db-diff`: 스키마 차이 확인 (hammer diff)
  - `make db-export`: 스키마 내보내기 (hammer export)
  - `make dropdb`: DB 삭제 (wrench drop)
  - `make resetdb`: DB 리셋
  - `make migrate-dml`: DML 마이그레이션 (wrench apply --dml)
  - `make generate-models`: 코드 생성 (이전의 generate-yo)
  - `make build/ext`: 외부 도구 빌드

- **제거된 명령어**:
  - `make install-tools` → `make build/ext`로 대체
  - `make migrate-up-wrench` → `make createdb` + `make db-apply`로 대체
  - `make migrate-down-wrench` → 더 이상 필요 없음
  - `make migrate-up-hammer` → 제거
  - `make migrate-down-hammer` → 제거
  - `make generate-yo` → `make generate-models`로 이름 변경

#### 3. schema/schema.sql 개선
- Spanner 주요 기능 주석 추가:
  ```sql
  -- Spanner 주요 기능:
  -- 1. DEFAULT 값: DEFAULT (값) 형식으로 괄호 필수
  -- 2. FOREIGN KEY: 기본 지원 (CASCADE 미지원)
  -- 3. INTERLEAVE: 부모-자식 관계 + CASCADE DELETE 지원 + 성능 최적화
  ```
- Posts 테이블에 DEFAULT 값 추가:
  ```sql
  published BOOL NOT NULL DEFAULT (false),
  ```
- FOREIGN KEY 제거 (INTERLEAVE 사용 권장)

#### 4. 문서 업데이트
- ✅ `README.md`: 전체 구조, 명령어, 워크플로우 업데이트
- ✅ `QUICK_START.md`: 새 명령어 반영
- ✅ `SPANNER.md`: Spanner 주요 기능 추가
- ✅ `TOOLS_GUIDE.md`: 새로 추가 (도구 사용법 상세 설명)
- ✅ `.gitignore`: exported_schema.sql 추가

### 📝 새로운 워크플로우

#### 초기 설정
```bash
make init  # Docker + 도구 빌드 + Instance + DB 생성 + 코드 생성
```

#### 스키마 변경
```bash
# 1. schema/schema.sql 수정
vim schema/schema.sql

# 2. 변경사항 확인
make db-diff

# 3. 변경사항 적용
make db-apply

# 4. 코드 재생성
make generate-models
```

#### 완전 리셋
```bash
make reset  # DB 재생성 + 코드 재생성
```

### 🎯 주요 개선사항

1. **단일 진실 공급원**
   - `schema/schema.sql`이 전체 스키마의 유일한 소스
   - 마이그레이션 히스토리 관리 불필요
   - 스키마 파일만 보면 전체 구조 파악 가능

2. **도구 역할 명확화**
   - hammer: DDL 전문
   - wrench: DB 관리 + DML 전문
   - yo: 코드 생성 전문

3. **로컬 빌드**
   - `go install` 대신 `go generate` 사용
   - `bin/ext/` 디렉토리에 도구 빌드
   - 버전 고정 및 재현 가능한 빌드

4. **워크플로우 단순화**
   - 복잡한 마이그레이션 관리 불필요
   - diff → apply → generate의 명확한 흐름

### 🔄 마이그레이션 가이드

#### 기존 방식 (제거됨)
```bash
# ❌ 더 이상 사용하지 않음
make migrate-up-wrench
make generate-yo
```

#### 새 방식
```bash
# ✅ 새로운 방식
make db-apply         # 스키마 적용
make generate-models  # 코드 생성
```

### 📚 참고 문서

- `TOOLS_GUIDE.md`: 도구 사용법 상세 가이드
- `README.md`: 프로젝트 개요 및 전체 가이드
- `QUICK_START.md`: 빠른 시작 가이드
- `SPANNER.md`: Spanner 기능 및 테스트 가이드

### 🚀 다음 단계

1. Spanner 주요 기능 활용
   - DEFAULT 값 활용
   - INTERLEAVE로 부모-자식 관계 구현
   - 적절한 인덱스 설계

2. DML 마이그레이션
   - `migrations/dml/` 디렉토리에 샘플 데이터 추가
   - `make migrate-dml`로 실행

3. Clean Architecture 적용
   - Repository 레이어에서 yo 모델 사용
   - 도메인 로직과 인프라 분리

