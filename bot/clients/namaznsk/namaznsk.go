package namaznsk

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"telegramBot/clients/entities"
	"telegramBot/internal/cache"
	"time"
)

type Namaz struct {
	logger *slog.Logger
	cache  *cache.Cache
	url    string
}

func New(logger *slog.Logger, cache *cache.Cache, url string) *Namaz {
	return &Namaz{
		logger: logger,
		cache:  cache,
		url:    url,
	}
}

// ежедневное обновление кеша в 00:00
func (n *Namaz) StartDailyUpdateCache() {
	now := time.Now()
	midnight := n.cache.CalculateMidnightUtc7()

	if !midnight.After(now) {
		midnight = midnight.Add(24 * time.Hour)
	}

	// вычисляю остаток времени до полуночи
	sleepDuration := midnight.Sub(now)
	n.logger.Info("ожидание следующей полуночи", "длительность", sleepDuration)
	time.Sleep(sleepDuration)
	n.logger.Info("обновление кэша в 00:00")

	data, err := n.todayScheduleRead()
	if err != nil {
		n.logger.Error("ошибка обновления кеша", "error", err)
		return
	}

	// cache update
	n.cache.Set(data)
	n.logger.Info("кэш успешно обновлен!", "time", now.Format("2006-01-02 15:04"))

}

func (n *Namaz) TodaySchedule() (entities.NamazData, error) {
	rd, err := n.todayDataCache()
	if err != nil {
		return entities.NamazData{}, err
	}

	var parsedData entities.NamazData
	err = json.Unmarshal(rd, &parsedData)
	if err != nil {
		msg := "не удалось спарсить json-ответ от сервера"
		n.logger.Error(msg, "error", err)
		return entities.NamazData{}, fmt.Errorf("%s: %w", msg, err)
	}

	return parsedData, nil
}

func (n *Namaz) todayDataCache() ([]byte, error) {
	if n.cache == nil {
		msg := "кэш не инициализирован"
		n.logger.Error(msg)
		return nil, errors.New(msg)
	}

	if data, ok := n.cache.Get(); ok {
		n.logger.Info("данные найдены в кеше", "длина", len(data))
		return data, nil
	}

	todayDataByte, err := n.todayScheduleRead()
	if err != nil {
		return nil, err
	}

	n.cache.Set(todayDataByte)
	n.logger.Info("данные сохранены в кэше", "длина", len(todayDataByte))

	return todayDataByte, nil

}

func (n *Namaz) todayScheduleRead() ([]byte, error) {

	resp, err := n.todayScheduleHttp(n.url)
	if err != nil {
		return nil, err
	}

	n.logger.Info("читаю полученный ответ от сервера", "ContentLength", resp.ContentLength)
	rd, err := io.ReadAll(resp.Body)
	if err != nil {
		msg := "не удалось прочитать ответ от сервера"
		n.logger.Error(msg, "error", err)
		return nil, fmt.Errorf("%s: %w", msg, err)
	}

	return rd, nil
}

func (n *Namaz) todayScheduleHttp(url string) (*http.Response, error) {
	const fn = "clients.namaznsk.TodaySchedule"
	n.logger.Info("получаю данные", "url", url)
	resp, err := http.Get(url)
	if err != nil {
		msg := "ошибка получения данных по url"
		n.logger.Error(msg, "error", err, "func", fn)
		return nil, fmt.Errorf("%s: %s, %w", msg, url, err)
	}

	if resp.StatusCode != http.StatusOK {
		msg := "ошибка получения данных"
		n.logger.Error(msg, "status code", resp.StatusCode, "error", err)
		return nil, fmt.Errorf("%s: %w", msg, err)
	}

	n.logger.Info("данные получены!", "status code", resp.StatusCode)
	return resp, nil
}
