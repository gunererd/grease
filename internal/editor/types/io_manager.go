package types

type Sink interface {
	Write(data []byte) error
	Flush() error
	Close() error
}

type Source interface {
	Read() ([]byte, error)
	Close() error
}

type IOManager interface {
	SetSource(source Source) error
	SetSink(sink Sink) error
}
