package handlers

import (
	"fmt"
	"log/slog"
	"telegramBot/clients/namaznsk"
	"telegramBot/internal/services"

	// "telegramBot/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	logger  *slog.Logger
	bot     tgbotapi.BotAPI
	namaz   namaznsk.Namaz
	service services.NamazService
}

func New(logger *slog.Logger, botKey tgbotapi.BotAPI, namaz namaznsk.Namaz, service services.NamazService) *Handler {
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
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Бот запущен!")
				msg2 := tgbotapi.NewMessage(update.Message.Chat.ID, "Для вызова справки используйте команду /help")
				h.bot.Send(msg)
				h.bot.Send(msg2)
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
			if update.Message.Command() == "notify" {
				text := "🔔 Вы подписались на уведомления о времени намаза!"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
				h.bot.Send(msg)

				// запускаю нотификатор для пользователя
				go h.service.StartNamazNotifier(update.Message)
				return
				// todo: сделать без команды /notify, а сразу по умолчанию после старта бота!
			}
		}

		// echo sms
	}
}
