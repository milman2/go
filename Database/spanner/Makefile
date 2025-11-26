# Makefile for Google Cloud Spanner + yo

# Spanner 설정
SPANNER_EMULATOR_HOST=localhost:9010
SPANNER_PROJECT_ID=test-project
SPANNER_INSTANCE_ID=test-instance
SPANNER_DATABASE_ID=test-database

# 경로 설정
MIGRATIONS_DIR=migrations
MODELS_DIR=models

# 색상 출력
GREEN=\033[0;32m
YELLOW=\033[1;33m
NC=\033[0m

.PHONY: help
help: ## 도움말 출력
	@echo "$(GREEN)Spanner + yo Makefile Commands$(NC)"
	@echo "=================================="
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-20s$(NC) %s\n", $$1, $$2}'

.PHONY: docker-up
docker-up: ## Docker Spanner emulator 시작
	@echo "$(GREEN)🐳 Starting Spanner emulator...$(NC)"
	docker-compose up -d spanner
	@echo "$(GREEN)✅ Spanner emulator started on :9010 (gRPC), :9020 (HTTP)$(NC)"

.PHONY: docker-down
docker-down: ## Docker Spanner emulator 중지
	@echo "$(YELLOW)🛑 Stopping Spanner emulator...$(NC)"
	docker-compose down

.PHONY: docker-ps
docker-ps: ## Docker 컨테이너 상태 확인
	docker-compose ps

.PHONY: setup-instance
setup-instance: ## Spanner Instance와 Database 생성
	@echo "$(GREEN)📦 Creating Spanner instance and database...$(NC)"
	@gcloud config set auth/disable_credentials true
	@gcloud config set project $(SPANNER_PROJECT_ID)
	@gcloud config set api_endpoint_overrides/spanner http://localhost:9020/
	@gcloud spanner instances create $(SPANNER_INSTANCE_ID) \
		--config=emulator-config \
		--description="Test Instance" \
		--nodes=1 || true
	@gcloud spanner databases create $(SPANNER_DATABASE_ID) \
		--instance=$(SPANNER_INSTANCE_ID) || true
	@echo "$(GREEN)✅ Instance and database ready$(NC)"

.PHONY: migrate-up-hammer
migrate-up-hammer: ## Hammer로 마이그레이션 UP
	@echo "$(GREEN)🔨 Running migrations with Hammer...$(NC)"
	@SPANNER_EMULATOR_HOST=$(SPANNER_EMULATOR_HOST) \
		hammer -p $(SPANNER_PROJECT_ID) -i $(SPANNER_INSTANCE_ID) -d $(SPANNER_DATABASE_ID) \
		-m $(MIGRATIONS_DIR) up
	@echo "$(GREEN)✅ Hammer migration completed$(NC)"

.PHONY: migrate-down-hammer
migrate-down-hammer: ## Hammer로 마이그레이션 DOWN
	@echo "$(YELLOW)⬇️ Rolling back migrations with Hammer...$(NC)"
	@SPANNER_EMULATOR_HOST=$(SPANNER_EMULATOR_HOST) \
		hammer -p $(SPANNER_PROJECT_ID) -i $(SPANNER_INSTANCE_ID) -d $(SPANNER_DATABASE_ID) \
		-m $(MIGRATIONS_DIR) down
	@echo "$(GREEN)✅ Hammer rollback completed$(NC)"

.PHONY: migrate-up-wrench
migrate-up-wrench: ## Wrench로 마이그레이션 UP
	@echo "$(GREEN)🔧 Running migrations with Wrench...$(NC)"
	@SPANNER_EMULATOR_HOST=$(SPANNER_EMULATOR_HOST) \
		wrench migrate up \
		--directory $(MIGRATIONS_DIR) \
		--database projects/$(SPANNER_PROJECT_ID)/instances/$(SPANNER_INSTANCE_ID)/databases/$(SPANNER_DATABASE_ID)
	@echo "$(GREEN)✅ Wrench migration completed$(NC)"

.PHONY: migrate-down-wrench
migrate-down-wrench: ## Wrench로 마이그레이션 DOWN
	@echo "$(YELLOW)⬇️ Rolling back migrations with Wrench...$(NC)"
	@SPANNER_EMULATOR_HOST=$(SPANNER_EMULATOR_HOST) \
		wrench migrate down \
		--directory $(MIGRATIONS_DIR) \
		--database projects/$(SPANNER_PROJECT_ID)/instances/$(SPANNER_INSTANCE_ID)/databases/$(SPANNER_DATABASE_ID)
	@echo "$(GREEN)✅ Wrench rollback completed$(NC)"

.PHONY: generate-yo
generate-yo: ## yo로 코드 생성 (스키마 → Go 코드)
	@echo "$(GREEN)🔨 Generating code with yo...$(NC)"
	@rm -rf $(MODELS_DIR)
	@mkdir -p $(MODELS_DIR)
	@SPANNER_EMULATOR_HOST=$(SPANNER_EMULATOR_HOST) \
		yo $(SPANNER_PROJECT_ID) $(SPANNER_INSTANCE_ID) $(SPANNER_DATABASE_ID) \
		-o $(MODELS_DIR) -p models
	@echo "$(GREEN)✅ Code generation completed: $(MODELS_DIR)/$(NC)"
	@ls -lh $(MODELS_DIR)/

.PHONY: install-tools
install-tools: ## 필요한 도구 설치
	@echo "$(GREEN)📦 Installing tools...$(NC)"
	@echo "Installing hammer..."
	@go install github.com/daichirata/hammer@v0.6.0
	@echo "Installing wrench..."
	@go install github.com/cloudspannerecosystem/wrench@v1.0.4
	@echo "Installing yo..."
	@go install go.mercari.io/yo@latest
	@echo "$(GREEN)✅ All tools installed$(NC)"
	@echo "  - hammer: $(shell which hammer)"
	@echo "  - wrench: $(shell which wrench)"
	@echo "  - yo: $(shell which yo)"

.PHONY: init
init: docker-up install-tools setup-instance migrate-up-wrench generate-yo ## 전체 초기화 (한번에)
	@echo ""
	@echo "$(GREEN)🎉 Spanner 환경 초기화 완료!$(NC)"
	@echo "=================================="
	@echo "  ✅ Docker Spanner emulator 실행"
	@echo "  ✅ Instance & Database 생성"
	@echo "  ✅ 마이그레이션 완료"
	@echo "  ✅ yo 코드 생성 완료"
	@echo ""
	@echo "$(YELLOW)📝 다음 단계:$(NC)"
	@echo "  make run          # 서버 실행"
	@echo "  make test         # API 테스트"

.PHONY: clean
clean: ## 생성된 파일 정리
	@echo "$(YELLOW)🧹 Cleaning up...$(NC)"
	@rm -rf $(MODELS_DIR)
	@echo "$(GREEN)✅ Cleaned$(NC)"

.PHONY: reset
reset: migrate-down-wrench migrate-up-wrench generate-yo ## DB 리셋 및 코드 재생성
	@echo "$(GREEN)✅ Database reset and code regenerated$(NC)"

.PHONY: run
run: ## 서버 실행
	@echo "$(GREEN)🚀 Starting server...$(NC)"
	@SPANNER_EMULATOR_HOST=$(SPANNER_EMULATOR_HOST) go run cmd/api/main.go

.PHONY: test
test: ## API 테스트
	@echo "$(GREEN)🧪 Running tests...$(NC)"
	@./test.sh

.PHONY: spanner-cli
spanner-cli: ## Spanner CLI 접속
	@echo "$(GREEN)💻 Connecting to Spanner CLI...$(NC)"
	docker-compose exec spanner-cli spanner-cli -p $(SPANNER_PROJECT_ID) -i $(SPANNER_INSTANCE_ID) -d $(SPANNER_DATABASE_ID)

.PHONY: show-schema
show-schema: ## 현재 스키마 확인
	@echo "$(GREEN)📋 Current schema:$(NC)"
	@SPANNER_EMULATOR_HOST=$(SPANNER_EMULATOR_HOST) \
		wrench migrate status \
		--directory $(MIGRATIONS_DIR) \
		--database projects/$(SPANNER_PROJECT_ID)/instances/$(SPANNER_INSTANCE_ID)/databases/$(SPANNER_DATABASE_ID)

.PHONY: info
info: ## 현재 설정 정보
	@echo "$(GREEN)ℹ️ Spanner Configuration$(NC)"
	@echo "=================================="
	@echo "Emulator Host: $(SPANNER_EMULATOR_HOST)"
	@echo "Project ID:    $(SPANNER_PROJECT_ID)"
	@echo "Instance ID:   $(SPANNER_INSTANCE_ID)"
	@echo "Database ID:   $(SPANNER_DATABASE_ID)"
	@echo ""
	@echo "Ports:"
	@echo "  gRPC: :9010"
	@echo "  HTTP: :9020"
	@echo ""
	@echo "Docker Status:"
	@docker ps --filter "name=spanner" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

