package utils

import (
	"errors"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"
)

const PATH_SCHEDULES = "./schedules/"

// ищет файл текущего месяца в каталоге PATH_SCHEDULES
func SearchCurrentMonthSchedule(path string) (string, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		slog.Error("Не удалось получить список файлов", "error", err)
		return "", err
	}

	currentMonthInt := int(time.Now().Month())
	for _, file := range files {
		if strings.Contains(file.Name(), strconv.Itoa(currentMonthInt)) {
			slog.Info("Файл найден", "filename", file.Name())
			return file.Name(), nil
		}
	}

	return "", errors.New("файл с расписанием за текущий месяц не найден")
}

// возвращает полный путь до файла
func FullPathToMonthSchedule() (string, error) {
	filename, err := SearchCurrentMonthSchedule(PATH_SCHEDULES)
	if err != nil {
		slog.Error("Файл не найден", "error", err)
		return "", err
	}
	return PATH_SCHEDULES + filename, nil
}
