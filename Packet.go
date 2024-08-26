package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
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
	Dest       string
	Src        string
	Size       int
	packetInfo PacketInfo
}

func CreateTCPPacket(packetSplit []string) (Packet, error) {
	dest := packetSplit[2]
	src := packetSplit[4]
	destSplit := strings.Split(dest, ".")
	srcSplit := strings.Split(src, ".")
	dest = ""
	src = ""

	packetSizeString := ""
	for i := 0; i < len(packetSplit); i++ {
		if packetSplit[i] == "length" {
			i++
			if i >= len(packetSplit) {
				break
			}

			packetSizeString = packetSplit[i]
			break
		}
	}

	if len(destSplit) < 4 {
		return Packet{}, errors.New("the destination IP is less then 4 octets")
	}

	if len(srcSplit) < 4 {
		return Packet{}, errors.New("the Source IP is less then 4 octets")
	}

	if packetSizeString == "" {
		log.Println("Mangeled Packet cannot get the size!")
	}

	for i := 0; i < len(destSplit); i++ {
		if i == len(destSplit)-1 {
			break
		}

		dest += destSplit[i] + "."
	}

	stringCmd := redisClient.Get(context.Background(), dest)
	if stringCmd == nil {
		for i := 0; i < len(srcSplit); i++ {
			if i == len(srcSplit)-1 {
				break
			}

			src += srcSplit[i] + "."
		}

		packetSize, err := strconv.Atoi(packetSizeString)
		if err != nil {
			return Packet{}, errors.New("failed to convert PacketSizeString to int, CreateTCPPacket")
		}

		packetInfo, _ := GetIpInfo(dest)
		p := Packet{
			dest,
			src,
			packetSize,
			packetInfo,
		}

		bytes, err := json.Marshal(p)
		if err != nil {
			return p, err
		}

		redisClient.Set(context.Background(), dest, bytes, 0)
		return p, nil
	}

	bytes, err := stringCmd.Bytes()
	if err != nil {
		return Packet{}, err
	}

	var packet Packet
	err = json.Unmarshal(bytes, &packet)
	if err != nil {
		return packet, err
	}

	packetSize, err := strconv.Atoi(packetSizeString)
	if err != nil {
		return Packet{}, errors.New("failed to convert PacketSizeString to int, CreateTCPPacket")
	}

	packet.Size = packetSize
	return packet, nil
}

func CreateUDPPacket(packetSplit []string) (Packet, error) {
	dest := packetSplit[2]
	src := packetSplit[4]
	destSplit := strings.Split(dest, ".")
	srcSplit := strings.Split(src, ".")
	dest = ""
	src = ""
	for i := 0; i < len(destSplit); i++ {
		if i == len(destSplit)-1 {
			break
		}

		dest += destSplit[i]
		if i != len(destSplit)-2 {
			dest += "."
		}
	}

	redisCmd := redisClient.Get(context.Background(), dest)
	if redisCmd != nil {
		for i := 0; i < len(srcSplit); i++ {
			if i == len(srcSplit)-1 {
				break
			}

			src += srcSplit[i]
			if i != len(srcSplit)-2 {
				src += "."
			}
		}

		packetInfo, err := GetIpInfo(dest)
		if err != nil {
			log.Println("Failed to get IP info due to ", err.Error())
		}

		packetSize, err := strconv.Atoi(packetSplit[7])
		if err != nil {
			log.Println("Cannot parse PacketSize!")
			return Packet{}, err
		}

		packet := Packet{
			dest,
			src,
			packetSize,
			packetInfo,
		}

		bytes, err := json.Marshal(packet)
		if err != nil {
			return Packet{}, err
		}

		redisClient.Set(context.Background(), dest, bytes, 0)
	}

	var packet Packet
	packetAsBytes, err := redisCmd.Bytes()
	if err != nil {
		return Packet{}, nil
	}

	err = json.Unmarshal(packetAsBytes, &packet)
	if err != nil {
		return Packet{}, err
	}

	packetSize, err := strconv.Atoi(packetSplit[7])
	if err != nil {
		return Packet{}, err
	}

	packet.Size = packetSize
	return packet, nil
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
