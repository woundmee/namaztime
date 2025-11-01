package services

import (
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"
	"namaztimeApi/internal/cache"
	"namaztimeApi/models"
	"namaztimeApi/pkg"
	"os"
	"strconv"
	"time"
)

// здесь будут сервисы (бизнес-логика).

type NamazData interface {
	NamazDataMonth() ([]models.NamazTime, error)
	NamazDataToday(day int, path string) (*models.NamazTime, error)
}

type Service struct {
	logger *slog.Logger
	cache  *cache.Cache
}

// конструктор
func New(log *slog.Logger, cache *cache.Cache) *Service {
	return &Service{
		logger: log,
		cache:  cache,
	}
}

func (s *Service) StartDailyUpdateCache() {
	s.logger.Info("запускаю горутину ежедневного обновления кэша")
	for {
		midnight := s.calculateMidnightUtc7()
		now := time.Now()

		if !midnight.After(now) {
			midnight = midnight.Add(24 * time.Hour)
		}

		// вычисляю остаток времени до полуночи
		sleepDuration := midnight.Sub(now)
		s.logger.Info("ожидание следующей полуночи", "длительность", sleepDuration)
		time.Sleep(sleepDuration)

		s.logger.Info("обновление кэша в 00:00")

		data, err := s.NamazDataMonth()
		if err != nil {
			s.logger.Error("ошибка получения месячных данных", "error", err)
			continue
		}

		s.cache.Set(data)
		s.logger.Info("кэш успешно обновлен!", "time", time.Now().Format(time.DateTime), "ContentLength", len(data))
	}
}

// получить расписание намазов за месяц
func (s *Service) NamazDataMonth() ([]models.NamazTime, error) {

	if cacheData, ok := s.cache.Get(); ok {
		return cacheData, nil
	}

	fullpath, err := pkg.FullPathToMonthScheduleFile()
	if err != nil {
		s.logger.Error("не удалось получить полный путь до файла с расписаниями", "error", err)
		return nil, err
	}

	file, err := os.Open(fullpath)
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

	s.cache.Set(data)
	s.logger.Info("данные записаны в кэш", "ContentLength", len(data))
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

	data, err := s.NamazDataMonth()
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

// from utc+7
func (s *Service) calculateMidnightUtc7() time.Time {
	// loc := time.FixedZone("UTC+7", 7*60*60)
	now := time.Now().In(s.timeZone())
	return time.Date(
		now.Year(), now.Month(), now.Day()+1,
		0, 0, 0, 0, s.timeZone(),
	)
}

func (s *Service) timeZone() *time.Location {
	timezone := os.Getenv("TIMEZONE")
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		s.logger.Error("не удалось задать timezone", "error", err)
		// loc = time.FixedZone("UTC+7", 7*3600)
		// n.logger.Warn("timezone задана принудительно", "timezone", loc)
		// return loc
		return nil
	}

	return loc
}
