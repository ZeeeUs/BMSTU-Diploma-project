package redis

type mockClient struct{}

func (m *mockClient) MGet(keys []string) (value []interface{}, err error) {
	_ = keys
	return
}

func (m *mockClient) ReadGroup(readProperty string) (messages []map[string]interface{}, err error) {
	_ = readProperty
	return
}

func (m *mockClient) GetValue(key string) (value interface{}, err error) {
	_ = key
	return
}

func (m *mockClient) Ack(id string) (err error) {
	_ = id
	return
}

func (m *mockClient) DeleteKeyValue(keys ...string) (err error) {
	_ = keys
	return
}

func (m *mockClient) IsNotFound(err error) bool {
	_ = err
	return false
}

func (m *mockClient) WriteKeyValues(pairs ...interface{}) (err error) {
	_ = pairs
	return
}

func (m *mockClient) WriteGroup(data map[string]interface{}) (err error) {
	_ = data
	return
}
func (m *mockClient) Set(key string, value interface{}) (err error) {
	_, _ = key, value
	return
}

func NewMockRedisClient() Client {
	return &mockClient{}
}
