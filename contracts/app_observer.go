package contracts

import (
	"context"
	"github.com/melodywen/docker-trace-log/package/logs"
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
}

// AppObserverInterface 定义app 关注者Interface
type AppObserverInterface interface {
	StartServerBeforeEvent(ctx context.Context, app AppAttributeInterface) error
	StartServerAfterEvent(ctx context.Context, app AppAttributeInterface) error
}
