package app

import (
	"fmt"
	"github.com/melodywen/docker-trace-log/contracts"
	"github.com/melodywen/docker-trace-log/helper"
	"github.com/spf13/viper"
	"log"
)

var app *Application

type Application struct {
	isInit bool
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
	defer helper.EnterExitFunc()
	for _, observerInterface := range app.observers {
		err := observerInterface.StartServerBeforeEvent()
		if err != nil {
			log.Fatalf("app notify start server before event ,found err:%s", err.Error())
		}
	}
}

func (a *Application) NotifyStartServerAfterEvent() {
	for _, observerInterface := range app.observers {
		err := observerInterface.StartServerAfterEvent()
		if err != nil {
			log.Fatalf("app notify start server after event ,found err:%s", err.Error())
		}
	}
}
