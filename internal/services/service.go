package services

import (
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"
	"namaztimeApi/internal/domain"
	"os"
	"strconv"
)

type Service struct {
	logger *slog.Logger
}

func NewService(logger *slog.Logger) *Service {
	return &Service{
		logger: logger,
	}
}

// получить расписание намазов за месяц
func (s *Service) NamazDataMonth(path string) ([]domain.NamazTime, error) {
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

	var data []domain.NamazTime

	for {
		record, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			s.logger.Error("ошибка чтения файла", "error", err)
			return nil, fmt.Errorf("ошибка чтения файла: %v", err)
		}

		nt := domain.NamazTime{
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
func (s *Service) NamazDataToday(day int, path string) (*domain.NamazTime, error) {
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
