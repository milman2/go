package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Item 구조체
type Item struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
}

// 메모리 저장소
type ItemStore struct {
	mu     sync.RWMutex
	items  map[int]Item
	nextID int
}

var store = &ItemStore{
	items:  make(map[int]Item),
	nextID: 1,
}

// Response 구조체
type Response struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
	Count   int         `json:"count,omitempty"`
}

func main() {
	r := chi.NewRouter()

	// 미들웨어 설정
	r.Use(middleware.Logger)          // 로깅
	r.Use(middleware.Recoverer)       // 패닉 복구
	r.Use(middleware.RequestID)       // 요청 ID 생성
	r.Use(middleware.RealIP)          // 실제 IP 추출
	r.Use(middleware.StripSlashes)    // URL 끝의 슬래시 제거

	// Health check
	r.Get("/health", healthCheck)

	// API v1 라우트
	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/items", func(r chi.Router) {
			r.Get("/", getItems)           // 목록 조회
			r.Post("/", createItem)        // 생성
			
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", getItem)        // 단일 조회
				r.Put("/", updateItem)     // 수정
				r.Delete("/", deleteItem)  // 삭제
			})
		})
	})

	// 서버 시작
	addr := ":8080"
	log.Printf("🚀 Chi 서버가 %s 포트에서 시작되었습니다\n", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("서버 시작 실패: %v", err)
	}
}

// respondJSON - JSON 응답 헬퍼 함수
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("JSON 인코딩 에러: %v", err)
	}
}

// healthCheck - 헬스 체크 엔드포인트
func healthCheck(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, Response{
		Message: "서버가 정상적으로 작동 중입니다",
		Data: map[string]string{
			"status": "ok",
		},
	})
}

// getItems - 모든 아이템 조회
func getItems(w http.ResponseWriter, r *http.Request) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	items := make([]Item, 0, len(store.items))
	for _, item := range store.items {
		items = append(items, item)
	}

	respondJSON(w, http.StatusOK, Response{
		Data:  items,
		Count: len(items),
	})
}

// getItem - 특정 아이템 조회
func getItem(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, Response{
			Error: "잘못된 ID 형식입니다",
		})
		return
	}

	store.mu.RLock()
	item, exists := store.items[id]
	store.mu.RUnlock()

	if !exists {
		respondJSON(w, http.StatusNotFound, Response{
			Error: "아이템을 찾을 수 없습니다",
		})
		return
	}

	respondJSON(w, http.StatusOK, Response{
		Data: item,
	})
}

// createItem - 새 아이템 생성
func createItem(w http.ResponseWriter, r *http.Request) {
	var newItem Item

	if err := json.NewDecoder(r.Body).Decode(&newItem); err != nil {
		respondJSON(w, http.StatusBadRequest, Response{
			Error: "잘못된 요청 데이터: " + err.Error(),
		})
		return
	}

	// 유효성 검증
	if newItem.Name == "" {
		respondJSON(w, http.StatusBadRequest, Response{
			Error: "이름은 필수 항목입니다",
		})
		return
	}

	if newItem.Price < 0 {
		respondJSON(w, http.StatusBadRequest, Response{
			Error: "가격은 0 이상이어야 합니다",
		})
		return
	}

	store.mu.Lock()
	newItem.ID = store.nextID
	store.items[newItem.ID] = newItem
	store.nextID++
	store.mu.Unlock()

	respondJSON(w, http.StatusCreated, Response{
		Message: "아이템이 생성되었습니다",
		Data:    newItem,
	})
}

// updateItem - 아이템 수정
func updateItem(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, Response{
			Error: "잘못된 ID 형식입니다",
		})
		return
	}

	var updatedItem Item
	if err := json.NewDecoder(r.Body).Decode(&updatedItem); err != nil {
		respondJSON(w, http.StatusBadRequest, Response{
			Error: "잘못된 요청 데이터: " + err.Error(),
		})
		return
	}

	// 유효성 검증
	if updatedItem.Name == "" {
		respondJSON(w, http.StatusBadRequest, Response{
			Error: "이름은 필수 항목입니다",
		})
		return
	}

	if updatedItem.Price < 0 {
		respondJSON(w, http.StatusBadRequest, Response{
			Error: "가격은 0 이상이어야 합니다",
		})
		return
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	if _, exists := store.items[id]; !exists {
		respondJSON(w, http.StatusNotFound, Response{
			Error: "아이템을 찾을 수 없습니다",
		})
		return
	}

	updatedItem.ID = id
	store.items[id] = updatedItem

	respondJSON(w, http.StatusOK, Response{
		Message: "아이템이 수정되었습니다",
		Data:    updatedItem,
	})
}

// deleteItem - 아이템 삭제
func deleteItem(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, Response{
			Error: "잘못된 ID 형식입니다",
		})
		return
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	if _, exists := store.items[id]; !exists {
		respondJSON(w, http.StatusNotFound, Response{
			Error: "아이템을 찾을 수 없습니다",
		})
		return
	}

	delete(store.items, id)

	respondJSON(w, http.StatusOK, Response{
		Message: fmt.Sprintf("아이템 ID %d가 삭제되었습니다", id),
	})
}

