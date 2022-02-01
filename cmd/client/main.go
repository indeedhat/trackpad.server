package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

var (
	lastX int
	lastY int
)

var conn *websocket.Conn

type Game struct {
	grids      map[image.Rectangle]ebiten.CursorShapeType
	gridColors map[image.Rectangle]color.Color
}

func (g *Game) Update() error {
	x, y := ebiten.CursorPosition()
	if x != lastX || y != lastY {
		conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("x(%d) y(%d)\n", x-lastX, y-lastY)))
		lastX = x
		lastY = y
	}

	pt := image.Pt(ebiten.CursorPosition())
	for r, c := range g.grids {
		if pt.In(r) {
			ebiten.SetCursorShape(c)
			return nil
		}
	}
	ebiten.SetCursorShape(ebiten.CursorShapeDefault)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for r, c := range g.gridColors {
		ebitenutil.DrawRect(screen, float64(r.Min.X), float64(r.Min.Y), float64(r.Dx()), float64(r.Dy()), c)
	}

	switch ebiten.CursorShape() {
	case ebiten.CursorShapeDefault:
		ebitenutil.DebugPrint(screen, "CursorShape: Default")
	case ebiten.CursorShapeText:
		ebitenutil.DebugPrint(screen, "CursorShape: Text")
	case ebiten.CursorShapeCrosshair:
		ebitenutil.DebugPrint(screen, "CursorShape: Crosshair")
	case ebiten.CursorShapePointer:
		ebitenutil.DebugPrint(screen, "CursorShape: Pointer")
	case ebiten.CursorShapeEWResize:
		ebitenutil.DebugPrint(screen, "CursorShape: EW Resize")
	case ebiten.CursorShapeNSResize:
		ebitenutil.DebugPrint(screen, "CursorShape: NS Resize")
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	var err error

	uri := url.URL{Scheme: "ws", Host: "localhost:8881", Path: "/ws"}
	conn, _, err = websocket.DefaultDialer.Dial(uri.String(), nil)
	if err != nil {
		log.Fatal("failed to connect to server")
	}
	defer conn.Close()

	g := &Game{
		grids: map[image.Rectangle]ebiten.CursorShapeType{
			image.Rect(100, 100, 200, 200): ebiten.CursorShapeDefault,
			image.Rect(200, 100, 300, 200): ebiten.CursorShapeText,
			image.Rect(300, 100, 400, 200): ebiten.CursorShapeCrosshair,
			image.Rect(100, 200, 200, 300): ebiten.CursorShapePointer,
			image.Rect(200, 200, 300, 300): ebiten.CursorShapeEWResize,
			image.Rect(300, 200, 400, 300): ebiten.CursorShapeNSResize,
		},
		gridColors: map[image.Rectangle]color.Color{},
	}
	for rect, c := range g.grids {
		clr := color.RGBA{0x40, 0x40, 0x40, 0xff}
		if c%2 == 0 {
			clr.R = 0x80
		} else {
			clr.G = 0x80
		}
		g.gridColors[rect] = clr
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Cursor (Ebiten Demo)")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
