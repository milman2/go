package main

import (
	"log"
	"net/http"

	httpDelivery "github.com/milman2/go-api/clean-architecture/internal/delivery/http"
	"github.com/milman2/go-api/clean-architecture/internal/repository/memory"
	"github.com/milman2/go-api/clean-architecture/internal/service"
	"github.com/milman2/go-api/clean-architecture/internal/usecase"
)

// 이 파일은 "Service" 용어를 사용한 예제입니다
// 실행: go run cmd/api/main_with_service.go

func main() {
	// 의존성 주입 (Service 용어 사용)

	// 1. Repository 생성
	userRepo := memory.NewUserRepository()

	// 2. Service 생성 (Use Case와 동일한 역할)
	userService := service.NewUserService(userRepo)

	// 3. Service를 사용하는 방법 보여주기
	// 실제로는 Service와 UseCase는 같은 인터페이스를 구현
	log.Printf("✅ UserService 생성됨: %T\n", userService)

	// 4. Handler 생성
	// 현재는 UseCase 기반 Handler를 재사용
	// (Service도 UseCase와 동일한 메서드를 가지므로 호환 가능)
	userUseCase := usecase.NewUserUseCase(userRepo)
	userHandler := httpDelivery.NewUserHandler(userUseCase)

	// 5. Router 설정
	router := httpDelivery.NewRouter(userHandler)

	// 6. 서버 시작
	addr := ":8080"
	log.Printf("\n🚀 Clean Architecture 서버 시작 (Service 용어 사용)\n")
	log.Printf("=" + "========================================" + "\n")
	log.Printf("📖 레이어 구조:\n")
	log.Printf("   - Domain: internal/domain/\n")
	log.Printf("   - Service: internal/service/ (= Use Case)\n")
	log.Printf("   - Repository: internal/repository/\n")
	log.Printf("   - Handler: internal/delivery/http/\n")
	log.Printf("\n")
	log.Printf("💡 핵심: Service = Use Case (역할 동일)\n")
	log.Printf("   - UserService: 비즈니스 로직 처리\n")
	log.Printf("   - UserUseCase: 비즈니스 로직 처리\n")
	log.Printf("   - 둘 다 같은 일을 합니다!\n")
	log.Printf("\n")
	log.Printf("🌐 서버 주소: http://localhost%s\n", addr)
	log.Printf("=" + "========================================" + "\n")

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("서버 시작 실패: %v", err)
	}
}
