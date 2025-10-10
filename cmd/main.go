package main

import (
	"namaztimeApi/internal/configs/slogger"
	"namaztimeApi/internal/handlers"
	"namaztimeApi/internal/services"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// todo: read from ENV
const ADDRESS = "localhost:8080"

func main() {
	logger := slogger.SetupLogger("local")

	service := services.NewService(logger)
	handler := handlers.NewHandler(logger, service)

	logger.Info("Создаю новый echo-сервер")
	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/namaztime/month", handler.GetNamazDataHandler)
	e.GET("/namaztime/today", handler.GetNamazDataFilteredHandler)

	// localhost:8080 - get from ENV
	logger.Info("Начинаю запуск приложения", "address", ADDRESS)
	e.Start(ADDRESS)
}
