package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"cloud.google.com/go/spanner"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
)

func main() {
	ctx := context.Background()
	os.Setenv("SPANNER_EMULATOR_HOST", "localhost:9010")

	database := "projects/test-project/instances/test-instance/databases/test-db"
	client, err := spanner.NewClient(ctx, database)
	if err != nil {
		log.Fatalf("❌ 연결 실패: %v", err)
	}
	defer client.Close()

	fmt.Println("🧪 Spanner CRUD 테스트")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// CREATE 테스트
	fmt.Println("\n1️⃣ CREATE 테스트")
	userID := uuid.New().String()
	testEmail := fmt.Sprintf("test-%d@example.com", time.Now().Unix())

	m := spanner.InsertMap("users", map[string]interface{}{
		"id":         userID,
		"email":      testEmail,
		"name":       "테스트 사용자",
		"created_at": spanner.CommitTimestamp,
		"updated_at": spanner.CommitTimestamp,
	})

	_, err = client.Apply(ctx, []*spanner.Mutation{m})
	if err != nil {
		log.Fatalf("❌ INSERT 실패: %v", err)
	}
	fmt.Printf("✅ 사용자 생성 성공: ID=%s, Email=%s\n", userID, testEmail)

	// READ 테스트 - Key 기반
	fmt.Println("\n2️⃣ READ 테스트 (Key 기반)")
	row, err := client.Single().ReadRow(ctx, "users",
		spanner.Key{userID},
		[]string{"id", "email", "name", "created_at"})
	if err != nil {
		log.Fatalf("❌ READ 실패: %v", err)
	}

	var id, email, name string
	var createdAt time.Time
	if err := row.Columns(&id, &email, &name, &createdAt); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("✅ 조회 성공:\n")
	fmt.Printf("   - ID: %s\n", id)
	fmt.Printf("   - Email: %s\n", email)
	fmt.Printf("   - Name: %s\n", name)
	fmt.Printf("   - Created: %s\n", createdAt.Format(time.RFC3339))

	// READ 테스트 - Query
	fmt.Println("\n3️⃣ READ 테스트 (Query)")
	stmt := spanner.Statement{
		SQL: `SELECT id, email, name FROM users 
		      WHERE email = @email`,
		Params: map[string]interface{}{
			"email": testEmail,
		},
	}

	iter := client.Single().Query(ctx, stmt)
	defer iter.Stop()

	foundCount := 0
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		var qID, qEmail, qName string
		row.Columns(&qID, &qEmail, &qName)
		fmt.Printf("✅ Query 결과: %s (%s)\n", qName, qEmail)
		foundCount++
	}
	fmt.Printf("   총 %d건 조회\n", foundCount)

	// UPDATE 테스트
	fmt.Println("\n4️⃣ UPDATE 테스트")
	newName := "수정된 테스트 사용자"
	m = spanner.UpdateMap("users", map[string]interface{}{
		"id":         userID,
		"name":       newName,
		"updated_at": spanner.CommitTimestamp,
	})

	_, err = client.Apply(ctx, []*spanner.Mutation{m})
	if err != nil {
		log.Fatalf("❌ UPDATE 실패: %v", err)
	}
	fmt.Printf("✅ 사용자 수정 성공: 새 이름='%s'\n", newName)

	// 수정 확인
	row, _ = client.Single().ReadRow(ctx, "users",
		spanner.Key{userID},
		[]string{"name"})
	var updatedName string
	row.Columns(&updatedName)
	fmt.Printf("✅ 수정 확인: %s\n", updatedName)

	// DELETE 테스트
	fmt.Println("\n5️⃣ DELETE 테스트")
	m = spanner.Delete("users", spanner.Key{userID})

	_, err = client.Apply(ctx, []*spanner.Mutation{m})
	if err != nil {
		log.Fatalf("❌ DELETE 실패: %v", err)
	}
	fmt.Printf("✅ 사용자 삭제 성공: ID=%s\n", userID)

	// 삭제 확인
	_, err = client.Single().ReadRow(ctx, "users",
		spanner.Key{userID},
		[]string{"id"})
	if err != nil {
		fmt.Println("✅ 삭제 확인: 사용자가 존재하지 않음 (정상)")
	} else {
		fmt.Println("❌ 삭제 확인 실패: 사용자가 여전히 존재함")
	}

	// BATCH CREATE 테스트
	fmt.Println("\n6️⃣ BATCH CREATE 테스트")
	mutations := []*spanner.Mutation{}
	batchIDs := []string{}

	for i := 0; i < 3; i++ {
		batchID := uuid.New().String()
		batchIDs = append(batchIDs, batchID)

		m := spanner.InsertMap("users", map[string]interface{}{
			"id":         batchID,
			"email":      fmt.Sprintf("batch-%d@example.com", i),
			"name":       fmt.Sprintf("배치 사용자 %d", i+1),
			"created_at": spanner.CommitTimestamp,
			"updated_at": spanner.CommitTimestamp,
		})
		mutations = append(mutations, m)
	}

	_, err = client.Apply(ctx, mutations)
	if err != nil {
		log.Fatalf("❌ BATCH INSERT 실패: %v", err)
	}
	fmt.Printf("✅ %d명의 사용자 일괄 생성 성공\n", len(batchIDs))

	// 전체 조회
	fmt.Println("\n7️⃣ 전체 조회 테스트")
	stmt = spanner.Statement{
		SQL: `SELECT id, email, name FROM users ORDER BY created_at DESC LIMIT 5`,
	}

	iter = client.Single().Query(ctx, stmt)
	defer iter.Stop()

	fmt.Println("   최근 사용자 5명:")
	count := 0
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		var qID, qEmail, qName string
		row.Columns(&qID, &qEmail, &qName)
		count++
		fmt.Printf("   %d. %s (%s)\n", count, qName, qEmail)
	}

	// 정리 (배치 생성한 사용자 삭제)
	fmt.Println("\n8️⃣ 정리 (테스트 데이터 삭제)")
	deleteMutations := []*spanner.Mutation{}
	for _, batchID := range batchIDs {
		deleteMutations = append(deleteMutations, spanner.Delete("users", spanner.Key{batchID}))
	}

	_, err = client.Apply(ctx, deleteMutations)
	if err != nil {
		log.Printf("⚠️ 정리 실패: %v", err)
	} else {
		fmt.Printf("✅ 테스트 데이터 정리 완료 (%d건)\n", len(batchIDs))
	}

	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🎉 CRUD 테스트 완료!")
}
