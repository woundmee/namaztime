package services

import (
	"encoding/csv"
	"fmt"
	"log"
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
		log.Fatalf("Не удалось открыть файл: %v", err)
		return []NamazTime{}, err
	}

	defer file.Close()

	csvReader := csv.NewReader(file)

	// пропускаю заголовки
	_, err = csvReader.Read()
	if err != nil {
		fmt.Printf("Ошибка чтения заголовка: %v", err)
		return []NamazTime{}, err
	}

	var data []NamazTime

	for {
		record, err := csvReader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			log.Fatalf("Ошибка чтения: %v", err)
			return []NamazTime{}, err
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
	data, err := NamazDataMonth(path)
	if err != nil {
		log.Fatalf("Ошибка получения данных: %v", err)
	}

	for _, d := range data {
		if d.Day == strconv.Itoa(day) {
			nt := NamazTime{
				Day:     d.Day,
				Fajr:    d.Fajr,
				Sunrise: d.Sunrise,
				Zuhr:    d.Zuhr,
				Asr:     d.Asr,
				Magrib:  d.Magrib,
				Isha:    d.Isha,
			}

			return nt, nil
		}
	}

	return NamazTime{}, err
}
