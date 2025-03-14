# slackbot-go

A basic scaffold for a slack bot written in go, using Echo as a web framework.

## Dev

Run `air` if installed, or can run `go run ./cmd/api` and manually reload after changes

## Prod

Can deploy with docker compose, or can build with `go build -o out ./cmd/api`, and then run the `out` file with `./out`.