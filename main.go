package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
)

func ParseData(validDataChannel chan []byte, packetChannel chan []Packet) {
	go func() {
		for {
			start := time.Now()
			packetInfo := string(<-validDataChannel)
			packetInfoSplit := strings.Split(packetInfo, "\n")
			packets := []Packet{}
			for _, packet := range packetInfoSplit {
				packetSplit := strings.Split(packet, " ")
				lengthPacketSplit := len(packetSplit)
				if lengthPacketSplit < 8 {
					continue
				}

				var packet Packet
				var err error
				if packetSplit[5] == "UDP," {
					packet, err = CreateUDPPacket(packetSplit)
				} else {
					packet, err = CreateTCPPacket(packetSplit)
				}

				if err != nil {
					continue
				}

				packets = append(packets, packet)
			}

			log.Println(time.Since(start))
			packetChannel <- packets
		}
	}()
}

func RecordMetrics(packetsChannel chan []Packet) {
	go func() {
		for {
			SetNetworkGuages(<-packetsChannel)
		}
	}()
}

var (
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
)

func main() {
	httpServer := HttpServer{
		"tcp",
		":5145",
	}

	listener := httpServer.StartHttpServer()
	conn := httpServer.HandleIncomingConnections(listener)
	validDataChannel := make(chan []byte, 128)
	packetChannel := make(chan []Packet, 128)
	httpServer.ReadAllBytesFromClient(validDataChannel, conn)
	ParseData(validDataChannel, packetChannel)
	RecordMetrics(packetChannel)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":2115", nil))
}
