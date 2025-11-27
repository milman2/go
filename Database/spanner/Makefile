# Makefile for Google Cloud Spanner + hammer + yo

# Makefile path information
MAKEFILE_PATH:=$(abspath $(lastword $(MAKEFILE_LIST)))
MAKEFILE_DIR:=$(patsubst %/,%,$(dir ${MAKEFILE_PATH}))

# Directories
CWD_DIR:=${MAKEFILE_DIR:%/=%}
BIN_DIR:=${CWD_DIR}/bin
EXT_BIN_DIR:=${BIN_DIR}/ext

# Spanner 설정
SPANNER_EMULATOR_HOST=localhost:9010
SPANNER_PROJECT_ID=test-project
SPANNER_INSTANCE_ID=test-instance
SPANNER_DATABASE_ID=test-db

# 경로 설정
SCHEMA_DIR=schema
SCHEMA_FILE=${SCHEMA_DIR}/schema.sql
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
setup-instance: ## Spanner Instance 생성 (데이터베이스는 createdb에서 생성)
	@echo "$(GREEN)📦 Creating Spanner instance...$(NC)"
	@gcloud config set auth/disable_credentials true
	@gcloud config set project $(SPANNER_PROJECT_ID)
	@gcloud config set api_endpoint_overrides/spanner http://localhost:9020/
	@gcloud spanner instances create $(SPANNER_INSTANCE_ID) \
		--config=emulator-config \
		--description="Test Instance" \
		--nodes=1 || true
	@echo "$(GREEN)✅ Instance ready$(NC)"

.PHONY: createdb
createdb: ${EXT_BIN_DIR}/wrench ${EXT_BIN_DIR}/hammer ## Spanner 데이터베이스 생성 (hammer create 사용)
	@echo "$(GREEN)🔨 Creating database with hammer...$(NC)"
	@SPANNER_EMULATOR_HOST=$(SPANNER_EMULATOR_HOST) \
		${EXT_BIN_DIR}/wrench drop --project $(SPANNER_PROJECT_ID) --instance $(SPANNER_INSTANCE_ID) --database $(SPANNER_DATABASE_ID) 2>/dev/null || true
	@SPANNER_EMULATOR_HOST=$(SPANNER_EMULATOR_HOST) \
		${EXT_BIN_DIR}/hammer create spanner://projects/$(SPANNER_PROJECT_ID)/instances/$(SPANNER_INSTANCE_ID)/databases/$(SPANNER_DATABASE_ID) $(SCHEMA_FILE)
	@echo "$(GREEN)✅ Database created from $(SCHEMA_FILE)$(NC)"

.PHONY: db-apply
db-apply: ${EXT_BIN_DIR}/hammer ## 스키마 변경사항 적용 (hammer apply 사용)
	@echo "$(GREEN)🔨 Applying schema changes with hammer...$(NC)"
	@SPANNER_EMULATOR_HOST=$(SPANNER_EMULATOR_HOST) \
		${EXT_BIN_DIR}/hammer apply spanner://projects/$(SPANNER_PROJECT_ID)/instances/$(SPANNER_INSTANCE_ID)/databases/$(SPANNER_DATABASE_ID) $(SCHEMA_FILE)
	@echo "$(GREEN)✅ Schema changes applied$(NC)"

.PHONY: db-diff
db-diff: ${EXT_BIN_DIR}/hammer ## 현재 DB와 스키마 파일 차이 확인 (hammer diff 사용)
	@echo "$(GREEN)🔍 Checking schema differences...$(NC)"
	@SPANNER_EMULATOR_HOST=$(SPANNER_EMULATOR_HOST) \
		${EXT_BIN_DIR}/hammer diff spanner://projects/$(SPANNER_PROJECT_ID)/instances/$(SPANNER_INSTANCE_ID)/databases/$(SPANNER_DATABASE_ID) $(SCHEMA_FILE)

.PHONY: db-export
db-export: ${EXT_BIN_DIR}/hammer ## 현재 스키마를 파일로 내보내기 (hammer export 사용)
	@echo "$(GREEN)📤 Exporting schema with hammer...$(NC)"
	@SPANNER_EMULATOR_HOST=$(SPANNER_EMULATOR_HOST) \
		${EXT_BIN_DIR}/hammer export spanner://projects/$(SPANNER_PROJECT_ID)/instances/$(SPANNER_INSTANCE_ID)/databases/$(SPANNER_DATABASE_ID) > exported_schema.sql
	@echo "$(GREEN)✅ Schema exported to exported_schema.sql$(NC)"

.PHONY: dropdb
dropdb: ${EXT_BIN_DIR}/wrench ## 데이터베이스 삭제 (wrench drop 사용)
	@echo "$(YELLOW)⚠️  데이터베이스를 삭제합니다!$(NC)"
	@read -p "계속하시겠습니까? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		echo "$(YELLOW)Dropping database...$(NC)"; \
		SPANNER_EMULATOR_HOST=$(SPANNER_EMULATOR_HOST) \
		${EXT_BIN_DIR}/wrench drop --project $(SPANNER_PROJECT_ID) --instance $(SPANNER_INSTANCE_ID) --database $(SPANNER_DATABASE_ID); \
		echo "$(GREEN)✅ Database dropped$(NC)"; \
	fi

.PHONY: resetdb
resetdb: createdb ## 데이터베이스 리셋 (삭제 후 재생성)
	@echo "$(GREEN)✅ Database reset completed$(NC)"

.PHONY: seed-data
seed-data: ## 샘플 데이터 삽입 (Go 스크립트 사용)
	@echo "$(GREEN)🌱 Seeding sample data...$(NC)"
	@SPANNER_EMULATOR_HOST=$(SPANNER_EMULATOR_HOST) \
		SPANNER_PROJECT_ID=$(SPANNER_PROJECT_ID) \
		SPANNER_INSTANCE_ID=$(SPANNER_INSTANCE_ID) \
		SPANNER_DATABASE_ID=$(SPANNER_DATABASE_ID) \
		go run scripts/seed_data.go
	@echo "$(GREEN)✅ Sample data seeded$(NC)"

.PHONY: clear-data
clear-data: ## 모든 데이터 삭제 (테이블 구조는 유지)
	@echo "$(YELLOW)⚠️  모든 데이터를 삭제합니다!$(NC)"
	@read -p "계속하시겠습니까? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		echo "$(YELLOW)🗑️  Clearing all data...$(NC)"; \
		SPANNER_EMULATOR_HOST=$(SPANNER_EMULATOR_HOST) \
		gcloud spanner databases execute-sql $(SPANNER_DATABASE_ID) \
			--instance=$(SPANNER_INSTANCE_ID) \
			--sql="DELETE FROM posts WHERE TRUE; DELETE FROM users WHERE TRUE;" || true; \
		echo "$(GREEN)✅ All data cleared$(NC)"; \
	fi

.PHONY: test-query
test-query: ## 샘플 쿼리 테스트
	@echo "$(GREEN)🔍 Running test queries...$(NC)"
	@echo ""
	@echo "$(YELLOW)📊 Users count:$(NC)"
	@SPANNER_EMULATOR_HOST=$(SPANNER_EMULATOR_HOST) \
		gcloud spanner databases execute-sql $(SPANNER_DATABASE_ID) \
		--instance=$(SPANNER_INSTANCE_ID) \
		--sql="SELECT COUNT(*) as user_count FROM users"
	@echo ""
	@echo "$(YELLOW)📊 Posts count:$(NC)"
	@SPANNER_EMULATOR_HOST=$(SPANNER_EMULATOR_HOST) \
		gcloud spanner databases execute-sql $(SPANNER_DATABASE_ID) \
		--instance=$(SPANNER_INSTANCE_ID) \
		--sql="SELECT COUNT(*) as post_count FROM posts"
	@echo ""
	@echo "$(YELLOW)📊 Published posts by user:$(NC)"
	@SPANNER_EMULATOR_HOST=$(SPANNER_EMULATOR_HOST) \
		gcloud spanner databases execute-sql $(SPANNER_DATABASE_ID) \
		--instance=$(SPANNER_INSTANCE_ID) \
		--sql="SELECT u.name, COUNT(p.id) as post_count FROM users u LEFT JOIN posts p ON u.id = p.user_id WHERE p.published = TRUE OR p.published IS NULL GROUP BY u.name ORDER BY post_count DESC"

.PHONY: generate-migration
generate-migration: ${EXT_BIN_DIR}/hammer ## 스키마 diff 기반 마이그레이션 템플릿 생성
	@./scripts/generate_migration.sh

.PHONY: migrate-dml
migrate-dml: ${EXT_BIN_DIR}/wrench ## DML 마이그레이션 실행 (wrench apply --dml) - 더 이상 사용 안 함
	@echo "$(YELLOW)⚠️  이 명령어는 더 이상 사용되지 않습니다.$(NC)"
	@echo "$(YELLOW)대신 'make seed-data'를 사용하세요.$(NC)"

.PHONY: generate-models
generate-models: ${EXT_BIN_DIR}/yo ## yo로 코드 생성 (스키마 → Go 코드)
	@echo "$(GREEN)🔨 Generating code with yo...$(NC)"
	@rm -rf $(MODELS_DIR)/*.yo.go
	@mkdir -p $(MODELS_DIR)
	@SPANNER_EMULATOR_HOST=$(SPANNER_EMULATOR_HOST) \
		${EXT_BIN_DIR}/yo $(SPANNER_PROJECT_ID) $(SPANNER_INSTANCE_ID) $(SPANNER_DATABASE_ID) \
		-o $(MODELS_DIR) -p models --ignore-tables SchemaMigrations
	@echo "$(GREEN)✅ Code generation completed: $(MODELS_DIR)/$(NC)"
	@ls -lh $(MODELS_DIR)/*.yo.go

.PHONY: build/ext
build/ext: ${EXT_BIN_DIR}/hammer ${EXT_BIN_DIR}/wrench ${EXT_BIN_DIR}/yo ## 외부 도구 빌드 (hammer, wrench, yo)
	@echo "$(GREEN)✅ External tools built$(NC)"
	@echo "  - hammer: ${EXT_BIN_DIR}/hammer"
	@echo "  - wrench: ${EXT_BIN_DIR}/wrench"
	@echo "  - yo: ${EXT_BIN_DIR}/yo"

${EXT_BIN_DIR}/%: go.mod go.sum ${CWD_DIR}/ext/%.go
	@mkdir -p $(dir $@)
	@echo "$(GREEN)📦 Building $(notdir $@)...$(NC)"
	@GO_ARTIFACT_PATH=$@ go generate -tags ext $(patsubst ${EXT_BIN_DIR}/%,${CWD_DIR}/ext/%.go,$@)

.PHONY: init
init: docker-up build/ext setup-instance createdb generate-models ## 전체 초기화 (한번에)
	@echo ""
	@echo "$(GREEN)🎉 Spanner 환경 초기화 완료!$(NC)"
	@echo "=================================="
	@echo "  ✅ Docker Spanner emulator 실행"
	@echo "  ✅ 외부 도구 빌드 (hammer, wrench, yo)"
	@echo "  ✅ Instance 생성"
	@echo "  ✅ Database 및 스키마 생성 (hammer)"
	@echo "  ✅ yo 코드 생성 완료"
	@echo ""
	@echo "$(YELLOW)📝 다음 단계:$(NC)"
	@echo "  make run          # 서버 실행"
	@echo "  make test         # API 테스트"
	@echo "  make db-diff      # 스키마 차이 확인"
	@echo "  make db-apply     # 스키마 변경 적용"

.PHONY: clean
clean: ## 생성된 파일 정리
	@echo "$(YELLOW)🧹 Cleaning up...$(NC)"
	@rm -rf $(MODELS_DIR)
	@echo "$(GREEN)✅ Cleaned$(NC)"

.PHONY: reset
reset: resetdb generate-models ## DB 리셋 및 코드 재생성
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
show-schema: ${EXT_BIN_DIR}/hammer ## 현재 스키마 확인 (hammer export)
	@echo "$(GREEN)📋 Current schema:$(NC)"
	@SPANNER_EMULATOR_HOST=$(SPANNER_EMULATOR_HOST) \
		${EXT_BIN_DIR}/hammer export spanner://projects/$(SPANNER_PROJECT_ID)/instances/$(SPANNER_INSTANCE_ID)/databases/$(SPANNER_DATABASE_ID)

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

.PHONY: test-connection
test-connection: ## Spanner 연결 테스트 (Go)
	@echo "$(GREEN)🔌 Spanner 연결 테스트...$(NC)"
	@SPANNER_EMULATOR_HOST=$(SPANNER_EMULATOR_HOST) go run test_connection.go

.PHONY: test-tables
test-tables: ## 테이블 정보 조회 (Go)
	@echo "$(GREEN)📊 테이블 정보 조회...$(NC)"
	@SPANNER_EMULATOR_HOST=$(SPANNER_EMULATOR_HOST) go run test_tables.go

.PHONY: test-crud
test-crud: ## CRUD 작업 테스트 (Go)
	@echo "$(GREEN)🧪 CRUD 테스트...$(NC)"
	@SPANNER_EMULATOR_HOST=$(SPANNER_EMULATOR_HOST) go run test_crud.go

.PHONY: test-all
test-all: ## 종합 테스트 (연결+테이블+CRUD)
	@./test_all.sh

.PHONY: sql
sql: ## SQL 실행 (사용: make sql SQL="SELECT * FROM users")
	@if [ -z "$(SQL)" ]; then \
		echo "$(YELLOW)사용법: make sql SQL=\"SELECT * FROM users\"$(NC)"; \
		exit 1; \
	fi
	@echo "$(GREEN)🔍 SQL 실행...$(NC)"
	@SPANNER_EMULATOR_HOST=$(SPANNER_EMULATOR_HOST) \
		gcloud spanner databases execute-sql $(SPANNER_DATABASE_ID) \
		--instance=$(SPANNER_INSTANCE_ID) \
		--sql='$(SQL)'

