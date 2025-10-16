package services

import (
	"fmt"
	"log/slog"
	"os"
	"telegramBot/clients/namaznsk"
)

type NamazService struct {
	logger *slog.Logger
}

func New(logger *slog.Logger) *NamazService {
	return &NamazService{
		logger: logger,
	}
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
	url := os.Getenv("URL_TODAY_SCHEDULE")

	// data, err := today.TodaySchedule(url)
	data, err := today.TodaySchedule(url)
	if err != nil {
		msg := "ошибка получения расписания"
		ns.logger.Error(msg, "error", err)
		return err.Error()
	}

	header := "🕗 Расписание на сегодня\n\n"
	// res := fmt.Sprintf("%s\t- Фаджр\n%s\t- Восход\n%s\t- Зухр", data.Fajr, data.Sunrise, data.Zuhr)
	res := fmt.Sprintf("%s   - Фаджр\n%s   - Восход\n%s - Зухр\n%s - Аср\n%s - Магриб\n%s - Иша",
		data.Fajr, data.Sunrise, data.Zuhr, data.Asr, data.Magrib, data.Isha)

	return header + res
}

// todo: вызвать в отдельной горутине
// todo: реализовать в for{} для проверки времени каждые 30сек времени намаза.
// todo: расписание за сегодня сохранить в файле или сохранить в кеше через отдельную функцию
func (ns *NamazService) NamazTimeNotify() {
	//
}
