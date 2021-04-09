package pymetadata

//扩展grpc的metadata，支持map[string]interface{}的存储

const (
	// Network
	RemoteIP = "remote_ip"
)

type defaultPyMetaDataKey struct{}

type PyMetaData map[string]interface{}

func NewPyMetaData(m map[string]interface{}) PyMetaData {
	if m == nil {
		return PyMetaData{}
	}
	nmd := new(PyMetaData)
	for k, v := range m {
		nmd[k] = v
	}
	return nmd
}

func (m *PyMetaData) Len() int {
	return len(m)
}

func (m *PyMetaData) Add(key string, value interface{}) {
	m[key] = value
	return
}

func (m *PyMetaData) Del(key string) interface{} {
	if _, exists := m[key]; exists {
		v := m[key]
		delete(m, key)
		return v
	}
	return nil
}
