package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/beego/beego/v2/server/web"
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

// HealthController - 헬스 체크 컨트롤러
type HealthController struct {
	web.Controller
}

// Get - 헬스 체크
func (c *HealthController) Get() {
	c.Data["json"] = Response{
		Message: "서버가 정상적으로 작동 중입니다",
		Data: map[string]string{
			"status": "ok",
		},
	}
	c.ServeJSON()
}

// ItemsController - 아이템 컨트롤러
type ItemsController struct {
	web.Controller
}

// GetAll - 모든 아이템 조회
func (c *ItemsController) GetAll() {
	store.mu.RLock()
	defer store.mu.RUnlock()

	items := make([]Item, 0, len(store.items))
	for _, item := range store.items {
		items = append(items, item)
	}

	c.Data["json"] = Response{
		Data:  items,
		Count: len(items),
	}
	c.ServeJSON()
}

// Get - 특정 아이템 조회
func (c *ItemsController) Get() {
	idStr := c.Ctx.Input.Param(":id")
	id := 0
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = Response{
			Error: "잘못된 ID 형식입니다",
		}
		c.ServeJSON()
		return
	}

	store.mu.RLock()
	item, exists := store.items[id]
	store.mu.RUnlock()

	if !exists {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = Response{
			Error: "아이템을 찾을 수 없습니다",
		}
		c.ServeJSON()
		return
	}

	c.Data["json"] = Response{
		Data: item,
	}
	c.ServeJSON()
}

// Post - 새 아이템 생성
func (c *ItemsController) Post() {
	var newItem Item

	if err := c.BindJSON(&newItem); err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = Response{
			Error: "잘못된 요청 데이터: " + err.Error(),
		}
		c.ServeJSON()
		return
	}

	// 유효성 검증
	if newItem.Name == "" {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = Response{
			Error: "이름은 필수 항목입니다",
		}
		c.ServeJSON()
		return
	}

	if newItem.Price < 0 {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = Response{
			Error: "가격은 0 이상이어야 합니다",
		}
		c.ServeJSON()
		return
	}

	store.mu.Lock()
	newItem.ID = store.nextID
	store.items[newItem.ID] = newItem
	store.nextID++
	store.mu.Unlock()

	c.Ctx.Output.SetStatus(201)
	c.Data["json"] = Response{
		Message: "아이템이 생성되었습니다",
		Data:    newItem,
	}
	c.ServeJSON()
}

// Put - 아이템 수정
func (c *ItemsController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id := 0
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = Response{
			Error: "잘못된 ID 형식입니다",
		}
		c.ServeJSON()
		return
	}

	var updatedItem Item
	if err := c.BindJSON(&updatedItem); err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = Response{
			Error: "잘못된 요청 데이터: " + err.Error(),
		}
		c.ServeJSON()
		return
	}

	// 유효성 검증
	if updatedItem.Name == "" {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = Response{
			Error: "이름은 필수 항목입니다",
		}
		c.ServeJSON()
		return
	}

	if updatedItem.Price < 0 {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = Response{
			Error: "가격은 0 이상이어야 합니다",
		}
		c.ServeJSON()
		return
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	if _, exists := store.items[id]; !exists {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = Response{
			Error: "아이템을 찾을 수 없습니다",
		}
		c.ServeJSON()
		return
	}

	updatedItem.ID = id
	store.items[id] = updatedItem

	c.Data["json"] = Response{
		Message: "아이템이 수정되었습니다",
		Data:    updatedItem,
	}
	c.ServeJSON()
}

// Delete - 아이템 삭제
func (c *ItemsController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id := 0
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = Response{
			Error: "잘못된 ID 형식입니다",
		}
		c.ServeJSON()
		return
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	if _, exists := store.items[id]; !exists {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = Response{
			Error: "아이템을 찾을 수 없습니다",
		}
		c.ServeJSON()
		return
	}

	delete(store.items, id)

	c.Data["json"] = Response{
		Message: fmt.Sprintf("아이템 ID %d가 삭제되었습니다", id),
	}
	c.ServeJSON()
}

func main() {
	// Beego 설정
	web.BConfig.RunMode = "prod"
	web.BConfig.CopyRequestBody = true

	// 라우트 설정
	web.Router("/health", &HealthController{})

	// API 네임스페이스 (v1)
	ns := web.NewNamespace("/api/v1",
		web.NSRouter("/items", &ItemsController{}, "get:GetAll;post:Post"),
		web.NSRouter("/items/:id", &ItemsController{}, "get:Get;put:Put;delete:Delete"),
	)
	web.AddNamespace(ns)

	// 서버 시작
	log.Println("🚀 Beego 서버가 :8080 포트에서 시작되었습니다")
	web.Run(":8080")
}

