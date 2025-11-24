package main

import (
	"embed"
	"log"
	"net/http"

	"github.com/milman2/go-api/pkg/api"
	"github.com/oklog/ulid/v2"

	"github.com/labstack/echo/v4"
)

//go:embed swagger
//go:embed openapi.yaml
var content embed.FS

func main() {
	log.Println("ulid:", ulid.Make())
	server := api.NewServer()

	e := echo.New()

	api.RegisterHandlers(e, api.NewStrictHandler(
		server,
		[]api.StrictMiddlewareFunc{},
	))

	// Swagger UI
	e.GET("/swagger/*", echo.WrapHandler(http.StripPrefix("/swagger/", http.FileServer(http.FS(content)))))

	e.Start("127.0.0.1:8080")
}
