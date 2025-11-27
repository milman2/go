# 스키마 변경 이력

schema.sql의 주요 변경 사항을 기록합니다.

## 2025-11-27 - 초기 스키마

### Users Table
```sql
CREATE TABLE users (
  id STRING(36) NOT NULL,
  email STRING(255) NOT NULL,
  name STRING(100) NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
) PRIMARY KEY (id);

CREATE UNIQUE INDEX users_email_idx ON users(email);
```

### Posts Table
```sql
CREATE TABLE posts (
  id STRING(36) NOT NULL,
  user_id STRING(36) NOT NULL,
  title STRING(200) NOT NULL,
  content STRING(MAX),
  published BOOL NOT NULL DEFAULT (false),
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
) PRIMARY KEY (id);

CREATE INDEX posts_user_id_idx ON posts(user_id);
CREATE INDEX posts_published_idx ON posts(published);
```

---

## 변경 템플릿

아래 형식으로 변경사항을 기록하세요:

```markdown
## YYYY-MM-DD - 변경 제목

### 변경 내용
- 추가: users.age 컬럼 (INT64)
- 삭제: posts.legacy_field
- 수정: users.email 크기 (255 → 320)

### SQL
```sql
ALTER TABLE users ADD COLUMN age INT64;
```

### 이유
- 사용자 나이 기반 필터링 기능 추가

### 영향
- 기존 데이터: NULL 값으로 채워짐
- 애플리케이션: User 모델에 Age 필드 추가 필요

### 롤백 (필요 시)
```sql
ALTER TABLE users DROP COLUMN age;
```
```

---

## 사용 방법

### 변경 전
1. 이 파일에 계획 작성
2. 팀원 리뷰 요청

### 변경 후
1. 실제 적용 내용 기록
2. 커밋 메시지에 참조

### 운영 배포 시
1. 이 히스토리 기반으로 운영용 마이그레이션 파일 생성
2. 단계별 테스트 후 적용

