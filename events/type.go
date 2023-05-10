package events

type Processor interface {
	Process(e Event) error
}

type Fetcher interface {
	Fetch(limit int) ([]Event, error)
}

type Type int

const (
	Unknown Type = iota
	Message
)

type Event struct {
	Type     Type
	ChatId   int
	Username string
	Text     string
}
