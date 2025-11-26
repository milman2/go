package main

import (
	"log"
	"net/http"

	httpDelivery "github.com/milman2/go-api/clean-architecture/internal/delivery/http"
	"github.com/milman2/go-api/clean-architecture/internal/repository/memory"
	"github.com/milman2/go-api/clean-architecture/internal/usecase"
)

func main() {
	// 의존성 주입 (Dependency Injection)
	// 외부 레이어에서 내부 레이어로 의존성 주입

	// 1. Repository 생성 (가장 바깥 레이어)
	userRepo := memory.NewUserRepository()

	// 2. Use Case 생성 (중간 레이어)
	userUseCase := usecase.NewUserUseCase(userRepo)

	// 3. Handler 생성 (프레젠테이션 레이어)
	userHandler := httpDelivery.NewUserHandler(userUseCase)

	// 4. Router 설정
	router := httpDelivery.NewRouter(userHandler)

	// 5. 서버 시작
	addr := ":8080"
	log.Printf("🚀 Clean Architecture 서버가 %s 포트에서 시작되었습니다\n", addr)
	log.Printf("📖 Clean Architecture 레이어:\n")
	log.Printf("   - Domain (Entities): internal/domain\n")
	log.Printf("   - Use Cases: internal/usecase\n")
	log.Printf("   - Interface Adapters: internal/repository\n")
	log.Printf("   - Frameworks & Drivers: internal/delivery/http\n")

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("서버 시작 실패: %v", err)
	}
}
