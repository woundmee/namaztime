package handlers

import (
	"log/slog"
	"namaztimeApi/internal/domain"
	"namaztimeApi/internal/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type NamazData interface {
	NamazDataMonth(path string) ([]domain.NamazTime, error)
	NamazDataToday(day int, path string) (*domain.NamazTime, error)
}

type Handler struct {
	logger    *slog.Logger
	namazData NamazData
}

func NewHandler(logger *slog.Logger, namazData NamazData) *Handler {
	return &Handler{
		logger:    logger,
		namazData: namazData,
	}
}

// Здесь будут обработчики (handlers) http-запросов.

// PATH_SCHEDULES add in gitignore and get from getenv

// note: HTTP: GET-request
func (h *Handler) GetNamazDataHandler(c echo.Context) error {
	fp, err := utils.FullPathToMonthSchedule()
	if err != nil {
		h.logger.Error("Не удалось получить полный путь до файла", "error", err)
		c.JSON(http.StatusInternalServerError, map[string]string{"error": "Не удалось получить путь до файла расписания"})
	}

	fr, err := h.namazData.NamazDataMonth(fp)
	if err != nil {
		h.logger.Error("Ошибка получения расписания намазов за месяц", "error", err, "http", http.StatusInternalServerError)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Не удалось получить данные",
		})
	}

	h.logger.Info("Расписание за месяц получено", "http", http.StatusOK)
	return c.JSON(http.StatusOK, fr)
}

func currentDayInt() (int, error) {
	loc := time.FixedZone("UTC+7", 7*60*60)
	dayStr := time.Now().In(loc).Format("02")
	return strconv.Atoi(dayStr)
}

// note: HTTP: GET-request (schedule from /today)
func (h *Handler) GetNamazDataFilteredHandler(c echo.Context) error {
	day, err := currentDayInt()
	if err != nil {
		h.logger.Error("Не удалось получить текущий день", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Не удалось получить текущий день"})
	}

	fp, err := utils.FullPathToMonthSchedule()
	if err != nil {
		h.logger.Error("Не удалось получить полный путь до файла", "error", err)
		c.JSON(http.StatusInternalServerError, map[string]string{"error": "Не удалось получить путь до файла расписания"})
	}

	res, err := h.namazData.NamazDataToday(day, fp)
	if err != nil {
		h.logger.Error("Не удалось получить расписание намазов за текущий день", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Данные за текущий день не найдены"})
	}

	h.logger.Info("Расписание на текущий день получено", "http", http.StatusOK)
	return c.JSON(http.StatusOK, res)
}
