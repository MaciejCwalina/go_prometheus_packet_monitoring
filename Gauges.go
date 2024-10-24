package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	serverNetworkInformation = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "proxmox_server_network_information",
			Help: "Sends the network information of the server",
		},

		[]string{"src", "dst"},
	)
)

func SetNetworkGuages(packet Packet) {
	serverNetworkInformation.WithLabelValues(packet.Src, packet.Dest).Add(float64(packet.Size))
}
