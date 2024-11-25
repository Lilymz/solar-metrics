package main

import (
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

var addr = flag.String("address", ":8080", "The address to listen on for HTTP requests.")

func main() {
	flag.Parse()
	reg := prometheus.NewRegistry()
	http.Handle("/metrics", promhttp.HandlerFor(
		reg,
		promhttp.HandlerOpts{
			EnableOpenMetrics: true,
		},
	))
	fmt.Println("Start collecting solar metrics into Prometheus!")
	log.Fatal(http.ListenAndServe(*addr, nil))
}
