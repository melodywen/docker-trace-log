package main

import (
	app2 "github.com/melodywen/docker-trace-log/app"
	"github.com/melodywen/docker-trace-log/package/logs"
)

func main() {
	app := app2.GetApp()
	ctx := logs.GetContextOfLog()

	app.NotifyStartServerBeforeEvent(ctx)

	app.NotifyStartServerAfterEvent(ctx)
}
