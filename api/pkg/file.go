package pkg

import (
	"errors"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"
)

// ищет файл текущего месяца в каталоге PATH_SCHEDULES
func searchCurrentMonthScheduleFile(path string) (string, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		slog.Error("Не удалось получить список файлов", "error", err)
		return "", err
	}

	currentMonthInt := int(time.Now().Month())
	for _, file := range files {
		if strings.Contains(file.Name(), strconv.Itoa(currentMonthInt)) {
			slog.Info("файл найден", "filename", file.Name())
			return file.Name(), nil
		}
	}

	return "", errors.New("файл с расписанием за текущий месяц не найден")
}

// возвращает полный путь до файла
func FullPathToMonthScheduleFile() (string, error) {
	path := os.Getenv("PATH_SCHEDULES")
	filename, err := searchCurrentMonthScheduleFile(path)
	if err != nil {
		slog.Error("Файл не найден", "error", err)
		return "", err
	}
	return path + filename, nil
}
