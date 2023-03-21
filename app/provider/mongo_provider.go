package provider

import (
	"github.com/melodywen/docker-trace-log/contracts"
)

type MongoProvider struct {
}

func NewMongoProvider() *MongoProvider {
	return &MongoProvider{}
}

func (m *MongoProvider) StartServerBeforeEvent(app contracts.AppAttributeInterface) error {
	defer app.EnterExitFunc()()
	return nil
}

func (m *MongoProvider) StartServerAfterEvent(app contracts.AppAttributeInterface) error {
	defer app.EnterExitFunc()()
	return nil
}
