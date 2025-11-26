-- Create posts table (yo 기능 테스트용)
CREATE TABLE posts (
  id STRING(36) NOT NULL,
  user_id STRING(36) NOT NULL,
  title STRING(200) NOT NULL,
  content STRING(MAX),
  published BOOL NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true),
  updated_at TIMESTAMP NOT NULL OPTIONS (allow_commit_timestamp=true),
  CONSTRAINT fk_posts_user FOREIGN KEY (user_id) REFERENCES users (id),
) PRIMARY KEY (id);

-- Index on user_id for efficient queries
CREATE INDEX posts_user_id_idx ON posts(user_id);

-- Index on published for filtering
CREATE INDEX posts_published_idx ON posts(published);

