package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/joho/godotenv"
	"github.com/lmittmann/tint"
	"github.com/nathanjms/slackbot-go/internal/application"
	"github.com/nathanjms/slackbot-go/internal/env"
	"github.com/nathanjms/slackbot-go/internal/version"
)

func main() {
	fmt.Printf("version: %s\n", version.Get())

	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: slog.LevelDebug}))

	err := run(logger)

	if err != nil {
		logger.Error(err.Error(), "trace", string(debug.Stack()))

		// Spin up a basic http server for the 500 error:
		httpErr := http.ListenAndServe(fmt.Sprintf(":%d", env.GetInt("PORT", 3000)), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Internal Server Error"))
		}))
		if httpErr != nil {
			// Well, we tried our best
			log.Fatal(httpErr)
			os.Exit(1)
		}

	}
}

func run(logger *slog.Logger) error {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file, proceeding with system environment variables")
	}

	app, err := application.New(logger)
	if err != nil {
		return err
	}

	return serveHttp(app)
}
