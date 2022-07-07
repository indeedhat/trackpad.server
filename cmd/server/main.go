package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/indeedhat/track-pad/internal/config"
	"github.com/indeedhat/track-pad/internal/env"
	"github.com/indeedhat/track-pad/internal/net"
	"github.com/micmonay/keybd_event"
)

var kb *keybd_event.KeyBonding

func main() {
	verbose := flag.Bool("v", false, "Verbose")
	flag.Parse()

	v := func(format string, args ...any) {
		if !*verbose {
			return
		}

		log.Printf(format, args...)
	}

	env.Load()
	if _kb, err := keybd_event.NewKeyBonding(); err == nil {
		kb = &_kb
	}

	done := make(chan struct{})
	go net.BroadcastExistence(done)

	http.HandleFunc("/ws", net.WebsocketHandler(kb, v))

	serverAddress := fmt.Sprintf(":%s", env.Get(env.ServerPort, config.HttpPort))
	http.ListenAndServe(serverAddress, nil)
}
