package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func ParseData(validDataChannel chan []byte, packetChannel chan []Packet) {
	go func() {
		for {
			packetInfo := string(<-validDataChannel)
			packetInfoSplit := strings.Split(packetInfo, "\n")
			packets := []Packet{}
			for _, packet := range packetInfoSplit {
				packetSplit := strings.Split(packet, " ")
				lengthPacketSplit := len(packetSplit)
				if lengthPacketSplit < 8 {
					continue
				}

				packet := Packet{}
				if packetSplit[5] == "UDP," {
					packet = CreateUDPPacket(packetSplit)
				} else {
					packet = CreateTCPPacket(packetSplit)
				}

				if packet.Size == -1 {
					continue
				}

				packets = append(packets, packet)
			}

			packetChannel <- packets
		}
	}()
}
func RecordMetrics(packetsChannel chan []Packet) {
	for {
		SetNetworkGuages(<-packetsChannel)
	}
}

func main() {
	httpServer := HttpServer{
		"tcp",
		":5145",
	}

	listener := httpServer.StartHttpServer()
	conn := httpServer.HandleIncomingConnections(listener)
	validDataChannel := make(chan []byte)
	packetChannel := make(chan []Packet)
	httpServer.ReadAllBytesFromClient(validDataChannel, conn)
	ParseData(validDataChannel, packetChannel)
	RecordMetrics(packetChannel)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":2115", nil))
}
