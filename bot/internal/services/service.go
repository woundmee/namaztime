package services

import (
	"fmt"
	"log/slog"
	"telegramBot/clients/namaznsk"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type NamazService struct {
	logger   *slog.Logger
	namaznsk *namaznsk.Namaz
	bot      *tgbotapi.BotAPI
}

func New(logger *slog.Logger, bot *tgbotapi.BotAPI) *NamazService {
	return &NamazService{
		logger: logger,
		bot:    bot,
	}
}

// Добавьте метод для установки namaznsk клиента
func (ns *NamazService) SetNamazClient(namazClient *namaznsk.Namaz) {
	ns.namaznsk = namazClient
}

func (ns *NamazService) CommandHelp() string {
	msg := `Ассаляму аляйкум!
Я бот для получения расписания намазов по г. Норильск.
	
	Что я умею?
	/help — получить справку
	/today — расписание намазов за сегодня
	
	добавить откуда беру расписание и что всё доступно на сайте. Также и распечатать можно на сайте`

	return msg
}

func (ns *NamazService) CommandToday(today namaznsk.Namaz) string {

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

func (ns *NamazService) StartNamazNotifier(botID int64) {
	for {
		namazTime, name, isExistData := ns.NamazNotify()
		if isExistData {
			msgText := fmt.Sprintf("%s - %s - время намаза", name, namazTime)
			msg := tgbotapi.NewMessage(botID, msgText)
			ns.bot.Send(msg)
			ns.logger.Info("Уведомление о наступлении времени намаза отправлено!", "botID", botID, "message", msgText, "msg", msg)
		}

		time.Sleep(time.Second * 25)
	}
}

func (ns *NamazService) NamazNotify() (string, string, bool) {
	const fn = "services.service.NamazNotify"

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
