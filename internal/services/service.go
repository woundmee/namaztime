package services

import (
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strconv"
)

// здесь будут сервисы (бизнес-логика).
type NamazTime struct {
	Day     string
	Fajr    string
	Sunrise string
	Zuhr    string
	Asr     string
	Magrib  string
	Isha    string
}

// получить расписание намазов за месяц
func NamazDataMonth(path string) ([]NamazTime, error) {
	file, err := os.Open(path)
	if err != nil {
		slog.Error("", "error", err)
		return nil, fmt.Errorf("не удалось открыть файл: %w", err)
	}

	defer file.Close()

	slog.Info("начинаю читать файл", "filename", file.Name())
	csvReader := csv.NewReader(file)

	// пропускаю заголовки
	_, err = csvReader.Read()
	if err != nil {
		slog.Error("ошибка чтения заголовка csv-файла", "error", err)
		return nil, fmt.Errorf("ошибка чтения заголовка csv-файла: %v", err)
	}

	var data []NamazTime

	for {
		record, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			slog.Error("ошибка чтения файла", "error", err)
			return nil, fmt.Errorf("ошибка чтения файла: %v", err)
		}

		nt := NamazTime{
			Day:     record[0],
			Fajr:    record[1],
			Sunrise: record[2],
			Zuhr:    record[3],
			Asr:     record[4],
			Magrib:  record[5],
			Isha:    record[6],
		}

		data = append(data, nt)
	}
	return data, nil
}

// получить расписание намазов за текущий день
func NamazDataToday(day int, path string) (NamazTime, error) {
	slog.Info("получаю расписание за месяц для дальнейшей обработки...")
	data, err := NamazDataMonth(path)
	if err != nil {
		slog.Error("ошибка получения данных", "error", err)
		return NamazTime{}, fmt.Errorf("ошибка получения данных: %v", err)
	}

	slog.Info("начинаю перебор расписания на получения расписания текущего дня")
	for _, d := range data {
		if d.Day == strconv.Itoa(day) {
			slog.Info("расписание на текущий день найдено", "day", d.Day)
			d = NamazTime{
				Day:     d.Day,
				Fajr:    d.Fajr,
				Sunrise: d.Sunrise,
				Zuhr:    d.Zuhr,
				Asr:     d.Asr,
				Magrib:  d.Magrib,
				Isha:    d.Isha,
			}

			return d, nil
		}
	}

	slog.Error("расписание на текущий день не найдено", "day", day)
	return NamazTime{}, fmt.Errorf("данные за день %d не найдены", day)
}
