package gomicro

import (
	"context"

	"base/broker"

	micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/server"
)

type gomicroBroker struct {
	client client.Client
	server server.Server
}

func NewGomicroBroker(client client.Client, server server.Server) *gomicroBroker {
	return &gomicroBroker{
		client: client,
		server: server,
	}
}

func (self *gomicroBroker) Publish(p broker.Message) error {
	return self.client.Publish(context.TODO(), p)
}

func (self *gomicroBroker) Subscribe(topic string, h interface{}) error {
	return micro.RegisterSubscriber(topic, self.server, h)
}
