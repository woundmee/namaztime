package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"api/services"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

// Здесь будут обработчики (handlers) http-запросов.

const PATH_SCHEDULES = "./schedules/"

// ищет файл текущего месяца в каталоге PATH_SCHEDULES
func searchCurrentMonthSchedule(path string) (string, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		log.Error("Не удалось получить список файлов", "error", err)
		return "", err
	}

	currentMonthInt := int(time.Now().Month())
	for _, file := range files {
		if strings.Contains(file.Name(), strconv.Itoa(currentMonthInt)) {
			log.Info("Файл найден ", "filename ", file.Name())
			return file.Name(), nil
		}
	}

	return "", errors.New("файл с расписанием за текущий месяц не найден")
}

// возвращает
func fullPathToMonthSchedule() string {
	filename, err := searchCurrentMonthSchedule(PATH_SCHEDULES)
	if err != nil {
		// log.Fatalf("Ошибка: %v", err)
		log.Error("ERROR: ", "message: ", err)
		os.Exit(1)
	}
	return PATH_SCHEDULES + filename
}

// note: HTTP: GET-request
func GetNamazDataHandler(cntx echo.Context) error {
	fr, err := services.NamazDataMonth(fullPathToMonthSchedule())
	fmt.Println(fr)

	if err != nil {
		return cntx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Не удалось получить данные",
		})
	}

	return cntx.JSON(http.StatusOK, fr)
}

// note: HTTP: GET-request (schedule from /today)
func GetNamazDataFilteredHandler(c echo.Context) error {
	currentDay := strings.Split(time.Now().Format("2006-01-02"), "-")[2]
	currentDayInt, err := strconv.Atoi(currentDay)
	if err != nil {
		log.Error("Не удалось сконвертировать тип в int", currentDay, err)
		return err
	}

	res, err := services.NamazDataToday(currentDayInt, fullPathToMonthSchedule())
	if err != nil {
		// log.Fatalf("Не удалось получить данные: %v", err)
		log.Error("Не удалось получить данные", "error", err)
		return err
	}

	return c.JSON(http.StatusOK, res)
}
