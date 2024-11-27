package mq

import (
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
	"solar-metrics/internal/model"
	"strconv"
)

// Delivery 获取消息投递
func Delivery(url, queue string) ([]<-chan amqp.Delivery, error) {
	var deliveries []<-chan amqp.Delivery
	_, channels, errors := setUp(url)
	if len(errors) > 0 {
		log.Fatal("获取消息投递者失败！", errors)
		return nil, nil
	}
	qos := model.GetConfig().Solar.Rabbitmq.Consumer.Qos
	consumerTag := model.GetConfig().Solar.Rabbitmq.Consumer.Name
	for i, channel := range channels {
		if err := channel.Qos(qos, 0, true); err != nil {
			return nil, err
		}
		consume, err := channel.Consume(queue, consumerTag+strconv.Itoa((i+1)), false, false, false, false, nil)
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
		}).Error("获取mq消息失败")
	}
	return power
}
