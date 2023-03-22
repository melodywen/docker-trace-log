package contracts

import (
    "context"
    "github.com/melodywen/docker-trace-log/package/logs"
    "github.com/spf13/viper"
)

// AppSubjectInterface 定义app被关注者Interface
type AppSubjectInterface interface {
    AttachAppObserver(observer AppObserverInterface)
    DetachAppObserver(observer AppObserverInterface)
    NotifyStartServerBeforeEvent(ctx context.Context)
    NotifyStartServerAfterEvent(ctx context.Context)
}

// AppAttributeInterface 定义app 关注者Interface
type AppAttributeInterface interface {
    SetLog(log *logs.Log)
    GetLog() *logs.Log
    SetConfig(config *viper.Viper)
    GetConfig() *viper.Viper
}

// AppObserverInterface 定义app 关注者Interface
type AppObserverInterface interface {
    StartServerBeforeEvent(ctx context.Context, app AppAttributeInterface) error
    StartServerAfterEvent(ctx context.Context, app AppAttributeInterface) error
}
