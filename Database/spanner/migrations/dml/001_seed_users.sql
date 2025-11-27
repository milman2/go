-- Seed data for users table
-- 사용법: make seed-data

-- 테스트 사용자 1
INSERT INTO users (id, email, name, created_at, updated_at)
VALUES (
  '550e8400-e29b-41d4-a716-446655440001',
  'john.doe@example.com',
  'John Doe',
  CURRENT_TIMESTAMP(),
  CURRENT_TIMESTAMP()
);

-- 테스트 사용자 2
INSERT INTO users (id, email, name, created_at, updated_at)
VALUES (
  '550e8400-e29b-41d4-a716-446655440002',
  'jane.smith@example.com',
  'Jane Smith',
  CURRENT_TIMESTAMP(),
  CURRENT_TIMESTAMP()
);

-- 테스트 사용자 3
INSERT INTO users (id, email, name, created_at, updated_at)
VALUES (
  '550e8400-e29b-41d4-a716-446655440003',
  'bob.johnson@example.com',
  'Bob Johnson',
  CURRENT_TIMESTAMP(),
  CURRENT_TIMESTAMP()
);

