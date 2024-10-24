package main

import (
	"context"
	"encoding/json"
	"time"

	"io"
	"log"
	"net/http"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/redis/go-redis/v9"
)

type PacketInfo struct {
	Query         string  `json:"query"`
	Status        string  `json:"status"`
	Continent     string  `json:"continent"`
	ContinentCode string  `json:"continentCode"`
	Country       string  `json:"country"`
	CountryCode   string  `json:"countryCode"`
	Region        string  `json:"region"`
	RegionName    string  `json:"regionName"`
	City          string  `json:"city"`
	District      string  `json:"district"`
	Zip           string  `json:"zip"`
	Lat           float64 `json:"lat"`
	Lon           float64 `json:"lon"`
	Timezone      string  `json:"timezone"`
	Offset        int     `json:"offset"`
	Currency      string  `json:"currency"`
	ISP           string  `json:"isp"`
	Org           string  `json:"org"`
	AS            string  `json:"as"`
	ASName        string  `json:"asname"`
	Mobile        bool    `json:"mobile"`
	Proxy         bool    `json:"proxy"`
	Hosting       bool    `json:"hosting"`
}

type Packet struct {
	Dest           string
	Src            string
	Size           int
	DestPacketInfo PacketInfo
}

func CapturePacketsAsync(packetChannel chan Packet) {
	handle, err := pcap.OpenLive(configInstance.InterfaceName, 68500, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}

	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		start := time.Now()
		ipLayer := packet.Layer(layers.LayerTypeIPv4)
		if ipLayer == nil {
			continue
		}

		ip, _ := ipLayer.(*layers.IPv4)
		packetInfo, err := GetPacketInfoFromRedis(ip.DstIP.String())

		if err != nil {
			log.Println(err.Error())
			continue
		}

		packetChannel <- Packet{
			Src:            ip.SrcIP.String(),
			Dest:           ip.DstIP.String(),
			Size:           packet.Metadata().Length,
			DestPacketInfo: packetInfo,
		}

		log.Println(time.Since(start))
	}
}

func GetPacketInfoFromRedis(ipAddress string) (PacketInfo, error) {
	redisCmd := redisClient.Get(context.Background(), ipAddress)
	if redisCmd.Err() == redis.Nil {
		packetInfo, err := GetIpInfo(ipAddress)
		if err != nil {
			return PacketInfo{}, err
		}

		bytes, err := json.Marshal(packetInfo)

		if err != nil {
			return PacketInfo{}, err
		}

		redisClient.Set(context.Background(), ipAddress, bytes, 0)
		return packetInfo, nil
	}

	bytes, err := redisCmd.Bytes()
	if err != nil {
		return PacketInfo{}, err
	}

	var packetInfo PacketInfo
	err = json.Unmarshal(bytes, &packetInfo)
	if err != nil {
		return PacketInfo{}, err
	}

	return packetInfo, nil
}

func GetIpInfo(ipAdrress string) (PacketInfo, error) {
	var packetInfo PacketInfo
	response, err := http.Get("http://ip-api.com/json/" + ipAdrress)
	if err != nil {
		log.Println("Response returned error", err.Error())
		return packetInfo, err
	}

	responseAsBytes, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("Failed to ReadAll the bytes from the response", err.Error())
		return packetInfo, err
	}

	json.Unmarshal(responseAsBytes, &packetInfo)
	return packetInfo, nil
}
