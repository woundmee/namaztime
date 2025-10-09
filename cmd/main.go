package main

import (
	"log/slog"
	"namaztimeApi/internal/configs/slogger"
	"namaztimeApi/internal/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// todo: read from ENV
const ADDRESS = "localhost:8080"

func main() {

	err := slogger.Init("namaztime.log")
	if err != nil {
		panic("Не удалось инициализировать логгер: " + err.Error())
	}

	slog.Info("Создаю новый echo-сервер")
	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/namaztime/month", handlers.GetNamazDataHandler)
	e.GET("/namaztime/today", handlers.GetNamazDataFilteredHandler)

	// localhost:8080 - get from ENV
	slog.Info("Начинаю запуск приложения", "address", ADDRESS)
	e.Start(ADDRESS)
}
