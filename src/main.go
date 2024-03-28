package main

import (
	"github.com/RubenPari/tracksByPopularity/src/utils"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func main() {
	server := echo.New()

	// Set up session
	server.Use(session.Middleware(sessions.NewCookieStore([]byte(utils.RandomString(64)))))

	// Run server
	server.Logger.Fatal(server.Start(":8080"))
}
