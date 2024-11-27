package main

import (
	"flag"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"solar-metrics/internal/consumer"
	"solar-metrics/internal/metric"
	"solar-metrics/internal/model"
	"strings"
	"syscall"
)

var (
	addr = flag.String("address", "127.0.0.1:9999", "The address to listen on for HTTP requests.")
)

func init() {
	logConfig := model.GetConfig().Solar.Logrus
	if logConfig.Format == "text" {
		log.SetFormatter(&log.TextFormatter{})
	} else if logConfig.Format == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		log.SetFormatter(&log.TextFormatter{})
	}
	if logConfig.RecordMethod {
		log.SetReportCaller(true)
	} else {
		log.SetReportCaller(false)
	}
	level, err := log.ParseLevel(strings.ToLower(logConfig.Level))
	if err != nil {
		level = log.InfoLevel
	}
	log.SetLevel(level)
	log.Info("solar-metrics完成日志配置初始化!")
	// todo rabbitmq初始化

}
func main() {
	flag.Parse()
	signals := make(chan os.Signal, 1)
	//  Prometheus采集端子
	go func() {
		reg := metric.Registry()
		http.Handle("/metrics", promhttp.HandlerFor(
			reg,
			promhttp.HandlerOpts{
				EnableOpenMetrics: true,
			},
		))
		log.Info("Start collecting solar metrics into Prometheus!")
		log.Fatal(http.ListenAndServe(*addr, nil))
	}()
	// 消息消费
	consumer.Start()
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals
	consumer.Stop()
}
