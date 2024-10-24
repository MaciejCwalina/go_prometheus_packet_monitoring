package main

import (
	"github.com/redis/go-redis/v9"
)

func ParseData(validDataChannel chan []byte, packetChannel chan []Packet) {

}

func RecordMetrics(packetsChannel chan []Packet) {
	go func() {
		for {
			if len(packetsChannel) < cap(packetsChannel) {
				SetNetworkGuages(<-packetsChannel)
			}
		}
	}()
}

var (
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "192.168.20.119:6379",
		Password: "",
		DB:       0,
	})
)

func main() {

	packetChannel := make(chan Packet, 256)
	go CapturePacketsAsync(packetChannel)
	for {
		_ = <-packetChannel
	}

	//RecordMetrics(packetChannel)
	//http.Handle("/metrics", promhttp.Handler())
	//log.Fatal(http.ListenAndServe(":2115", nil))
}
