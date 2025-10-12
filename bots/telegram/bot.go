package telegram

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"namaztimeApi/models"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

// todo: добавить структуру + логирование
// todo: сделать функцию, которая оповещает по каждому намазу. Проверка - каждые 30сек.

func Bot() {

	godotenv.Load()

	token := os.Getenv("TG_BOT_TOKEN")
	if token == "" {
		log.Fatal("TG_BOT_TOKEN is not set")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Бот %s запущен!", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "help":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, helper())
				bot.Send(msg)
			case "today":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, namazToday())
				bot.Send(msg)
			}
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		// msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		// bot.Send(msg)
	}
}

func helper() string {
	return `👋 Ассаляму аляйкум! Я telegram-бот для отправки расписания намазов по г. Норильск.

	ℹ️ Используйте команду /help для получения доп. информации.
	
	Источник: @nurdkamal
	Веб-версия: namaznsk.ru`
}

func namazTodayJson() models.NamazTime {
	url := "https://namaznsk.ru/namaztime/today"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Ошибка получения расписания через API: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Неверный статус ответа: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Ошибка при чтении тела ответа: %v", err)
	}

	var data models.NamazTime

	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatalf("Ошибка при парсинге JSON: %v", err)
	}

	return data
}

func namazToday() string {
	namaz := namazTodayJson()

	res := fmt.Sprintf("🕘 Расписание на сегодня\n\n%s\tФаджр\n%s\tВосход\n%s\tЗухр\n%s\tАср\n%s\tМагриб\n%s\tИша",
		namaz.Fajr, namaz.Sunrise, namaz.Zuhr, namaz.Asr, namaz.Magrib, namaz.Isha)
	return res
}
