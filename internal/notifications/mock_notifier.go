package notifications

type MockNotifier struct {
	Message string
}

func (m *MockNotifier) Send(message string) error {
	m.Message = message
	return nil
}
