package app

import (
	"fmt"
	"github.com/melodywen/docker-trace-log/app/provider"
	"github.com/melodywen/docker-trace-log/contracts"
	"github.com/melodywen/docker-trace-log/helper"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

var app *Application

type Application struct {
	isInit bool
	Log    *logrus.Logger
	Config *viper.Viper

	observers []contracts.AppObserverInterface
}

func newApplication() *Application {
	return &Application{}
}

func GetApp() *Application {
	if app == nil {
		app = newApplication()

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

func (a *Application) NotifyStartServerBeforeEvent() {
	for _, observerInterface := range app.observers {
		err := observerInterface.StartServerBeforeEvent(a)
		if err != nil {
			a.Log.Fatalf("app notify start server before event ,found err:%s", err.Error())
		}
	}
}

func (a *Application) NotifyStartServerAfterEvent() {
	for _, observerInterface := range app.observers {
		err := observerInterface.StartServerAfterEvent(a)
		if err != nil {
			a.Log.Fatalf("app notify start server after event ,found err:%s", err.Error())
		}
	}
}

func (a *Application) SetLog(log *logrus.Logger) {
	a.Log = log
}

// EnterExitFunc 打印函数进出日志:
//     使用方法
// 	defer EnterExitFunc()()
func (a *Application) EnterExitFunc() func() {
	funcName, file, line := helper.GetCallerInfo(true)
	start := time.Now()

	a.Log.Debugf("enter %s func (%s:%d)", funcName, file, line)
	return func() {
		_, file, line = helper.GetCallerInfo(false)
		a.Log.Debugf("exit %s (%s) func (%s:%d)", funcName, time.Since(start), file, line)
	}
}
