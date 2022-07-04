package net

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/indeedhat/track-pad/internal/config"
	"github.com/indeedhat/track-pad/internal/env"
)

// BroadcastExistence of the server via udp multicast messages
func BroadcastExistence(done chan struct{}) {
	httpAddress := fmt.Sprintf(":%s", env.Get(env.ServerPort, config.HttpPort))

	addr, err := net.ResolveUDPAddr("udp4", config.MultiCastAddress)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		log.Fatal(err)
	}

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "Unknown"
	}

	interval := env.GetInt(env.DiscoveryInterval, config.DiscoveryInterval)
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	for {
		select {
		case <-ticker.C:
			conn.Write([]byte(fmt.Sprintf("trackpad.server;%s;%s;", httpAddress[1:], hostname)))
		case <-done:
			break
		}
	}
}
