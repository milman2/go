package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Item 구조체
type Item struct {
	ID          int    `json:"id"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Price       int    `json:"price" validate:"required,gte=0"`
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
	e := echo.New()

	// 배너 숨기기
	e.HideBanner = true

	// 미들웨어 설정
	e.Use(middleware.Logger())                                       // 로깅
	e.Use(middleware.Recover())                                      // 패닉 복구
	e.Use(middleware.RequestID())                                    // 요청 ID
	e.Use(middleware.CORS())                                         // CORS
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{          // Gzip 압축
		Level: 5,
	}))

	// Health check
	e.GET("/health", healthCheck)

	// API v1 그룹
	v1 := e.Group("/api/v1")
	{
		// 아이템 라우트
		items := v1.Group("/items")
		{
			items.GET("", getItems)           // 목록 조회
			items.GET("/:id", getItem)        // 단일 조회
			items.POST("", createItem)        // 생성
			items.PUT("/:id", updateItem)     // 수정
			items.DELETE("/:id", deleteItem)  // 삭제
		}
	}

	// 서버 시작
	addr := ":8080"
	log.Printf("🚀 Echo 서버가 %s 포트에서 시작되었습니다\n", addr)
	if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
		log.Fatalf("서버 시작 실패: %v", err)
	}
}

// healthCheck - 헬스 체크 엔드포인트
func healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, Response{
		Message: "서버가 정상적으로 작동 중입니다",
		Data: map[string]string{
			"status": "ok",
		},
	})
}

// getItems - 모든 아이템 조회
func getItems(c echo.Context) error {
	store.mu.RLock()
	defer store.mu.RUnlock()

	items := make([]Item, 0, len(store.items))
	for _, item := range store.items {
		items = append(items, item)
	}

	return c.JSON(http.StatusOK, Response{
		Data:  items,
		Count: len(items),
	})
}

// getItem - 특정 아이템 조회
func getItem(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Error: "잘못된 ID 형식입니다",
		})
	}

	store.mu.RLock()
	item, exists := store.items[id]
	store.mu.RUnlock()

	if !exists {
		return c.JSON(http.StatusNotFound, Response{
			Error: "아이템을 찾을 수 없습니다",
		})
	}

	return c.JSON(http.StatusOK, Response{
		Data: item,
	})
}

// createItem - 새 아이템 생성
func createItem(c echo.Context) error {
	var newItem Item

	if err := c.Bind(&newItem); err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Error: "잘못된 요청 데이터: " + err.Error(),
		})
	}

	// 유효성 검증
	if newItem.Name == "" {
		return c.JSON(http.StatusBadRequest, Response{
			Error: "이름은 필수 항목입니다",
		})
	}

	if newItem.Price < 0 {
		return c.JSON(http.StatusBadRequest, Response{
			Error: "가격은 0 이상이어야 합니다",
		})
	}

	store.mu.Lock()
	newItem.ID = store.nextID
	store.items[newItem.ID] = newItem
	store.nextID++
	store.mu.Unlock()

	return c.JSON(http.StatusCreated, Response{
		Message: "아이템이 생성되었습니다",
		Data:    newItem,
	})
}

// updateItem - 아이템 수정
func updateItem(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Error: "잘못된 ID 형식입니다",
		})
	}

	var updatedItem Item
	if err := c.Bind(&updatedItem); err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Error: "잘못된 요청 데이터: " + err.Error(),
		})
	}

	// 유효성 검증
	if updatedItem.Name == "" {
		return c.JSON(http.StatusBadRequest, Response{
			Error: "이름은 필수 항목입니다",
		})
	}

	if updatedItem.Price < 0 {
		return c.JSON(http.StatusBadRequest, Response{
			Error: "가격은 0 이상이어야 합니다",
		})
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	if _, exists := store.items[id]; !exists {
		return c.JSON(http.StatusNotFound, Response{
			Error: "아이템을 찾을 수 없습니다",
		})
	}

	updatedItem.ID = id
	store.items[id] = updatedItem

	return c.JSON(http.StatusOK, Response{
		Message: "아이템이 수정되었습니다",
		Data:    updatedItem,
	})
}

// deleteItem - 아이템 삭제
func deleteItem(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Error: "잘못된 ID 형식입니다",
		})
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	if _, exists := store.items[id]; !exists {
		return c.JSON(http.StatusNotFound, Response{
			Error: "아이템을 찾을 수 없습니다",
		})
	}

	delete(store.items, id)

	return c.JSON(http.StatusOK, Response{
		Message: fmt.Sprintf("아이템 ID %d가 삭제되었습니다", id),
	})
}

