package mock

type InfoWriter struct {
	WriteFunc func(p []byte) (n int, err error)
	FlushFunc func() (string, error)
}

func (m *InfoWriter) Write(p []byte) (n int, err error) {
	return m.WriteFunc(p)
}

func (m *InfoWriter) Flush() (string, error) {
	return m.FlushFunc()
}

