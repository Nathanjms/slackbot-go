package SlackHandler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/nathanjms/slackbot-go/internal/application"
	"github.com/slack-go/slack"
)

func CommandHandler(app *application.Application) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Return text "test"

		attachment := slack.Attachment{
			Pretext: "Super Bot Message",
			Text:    "some text",
			Color:   "4af030",
			Fields: []slack.AttachmentField{
				{
					Title: "Date",
					Value: time.Now().String(),
				},
			},
		}

		_, timestamp, err := app.SlackClient.PostMessage(
			app.Config.SlackConfig.ChannelID,
			slack.MsgOptionAttachments(attachment),
		)

		if err != nil {
			return err
		}

		fmt.Printf("Message sent at %s", timestamp)

		return c.NoContent(http.StatusOK)
	}
}
