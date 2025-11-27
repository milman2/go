package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"cloud.google.com/go/spanner"
)

func main() {
	ctx := context.Background()

	// 환경변수에서 설정 읽기
	projectID := getEnv("SPANNER_PROJECT_ID", "test-project")
	instanceID := getEnv("SPANNER_INSTANCE_ID", "test-instance")
	databaseID := getEnv("SPANNER_DATABASE_ID", "test-db")

	databaseName := fmt.Sprintf("projects/%s/instances/%s/databases/%s",
		projectID, instanceID, databaseID)

	fmt.Printf("🌱 Connecting to: %s\n", databaseName)

	// Spanner 클라이언트 생성
	client, err := spanner.NewClient(ctx, databaseName)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// 트랜잭션으로 데이터 삽입
	_, err = client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		// 1. Users 삽입
		fmt.Println("📝 Inserting users...")
		userMutations := []*spanner.Mutation{
			spanner.Insert("users",
				[]string{"id", "email", "name", "created_at", "updated_at"},
				[]interface{}{
					"550e8400-e29b-41d4-a716-446655440001",
					"john.doe@example.com",
					"John Doe",
					spanner.CommitTimestamp,
					spanner.CommitTimestamp,
				}),
			spanner.Insert("users",
				[]string{"id", "email", "name", "created_at", "updated_at"},
				[]interface{}{
					"550e8400-e29b-41d4-a716-446655440002",
					"jane.smith@example.com",
					"Jane Smith",
					spanner.CommitTimestamp,
					spanner.CommitTimestamp,
				}),
			spanner.Insert("users",
				[]string{"id", "email", "name", "created_at", "updated_at"},
				[]interface{}{
					"550e8400-e29b-41d4-a716-446655440003",
					"bob.johnson@example.com",
					"Bob Johnson",
					spanner.CommitTimestamp,
					spanner.CommitTimestamp,
				}),
		}

		if err := txn.BufferWrite(userMutations); err != nil {
			return fmt.Errorf("failed to insert users: %w", err)
		}
		fmt.Println("  ✅ 3 users inserted")

		// 2. Posts 삽입
		fmt.Println("📝 Inserting posts...")
		postMutations := []*spanner.Mutation{
			spanner.Insert("posts",
				[]string{"id", "user_id", "title", "content", "published", "created_at", "updated_at"},
				[]interface{}{
					"660e8400-e29b-41d4-a716-446655440001",
					"550e8400-e29b-41d4-a716-446655440001",
					"Getting Started with Cloud Spanner",
					"Cloud Spanner is a fully managed, mission-critical, relational database service that offers transactional consistency at global scale...",
					true,
					spanner.CommitTimestamp,
					spanner.CommitTimestamp,
				}),
			spanner.Insert("posts",
				[]string{"id", "user_id", "title", "content", "published", "created_at", "updated_at"},
				[]interface{}{
					"660e8400-e29b-41d4-a716-446655440002",
					"550e8400-e29b-41d4-a716-446655440001",
					"Advanced Spanner Features",
					"This is a draft post about advanced features...",
					false,
					spanner.CommitTimestamp,
					spanner.CommitTimestamp,
				}),
			spanner.Insert("posts",
				[]string{"id", "user_id", "title", "content", "published", "created_at", "updated_at"},
				[]interface{}{
					"660e8400-e29b-41d4-a716-446655440003",
					"550e8400-e29b-41d4-a716-446655440002",
					"Building Scalable Applications",
					"Learn how to build applications that can scale globally with Cloud Spanner...",
					true,
					spanner.CommitTimestamp,
					spanner.CommitTimestamp,
				}),
			spanner.Insert("posts",
				[]string{"id", "user_id", "title", "content", "published", "created_at", "updated_at"},
				[]interface{}{
					"660e8400-e29b-41d4-a716-446655440004",
					"550e8400-e29b-41d4-a716-446655440002",
					"Database Design Best Practices",
					"Here are some best practices for designing your Spanner schema...",
					true,
					spanner.CommitTimestamp,
					spanner.CommitTimestamp,
				}),
			spanner.Insert("posts",
				[]string{"id", "user_id", "title", "content", "published", "created_at", "updated_at"},
				[]interface{}{
					"660e8400-e29b-41d4-a716-446655440005",
					"550e8400-e29b-41d4-a716-446655440003",
					"Work in Progress",
					"This is still being written...",
					false,
					spanner.CommitTimestamp,
					spanner.CommitTimestamp,
				}),
		}

		if err := txn.BufferWrite(postMutations); err != nil {
			return fmt.Errorf("failed to insert posts: %w", err)
		}
		fmt.Println("  ✅ 5 posts inserted")

		return nil
	})

	if err != nil {
		log.Fatalf("❌ Transaction failed: %v", err)
	}

	// 결과 확인
	time.Sleep(100 * time.Millisecond)
	fmt.Println("\n📊 Verification:")
	
	// Users 수 확인
	stmt := spanner.Statement{SQL: "SELECT COUNT(*) FROM users"}
	iter := client.Single().Query(ctx, stmt)
	defer iter.Stop()
	row, err := iter.Next()
	if err == nil {
		var count int64
		if err := row.Columns(&count); err == nil {
			fmt.Printf("  👥 Users: %d\n", count)
		}
	}

	// Posts 수 확인
	stmt = spanner.Statement{SQL: "SELECT COUNT(*) FROM posts"}
	iter = client.Single().Query(ctx, stmt)
	defer iter.Stop()
	row, err = iter.Next()
	if err == nil {
		var count int64
		if err := row.Columns(&count); err == nil {
			fmt.Printf("  📝 Posts: %d\n", count)
		}
	}

	fmt.Println("\n✅ Sample data seeded successfully!")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

