package consumer

import (
	"github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
	"solar-metrics/internal/metric"
	"solar-metrics/internal/model"
	"solar-metrics/internal/mq"
)

// Start 开启消费
func Start() {
	doStart()
}

// Stop 结束消费
func Stop() {
	// todo 未确认数据回写mq,且不在生产任务
}

func doStart() {
	// 创建连接
	url := model.GetConfig().Solar.Rabbitmq.Dsl
	for _, queue := range model.GetConfig().Solar.Rabbitmq.Consumer.Queues {
		deliveries, err := mq.Delivery(url, queue)
		if err != nil {
			log.WithFields(log.Fields{
				"queue": queue,
				"error": err,
			}).Error("创建mq的消费者出现错误")
		}
		for _, delivery := range deliveries {
			go handler(delivery)
		}
	}
}
func handler(delivery <-chan amqp091.Delivery) {
	metrics := metric.GetMetrics()
	for message := range delivery {
		power := mq.GetMessage(message)
		metrics.IncSolarTotalPower(power)
		// 消息确认
		_ = message.Ack(false)
	}
}
