package consumer

import (
	log "github.com/sirupsen/logrus"
	"solar-metrics/internal/model"
	"solar-metrics/internal/mq"
)

var consumer mq.Consumer

func Start() {
	doStart()
}
func Stop() {
	defer func(consumer mq.Consumer) {
		err := consumer.Close()
		if err != nil {
			log.Error(err)
		}
	}(consumer)
}
func doStart() {
	dsl := model.GetConfig().Solar.Rabbitmq.Dsl
	consumer = mq.NewConsumer(dsl)
	for _, queue := range model.GetConfig().Solar.Rabbitmq.Consumer.Queues {
		deliveries, err := consumer.Delivery(dsl, queue)
		if err != nil {
			log.WithFields(log.Fields{
				"queue": queue,
				"error": err,
			}).Error("creat mq consumer error")
		}
		for _, delivery := range deliveries {
			go consumer.Consumer(delivery)
		}
	}
}
