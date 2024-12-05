package consumer

import (
	"github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
	"solar-metrics/internal/model"
	"solar-metrics/internal/mq"
	"sync"
	"time"
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
	url := model.GetConfig().Solar.Rabbitmq.Dsl
	consumer = mq.NewConsumer(url)
	queues := model.GetConfig().Solar.Rabbitmq.Consumer.Queues
	var wg sync.WaitGroup
	for _, queue := range queues {
		wg.Add(1)
		go func(q string) {
			defer wg.Done()
			startQueueConsumer(url, q)
		}(queue)
	}
	wg.Wait()
}
func startQueueConsumer(url, queue string) {
	for {
		deliveries, err := consumer.Delivery(url, queue)
		if err != nil {
			log.WithFields(log.Fields{
				"queue": queue,
				"error": err,
			}).Error("creat mq consumer error")
			time.Sleep(5 * time.Second)
			continue
		}
		// 处理单个队列的消费
		var innerWg sync.WaitGroup
		for _, delivery := range deliveries {
			innerWg.Add(1)
			go func(d <-chan amqp091.Delivery) {
				defer innerWg.Done()
				consumer.Consumer(delivery)
			}(delivery)
		}
		innerWg.Wait()
	}
}
