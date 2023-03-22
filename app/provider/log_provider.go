package provider

import (
	"context"
	"github.com/melodywen/docker-trace-log/contracts"
	"github.com/melodywen/docker-trace-log/package/logs"
)

type LogProvider struct {
}

func NewLogProvider() *LogProvider {
	return &LogProvider{}
}

func (l *LogProvider) StartServerBeforeEvent(ctx context.Context, app contracts.AppAttributeInterface) error {
	instance := logs.NewLog()
	app.SetLog(instance)

	instance.Debug(ctx, "成功注册了logrus")

	return nil
}

func (l *LogProvider) StartServerAfterEvent(ctx context.Context, app contracts.AppAttributeInterface) error {
	defer app.GetLog().EnterExitFunc(ctx)()
	return nil
}
