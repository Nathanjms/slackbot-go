package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nathanjms/slackbot-go/cmd/api/handlers/SlackHandler"
	"github.com/nathanjms/slackbot-go/internal/application"
)

func InitRoutes(e *echo.Echo, app *application.Application) {
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": true,
			"message": "Hello, World!",
		})
	})
	e.GET("status", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": true,
			"message": "OK",
		})
	})

	// --- SLACK ROUTES ---
	authed := e.Group("")
	// TODO: ADD SLACK MIDDLEWARE
	// authed.Use(middleware.JWTAuthMiddleware(app))

	// // User Routes
	authed.POST("slack/harvest", SlackHandler.CommandHandler(app))
	// authed.DELETE("user", UserHandler.DeleteAccountHandler(app))

}
