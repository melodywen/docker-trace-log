package contracts

import "github.com/sirupsen/logrus"

// AppSubjectInterface 定义app被关注者Interface
type AppSubjectInterface interface {
	AttachAppObserver(observer AppObserverInterface)
	DetachAppObserver(observer AppObserverInterface)
	NotifyStartServerBeforeEvent()
	NotifyStartServerAfterEvent()
}

// AppAttributeInterface 定义app 关注者Interface
type AppAttributeInterface interface {
	SetLog(log *logrus.Logger)
	EnterExitFunc() func()
}

// AppObserverInterface 定义app 关注者Interface
type AppObserverInterface interface {
	StartServerBeforeEvent(app AppAttributeInterface) error
	StartServerAfterEvent(app AppAttributeInterface) error
}
