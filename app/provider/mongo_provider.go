package provider

import (
    "context"
    "github.com/melodywen/docker-trace-log/contracts"
)

type MongoProvider struct {
}

func NewMongoProvider() *MongoProvider {
    return &MongoProvider{}
}

func (m *MongoProvider) StartServerBeforeEvent(ctx context.Context, app contracts.AppAttributeInterface) error {

    return nil
}

func (m *MongoProvider) StartServerAfterEvent(ctx context.Context, app contracts.AppAttributeInterface) error {

    return nil
}
