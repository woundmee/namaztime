package services

import (
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"
	"namaztimeApi/models"
	"os"
	"strconv"
	"time"
)

// здесь будут сервисы (бизнес-логика).

type NamazData interface {
	NamazDataMonth(path string) ([]models.NamazTime, error)
	NamazDataToday(day int, path string) (*models.NamazTime, error)
}

type Service struct {
	logger *slog.Logger
}

// конструктор
func New(log *slog.Logger) *Service {
	return &Service{
		logger: log,
	}
}

// получить расписание намазов за месяц
func (s *Service) NamazDataMonth(path string) ([]models.NamazTime, error) {
	file, err := os.Open(path)
	if err != nil {
		s.logger.Error("", "error", err)
		return nil, fmt.Errorf("не удалось открыть файл: %w", err)
	}

	defer file.Close()

	s.logger.Info("начинаю читать файл", "filename", file.Name())
	csvReader := csv.NewReader(file)

	// пропускаю заголовки
	_, err = csvReader.Read()
	if err != nil {
		s.logger.Error("ошибка чтения заголовка csv-файла", "error", err)
		return nil, fmt.Errorf("ошибка чтения заголовка csv-файла: %v", err)
	}

	var data []models.NamazTime

	for {
		record, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			s.logger.Error("ошибка чтения файла", "error", err)
			return nil, fmt.Errorf("ошибка чтения файла: %v", err)
		}

		data = append(data, models.NamazTime{
			Day:     record[0],
			Fajr:    s.timeFormatted(record[1]),
			Sunrise: s.timeFormatted(record[2]),
			Zuhr:    s.timeFormatted(record[3]),
			Asr:     s.timeFormatted(record[4]),
			Magrib:  s.timeFormatted(record[5]),
			Isha:    s.timeFormatted(record[6]),
		})
	}
	return data, nil
}

func (s *Service) timeFormatted(t string) string {
	timeCustomFormat, err := time.Parse("15:04", t)
	if err != nil {
		s.logger.Error("не удалось сконвертировать время", "error", err)
		return err.Error()
	}

	return timeCustomFormat.Format("15:04")
}

// получить расписание намазов за текущий день
func (s *Service) NamazDataToday(day int, path string) (*models.NamazTime, error) {
	s.logger.Info("получаю расписание за месяц для дальнейшей обработки...")

	data, err := s.NamazDataMonth(path)
	if err != nil {
		s.logger.Error("ошибка получения данных", "error", err)
		return nil, fmt.Errorf("ошибка получения данных: %v", err)
	}

	s.logger.Info("начинаю перебор расписания на получения расписания текущего дня")
	for _, d := range data {
		if d.Day == strconv.Itoa(day) {
			s.logger.Info("расписание на текущий день найдено", "day", d.Day)
			return &d, nil
		}
	}

	s.logger.Error("расписание на текущий день не найдено", "day", day)
	return nil, fmt.Errorf("данные за день %d не найдены", day)
}
