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

func (s *Service) SendAll(message string) {

	users, err := s.storage.GetUsers()
	if err != nil {
		s.logger.Error("не удалось получить список пользователей", "error", err)
		return
	}

	for chatID, username := range users {
		msg := tgbotapi.NewMessage(chatID, message)
		s.bot.Send(msg)
		s.logger.Info("пользователею отправлено сообщение", "username", username, "chatID", chatID, "message", message)
	}
}

func (s *Service) StartNamazNotifier() {
	s.logger.Info("запускаю нотификатор времени намазов")

	for {
		_, name, isExistData := s.isNamazTime()

		// пятничные уведомления
		if int(time.Now().Weekday()) == int(time.Friday) {

			now := time.Now().Format("15:04")

			sunnahText, isBefore2HourZuhr := s.friday(now)
			duaText, isBefore1HourMagrib := s.beforeMagrib1Hour(now)

			// сунны и адабы пятницы
			if isBefore2HourZuhr {
				users, err := s.storage.GetUsers()
				if err != nil {
					s.logger.Error("не удалось получить список пользователей", "error", err)
					return
				}
				for userChatID := range users {
					msg := tgbotapi.NewMessage(userChatID, sunnahText)
					s.bot.Send(msg)
				}
			}

			// последний час асра
			if isBefore1HourMagrib {
				users, err := s.storage.GetUsers()
				if err != nil {
					s.logger.Error("не удалось получить список пользователей", "error", err)
					return
				}

				for chatID := range users {
					msg := tgbotapi.NewMessage(chatID, duaText)
					msg.ParseMode = "HTML"
					_, err = s.bot.Send(msg)
					if err != nil {
						s.logger.Error("ошибка отправки смс в telegram", "error", err)
					}
				}
			}
		}

		if isExistData {
			users, err := s.storage.GetUsers()
			if err != nil {
				s.logger.Error("не удалось получить список пользователей", "error", err)
				return
			}

			msgText := fmt.Sprintf("%s — время намаза", strings.ToLower(name))

			if name == "Восход" {
				msgText = "время восхода"
			}

			for userChatID, userName := range users {
				msg := tgbotapi.NewMessage(userChatID, msgText)
				_, err := s.bot.Send(msg)
				if err != nil {
					if strings.Contains(err.Error(), "Forbidden: bot was blocked by the user") {
						s.logger.Info("пользователь заблокировал бота", "username", userName, "chatID", userChatID, "error", err)
						s.storage.DeleteUser(userChatID)
						s.logger.Info("пользователь удален из БД", "username", userName, "chatID", userChatID)
						continue
					}
					s.logger.Error("ошибка отправки уведомления", "chatID", userChatID, "username", userName, "error", err)
				} else {
					s.logger.Info("уведомление о наступлении времени намаза отправлено!", "namaz", name, "chatID", userChatID, "username", userName)
				}
			}

			// пауза после отправки времени намаза
			time.Sleep(time.Second * 60)
			// time.Sleep(time.Minute * 20)

		}

		time.Sleep(time.Second * 50)
	}
}

func (s *Service) isNamazTime() (string, string, bool) {
	now := time.Now().Format("15:04")
	namazTimes, err := s.namazDataMap()
	if err != nil {
		s.logger.Error("не удалось получить расписание за текущий день", "error", err)
		return "", "", false
	}

	if name, ok := namazTimes[now]; ok {
		return now, name, true
	}

	return "", "", false
}

func (s *Service) friday(now string) (string, bool) {
	isFriday := int(time.Now().Weekday()) == int(time.Friday)

	namazData, err := s.namazDataMap()
	if err != nil {
		s.logger.Error("не удалось получить время намазов")
		return "", false
	}

	var timeZuhr string

	for t, n := range namazData {
		if n == "Зухр" {
			timeZuhr = t
			break
		}
	}

	timeZuhrParsed, err := time.Parse("15:04", timeZuhr)
	if err != nil {
		s.logger.Error("не удалось спарсить время", "error", err)
		return "", false
	}

	before2HoursZuhr := timeZuhrParsed.Add(-2 * time.Hour)

	text := `🌙 Благословенная пятница!
	
	Сунны и адабы:
	١. полное омовение (гусль)
	٢. использование благовоний
	٣. благословлять Пророка ﷺ
	٤. использование сивака
	٥. облачение в красивую одежду
	٦. чтение суры аль-Кахф
	٧. рано отправиться в мечеть
	٨. мольба (ду'а)

	اللَّهُمَّ صَلِّ عَلَى مُحَمَّدٍ وَعَلَى آلِ مُحَمَّدٍ	
	«О Аллах, благослови Мухаммада и семейство Мухаммада»`

	if isFriday && now == before2HoursZuhr.Format("15:04") {
		s.logger.Info("уведомление о суннах пятницы отправлено", "time", before2HoursZuhr.Format("15:04"))
		return text, true
	}

	return "", false
}

func (s *Service) beforeMagrib1Hour(now string) (duaText string, isTime bool) {
	namazTimes, err := s.namazDataMap()
	if err != nil {
		s.logger.Error("не удалось получить время намазов", "error", err)
		return duaText, isTime
	}

	var timeMagrib string

	for t, n := range namazTimes {
		if n == "Магриб" {
			timeMagrib = t
			break
		}
	}

	timeMagribParsed, err := time.Parse("15:04", timeMagrib)
	if err != nil {
		s.logger.Error("не удалось спарсить время намаза Магриб", "error", err)
		return duaText, true
	}

	before1HourMagrib := timeMagribParsed.Add(-1 * time.Hour)

	duaText = "🕘 Последний час Асра!\n\n" +
		"Передается от посланника Аллаха, да благословит его Аллах и приветствует, что он сказал: " +
		"<b>«В пятнице двенадцать часов,</b> (и среди них есть такой час), <b>что если мусульманин " +
		"попросит в нём Аллаха о чём-либо, то ему будет даровано это. Так ищите же его в последний час времени Асра»</b>.\n\n" +
		"<i>см. ан-Насаи 1387, Абу Дауд 1048, Фатх аль-Бари 2/420 и др.</i>"

	if now == before1HourMagrib.Format("15:04") {
		s.logger.Info("уведомление о последнем часе асра отправлено", "time", now)
		return duaText, true
	}

	return duaText, false
}

func (s *Service) namazDataMap() (map[string]string, error) {
	todayData, err := s.namaznsk.TodaySchedule()
	if err != nil {
		return nil, err
	}

	namazTimes := map[string]string{
		todayData.Fajr:    "Фаджр",
		todayData.Sunrise: "Восход",
		todayData.Zuhr:    "Зухр",
		todayData.Asr:     "Аср",
		todayData.Magrib:  "Магриб",
		todayData.Isha:    "Иша",
	}

	return namazTimes, nil
}

func (s *Service) CommandStart(chatID int64, username string) string {
	text := "🚀 Бот запущен!\n\n" +
		"🔔 Вы будете получать уведомления при наступлении времени намаза. " +
		"Для получения справочной информации используйте команду /help"

	s.storage.AddUser(chatID, username)
	return text
}

func (s *Service) CommandUnsubscribe(chatID int64) string {
	text := "🔕 Вы отписались от всех уведомлений, однако можете дальше пользоваться ботом в ручном режиме, " +
		"вызывая команды по необходимости. Список доступных команд можно получить командой /help"

	s.storage.DeleteUser(chatID)
	return text
}

func (s *Service) CommandHelp() string {
	msg := `Ассаляму аляйкум!
Я бот для получения расписания намазов по г. Норильск.
	
	Что я умею?
	/help — получить справку
	/today — расписание намазов за сегодня
	
	добавить откуда беру расписание и что всё доступно на сайте. Также и распечатать можно на сайте`

	return msg
}

func (s *Service) CommandToday(today namaznsk.Namaz) string {

	// data, err := today.TodaySchedule(url)
	data, err := today.TodaySchedule()
	if err != nil {
		msg := "ошибка получения расписания"
		s.logger.Error(msg, "error", err)
		return err.Error()
	}

	icon := s.monthIcon(time.Now().Month())

	header := icon + " День " + data.Day + "\n" +
		"🕌 Норильск\n\n"
	// res := fmt.Sprintf("%s\t- Фаджр\n%s\t- Восход\n%s\t- Зухр", data.Fajr, data.Sunrise, data.Zuhr)
	res := fmt.Sprintf("%s   - Фаджр\n%s   - Восход\n%s - Зухр\n%s - Аср\n%s - Магриб\n%s - Иша",
		data.Fajr, data.Sunrise, data.Zuhr, data.Asr, data.Magrib, data.Isha)

	return header + res
}

func (s *Service) monthIcon(month time.Month) string {
	switch month {
	case time.May, time.June, time.July, time.August, time.September:
		return "☀"
	default:
		return "❄"
	}
}
