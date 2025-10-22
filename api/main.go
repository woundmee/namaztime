package main

import (
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
		msg := "Не удалось загрузить переменные окружения"
		// logger.Error(msg)
		panic(msg)
	}

	logFile := os.Getenv("API_LOG_FILE")

	logger, err := slogger.Init(logFile)
	if err != nil {
		panic("Не удалось инициализировать логгер: " + err.Error())
	}
	service := services.New(logger)
	handler := handlers.New(logger, service)

	logger.Info("Создаю новый echo-сервер")
	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/namaztime/month", handler.GetNamazDataHandler)
	e.GET("/namaztime/today", handler.GetNamazDataFilteredHandler)

	logger.Info("Начинаю запуск приложения", "address", os.Getenv("ADDRESS"))
	e.Start(os.Getenv("ADDRESS"))
}
