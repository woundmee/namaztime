package main

import (
	"log"
	"os"
	"strconv"
	"sync"

	"telegramBot/clients/namaznsk"
	"telegramBot/internal/cache"
	"telegramBot/internal/handlers"
	"telegramBot/internal/services"
	"telegramBot/slogger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {

	// TODO START
	// + 1. Написать функцию, которая в 00:00 ежедневно обновляет кеш
	// 2. Написать функцию, которая по времени отправляет уведомление. Пример: наступил Зухр - оповещает.
	// 3. Обновить readme.md: добавить информацию как поднять API и Bot (подробно).
	// TODO END

	var wg sync.WaitGroup

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

	urlTodaySchedule := os.Getenv("URL_TODAY_SCHEDULE")

	// cache init
	cache := cache.New(logger)

	// client init
	clientNamaznsk := namaznsk.New(logger, cache, urlTodaySchedule)

	wg.Add(1)
	go clientNamaznsk.StartDailyUpdateCache()

	// bot
	token := os.Getenv("TG_BOT_TOKEN")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic("Не удалось загрузить переменные окружения")
	}

	// bot.Debug = true

	service := services.New(logger, bot)
	service.SetNamazClient(clientNamaznsk)
	botID, err := strconv.Atoi(os.Getenv("TG_BOT_ID"))
	if err != nil {
		logger.Error("Не удалось получить/сконвертировать переменную окружения TG_BOT_ID", "error", err)
	}

	go service.StartNamazNotifier(int64(botID))

	// bot init & start
	botHandler := handlers.New(logger, *bot, *clientNamaznsk, *service)
	botHandler.Start()

	wg.Wait()

}
