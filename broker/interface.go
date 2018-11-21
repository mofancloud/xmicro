package broker

// Publication is the interface for a message published asynchronously
type Message interface {
	Topic() string
	Payload() interface{}
	Message() interface{}
	ContentType() string
}

type Subscriber interface {
	Topic() string
	Subscriber() interface{}
}

type Broker interface {
	Publish(p Message) error
	Subscribe(topic string, h interface{}) error
}

const ContentTypeJson = "application/json"

type Handler func(Message) error

type defaultPublication struct {
	topic       string
	message     interface{}
	contentType string
}

func NewDefaultPublication(topic string, message interface{}, contentType string) *defaultPublication {
	return &defaultPublication{
		topic:       topic,
		message:     message,
		contentType: contentType,
	}
}

func (self *defaultPublication) ContentType() string {
	return self.contentType
}

func (self *defaultPublication) Topic() string {
	return self.topic
}

func (self *defaultPublication) Message() interface{} {
	return self.message
}

func (self *defaultPublication) Payload() interface{} {
	return self.message
}
