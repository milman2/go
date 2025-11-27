-- Seed data for posts table
-- 사용법: make seed-data

-- John Doe의 게시글 (published)
INSERT INTO posts (id, user_id, title, content, published, created_at, updated_at)
VALUES (
  '660e8400-e29b-41d4-a716-446655440001',
  '550e8400-e29b-41d4-a716-446655440001',
  'Getting Started with Cloud Spanner',
  'Cloud Spanner is a fully managed, mission-critical, relational database service that offers transactional consistency at global scale...',
  TRUE,
  CURRENT_TIMESTAMP(),
  CURRENT_TIMESTAMP()
);

-- John Doe의 게시글 (draft)
INSERT INTO posts (id, user_id, title, content, published, created_at, updated_at)
VALUES (
  '660e8400-e29b-41d4-a716-446655440002',
  '550e8400-e29b-41d4-a716-446655440001',
  'Advanced Spanner Features',
  'This is a draft post about advanced features...',
  FALSE,
  CURRENT_TIMESTAMP(),
  CURRENT_TIMESTAMP()
);

-- Jane Smith의 게시글 (published)
INSERT INTO posts (id, user_id, title, content, published, created_at, updated_at)
VALUES (
  '660e8400-e29b-41d4-a716-446655440003',
  '550e8400-e29b-41d4-a716-446655440002',
  'Building Scalable Applications',
  'Learn how to build applications that can scale globally with Cloud Spanner...',
  TRUE,
  CURRENT_TIMESTAMP(),
  CURRENT_TIMESTAMP()
);

-- Jane Smith의 게시글 (published)
INSERT INTO posts (id, user_id, title, content, published, created_at, updated_at)
VALUES (
  '660e8400-e29b-41d4-a716-446655440004',
  '550e8400-e29b-41d4-a716-446655440002',
  'Database Design Best Practices',
  'Here are some best practices for designing your Spanner schema...',
  TRUE,
  CURRENT_TIMESTAMP(),
  CURRENT_TIMESTAMP()
);

-- Bob Johnson의 게시글 (draft)
INSERT INTO posts (id, user_id, title, content, published, created_at, updated_at)
VALUES (
  '660e8400-e29b-41d4-a716-446655440005',
  '550e8400-e29b-41d4-a716-446655440003',
  'Work in Progress',
  'This is still being written...',
  FALSE,
  CURRENT_TIMESTAMP(),
  CURRENT_TIMESTAMP()
);

