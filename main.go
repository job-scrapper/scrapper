package main

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/job-scrapper/scrapper/scrapper"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Scrapper!\n")
	})

	e.GET("/alba", func(c echo.Context) error {
		job := c.QueryParam("job")
		area := c.QueryParam("area")
		volume, err :=  strconv.Atoi(c.QueryParam("volume"))
		if err != nil {
			echo.NewHTTPError(http.StatusBadRequest, err.Error)
		}

		response, err := scrapper.GetAlbaPages(job, area, volume)
		if err != nil {
			echo.NewHTTPError(http.StatusBadRequest, err.Error)
		}
		return c.JSON(http.StatusOK, response)
	})

	e.Logger.Fatal(e.Start(":2222"))
}
