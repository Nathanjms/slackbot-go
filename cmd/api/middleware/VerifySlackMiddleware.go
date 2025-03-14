package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/nathanjms/slackbot-go/internal/application"
)

const MaxPermittedRequestAge time.Duration = 100 * time.Second

// Heavily inspired from https://github.com/coro/verifyslack
func VerifySlackMiddleware(app *application.Application) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var timestamp string
			if timestamp = c.Request().Header.Get("X-Slack-Request-Timestamp"); timestamp == "" {
				app.ReportError(fmt.Errorf("request did not contain a request timestamp"))
				return c.String(http.StatusForbidden, "request did not contain a request timestamp")
			}

			intTimestamp, err := strconv.ParseInt(timestamp, 10, 64)
			if err != nil {
				app.ReportError(fmt.Errorf("request did not contain a request timestamp"))
				return c.String(http.StatusForbidden, "request did not contain a request timestamp")
			}

			timeNow := time.Now()
			if timeNow.After(time.Unix(intTimestamp, 0).Add(MaxPermittedRequestAge)) {
				app.ReportError(fmt.Errorf("request did not contain a request timestamp"))
				return c.String(http.StatusForbidden, "request did not contain a request timestamp")
			}

			var slackSignature string
			if slackSignature = c.Request().Header.Get("X-Slack-Signature"); slackSignature == "" {
				app.ReportError(fmt.Errorf("request does not provide a Slack-signed signature"))
				return c.String(http.StatusForbidden, "request does not provide a Slack-signed signature")
			}

			requestBody, err := io.ReadAll(c.Request().Body)
			if err != nil {
				return err
			}

			// Rewrite the request body back
			c.Request().Body = io.NopCloser(bytes.NewBuffer(requestBody))

			expectedSignature := GenerateExpectedSignature(timestamp, requestBody, app.Config.SlackConfig.SigningSecret)

			if !hmac.Equal([]byte(expectedSignature), []byte(slackSignature)) {
				app.ReportError(fmt.Errorf("request is not signed with a valid Slack signature"))
				return c.String(http.StatusForbidden, "request is not signed with a valid Slack signature")
			}

			return next(c)
		}
	}
}

func GenerateExpectedSignature(timestamp string, requestBody []byte, signingSecret string) string {
	baseSignature := append([]byte(fmt.Sprintf("v0:%s:", timestamp)), requestBody...)
	mac := hmac.New(sha256.New, []byte(signingSecret))
	mac.Write(baseSignature)

	expectedSignature := fmt.Sprintf("v0=%s", hex.EncodeToString(mac.Sum(nil)))
	return expectedSignature
}
