-- ============================================================================
-- Google Cloud Spanner Database Schema
-- ============================================================================
-- 
-- 주요 특징:
-- 1. DEFAULT 값: DEFAULT (값) 형식으로 괄호 필수
-- 2. FOREIGN KEY: 기본 지원 (CASCADE 미지원)
-- 3. INTERLEAVE: 부모-자식 관계 + CASCADE DELETE 지원 + 성능 최적화
--
-- ============================================================================

-- ============================================================================
-- Users Table
-- ============================================================================
CREATE TABLE users (
  id STRING(36) NOT NULL,
  email STRING(255) NOT NULL,
  name STRING(100) NOT NULL,
  created_at TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true),
  updated_at TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true),
) PRIMARY KEY (id);

CREATE UNIQUE INDEX users_email_idx ON users(email);

-- ============================================================================
-- Posts Table
-- ============================================================================
-- 
-- 💡 두 가지 구현 방식 선택 가능:
--
-- 방식 1: FOREIGN KEY (현재 사용 중) - 일반적인 참조 관계
-- 방식 2: INTERLEAVE (주석 참고) - 강한 부모-자식 관계 + CASCADE DELETE
--
-- ============================================================================

-- [방식 1] FOREIGN KEY 버전 (현재 활성화)
CREATE TABLE posts (
  id STRING(36) NOT NULL,
  user_id STRING(36) NOT NULL,
  title STRING(200) NOT NULL,
  content STRING(MAX),
  published BOOL NOT NULL DEFAULT (false),
  created_at TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true),
  updated_at TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true),
) PRIMARY KEY (id);

CREATE INDEX posts_user_id_idx ON posts(user_id);
CREATE INDEX posts_published_idx ON posts(published);

-- ============================================================================
-- [방식 2] INTERLEAVE 버전 (주석 처리 - 필요 시 활성화)
-- ============================================================================
-- 
-- INTERLEAVE 사용 시 장점:
-- ✅ CASCADE DELETE 자동 (user 삭제 시 posts도 자동 삭제)
-- ✅ 성능 최적화 (부모-자식이 같은 노드에 저장)
-- ✅ 강력한 참조 무결성
--
-- 주의사항:
-- ⚠️ PRIMARY KEY에 부모 키(user_id) 포함 필수
-- ⚠️ 부모-자식 관계가 명확한 1:N 구조에만 사용
--
-- 활성화 방법:
-- 1. 위의 [방식 1] 테이블 생성 부분을 주석 처리
-- 2. 아래 주석을 해제
--
-- CREATE TABLE posts (
--   user_id STRING(36) NOT NULL,      -- 부모 키 (첫 번째)
--   id STRING(36) NOT NULL,           -- 자식 키
--   title STRING(200) NOT NULL,
--   content STRING(MAX),
--   published BOOL NOT NULL DEFAULT (false),
--   created_at TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true),
--   updated_at TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true),
-- ) PRIMARY KEY (user_id, id),        -- 복합 키: 부모 + 자식
--   INTERLEAVE IN PARENT users ON DELETE CASCADE;
--
-- CREATE INDEX posts_published_idx ON posts(published);
-- 
-- ============================================================================

