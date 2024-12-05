package mq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	RABBITMQ = "rabbitmq"
	KAFKA    = "kafka"
	REDIS    = "redis"
)

type Consumer interface {
	// SetUp 消费者连接建立
	SetUp(url string) []error
	// GetConnections 获取所有已建立的连接
	GetConnections() []*amqp.Connection
	// GetChannels 获取所有正在交互的通道
	GetChannels() []*amqp.Channel
	// Consumer 消费逻辑
	Consumer(delivery <-chan amqp.Delivery)
	// Delivery 投递者
	Delivery(url, queue string) ([]<-chan amqp.Delivery, error)
	// Close 结束消费
	Close() error
}

// NewConsumer 创建消费者
func NewConsumer(brokerName string) Consumer {
	if brokerName == KAFKA {
		// todo
	} else if brokerName == REDIS {
		// todo
	} else if brokerName == RABBITMQ {
		return &rabbitMqConsumer{}
	}
	return &rabbitMqConsumer{}
}
