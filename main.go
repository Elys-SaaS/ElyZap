package main

import (
	mw "github.com/Elys-SaaS/ElyZap/middleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type VideoPackage struct {
	VideoBytes  []byte
	VideoName   string
	RandomId    string
	ChunkNumber int
	ChunkCount  int
}

type App struct {
	videoMap    map[string][][]byte
	assignedIds map[string]string
}

func main() {
	app := App{
		videoMap:    make(map[string][][]byte),
		assignedIds: make(map[string]string),
	}

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello, World!!!")
	})
	e.Use(middleware.CORS())
	e.POST("/upload-video", app.UploadVideo, mw.CheckJWT)
	e.Logger.Fatal(e.Start(":8080"))
}
