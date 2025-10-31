package main

import (
	"os"
	"telegramBot/clients/namaznsk"
	"telegramBot/internal/cache"
	"telegramBot/internal/handlers"
	"telegramBot/internal/services"
	storage "telegramBot/internal/storage/sqlite"
	"telegramBot/slogger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		panic("Не удалось загрузить переменные окружения")
	}

	// log file
	botLogFile := os.Getenv("LOG_FILE")
	logger, err := slogger.Init(botLogFile)
	if err != nil {
		panic("не удалось инициализировать логгер")
	}

	// cache init
	cache := cache.New(logger)
	if cache == nil {
		logger.Error("не удалось инициализировать кэш", "cache", cache)
	}

	// storage init
	storage := storage.New(logger)
	if storage == nil {
		logger.Error("не удалось инициализировать storage", "storage", storage)
	}
	storage.Create()

	urlTodaySchedule := os.Getenv("URL_TODAY_SCHEDULE")

	// client init
	clientNamaznsk := namaznsk.New(logger, cache, urlTodaySchedule)
	if clientNamaznsk == nil {
		logger.Error("не удалось инициализировать clientNamaznsk", "clientNamaznsk", clientNamaznsk)
	}
	go clientNamaznsk.StartDailyUpdateCache()

	// bot
	token := os.Getenv("TG_BOT_TOKEN")
	bot, err := tgbotapi.NewBotAPI(token)
	if bot == nil {
		logger.Error("не удалось инициализировать bot", "bot", bot)
	}
	if err != nil {
		logger.Error("не удалось загрузить переменные окружения", "error", err)
		panic("не удалось загрузить переменные окружения")
		// log.Panic("Не удалось загрузить переменные окружения", "error", err, "bot", bot)
	}
	// bot.Debug = true

	// services init
	service := services.New(logger, clientNamaznsk, bot, storage)
	if service == nil {
		logger.Error("не удалось инициализировать service", "service", service)
	}
	// service.SetNamazClient(clientNamaznsk)
	go service.StartNamazNotifier()

	// bot init & start
	botHandler := handlers.New(logger, *bot, *clientNamaznsk, *service, storage)
	if botHandler == nil {
		logger.Error("не удалось инициализировать botHandler", "botHandler", botHandler)
	}
	botHandler.Start()

}
