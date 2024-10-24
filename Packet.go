package main

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"io"
	"log"
	"net/http"
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

func ParseDataToPacketsAsync(unParsedDataChannel chan string, packetsChannel chan []Packet) {
	for {
		var data string
		select {
		case d, ok := <-unParsedDataChannel:
			if ok {
				data = d
			} else {
				log.Fatal("Failed to read from channel, this should never happen")
			}

		default:
			continue
		}

		dataSplit := strings.Split(data, "\n")
		var packets []Packet
		for _, dataEntry := range dataSplit {
			dataSplitSpace := strings.Split(dataEntry, " ")
			ipAddresses := strings.Split(dataSplitSpace[0], ",")
			size, err := strconv.Atoi(dataSplitSpace[1])
			if err != nil {
				log.Println("Cannot get the size of the packet reason: ", err.Error())
				continue
			}

			source := ipAddresses[0]
			dest := ipAddresses[1]
			packetInfo, err := GetPacketInfoFromRedis(dest)
			if err != nil {
				log.Println("Failed to get the packet info from redis reason: ", err.Error())
				continue
			}

			packet := Packet{
				Src:            source,
				Dest:           dest,
				Size:           size,
				DestPacketInfo: packetInfo,
			}

			packets = append(packets, packet)
			if len(packetsChannel) < cap(packetsChannel) {
				packetsChannel <- packets
			}
		}
	}
}

func GetPacketInfoFromRedis(ipAddress string) (PacketInfo, error) {
	redisCmd := redisClient.Get(context.Background(), ipAddress)
	if redisCmd == nil {
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
	response, err := http.Get("https://ip-api.com/json/" + ipAdrress)
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
