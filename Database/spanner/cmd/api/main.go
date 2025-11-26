package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/spanner"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

// yo로 생성된 모델을 임포트 (생성 후 주석 해제)
// import "github.com/milman2/go-api/spanner-yo/models"

var (
	projectID  = getEnv("SPANNER_PROJECT_ID", "test-project")
	instanceID = getEnv("SPANNER_INSTANCE_ID", "test-instance")
	databaseID = getEnv("SPANNER_DATABASE_ID", "test-database")
)

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	ctx := context.Background()

	// Spanner 클라이언트 생성
	database := "projects/" + projectID + "/instances/" + instanceID + "/databases/" + databaseID
	client, err := spanner.NewClient(ctx, database)
	if err != nil {
		log.Fatalf("Spanner 클라이언트 생성 실패: %v", err)
	}
	defer client.Close()

	log.Printf("✅ Spanner 연결 성공: %s\n", database)

	// HTTP 라우터 설정
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","database":"spanner"}`))
	})

	// API 라우트
	r.Route("/api/v1", func(r chi.Router) {
		// Users
		r.Route("/users", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				getUsers(w, r, client)
			})
			r.Post("/", func(w http.ResponseWriter, r *http.Request) {
				createUser(w, r, client)
			})
			r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
				getUser(w, r, client)
			})
			r.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
				deleteUser(w, r, client)
			})
		})
	})

	// 서버 시작
	addr := ":8080"
	log.Printf("\n🚀 Spanner + yo 서버 시작\n")
	log.Printf("=========================================\n")
	log.Printf("📦 Database: Google Cloud Spanner\n")
	log.Printf("🔨 Code Generator: yo (go.mercari.io/yo)\n")
	log.Printf("🔧 Migration: Hammer + Wrench\n")
	log.Printf("\n")
	log.Printf("🌐 서버 주소: http://localhost%s\n", addr)
	log.Printf("=========================================\n")

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("서버 시작 실패: %v", err)
	}
}

// Handlers (yo 생성 후 모델 사용)
func getUsers(w http.ResponseWriter, r *http.Request, client *spanner.Client) {
	ctx := r.Context()

	// yo 생성 코드를 사용한 조회
	// users, err := models.FindAllUsers(ctx, client)

	// 임시: Raw SQL 사용
	stmt := spanner.Statement{SQL: `SELECT id, email, name, created_at, updated_at FROM users ORDER BY created_at DESC`}
	iter := client.Single().Query(ctx, stmt)
	defer iter.Stop()

	type User struct {
		ID        string `json:"id"`
		Email     string `json:"email"`
		Name      string `json:"name"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	var users []User
	for {
		row, err := iter.Next()
		if err != nil {
			break
		}
		var user User
		if err := row.Columns(&user.ID, &user.Email, &user.Name, &user.CreatedAt, &user.UpdatedAt); err != nil {
			continue
		}
		users = append(users, user)
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"data":  users,
		"count": len(users),
	})
}

func createUser(w http.ResponseWriter, r *http.Request, client *spanner.Client) {
	ctx := r.Context()

	type CreateUserRequest struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	var req CreateUserRequest
	if err := parseJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "잘못된 요청")
		return
	}

	// yo 생성 코드를 사용한 INSERT
	// user := &models.User{
	//     ID:    uuid.New().String(),
	//     Email: req.Email,
	//     Name:  req.Name,
	// }
	// _, err := client.Apply(ctx, []*spanner.Mutation{user.Insert(ctx)})

	// 임시: Raw Mutation 사용
	id := uuid.New().String()
	m := spanner.InsertMap("users", map[string]interface{}{
		"id":         id,
		"email":      req.Email,
		"name":       req.Name,
		"created_at": spanner.CommitTimestamp,
		"updated_at": spanner.CommitTimestamp,
	})

	_, err := client.Apply(ctx, []*spanner.Mutation{m})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "사용자 생성 실패: "+err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "사용자가 생성되었습니다",
		"id":      id,
		"email":   req.Email,
		"name":    req.Name,
	})
}

func getUser(w http.ResponseWriter, r *http.Request, client *spanner.Client) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	// yo 생성 코드 사용
	// user, err := models.FindUserByID(ctx, client, id)

	// 임시: Raw SQL
	stmt := spanner.Statement{
		SQL:    `SELECT id, email, name, created_at, updated_at FROM users WHERE id = @id`,
		Params: map[string]interface{}{"id": id},
	}
	iter := client.Single().Query(ctx, stmt)
	defer iter.Stop()

	row, err := iter.Next()
	if err != nil {
		respondError(w, http.StatusNotFound, "사용자를 찾을 수 없습니다")
		return
	}

	type User struct {
		ID        string `json:"id"`
		Email     string `json:"email"`
		Name      string `json:"name"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	var user User
	if err := row.Columns(&user.ID, &user.Email, &user.Name, &user.CreatedAt, &user.UpdatedAt); err != nil {
		respondError(w, http.StatusInternalServerError, "데이터 파싱 실패")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{"data": user})
}

func deleteUser(w http.ResponseWriter, r *http.Request, client *spanner.Client) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	// yo 생성 코드 사용
	// user, _ := models.FindUserByID(ctx, client, id)
	// _, err := client.Apply(ctx, []*spanner.Mutation{user.Delete(ctx)})

	// 임시: Raw Mutation
	m := spanner.Delete("users", spanner.Key{id})
	_, err := client.Apply(ctx, []*spanner.Mutation{m})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "삭제 실패")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Helper functions
func parseJSON(r *http.Request, v interface{}) error {
	return nil // 간단히 구현
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	// JSON 인코딩 생략 (간단히)
}

func respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
}
