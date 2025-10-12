package pkg

import (
	"strconv"
	"time"
)

// получить текущий день с локалью +7ч в Int
func CurrentDayUtc7Int() (int, error) {
	loc := time.FixedZone("UTC+7", 7*60*60)
	dayStr := time.Now().In(loc).Format("02")
	return strconv.Atoi(dayStr)
}
