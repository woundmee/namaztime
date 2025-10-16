package namaznsk

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"telegramBot/clients/entities"
	"telegramBot/internal/cache"
)

type Namaz struct {
	logger *slog.Logger
	cache  *cache.Cache
}

func New(logger *slog.Logger, cache *cache.Cache) *Namaz {
	return &Namaz{
		logger: logger,
		cache:  cache,
	}
}

func (n *Namaz) TodaySchedule(url string) (entities.NamazData, error) {
	rd, err := n.todayDataCache(url)
	if err != nil {
		return entities.NamazData{}, err
	}

	var jsonData entities.NamazData
	err = json.Unmarshal(rd, &jsonData)
	if err != nil {
		msg := "не удалось спарсить json-ответ от сервера"
		n.logger.Error(msg, "error", err)
		return entities.NamazData{}, fmt.Errorf("%s: %w", msg, err)
	}

	return jsonData, nil
}

func (n *Namaz) todayDataCache(url string) ([]byte, error) {
	if data, ok := n.cache.Get(); ok {
		n.logger.Info("данные найдены в кеше", "длина", len(data))
		return data, nil
	}

	todayDataByte, err := n.todayScheduleRead(url)
	if err != nil {
		return nil, err
	}

	n.cache.Set(todayDataByte)
	n.logger.Info("данные сохранены в кэше", "длина", len(todayDataByte))

	return todayDataByte, nil

}

func (n *Namaz) todayScheduleRead(url string) ([]byte, error) {

	resp, err := n.todayScheduleHttp(url)
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
