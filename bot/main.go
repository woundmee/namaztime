package main

import (
	"log"
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
		log.Panic("Не удалось загрузить переменные окружения")
	}

	// log file
	botLogFile := os.Getenv("LOG_FILE")
	logger, err := slogger.Init(botLogFile)
	if err != nil {
		panic("не удалось инициализировать логгер")
	}

	// cache init
	cache := cache.New(logger)

	// storage init
	storage := storage.New(logger)
	storage.Create()

	urlTodaySchedule := os.Getenv("URL_TODAY_SCHEDULE")

	// client init
	clientNamaznsk := namaznsk.New(logger, cache, urlTodaySchedule)
	go clientNamaznsk.StartDailyUpdateCache()

	// bot
	token := os.Getenv("TG_BOT_TOKEN")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic("Не удалось загрузить переменные окружения")
	}
	// bot.Debug = true

	// services init
	service := services.New(logger, clientNamaznsk, bot, storage)
	// service.SetNamazClient(clientNamaznsk)
	go service.StartNamazNotifier()

	// bot init & start
	botHandler := handlers.New(logger, *bot, *clientNamaznsk, *service, storage)
	botHandler.Start()

}
