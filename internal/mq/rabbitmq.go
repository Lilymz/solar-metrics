package mq

import "C"
import (
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
	"solar-metrics/internal/metric"
	"solar-metrics/internal/model"
	"strconv"
)

type rabbitMqConsumer struct {
	connections []*amqp.Connection
	channels    []*amqp.Channel
}

// SetUp rabbitmq 消费者连接建立逻辑
func (c *rabbitMqConsumer) SetUp(url string) []error {
	var errors []error
	for range model.GetConfig().Solar.Rabbitmq.Consumer.ConnectionNums {
		conn, err := amqp.Dial(url)
		if err != nil {
			errors = append(errors, err)
			continue
		}
		c.connections = append(c.connections, conn)
		for range model.GetConfig().Solar.Rabbitmq.Consumer.ChannelNums {
			ch, channelError := conn.Channel()
			if channelError != nil {
				errors = append(errors, channelError)
				continue
			}
			c.channels = append(c.channels, ch)
		}
	}
	return errors
}
func (c *rabbitMqConsumer) GetConnections() []*amqp.Connection {
	return c.connections
}

func (c *rabbitMqConsumer) GetChannels() []*amqp.Channel {
	return c.channels
}

func (c *rabbitMqConsumer) Consumer(delivery <-chan amqp.Delivery) {
	metrics := metric.GetMetrics()
	for message := range delivery {
		power := GetMessage(message)
		metrics.IncSolarTotalPower(power)
		// 消息确认
		_ = message.Ack(false)
	}
}

// Delivery 获取消息投递
func (c *rabbitMqConsumer) Delivery(url, queue string) ([]<-chan amqp.Delivery, error) {
	var deliveries []<-chan amqp.Delivery
	errors := c.SetUp(url)
	if len(errors) > 0 {
		log.Fatal("failed get rabbitmq delivery！", errors)
		return nil, nil
	}
	qos := model.GetConfig().Solar.Rabbitmq.Consumer.Qos
	consumerTag := model.GetConfig().Solar.Rabbitmq.Consumer.Name
	for i, channel := range c.GetChannels() {
		if err := channel.Qos(qos, 0, true); err != nil {
			return nil, err
		}
		consume, err := channel.Consume(queue, consumerTag+strconv.Itoa(i+1), false, false, false, false, nil)
		if err != nil {
			return nil, err
		}
		deliveries = append(deliveries, consume)
	}
	return deliveries, nil
}

// GetMessage 将当前消息解析为指标消息
func GetMessage(delivery amqp.Delivery) (power model.EquipmentOriginalPower) {
	body := delivery.Body
	err := json.Unmarshal(body, &power)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("get message failed ")
	}
	return power
}
func (c *rabbitMqConsumer) Close() (err error) {
	for _, channel := range c.GetChannels() {
		err = channel.Close()
	}
	for _, connection := range c.GetConnections() {
		err = connection.Close()
	}
	return err
}
