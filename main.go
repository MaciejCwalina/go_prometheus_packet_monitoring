package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
)

func RecordMetricsAsync(packetChannel chan Packet) {
	for {
		if len(packetChannel) < cap(packetChannel) {
			SetNetworkGuages(<-packetChannel)
		}
	}
}

var redisClient *redis.Client

func main() {
	ReadConfigFromFile()
	redisClient = redis.NewClient(&redis.Options{
		Addr:     configInstance.RedisDataBaseAddress,
		Password: "",
		DB:       0,
	})

	packetChannel := make(chan Packet, 256)
	go CapturePacketsAsync(packetChannel)
	go RecordMetricsAsync(packetChannel)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":2115", nil))
}
