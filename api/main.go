package main

import (
	"api/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/namaztime/month", handlers.GetNamazDataHandler)
	e.GET("/namaztime/today", handlers.GetNamazDataFilteredHandler)

	e.Start("localhost:8080")

}
