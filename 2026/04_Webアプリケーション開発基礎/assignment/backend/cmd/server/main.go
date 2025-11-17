package main

import (
	"path/filepath"

	"backend/internal/handlers"
	"backend/internal/store"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Echoインスタンスを作成
	e := echo.New()

	// ミドルウェアを設定
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// TODOストアを初期化
	todoStore, err := store.NewTodoStore(filepath.Join("tmp", "todos.json"))
	if err != nil {
		e.Logger.Fatal(err)
	}

	// TODOハンドラーを作成してルートを登録
	todoHandler := handlers.NewTodoHandler(todoStore)
	todoHandler.RegisterRoutes(e)

	// サーバーを起動
	e.Logger.Fatal(e.Start(":8080"))
}
