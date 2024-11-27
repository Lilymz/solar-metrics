package mq

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"solar-metrics/internal/model"
)

var (
	connections []*amqp.Connection
	channels    []*amqp.Channel
	errors      []error
)

// setUp 根据传入的amqp091定义的url，建立
func setUp(url string) ([]*amqp.Connection, []*amqp.Channel, []error) {
	for range model.GetConfig().Solar.Rabbitmq.Consumer.ConnectionNums {
		conn, err := amqp.Dial(url)
		if err != nil {
			errors = append(errors, err)
			continue
		}
		connections = append(connections, conn)
		for range model.GetConfig().Solar.Rabbitmq.Consumer.ChannelNums {
			ch, channelError := conn.Channel()
			if channelError != nil {
				errors = append(errors, channelError)
			}
			channels = append(channels, ch)
		}
	}
	return connections, channels, nil
}

// GetConnections 获取连接
func GetConnections() []*amqp.Connection {
	return connections
}

// GetChannels 获取通道
func GetChannels() []*amqp.Connection {
	return connections
}
