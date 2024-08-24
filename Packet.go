package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type PacketInfo struct {
	IP       string `json:"ip"`
	Hostname string `json:"hostname"`
	City     string `json:"city"`
	Region   string `json:"region"`
	Country  string `json:"country"`
	Location string `json:"loc"`
	Org      string `json:"org"`
	Postal   string `json:"postal"`
	Timezone string `json:"timezone"`
	Readme   string `json:"readme"`
}

type Packet struct {
	Dest       string
	Src        string
	Size       int
	packetInfo PacketInfo
}

func CreateTCPPacket(packetSplit []string) Packet {
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
			packetSizeString = packetSplit[i]
			break
		}
	}

	if len(destSplit) < 4 {
		log.Println("Destination IP is mangeled")
		return Packet{
			"null",
			"null",
			-1,
			PacketInfo{},
		}
	}

	if len(srcSplit) < 4 {
		log.Println("Source ip is mangeled")
		return Packet{
			"null",
			"null",
			-1,
			PacketInfo{},
		}
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

	for i := 0; i < len(srcSplit); i++ {
		if i == len(srcSplit)-1 {
			break
		}

		src += srcSplit[i] + "."
	}

	packetSize, err := strconv.Atoi(packetSizeString)
	if err != nil {
		log.Println("Failed to convert PacketSizeString to int, CreateTCPPacket")
		return Packet{
			"null",
			"null",
			-1,
			PacketInfo{},
		}
	}

	return Packet{
		dest,
		src,
		packetSize,
		PacketInfo{},
	}
}

func CreateUDPPacket(packetSplit []string) Packet {
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

	for i := 0; i < len(srcSplit); i++ {
		if i == len(srcSplit)-1 {
			break
		}

		src += srcSplit[i]
		if i != len(srcSplit)-2 {
			src += "."
		}
	}

	packetInfo := GetIpInfo(dest)
	packetSize, err := strconv.Atoi(packetSplit[7])
	if err != nil {
		log.Println("Cannot parse PacketSize!")
		return Packet{
			"null",
			"null",
			-1,
			PacketInfo{},
		}
	}

	return Packet{
		dest,
		src,
		packetSize,
		packetInfo,
	}
}

func GetIpInfo(ipAdrress string) PacketInfo {
	response, err := http.Get("http://ipinfo.io/" + ipAdrress)
	if err != nil {
		log.Println("Response returned error", err.Error())
		return PacketInfo{}
	}

	responseAsBytes, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("Failed to ReadAll the bytes from the response", err.Error())
		return PacketInfo{}
	}

	var packetInfo PacketInfo
	json.Unmarshal(responseAsBytes, &packetInfo)
	return packetInfo
}
