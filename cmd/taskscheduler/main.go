package main

import (
	"os"

	app "github.com/futig/task-scheduler/internal/app"
)

func main() {
	app.Run(os.Getenv("BOT_TOKEN"))
}