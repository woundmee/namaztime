package handlers

import (
	"log/slog"
	"net/http"

	"namaztimeApi/internal/services"
	"namaztimeApi/pkg"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	logger    *slog.Logger
	namazData services.NamazData
}

// конструктор
func New(log *slog.Logger, namazData services.NamazData) *Handler {
	return &Handler{
		logger:    log,
		namazData: namazData,
	}
}

func sendHttpJsonResponse(c echo.Context, httpStatus int, key, value string) error {
	return c.JSON(httpStatus, map[string]string{
		key: value,
	})
}

// note: HTTP: GET-request
func (h *Handler) GetNamazDataHandler(c echo.Context) error {
	fp, err := pkg.FullPathToMonthScheduleFile()
	if err != nil {
		h.logger.Error("Не удалось получить полный путь до файла", "error", err)
		// return c.JSON(http.StatusInternalServerError, map[string]string{
		// 	"error": "Не удалось получить полный путь до файла",
		// })
		sendHttpJsonResponse(c, http.StatusInternalServerError, "error", "Не удалось получить полный путь до файла")
	}

	fr, err := h.namazData.NamazDataMonth(fp)
	if err != nil {
		h.logger.Error("Ошибка получения расписания намазов за месяц", "error", err, "http", http.StatusInternalServerError)
		// return c.JSON(http.StatusInternalServerError, map[string]string{
		// 	"error": "Ошибка получения расписания намазов за месяц",
		// })
		sendHttpJsonResponse(c, http.StatusInternalServerError, "error", "Ошибка получения расписания намазов за месяц")
	}

	h.logger.Info("Расписание за месяц получено", "http", http.StatusOK)
	return c.JSON(http.StatusOK, fr)
}

// note: HTTP: GET-request (schedule from /today)
func (h *Handler) GetNamazDataFilteredHandler(c echo.Context) error {

	day, err := pkg.CurrentDayUtc7Int()
	if err != nil {
		h.logger.Error("Не удалось получить текущий день", "error", err)
		// return c.JSON(http.StatusInternalServerError, map[string]string{
		// 	"error": "Не удалось получить текущий день",
		// })
		sendHttpJsonResponse(c, http.StatusInternalServerError, "error", "Не удалось получить текущий день")
	}

	fp, err := pkg.FullPathToMonthScheduleFile()
	if err != nil {
		h.logger.Error("Не удалось получить полный путь до файла", "error", err)
		// return c.JSON(http.StatusInternalServerError, map[string]string{
		// 	"error": "Не удалось получить полный путь до файлаь",
		// })
		sendHttpJsonResponse(c, http.StatusInternalServerError, "error", "Не удалось получить полный путь до файла")
	}

	res, err := h.namazData.NamazDataToday(day, fp)
	if err != nil {
		h.logger.Error("Не удалось получить расписание намазов за текущий день", "error", err)
		// return c.JSON(http.StatusInternalServerError, map[string]string{
		// 	"error": "Не удалось получить расписание намазов за текущий день",
		// })
		sendHttpJsonResponse(c, http.StatusInternalServerError, "error", "Не удалось получить расписание намазов за текущий день")
	}

	h.logger.Info("Расписание на текущий день получено", "http", http.StatusOK)
	return c.JSON(http.StatusOK, res)
}
