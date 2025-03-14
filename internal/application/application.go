package application

import (
	"log/slog"

	"github.com/nathanjms/slackbot-go/internal/env"
	"github.com/slack-go/slack"
)

type Config struct {
	BaseURL     string
	HTTPPort    int
	Env         string
	SlackConfig struct {
		ChannelID     string
		Token         string
		SigningSecret string
	}
}

type Application struct {
	Config      Config
	Logger      *slog.Logger
	SlackClient *slack.Client
}

func New(logger *slog.Logger) (*Application, error) {
	app := &Application{}
	// --- Config ---
	cfg := initConfig()

	app.Config = cfg
	app.Logger = logger

	// --- Slack ---

	app.SlackClient = slack.New(cfg.SlackConfig.Token, slack.OptionDebug(cfg.Env == "development"))

	return app, nil
}

func initConfig() Config {
	var cfg Config

	cfg.Env = env.GetString("ENV", "production")
	cfg.BaseURL = env.GetString("BASE_URL", "http://localhost")
	cfg.HTTPPort = env.GetInt("PORT", 3000)
	cfg.SlackConfig.ChannelID = env.GetString("SLACK_CHANNEL_ID", "")
	cfg.SlackConfig.Token = env.GetString("SLACK_AUTH_TOKEN", "")
	cfg.SlackConfig.SigningSecret = env.GetString("SLACK_SIGNING_SECRET", "")

	return cfg
}

// reportError reports the error to Sentry and logs it
func (app *Application) ReportError(err error) {
	// 1. Log the error
	app.Logger.Error(err.Error())
}
