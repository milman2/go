package main

import (
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

// Item 구조체
type Item struct {
	ID          int    `json:"id"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Price       int    `json:"price" binding:"required,min=0"`
}

// 메모리 저장소
type ItemStore struct {
	mu      sync.RWMutex
	items   map[int]Item
	nextID  int
}

var store = &ItemStore{
	items:  make(map[int]Item),
	nextID: 1,
}

func main() {
	// Gin 라우터 생성 (릴리즈 모드)
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Health check
	r.GET("/health", healthCheck)

	// API v1 그룹
	v1 := r.Group("/api/v1")
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
	log.Println("🚀 Gin 서버가 :8080 포트에서 시작되었습니다")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("서버 시작 실패: %v", err)
	}
}

// healthCheck - 헬스 체크 엔드포인트
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"message": "서버가 정상적으로 작동 중입니다",
	})
}

// getItems - 모든 아이템 조회
func getItems(c *gin.Context) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	items := make([]Item, 0, len(store.items))
	for _, item := range store.items {
		items = append(items, item)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": items,
		"count": len(items),
	})
}

// getItem - 특정 아이템 조회
func getItem(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "잘못된 ID 형식입니다",
		})
		return
	}

	store.mu.RLock()
	item, exists := store.items[id]
	store.mu.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "아이템을 찾을 수 없습니다",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": item,
	})
}

// createItem - 새 아이템 생성
func createItem(c *gin.Context) {
	var newItem Item

	if err := c.ShouldBindJSON(&newItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "잘못된 요청 데이터: " + err.Error(),
		})
		return
	}

	store.mu.Lock()
	newItem.ID = store.nextID
	store.items[newItem.ID] = newItem
	store.nextID++
	store.mu.Unlock()

	c.JSON(http.StatusCreated, gin.H{
		"message": "아이템이 생성되었습니다",
		"data": newItem,
	})
}

// updateItem - 아이템 수정
func updateItem(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "잘못된 ID 형식입니다",
		})
		return
	}

	var updatedItem Item
	if err := c.ShouldBindJSON(&updatedItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "잘못된 요청 데이터: " + err.Error(),
		})
		return
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	if _, exists := store.items[id]; !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "아이템을 찾을 수 없습니다",
		})
		return
	}

	updatedItem.ID = id
	store.items[id] = updatedItem

	c.JSON(http.StatusOK, gin.H{
		"message": "아이템이 수정되었습니다",
		"data": updatedItem,
	})
}

// deleteItem - 아이템 삭제
func deleteItem(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "잘못된 ID 형식입니다",
		})
		return
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	if _, exists := store.items[id]; !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "아이템을 찾을 수 없습니다",
		})
		return
	}

	delete(store.items, id)

	c.JSON(http.StatusOK, gin.H{
		"message": "아이템이 삭제되었습니다",
	})
}

