package consumer

import (
	"github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
	"solar-metrics/internal/metric"
	"solar-metrics/internal/model"
	"solar-metrics/internal/mq"
	"sync"
	"time"
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
	url := model.GetConfig().Solar.Rabbitmq.Dsl
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
		deliveries, err := mq.Delivery(url, queue)
		if err != nil {
			log.WithFields(log.Fields{
				"queue": queue,
				"error": err,
			}).Error("创建mq消费者出错")
			time.Sleep(5 * time.Second)
			continue
		}

		// 处理单个队列的消费
		var innerWg sync.WaitGroup
		for _, delivery := range deliveries {
			innerWg.Add(1)
			go func(d <-chan amqp091.Delivery) {
				defer innerWg.Done()
				handler(d)
			}(delivery)
		}

		innerWg.Wait()
	}
}
func handler(delivery <-chan amqp091.Delivery) {
	metrics := metric.GetMetrics()
	for message := range delivery {
		power := mq.GetMessage(message)
		metrics.IncSolarTotalPower(power)
		// 消息确认go env
		_ = message.Ack(false)
	}
}
