package net

import (
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-vgo/robotgo"
	"github.com/gorilla/websocket"
	"github.com/indeedhat/track-pad/internal/env"
	"github.com/micmonay/keybd_event"
)

var upgrader = websocket.Upgrader{}

// WebsocketHandler for managing main connection to app
func WebsocketHandler(
	kb *keybd_event.KeyBonding,
	v func(string, ...any),
) func(rw http.ResponseWriter, r *http.Request) {
	serverPass := env.Get(env.ConnetPass)

	return func(rw http.ResponseWriter, r *http.Request) {
		con, err := upgrader.Upgrade(rw, r, nil)
		if err != nil {
			return
		}
		defer con.Close()

		if err := awaitPassword(con, serverPass, v); err != nil {
			return
		}

		for {
			cmdParts, err := readMessage(con)
			if err != nil {
				break
			}

			if len(cmdParts) == 0 {
				continue
			}

			v("%v", cmdParts)
			switch cmdParts[0] {
			case "move":
				processMoveMessage(cmdParts)
			case "scroll":
				processScrollMessage(cmdParts)
			case "click":
				processClickMessage(cmdParts)
			case "keeb":
				processKeebMessage(cmdParts, kb)
			case "zoom":
				processZoomMessage(cmdParts)

			}
		}
	}
}

func readMessage(con *websocket.Conn) ([]string, error) {
	_, message, err := con.ReadMessage()
	if err != nil {
		return nil, err
	}

	return strings.Split(string(message), ";"), nil
}

func awaitPassword(con *websocket.Conn, serverPass string, v func(string, ...any)) error {
	if serverPass == "" {
		con.WriteMessage(websocket.TextMessage, []byte("auth;true"))
		v("auth;true")
		return nil
	}

	con.WriteMessage(websocket.TextMessage, []byte("auth;false"))
	v("auth;false")

	for {
		cmdParts, err := readMessage(con)
		if err != nil {
			return err
		}

		if len(cmdParts) != 2 || cmdParts[0] != "pass" {
			continue
		}

		if cmdParts[1] == serverPass {
			con.WriteMessage(websocket.TextMessage, []byte("auth;true"))
			v("auth;true")
			return nil
		}

		con.WriteMessage(websocket.TextMessage, []byte("auth;false"))
		v("auth;false")
	}
}

func processMoveMessage(cmdParts []string) {
	if len(cmdParts) != 3 {
		return
	}
	x, _ := strconv.ParseFloat(cmdParts[1], 32)
	y, _ := strconv.ParseFloat(cmdParts[2], 32)

	robotgo.MoveRelative(int(x), int(y))
}

func processScrollMessage(cmdParts []string) {
	if len(cmdParts) != 3 {
		return
	}

	var (
		x int
		y int
	)
	ix, _ := strconv.ParseFloat(cmdParts[1], 32)
	iy, _ := strconv.ParseFloat(cmdParts[2], 32)
	if math.Abs(ix) > math.Abs(iy) {
		x = scrollDistance(ix)
	} else {
		y = scrollDistance(iy)
	}

	robotgo.Scroll(int(x), int(y), 0)
}

func scrollDistance(input float64) int {
	if input > 0 {
		return 1
	} else {
		return -1
	}
}

func processClickMessage(cmdParts []string) {
	if len(cmdParts) != 3 {
		return
	}
	if cmdParts[2] == "true" {
		robotgo.Toggle(cmdParts[1], "down")
	} else {
		robotgo.Toggle(cmdParts[1], "up")
	}
}

func processKeebMessage(cmdParts []string, kb *keybd_event.KeyBonding) {
	if len(cmdParts) != 2 {
		return
	}

	code, err := strconv.Atoi(cmdParts[1])

	if err != nil {
		return
	}

	if code == 2408 {
		kb.SetKeys(keybd_event.VK_BACKSPACE)
		_ = kb.Launching()
	} else {
		robotgo.UnicodeType(uint32(code))
	}
}

func processZoomMessage(cmdParts []string) {
	if len(cmdParts) != 2 {
		return
	}

	direction, err := strconv.ParseFloat(cmdParts[1], 32)
	if err != nil {
		return
	}

	key := "-"
	if direction > 1 {
		// this seems to work in most places that + deos and some that + deos not
		key = "="
	}
	robotgo.KeyDown("ctrl")
	robotgo.KeyTap(key)
	robotgo.KeyUp("ctrl")
}
