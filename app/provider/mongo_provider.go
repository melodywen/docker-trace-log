package provider

type MongoProvider struct {
}

func (m *MongoProvider) StartServerBeforeEvent() error {

	return nil
}

func (m *MongoProvider) StartServerAfterEvent() error {

	return nil
}
