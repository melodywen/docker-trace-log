package app

import (
	"context"
	"fmt"
	"github.com/melodywen/docker-trace-log/app/provider"
	"github.com/melodywen/docker-trace-log/contracts"
	"github.com/melodywen/docker-trace-log/helper"
	"github.com/melodywen/docker-trace-log/package/logs"
	"github.com/spf13/viper"
	"time"
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

	a.Config = a.loadConfig()

	a.AttachAppObserver(provider.NewLogProvider())
	a.AttachAppObserver(provider.NewMongoProvider())
}

func (a *Application) loadConfig() *viper.Viper {
	viper.SetConfigName("config")   // name of config file (without extension)
	viper.SetConfigType("json")     // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("./config") // optionally look for config in the working directory
	err := viper.ReadInConfig()     // Find and read the config file
	if err != nil {                 // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	return viper.GetViper()
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

// EnterExitFunc 打印函数进出日志:
//     使用方法
// 	defer EnterExitFunc()()
func (a *Application) EnterExitFunc(ctx context.Context) func() {
	funcName, file, line := helper.GetCallerInfo(true)
	start := time.Now()

	a.Log.Debug(ctx, "enter %s func (%s:%d)", funcName, file, line)
	return func() {
		_, file, line = helper.GetCallerInfo(false)
		a.Log.Debug(ctx, "exit %s (%s) func (%s:%d)", funcName, time.Since(start), file, line)
	}
}
