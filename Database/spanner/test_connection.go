package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"
)

func main() {
	ctx := context.Background()

	// 환경 변수 설정
	os.Setenv("SPANNER_EMULATOR_HOST", "localhost:9010")

	// Spanner 클라이언트 생성
	database := "projects/test-project/instances/test-instance/databases/test-db"
	client, err := spanner.NewClient(ctx, database)
	if err != nil {
		log.Fatalf("❌ 연결 실패: %v", err)
	}
	defer client.Close()

	fmt.Println("✅ Spanner 연결 성공!")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("Database: %s\n", database)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// 간단한 쿼리 실행
	fmt.Println("\n🔍 연결 테스트 쿼리 실행...")
	stmt := spanner.Statement{SQL: `SELECT 1 as test, 'Hello Spanner' as message`}
	iter := client.Single().Query(ctx, stmt)
	defer iter.Stop()

	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		var test int64
		var message string
		if err := row.Columns(&test, &message); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("✅ 쿼리 결과: test=%d, message='%s'\n", test, message)
	}

	fmt.Println("\n🎉 테스트 완료!")
}
