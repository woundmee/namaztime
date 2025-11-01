package main

import (
	"namaztimeApi/internal/cache"
	"namaztimeApi/internal/configs/slogger"
	"namaztimeApi/internal/handlers"
	"namaztimeApi/internal/services"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		panic("не удалось загрузить переменные окружения")
	}

	logFile := os.Getenv("API_LOG_FILE")

	// init logger
	logger, err := slogger.Init(logFile)
	if err != nil {
		panic("Не удалось инициализировать логгер: " + err.Error())
	}

	// init cache
	cache := cache.New(logger)

	// init service
	service := services.New(logger, cache)
	go service.StartDailyUpdateCache()

	// init handler
	handler := handlers.New(logger, service)

	// init echo server
	logger.Info("cоздаю новый echo-сервер")
	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/namaztime/month", handler.GetNamazDataHandler)
	e.GET("/namaztime/today", handler.GetNamazDataFilteredHandler)

	logger.Info("начинаю запуск приложения", "address", os.Getenv("ADDRESS"))
	e.Start(os.Getenv("ADDRESS"))
}
