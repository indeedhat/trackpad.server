package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

const addr = "localhost:8881"

var upgrader = websocket.Upgrader{}

func main() {
	http.HandleFunc("/ws", func(rw http.ResponseWriter, r *http.Request) {
		con, err := upgrader.Upgrade(rw, r, nil)
		if err != nil {
			return
		}
		defer con.Close()

		for {
			_, message, err := con.ReadMessage()
			if err != nil {
				break
			}

			log.Println(string(message))
		}
	})

	http.ListenAndServe(addr, nil)
}
