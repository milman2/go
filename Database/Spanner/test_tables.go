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
	os.Setenv("SPANNER_EMULATOR_HOST", "localhost:9010")

	database := "projects/test-project/instances/test-instance/databases/test-db"
	client, err := spanner.NewClient(ctx, database)
	if err != nil {
		log.Fatalf("❌ 연결 실패: %v", err)
	}
	defer client.Close()

	fmt.Println("📊 Spanner 데이터베이스 정보")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// 테이블 목록 조회
	fmt.Println("\n1️⃣ 테이블 목록:")
	stmt := spanner.Statement{
		SQL: `SELECT table_name, parent_table_name
		      FROM INFORMATION_SCHEMA.TABLES
		      WHERE table_catalog = '' AND table_schema = ''
		      ORDER BY table_name`,
	}

	iter := client.Single().Query(ctx, stmt)
	defer iter.Stop()

	tableCount := 0
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		var tableName, parentTable spanner.NullString
		if err := row.Columns(&tableName, &parentTable); err != nil {
			log.Fatal(err)
		}

		tableCount++
		fmt.Printf("   ├─ %s", tableName.StringVal)
		if parentTable.Valid {
			fmt.Printf(" (parent: %s)", parentTable.StringVal)
		}
		fmt.Println()
	}
	fmt.Printf("   └─ 총 %d개의 테이블\n", tableCount)

	// 인덱스 목록 조회
	fmt.Println("\n2️⃣ 인덱스 목록:")
	stmt2 := spanner.Statement{
		SQL: `SELECT 
		        index_name,
		        table_name,
		        index_type,
		        is_unique,
		        is_null_filtered
		      FROM INFORMATION_SCHEMA.INDEXES
		      WHERE table_catalog = '' AND table_schema = ''
		        AND index_name != 'PRIMARY_KEY'
		      ORDER BY table_name, index_name`,
	}

	iter2 := client.Single().Query(ctx, stmt2)
	defer iter2.Stop()

	indexCount := 0
	for {
		row, err := iter2.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		var indexName, tableName, indexType string
		var isUnique, isNullFiltered bool
		if err := row.Columns(&indexName, &tableName, &indexType, &isUnique, &isNullFiltered); err != nil {
			log.Fatal(err)
		}

		indexCount++
		uniqueStr := ""
		if isUnique {
			uniqueStr = " [UNIQUE]"
		}
		fmt.Printf("   ├─ %s.%s (%s)%s\n", tableName, indexName, indexType, uniqueStr)
	}
	fmt.Printf("   └─ 총 %d개의 인덱스\n", indexCount)

	// 각 테이블의 컬럼 정보
	fmt.Println("\n3️⃣ 테이블 상세 정보:")

	// 먼저 테이블 이름 목록 가져오기
	iter3 := client.Single().Query(ctx, spanner.Statement{
		SQL: `SELECT table_name FROM INFORMATION_SCHEMA.TABLES
		      WHERE table_catalog = '' AND table_schema = ''
		      ORDER BY table_name`,
	})

	var tables []string
	for {
		row, err := iter3.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		var tableName string
		row.Columns(&tableName)
		tables = append(tables, tableName)
	}
	iter3.Stop()

	// 각 테이블의 컬럼 정보 출력
	for i, tableName := range tables {
		if i > 0 {
			fmt.Println()
		}
		fmt.Printf("   📋 테이블: %s\n", tableName)

		stmt4 := spanner.Statement{
			SQL: `SELECT column_name, spanner_type, is_nullable
			      FROM INFORMATION_SCHEMA.COLUMNS
			      WHERE table_name = @tableName
			      ORDER BY ordinal_position`,
			Params: map[string]interface{}{
				"tableName": tableName,
			},
		}

		iter4 := client.Single().Query(ctx, stmt4)

		columnCount := 0
		for {
			row, err := iter4.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Fatal(err)
			}

			var columnName, spannerType, isNullable string
			if err := row.Columns(&columnName, &spannerType, &isNullable); err != nil {
				log.Fatal(err)
			}

			columnCount++
			nullable := ""
			if isNullable == "YES" {
				nullable = " (nullable)"
			}

			prefix := "      ├─"
			fmt.Printf("%s %-25s %s%s\n", prefix, columnName, spannerType, nullable)
		}
		fmt.Printf("      └─ %d개 컬럼\n", columnCount)
		iter4.Stop()
	}

	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🎉 조회 완료!")
}
