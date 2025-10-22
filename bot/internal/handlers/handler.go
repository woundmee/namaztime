package handlers

import (
	"fmt"
	"log/slog"
	"os"
	"telegramBot/clients/namaznsk"
	"telegramBot/internal/services"
	storage "telegramBot/internal/storage/sqlite"

	// "telegramBot/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	logger  *slog.Logger
	bot     tgbotapi.BotAPI
	namaz   namaznsk.Namaz
	service services.Service
	storage storage.Database
}

func New(logger *slog.Logger, botKey tgbotapi.BotAPI, namaz namaznsk.Namaz, service services.Service) *Handler {
	return &Handler{
		logger:  logger,
		bot:     botKey,
		namaz:   namaz,
		service: service,
	}
}

func (h *Handler) Start() {
	fmt.Printf("Бот @%s запущен!\n", h.bot.Self.UserName)
	h.logger.Info("Бот запущен", "name", "@"+h.bot.Self.UserName)

	go h.service.StartNamazNotifier()

	discardData := h.DiscardOfflineUpdates()
	u := tgbotapi.NewUpdate(discardData + 1)
	u.Timeout = 60

	updates := h.bot.GetUpdatesChan(u)
	for update := range updates {
		h.handlerUpdate(update)
	}

}

// отбрасывает смс, которые были получены в оффлайне
func (h *Handler) DiscardOfflineUpdates() int {
	updates, err := h.bot.GetUpdates(tgbotapi.UpdateConfig{
		Offset:  0,
		Limit:   100,
		Timeout: 0,
	})

	if err != nil {
		h.logger.Error("Ошибка при получении данных", "error", err)
	}

	maxUpdateID := 0
	for _, update := range updates {
		if update.UpdateID > maxUpdateID {
			maxUpdateID = update.UpdateID
		}
	}

	return maxUpdateID
}

func (h *Handler) handlerUpdate(update tgbotapi.Update) {
	if update.Message != nil {
		h.logger.Info("Вызывана команда", "user", "@"+update.Message.From.UserName, "command", update.Message.Text, "group", "@"+update.Message.Chat.UserName, "groupName", update.Message.Chat.Title)

		if update.Message.IsCommand() {
			if update.Message.Command() == "start" {
				text := "🚀 Бот запущен!\n\n" +
					"🔔 Теперь вы будете получать уведомления о наступлении времени намазов."
				text2 := "Для получения справочной информации, используйте команду /help"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
				msg2 := tgbotapi.NewMessage(update.Message.Chat.ID, text2)
				h.bot.Send(msg)
				h.bot.Send(msg2)

				// сохраняю пользователя в БД
				h.storage.Insert(update.Message.Chat.ID, update.Message.Chat.UserName)
				// h.service.AddUser(update.Message.Chat.ID, update.Message.Chat.UserName)

				return
			}
			if update.Message.Command() == "help" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, h.service.CommandHelp())
				h.bot.Send(msg)
				return
			}
			if update.Message.Command() == "today" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, h.service.CommandToday(h.namaz))
				h.bot.Send(msg)
				return
			}
			if update.Message.Command() == "stop" {
				// h.service.DeleteUser(update.Message.Chat.ID)
				h.storage.Delete(update.Message.Chat.ID)
				return
			}
			if update.Message.Command() == "all" {
				admin := os.Getenv("ADMIN")
				if update.Message.Chat.UserName == admin {
					h.service.SendAll(update.Message.CommandArguments())
					return
				}
			}
		}

		// echo sms
	}

}
