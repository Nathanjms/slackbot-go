package SlackHandler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/nathanjms/slackbot-go/internal/application"
	"github.com/slack-go/slack"
)

// Command types are "test", "list", "worst":
const (
	CommandTypeTest  = "test"
	CommandTypeList  = "list"
	CommandTypeWorst = "worst"
)

func CommandHandler(app *application.Application) echo.HandlerFunc {
	return func(c echo.Context) error {
		command := c.FormValue("command")
		if command != "/harvest" {
			return c.String(http.StatusOK, "Invalid Command Option!")
		}

		text := strings.TrimSpace(c.FormValue("text"))
		fmt.Println(text)

		// Handle the commands
		if text == CommandTypeTest {
			_, _, err := app.SlackClient.PostMessage(
				app.Config.SlackConfig.ChannelID,
				slack.MsgOptionText("Testing...", false),
			)

			if err != nil {
				return err
			}

			return c.NoContent(http.StatusOK)
		}

		if text == CommandTypeList {
			_, _, err := app.SlackClient.PostMessage(
				app.Config.SlackConfig.ChannelID,
				slack.MsgOptionText("List", false),
			)

			if err != nil {
				return err
			}

			return c.NoContent(http.StatusOK)
		}

		if text == CommandTypeWorst {
			_, _, err := app.SlackClient.PostMessage(
				app.Config.SlackConfig.ChannelID,
				slack.MsgOptionText("Worst", false),
			)

			if err != nil {
				return err
			}

			return c.NoContent(http.StatusOK)
		}

		return c.String(http.StatusOK, "Invalid Command Option!")
	}
}
