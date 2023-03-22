package app

import (
	"context"
	"github.com/melodywen/docker-trace-log/app/provider"
	"github.com/melodywen/docker-trace-log/contracts"
	"github.com/melodywen/docker-trace-log/package/logs"
	"github.com/spf13/viper"
)

var app *Application

type Application struct {
	isInit bool
	Log    *logs.Log
	Config *viper.Viper

	observers []contracts.AppObserverInterface
}

func newApplication() *Application {
	return &Application{}
}

func GetApp() *Application {
	if app == nil {
		app = newApplication()
		app.Init()
	}
	return app
}

func (a *Application) Init() {
	if app.isInit {
		return
	}
	a.AttachAppObserver(provider.NewLogProvider())
	a.AttachAppObserver(provider.NewConfigProvider())
	a.AttachAppObserver(provider.NewMongoProvider())
}

func (a *Application) AttachAppObserver(observer contracts.AppObserverInterface) {
	app.observers = append(app.observers, observer)
}

func (a *Application) DetachAppObserver(observer contracts.AppObserverInterface) {
	for index, observerInterface := range app.observers {
		if observerInterface == observer {
			app.observers = append(app.observers[:index], app.observers[index+1:]...)
		}
	}
}

func (a *Application) NotifyStartServerBeforeEvent(ctx context.Context) {
	for _, observerInterface := range app.observers {
		err := observerInterface.StartServerBeforeEvent(ctx, a)
		if err != nil {
			a.Log.Fatal(ctx, "app notify start server before event ,found err:%s", err.Error())
		}
	}
}

func (a *Application) NotifyStartServerAfterEvent(ctx context.Context) {
	for _, observerInterface := range app.observers {
		err := observerInterface.StartServerAfterEvent(ctx, a)
		if err != nil {
			a.Log.Fatal(ctx, "app notify start server after event ,found err:%s", err.Error())
		}
	}
}

func (a *Application) SetLog(log *logs.Log) {
	a.Log = log
}
func (a *Application) GetLog() *logs.Log {
	return a.Log
}

func (a *Application) SetConfig(config *viper.Viper) {
	a.Config = config
}

func (a *Application) GetConfig() *viper.Viper {
	return a.Config
}
