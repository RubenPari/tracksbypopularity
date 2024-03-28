package main

import "github.com/labstack/echo/v4"

func main() {
	server := echo.New()

	server.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello World!")
	})

	server.Logger.Fatal(server.Start(":8080"))
}
