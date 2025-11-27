# Spanner 도구 사용 가이드

## 도구 개요

이 프로젝트는 Google Cloud Spanner와 함께 작동하는 세 가지 주요 도구를 사용합니다:

### 1. **hammer** - DDL 스키마 관리 도구

**용도**: 데이터베이스 스키마(테이블, 인덱스 등) 관리

**주요 기능**:
- `hammer create`: 스키마 파일로부터 데이터베이스 생성
- `hammer apply`: 스키마 변경사항을 기존 데이터베이스에 적용
- `hammer diff`: 현재 DB 스키마와 파일의 차이점 확인
- `hammer export`: 현재 DB 스키마를 SQL 파일로 내보내기

**사용 예시**:
```bash
# 데이터베이스 생성
make createdb

# 스키마 변경사항 적용
make db-apply

# 스키마 차이 확인
make db-diff

# 스키마 내보내기
make db-export
```

### 2. **wrench** - 데이터베이스 관리 도구

**용도**: 데이터베이스 삭제 및 DML(데이터 조작어) 실행

**주요 기능**:
- `wrench drop`: 데이터베이스 완전 삭제
- `wrench apply --dml`: DML 파일 실행 (INSERT, UPDATE, DELETE 등)

**사용 예시**:
```bash
# 데이터베이스 삭제 (확인 필요)
make dropdb

# DML 마이그레이션 실행
make migrate-dml
```

### 3. **yo** - Go 코드 생성 도구

**용도**: Spanner 데이터베이스 스키마로부터 Go 코드 자동 생성

**주요 기능**:
- 테이블 구조를 Go struct로 변환
- CRUD 메서드 자동 생성
- 타입 안전한 데이터베이스 접근 코드 생성

**사용 예시**:
```bash
# 모델 코드 생성
make generate-models
```

## 워크플로우

### 초기 설정
```bash
# 전체 환경 초기화 (Docker + 도구 빌드 + DB 생성 + 코드 생성)
make init
```

### 스키마 변경 워크플로우

1. **스키마 파일 수정** (`schema/schema.sql`)
   ```sql
   -- 예: 새 테이블 추가
   CREATE TABLE comments (
     id STRING(36) NOT NULL,
     post_id STRING(36) NOT NULL,
     content STRING(MAX),
   ) PRIMARY KEY (id);
   ```

2. **변경사항 확인**
   ```bash
   make db-diff
   ```

3. **변경사항 적용**
   ```bash
   make db-apply
   ```

4. **Go 코드 재생성**
   ```bash
   make generate-models
   ```

### 데이터베이스 리셋
```bash
# DB 완전 리셋 + 코드 재생성
make reset
```

## 디렉토리 구조

```
Database/spanner/
├── schema/
│   └── schema.sql          # 전체 DDL 스키마 정의
├── migrations/
│   ├── dml/                # DML 마이그레이션 (선택사항)
│   └── *.{up,down}.sql     # 참고용 (더 이상 사용하지 않음)
├── models/
│   └── *.yo.go             # yo로 생성된 코드
├── ext/
│   ├── hammer.go           # hammer 빌드 설정
│   ├── wrench.go           # wrench 빌드 설정
│   └── yo.go               # yo 빌드 설정
└── bin/
    └── ext/                # 빌드된 도구들
        ├── hammer
        ├── wrench
        └── yo
```

## 주요 명령어

| 명령어 | 설명 |
|--------|------|
| `make help` | 모든 명령어 목록 보기 |
| `make init` | 전체 초기화 |
| `make createdb` | 데이터베이스 생성 (hammer) |
| `make db-apply` | 스키마 변경 적용 (hammer) |
| `make db-diff` | 스키마 차이 확인 (hammer) |
| `make db-export` | 스키마 내보내기 (hammer) |
| `make dropdb` | 데이터베이스 삭제 (wrench) |
| `make resetdb` | 데이터베이스 리셋 |
| `make reset` | DB 리셋 + 코드 재생성 |
| `make generate-models` | Go 코드 생성 (yo) |
| `make build/ext` | 외부 도구 빌드 |
| `make show-schema` | 현재 스키마 확인 |

## 참고

- **hammer vs wrench**: hammer는 스키마(DDL) 관리용, wrench는 데이터베이스 삭제 및 DML 실행용
- **migrations 폴더**: 기존 `.up.sql`, `.down.sql` 파일들은 참고용으로 남겨두고, 실제로는 `schema/schema.sql`을 단일 진실 공급원(single source of truth)으로 사용
- **DML 마이그레이션**: 샘플 데이터나 마스터 데이터 삽입이 필요한 경우 `migrations/dml/` 폴더에 SQL 파일을 추가하고 `make migrate-dml` 실행

