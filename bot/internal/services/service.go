package services

import (
	"fmt"
	"log/slog"
	"strings"
	"telegramBot/clients/namaznsk"
	storage "telegramBot/internal/storage/sqlite"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Service struct {
	logger   *slog.Logger
	namaznsk *namaznsk.Namaz
	bot      *tgbotapi.BotAPI
	storage  *storage.Database
}

func New(logger *slog.Logger, namaznsk *namaznsk.Namaz, bot *tgbotapi.BotAPI, storage *storage.Database) *Service {
	return &Service{
		logger:   logger,
		namaznsk: namaznsk,
		bot:      bot,
		storage:  storage,
	}
}

// // fixme: метод для установки namaznsk клиента. Нужна ли?!
// func (ns *NamazService) SetNamazClient(namazClient *namaznsk.Namaz) {
// 	ns.namaznsk = namazClient
// }

// func (ns *NamazService) AddUser(chatID int64, username string) {
// 	ns.storage.Insert(chatID, username)
// }

// func (ns *NamazService) DeleteUser(chatID int64) {
// 	ns.storage.Delete(chatID)
// }

func (ns *Service) SendAll(message string) {

	users, err := ns.storage.GetUsers()
	if err != nil {
		ns.logger.Error("не удалось получить список пользователей", "error", err)
		return
	}

	for chatID, username := range users {
		msg := tgbotapi.NewMessage(chatID, message)
		ns.bot.Send(msg)
		ns.logger.Info("пользователею отправлено сообщение", "username", username, "chatID", chatID, "message", message)
	}
}

func (ns *Service) StartNamazNotifier() {
	ns.logger.Info("запускаю нотификатор времени намазов")

	for {
		_, name, isExistData := ns.IsNamazTime()

		if isExistData {
			users, err := ns.storage.GetUsers()
			if err != nil {
				ns.logger.Error("не удалось получить список пользователей", "error", err)
				return
			}

			msgText := fmt.Sprintf("%s — время намаза", strings.ToLower(name))

			if name == "Восход" {
				msgText = "время восхода"
			}

			for userChatID, userName := range users {
				msg := tgbotapi.NewMessage(userChatID, msgText)
				_, err := ns.bot.Send(msg)
				if err != nil {
					if strings.Contains(err.Error(), "Forbidden: bot was blocked by the user") {
						ns.logger.Info("пользователь заблокировал бота", "username", userName, "chatID", userChatID, "error", err)
						ns.storage.DeleteUser(userChatID)
						ns.logger.Info("пользователь удален из БД", "username", userName, "chatID", userChatID)
						continue
					}
					ns.logger.Error("ошибка отправки уведомления", "chatID", userChatID, "username", userName, "error", err)
				} else {
					ns.logger.Info("уведомление о наступлении времени намаза отправлено!", "namaz", name, "chatID", userChatID, "username", userName)
				}
			}

			// пауза после отправки времени намаза
			time.Sleep(time.Second * 60)
			// time.Sleep(time.Minute * 20)

		}

		time.Sleep(time.Second * 50)
	}
}

func (ns *Service) IsNamazTime() (string, string, bool) {
	const fn = "services.service.IsNamazTime"

	now := time.Now().Format("15:04")
	todayData, err := ns.namaznsk.TodaySchedule()
	if err != nil {
		ns.logger.Error("не удалось получить расписание за текущий день", "error", err, "fn", fn)
		return "", "", false
	}

	namazTimes := map[string]string{
		todayData.Fajr:    "Фаджр",
		todayData.Sunrise: "Восход",
		todayData.Zuhr:    "Зухр",
		todayData.Asr:     "Аср",
		todayData.Magrib:  "Магриб",
		todayData.Isha:    "Иша",
	}

	if name, ok := namazTimes[now]; ok {
		return now, name, true
	}

	return "", "", false
}

func (ns *Service) CommandNotify(chatID int64, username string) string {
	text := "🔔 Вы подписались на уведомления о времени намаза!"
	ns.storage.AddUser(chatID, username)
	ns.logger.Info("в БД добавлен новый пользователь", "chatID", chatID, "username", username)
	return text
}

func (ns *Service) CommandHelp() string {
	msg := `Ассаляму аляйкум!
Я бот для получения расписания намазов по г. Норильск.
	
	Что я умею?
	/help — получить справку
	/today — расписание намазов за сегодня
	
	добавить откуда беру расписание и что всё доступно на сайте. Также и распечатать можно на сайте`

	return msg
}

func (ns *Service) CommandToday(today namaznsk.Namaz) string {

	// data, err := today.TodaySchedule(url)
	data, err := today.TodaySchedule()
	if err != nil {
		msg := "ошибка получения расписания"
		ns.logger.Error(msg, "error", err)
		return err.Error()
	}

	header := "🌙 День " + data.Day + "\n" +
		"🕌 Норильск\n\n"
	// res := fmt.Sprintf("%s\t- Фаджр\n%s\t- Восход\n%s\t- Зухр", data.Fajr, data.Sunrise, data.Zuhr)
	res := fmt.Sprintf("%s   - Фаджр\n%s   - Восход\n%s - Зухр\n%s - Аср\n%s - Магриб\n%s - Иша",
		data.Fajr, data.Sunrise, data.Zuhr, data.Asr, data.Magrib, data.Isha)

	return header + res
}
