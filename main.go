package main

import (
	"fmt"
	"log"
	"sort"

	"github.com/laranc/emuNES/bus"
	"github.com/laranc/emuNES/cartridge"
	"github.com/laranc/emuNES/mos6502"
	"github.com/laranc/emuNES/rp2C02"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

// Constants
const (
	title  = "emuNES"
	width  = 680
	height = 480
	scale  = 2
)

// Global State
var (
	debugWindow   *sdl.Window   = nil
	gameWindow    *sdl.Window   = nil
	debugRenderer *sdl.Renderer = nil
	gameRenderer  *sdl.Renderer = nil
	font          *ttf.Font     = nil
	nes           *bus.Bus      = nil
	stepMode      bool          = false
	step          bool          = false
	asm           map[uint16]string
	asmAddrs      []uint16
)

// Colors
var (
	black      = sdl.Color{R: 0, G: 0, B: 0, A: 0}
	white      = sdl.Color{R: 255, G: 255, B: 255, A: 255}
	red        = sdl.Color{R: 255, G: 0, B: 0, A: 255}
	green      = sdl.Color{R: 0, G: 255, B: 0, A: 255}
	cyan       = sdl.Color{R: 0, G: 255, B: 255, A: 255}
	background = sdl.Color{R: 0, G: 0, B: 128, A: 255}
)

func main() {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic(err)
	}
	defer sdl.Quit()

	err = ttf.Init()
	if err != nil {
		panic(err)
	}
	defer ttf.Quit()

	debugWindow, err = sdl.CreateWindow(title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, width*scale, height*scale, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer debugWindow.Destroy()

	gameWindow, err = sdl.CreateWindow("Game Title", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, rp2C02.ResX*rp2C02.Scale, rp2C02.ResY*rp2C02.Scale, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer gameWindow.Destroy()

	debugRenderer, err = sdl.CreateRenderer(debugWindow, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	defer debugRenderer.Destroy()
	debugRenderer.SetScale(scale, scale)

	gameRenderer, err = sdl.CreateRenderer(gameWindow, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	defer gameRenderer.Destroy()
	gameRenderer.SetScale(rp2C02.Scale, rp2C02.Scale)

	font, err = ttf.OpenFont("./assets/nes.ttf", 8)
	if err != nil {
		panic(err)
	}
	defer font.Close()

	nes = bus.NewBus()
	nes.ConnectRenderer(gameRenderer)
	cart := cartridge.NewROM("./nestest.nes")
	if !cart.ImageValid() {
		log.Fatal("reading from rom failed")
	}
	nes.InsertCartridge(cart)
	asm = nes.Disassemble(0x0000, 0xFFFF)
	asmAddrs = make([]uint16, 0, len(asm))
	for k := range asm {
		asmAddrs = append(asmAddrs, k)
	}
	sort.Slice(asmAddrs, func(i int, j int) bool {
		return asmAddrs[i] < asmAddrs[j]
	})
	nes.Reset()
	run()
}

func run() {
	go func() {
		for {
			if !stepMode {
				nes.Clock()
			} else {
				if step {
					nes.Clock()
					step = false
				}
			}
		}

	}()
	running := true
	for running {
		for e := sdl.PollEvent(); e != nil; e = sdl.PollEvent() {
			switch t := e.(type) {
			case sdl.QuitEvent:
				running = false
			case sdl.KeyboardEvent:
				switch t.Keysym.Sym {
				case sdl.K_TAB:
					stepMode = !stepMode
					if stepMode {
						fmt.Println("Step by step mode on")
					} else {
						fmt.Println("Step by step mode off")
					}
				case sdl.K_SPACE:
					step = true
				default:
					break
				}
			}
		}

		debugRenderer.SetDrawColor(background.R, background.G, background.B, background.A)
		debugRenderer.Clear()
		drawRAM(2, 2, 0x0000, 16, 16)
		drawRAM(2, 182, 0x8000, 16, 16)
		drawCPU(448, 2)
		drawCode(448, 72, 26)
		drawPalettes()
		//drawPatternTable(x, y)

		gameRenderer.SetDrawColor(0, 0, 0, 255)
		gameRenderer.Clear()
		for !nes.PPUFrameComplete() {
			nes.PPUClock()
		}

		debugRenderer.Present()
		gameRenderer.Present()
		sdl.Delay(16)
	}
}

func drawText(str string, x int32, y int32, color sdl.Color) {
	text, err := font.RenderUTF8Blended(str, color)
	if err != nil {
		log.Fatal(err)
	}
	defer text.Free()
	src := sdl.Rect{X: 0, Y: 0, W: text.W, H: text.H}
	dst := sdl.Rect{X: x, Y: y, W: text.W, H: text.H}
	texture, err := debugRenderer.CreateTextureFromSurface(text)
	if err != nil {
		log.Fatal(err)
	}
	defer texture.Destroy()
	debugRenderer.Copy(texture, &src, &dst)
}

func drawRAM(x int32, y int32, addr uint16, rows int, columns int) {
	ramX := x
	ramY := y
	for range rows {
		offset := "$" + fmt.Sprintf("%04X", addr) + ":"
		for range columns {
			offset += " " + fmt.Sprintf("%02X", nes.Read(addr, true))
			addr += 1
		}
		drawText(offset, ramX, ramY, white)
		ramY += 10
	}
}

func statusColor(flag uint8) sdl.Color {
	if nes.CPUGetStatus()&flag != 0 {
		return green
	} else {
		return red
	}
}

func drawCPU(x int32, y int32) {
	drawText("STATUS: ", x, y, white)
	drawText("N", x+64, y, statusColor(mos6502.N))
	drawText("V", x+80, y, statusColor(mos6502.V))
	drawText("-", x+96, y, statusColor(mos6502.U))
	drawText("B", x+112, y, statusColor(mos6502.B))
	drawText("D", x+128, y, statusColor(mos6502.D))
	drawText("I", x+144, y, statusColor(mos6502.I))
	drawText("Z", x+160, y, statusColor(mos6502.Z))
	drawText("C", x+178, y, statusColor(mos6502.C))
	drawText("PC: $"+fmt.Sprintf("%04X", nes.CPUGetPC()), x, y+10, white)
	drawText("SP: $"+fmt.Sprintf("%02X", nes.CPUGetSP()), x, y+50, white)
	drawText("A: $"+fmt.Sprintf("%02X", nes.CPUGetA()), x, y+20, white)
	drawText("X: $"+fmt.Sprintf("%02X", nes.CPUGetX()), x, y+30, white)
	drawText("Y: $"+fmt.Sprintf("%02X", nes.CPUGetY()), x, y+40, white)
}

func drawCode(x int32, y int32, lines int32) {
	pc := nes.CPUGetPC()
	addr := sort.Search(len(asmAddrs), func(i int) bool {
		return asmAddrs[i] >= pc
	})
	lineY := (lines>>1)*10 + y
	drawText(asm[asmAddrs[addr]], x, lineY, cyan)
	for i := addr + 1; i < len(asmAddrs) && lineY < (lines*10)+y; i++ {
		lineY += 10
		drawText(asm[asmAddrs[i]], x, lineY, white)
	}
	lineY = (lines>>1)*10 + y
	for i := addr - 1; i >= 0 && lineY > y; i-- {
		lineY -= 10
		drawText(asm[asmAddrs[i]], x, lineY, white)
	}

}

func drawPalettes() {

}

func drawPatternTable(x int32, y int32) {

}
