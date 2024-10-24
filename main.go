package main

import (
	"log"
	"net/http"

	"github.com/coreos/go-iptables/iptables"
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

func BlockIpAddress(packet Packet) error {
	shouldBlock := false
	log.Println(packet.DestPacketInfo.CountryCode)
	for _, countryCode := range configInstance.ListOfBlockedCountryCodes {
		if countryCode == packet.DestPacketInfo.CountryCode {
			shouldBlock = true
			break
		}
	}

	if !shouldBlock {
		return nil
	}

	err := ipTables.Append("filter", "INPUT", "-s", packet.Dest, "-j", "DROP")
	if err != nil {
		return err
	}

	err = ipTables.Append("filter", "OUTPUT", "-d", packet.Dest, "-j", "DROP")
	if err != nil {
		return err
	}

	return nil
}

var (
	redisClient *redis.Client
	ipTables    *iptables.IPTables
)

func main() {
	ReadConfigFromFile()
	redisClient = redis.NewClient(&redis.Options{
		Addr:     configInstance.RedisDataBaseAddress,
		Password: "",
		DB:       0,
	})

	packetChannel := make(chan Packet, 256)
	var err error

	ipTables, err = iptables.New(iptables.IPFamily(iptables.ProtocolIPv4), iptables.Timeout(0))
	if err != nil {
		log.Fatal(err.Error())
	}

	go CapturePacketsAsync(packetChannel)
	go RecordMetricsAsync(packetChannel)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":2115", nil))
}
