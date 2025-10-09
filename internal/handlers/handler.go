package handlers

import (
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"namaztimeApi/internal/configs/slogger"
	"namaztimeApi/internal/services"

	"github.com/labstack/echo/v4"
)

// Здесь будут обработчики (handlers) http-запросов.

// PATH_SCHEDULES add in gitignore and get from getenv
const PATH_SCHEDULES = "./schedules/"

// ищет файл текущего месяца в каталоге PATH_SCHEDULES
func searchCurrentMonthSchedule(path string) (string, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		slogger.Log.Error("Не удалось получить список файлов", "error", err)
		return "", err
	}

	currentMonthInt := int(time.Now().Month())
	for _, file := range files {
		if strings.Contains(file.Name(), strconv.Itoa(currentMonthInt)) {
			slogger.Log.Info("Файл найден", "filename", file.Name())
			return file.Name(), nil
		}
	}

	return "", errors.New("файл с расписанием за текущий месяц не найден")
}

// возвращает полный путь до файла
func fullPathToMonthSchedule() (string, error) {
	filename, err := searchCurrentMonthSchedule(PATH_SCHEDULES)
	if err != nil {
		slogger.Log.Error("Файл не найден", "error", err)
		return "", err
	}
	return PATH_SCHEDULES + filename, nil
}

// note: HTTP: GET-request
func GetNamazDataHandler(c echo.Context) error {
	fp, err := fullPathToMonthSchedule()
	if err != nil {
		slogger.Log.Error("Не удалось получить полный путь до файла", "error", err)
	}

	fr, err := services.NamazDataMonth(fp)
	if err != nil {
		slogger.Log.Error("Ошибка получения расписания намазов за месяц", "error", err, "http", http.StatusInternalServerError)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Не удалось получить данные",
		})
	}

	slogger.Log.Info("Расписание за месяц получено", "http", http.StatusOK)
	return c.JSON(http.StatusOK, fr)
}

func currentDayInt() (int, error) {
	loc := time.FixedZone("UTC+7", 7*60*60)
	dayStr := time.Now().In(loc).Format("02")
	return strconv.Atoi(dayStr)
}

// note: HTTP: GET-request (schedule from /today)
func GetNamazDataFilteredHandler(c echo.Context) error {

	day, err := currentDayInt()
	if err != nil {
		slogger.Log.Error("Не удалось получить текущий день", "error", err)
		return err
	}

	fp, err := fullPathToMonthSchedule()
	if err != nil {
		slogger.Log.Error("Не удалось получить полный путь до файла", "error", err)
	}

	res, err := services.NamazDataToday(day, fp)
	if err != nil {
		slogger.Log.Error("Не удалось получить расписание намазов за текущий день", "error", err)
		return err
	}

	slogger.Log.Info("Расписание на текущий день получено", "http", http.StatusOK)
	return c.JSON(http.StatusOK, res)
}
