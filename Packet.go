package main

import (
	"encoding/json"

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
	Dest       string
	Src        string
	Size       int
	packetInfo PacketInfo
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
