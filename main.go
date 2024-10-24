package main

import (
	"go_prometheus_packet_monitoring/TSharkWrapper"

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
		Addr:     "apollo-3-prot.local:6379",
		Password: "",
		DB:       0,
	})
)

func main() {
	tshark := TSharkWrapper.NewTShark()
	tshark.Run()
	unParsedDataChannel := make(chan string, 256)
	packetChannel := make(chan []Packet, 256)
	go tshark.RedirectOutputToChannelAsync(unParsedDataChannel)
	go ParseDataToPacketsAsync(unParsedDataChannel, packetChannel)
	//RecordMetrics(packetChannel)
	//http.Handle("/metrics", promhttp.Handler())
	//log.Fatal(http.ListenAndServe(":2115", nil))
}
