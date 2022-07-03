package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-vgo/robotgo"
	"github.com/gorilla/websocket"
	"github.com/micmonay/keybd_event"
)

const addr = ":8881"

var upgrader = websocket.Upgrader{}
var kb *keybd_event.KeyBonding

func main() {
	if _kb, err := keybd_event.NewKeyBonding(); err == nil {
		kb = &_kb
	}

	http.HandleFunc("/ws", websocketHandler)

	http.ListenAndServe(addr, nil)
}

func websocketHandler(rw http.ResponseWriter, r *http.Request) {
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

		cmdParts := strings.Split(string(message), ",")
		if len(cmdParts) == 0 {
			continue
		}

		switch cmdParts[0] {
		case "move":
			if len(cmdParts) != 3 {
				break
			}
			x, _ := strconv.ParseFloat(cmdParts[1], 32)
			y, _ := strconv.ParseFloat(cmdParts[2], 32)

			robotgo.MoveRelative(int(x), int(y))

		case "scroll":
			if len(cmdParts) != 3 {
				break
			}
			x, _ := strconv.ParseFloat(cmdParts[1], 32)
			y, _ := strconv.ParseFloat(cmdParts[2], 32)
			if x > y {
				y = 0
				x = 1
			} else {
				x = 0
				y = 1
			}

			// this continues infinately
			// robotgo.ScrollRelative(int(x), int(y))

		case "click":
			if len(cmdParts) != 2 {
				break
			}
			robotgo.Click(cmdParts[1])

		case "keeb":
			if len(cmdParts) != 2 {
				break
			}

			code, err := strconv.Atoi(cmdParts[1])

			if err == nil {
				if code == 2408 {
					kb.SetKeys(keybd_event.VK_BACKSPACE)
					_ = kb.Launching()
				} else {
					robotgo.UnicodeType(uint32(code))
				}
			}
		}

	}
}
