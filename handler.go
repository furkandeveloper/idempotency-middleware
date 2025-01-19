package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func ExampleHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"message": "This is an idempotent response",
	})
}

func ExampleHandlerWithoutIdempotency(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"message": "This is a non-idempotent response",
	})
}
