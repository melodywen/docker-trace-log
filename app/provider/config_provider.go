package provider

import (
	"context"
	"fmt"
	"github.com/melodywen/docker-trace-log/contracts"
	"github.com/spf13/viper"
)

type ConfigProvider struct {
}

func NewConfigProvider() *ConfigProvider {
	return &ConfigProvider{}
}

func (c ConfigProvider) StartServerBeforeEvent(ctx context.Context, app contracts.AppAttributeInterface) error {
	defer app.GetLog().EnterExitFunc(ctx)()

	viper.SetConfigName("config")   // name of config file (without extension)
	viper.SetConfigType("json")     // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("./config") // optionally look for config in the working directory
	err := viper.ReadInConfig()     // Find and read the config file
	if err != nil {                 // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	instance := viper.GetViper()
	app.SetConfig(instance)
	return nil
}

func (c ConfigProvider) StartServerAfterEvent(ctx context.Context, app contracts.AppAttributeInterface) error {
	defer app.GetLog().EnterExitFunc(ctx)()
	return nil
}
