package contracts

// AppSubjectInterface 定义app被关注者Interface
type AppSubjectInterface interface {
	AttachAppObserver(observer AppObserverInterface)
	DetachAppObserver(observer AppObserverInterface)
	NotifyStartServerBeforeEvent()
	NotifyStartServerAfterEvent()
}

// AppObserverInterface 定义app 关注者Interface
type AppObserverInterface interface {
	StartServerBeforeEvent() error
	StartServerAfterEvent() error
}
