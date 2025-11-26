package main

import (
	"log"
	"net/http"

	httpDelivery "github.com/milman2/go-api/clean-architecture/internal/delivery/http"
	gormRepo "github.com/milman2/go-api/clean-architecture/internal/repository/gorm"
	"github.com/milman2/go-api/clean-architecture/internal/usecase"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 이 파일은 GORM을 사용한 예제입니다
// 실행: go run cmd/api/main_with_gorm.go

func main() {
	// 1. GORM DB 연결 (SQLite 사용)
	db, err := gorm.Open(sqlite.Open("users.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // SQL 로그 출력
	})
	if err != nil {
		log.Fatalf("DB 연결 실패: %v", err)
	}

	// 2. 테이블 자동 마이그레이션
	if err := db.AutoMigrate(&gormRepo.UserModel{}); err != nil {
		log.Fatalf("마이그레이션 실패: %v", err)
	}
	log.Println("✅ 데이터베이스 마이그레이션 완료")

	// 3. Repository 생성 (GORM 구현)
	userRepo := gormRepo.NewUserRepository(db)

	// 4. Use Case 생성 (Repository 인터페이스 사용)
	userUseCase := usecase.NewUserUseCase(userRepo)

	// 5. Handler 생성
	userHandler := httpDelivery.NewUserHandler(userUseCase)

	// 6. Router 설정
	router := httpDelivery.NewRouter(userHandler)

	// 7. 서버 시작
	addr := ":8080"
	log.Printf("\n🚀 Clean Architecture + GORM 서버 시작\n")
	log.Printf("=" + "========================================" + "\n")
	log.Printf("📖 아키텍처 레이어:\n")
	log.Printf("   - Domain: internal/domain/ (순수 비즈니스 로직)\n")
	log.Printf("   - Use Case: internal/usecase/ (애플리케이션 로직)\n")
	log.Printf("   - Repository: internal/repository/gorm/ (GORM 어댑터)\n")
	log.Printf("   - Handler: internal/delivery/http/ (HTTP 어댑터)\n")
	log.Printf("\n")
	log.Printf("💾 데이터베이스:\n")
	log.Printf("   - ORM: GORM (Yet Another ORM)\n")
	log.Printf("   - Driver: SQLite\n")
	log.Printf("   - File: users.db\n")
	log.Printf("\n")
	log.Printf("✨ GORM 기능:\n")
	log.Printf("   - ✅ 자동 마이그레이션 (AutoMigrate)\n")
	log.Printf("   - ✅ 관계 매핑 (Associations)\n")
	log.Printf("   - ✅ 트랜잭션 (Transactions)\n")
	log.Printf("   - ✅ Hooks (Before/After)\n")
	log.Printf("   - ✅ 프리로드 (Preload)\n")
	log.Printf("   - ✅ SQL 로깅\n")
	log.Printf("\n")
	log.Printf("🌐 서버 주소: http://localhost%s\n", addr)
	log.Printf("=" + "========================================" + "\n")

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("서버 시작 실패: %v", err)
	}
}

