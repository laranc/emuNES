package rp2C02

import (
	"math/rand"

	"github.com/laranc/emuNES/cartridge"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	ResX  = 320
	ResY  = 240
	Scale = 4
)

type PPU struct {
	nameTable       [2][1024]uint8
	paletteTable    [32]uint8
	rom             *cartridge.ROM
	palScreen       [64]sdl.Color
	sprScreen       *sdl.Surface
	sprNameTable    [2]*sdl.Surface
	sprPatternTable [2]*sdl.Surface
	frameComplete   bool
	scanLine        uint16
	cycle           uint16
	status          StatusRegister
	mask            MaskRegister
	control         ControlRegister
	addressLatch    uint8
	dataBuffer      uint8
	address         uint16
	renderer        *sdl.Renderer
}

func NewPPU() *PPU {
	ppu := &PPU{
		nameTable:       [2][1024]uint8{},
		paletteTable:    [32]uint8{},
		rom:             nil,
		palScreen:       [64]sdl.Color{},
		sprScreen:       nil,                       // W: 256 H: 240
		sprNameTable:    [2]*sdl.Surface{nil, nil}, // W: 256 H: 240
		sprPatternTable: [2]*sdl.Surface{nil, nil}, // W: 128 H: 128
		frameComplete:   false,
		scanLine:        0,
		cycle:           0,
		status:          MakeStatusRegister(),
		mask:            MakeMaskRegister(),
		control:         MakeControlRegister(),
		addressLatch:    0,
		dataBuffer:      0x00,
		address:         0x0000,
		renderer:        nil,
	}
	ppu.sprScreen, _ = sdl.CreateRGBSurface(0, 256, 240, 0, 0, 0, 0, 0)
	ppu.sprNameTable[0], _ = sdl.CreateRGBSurface(0, 256, 240, 0, 0, 0, 0, 0)
	ppu.sprNameTable[1], _ = sdl.CreateRGBSurface(0, 256, 240, 0, 0, 0, 0, 0)
	ppu.sprPatternTable[0], _ = sdl.CreateRGBSurface(0, 128, 128, 0, 0, 0, 0, 0)
	ppu.sprPatternTable[1], _ = sdl.CreateRGBSurface(0, 128, 128, 0, 0, 0, 0, 0)
	ppu.palScreen[0x00] = sdl.Color{R: 84, G: 84, B: 84, A: 255}
	ppu.palScreen[0x01] = sdl.Color{R: 0, G: 30, B: 116, A: 255}
	ppu.palScreen[0x02] = sdl.Color{R: 8, G: 16, B: 144, A: 255}
	ppu.palScreen[0x03] = sdl.Color{R: 48, G: 0, B: 136, A: 255}
	ppu.palScreen[0x04] = sdl.Color{R: 68, G: 0, B: 100, A: 255}
	ppu.palScreen[0x05] = sdl.Color{R: 92, G: 0, B: 48, A: 255}
	ppu.palScreen[0x06] = sdl.Color{R: 84, G: 4, B: 0, A: 255}
	ppu.palScreen[0x07] = sdl.Color{R: 60, G: 24, B: 0, A: 255}
	ppu.palScreen[0x08] = sdl.Color{R: 32, G: 42, B: 0, A: 255}
	ppu.palScreen[0x09] = sdl.Color{R: 8, G: 58, B: 0, A: 255}
	ppu.palScreen[0x0A] = sdl.Color{R: 0, G: 64, B: 0, A: 255}
	ppu.palScreen[0x0B] = sdl.Color{R: 0, G: 60, B: 0, A: 255}
	ppu.palScreen[0x0C] = sdl.Color{R: 0, G: 50, B: 60, A: 255}
	ppu.palScreen[0x0D] = sdl.Color{R: 0, G: 0, B: 0, A: 255}
	ppu.palScreen[0x0E] = sdl.Color{R: 0, G: 0, B: 0, A: 255}
	ppu.palScreen[0x0F] = sdl.Color{R: 0, G: 0, B: 0, A: 255}
	ppu.palScreen[0x10] = sdl.Color{R: 152, G: 150, B: 152, A: 255}
	ppu.palScreen[0x11] = sdl.Color{R: 8, G: 76, B: 196, A: 255}
	ppu.palScreen[0x12] = sdl.Color{R: 48, G: 50, B: 236, A: 255}
	ppu.palScreen[0x13] = sdl.Color{R: 92, G: 30, B: 228, A: 255}
	ppu.palScreen[0x14] = sdl.Color{R: 136, G: 20, B: 176, A: 255}
	ppu.palScreen[0x15] = sdl.Color{R: 160, G: 20, B: 100, A: 255}
	ppu.palScreen[0x16] = sdl.Color{R: 152, G: 34, B: 32, A: 255}
	ppu.palScreen[0x17] = sdl.Color{R: 120, G: 60, B: 0, A: 255}
	ppu.palScreen[0x18] = sdl.Color{R: 84, G: 90, B: 0, A: 255}
	ppu.palScreen[0x19] = sdl.Color{R: 40, G: 114, B: 0, A: 255}
	ppu.palScreen[0x1A] = sdl.Color{R: 8, G: 124, B: 0, A: 255}
	ppu.palScreen[0x1B] = sdl.Color{R: 0, G: 118, B: 40, A: 255}
	ppu.palScreen[0x1C] = sdl.Color{R: 0, G: 102, B: 120, A: 255}
	ppu.palScreen[0x1D] = sdl.Color{R: 0, G: 0, B: 0, A: 255}
	ppu.palScreen[0x1E] = sdl.Color{R: 0, G: 0, B: 0, A: 255}
	ppu.palScreen[0x1F] = sdl.Color{R: 0, G: 0, B: 0, A: 255}
	ppu.palScreen[0x20] = sdl.Color{R: 236, G: 238, B: 236, A: 255}
	ppu.palScreen[0x21] = sdl.Color{R: 76, G: 154, B: 236, A: 255}
	ppu.palScreen[0x22] = sdl.Color{R: 120, G: 124, B: 236, A: 255}
	ppu.palScreen[0x23] = sdl.Color{R: 176, G: 98, B: 236, A: 255}
	ppu.palScreen[0x24] = sdl.Color{R: 228, G: 84, B: 236, A: 255}
	ppu.palScreen[0x25] = sdl.Color{R: 236, G: 88, B: 180, A: 255}
	ppu.palScreen[0x26] = sdl.Color{R: 236, G: 106, B: 100, A: 255}
	ppu.palScreen[0x27] = sdl.Color{R: 212, G: 136, B: 32, A: 255}
	ppu.palScreen[0x28] = sdl.Color{R: 160, G: 170, B: 0, A: 255}
	ppu.palScreen[0x29] = sdl.Color{R: 116, G: 196, B: 0, A: 255}
	ppu.palScreen[0x2A] = sdl.Color{R: 76, G: 208, B: 32, A: 255}
	ppu.palScreen[0x2B] = sdl.Color{R: 56, G: 204, B: 108, A: 255}
	ppu.palScreen[0x2C] = sdl.Color{R: 56, G: 180, B: 204, A: 255}
	ppu.palScreen[0x2D] = sdl.Color{R: 60, G: 60, B: 60, A: 255}
	ppu.palScreen[0x2E] = sdl.Color{R: 0, G: 0, B: 0, A: 255}
	ppu.palScreen[0x2F] = sdl.Color{R: 0, G: 0, B: 0, A: 255}
	ppu.palScreen[0x30] = sdl.Color{R: 236, G: 238, B: 236, A: 255}
	ppu.palScreen[0x31] = sdl.Color{R: 168, G: 204, B: 236, A: 255}
	ppu.palScreen[0x32] = sdl.Color{R: 188, G: 188, B: 236, A: 255}
	ppu.palScreen[0x33] = sdl.Color{R: 212, G: 178, B: 236, A: 255}
	ppu.palScreen[0x34] = sdl.Color{R: 236, G: 174, B: 236, A: 255}
	ppu.palScreen[0x35] = sdl.Color{R: 236, G: 174, B: 212, A: 255}
	ppu.palScreen[0x36] = sdl.Color{R: 236, G: 180, B: 176, A: 255}
	ppu.palScreen[0x37] = sdl.Color{R: 228, G: 196, B: 144, A: 255}
	ppu.palScreen[0x38] = sdl.Color{R: 204, G: 210, B: 120, A: 255}
	ppu.palScreen[0x39] = sdl.Color{R: 180, G: 222, B: 120, A: 255}
	ppu.palScreen[0x3A] = sdl.Color{R: 168, G: 226, B: 144, A: 255}
	ppu.palScreen[0x3B] = sdl.Color{R: 152, G: 226, B: 180, A: 255}
	ppu.palScreen[0x3C] = sdl.Color{R: 160, G: 214, B: 228, A: 255}
	ppu.palScreen[0x3D] = sdl.Color{R: 160, G: 162, B: 160, A: 255}
	ppu.palScreen[0x3E] = sdl.Color{R: 0, G: 0, B: 0, A: 255}
	ppu.palScreen[0x3F] = sdl.Color{R: 0, G: 0, B: 0, A: 255}

	return ppu
}

func (ppu *PPU) Read(addr uint16, readOnly bool) uint8 {
	var data uint8 = 0x00
	addr &= 0x3FFF
	if ppu.rom.PPURead(addr, &data) {
		// Read from the rom or pass and read from PPU memory
	} else if addr >= 0x2000 && addr <= 0x3EFF {
		addr &= 0x0FFF
		switch ppu.rom.GetMirror() {
		case cartridge.MirrorVertical:
			if addr <= 0x03FF {
				data = ppu.nameTable[0][addr&0x03FF]
			} else if addr >= 0x0400 && addr <= 0x07FF {
				data = ppu.nameTable[1][addr&0x03FF]
			} else if addr >= 0x0800 && addr <= 0x0BFF {
				data = ppu.nameTable[0][addr&0x03FF]
			} else if addr >= 0x0C00 && addr <= 0x0FFF {
				data = ppu.nameTable[1][addr&0x03FF]
			}
		case cartridge.MirrorHorizontal:
			if addr <= 0x03FF {
				data = ppu.nameTable[0][addr&0x03FF]
			} else if addr >= 0x0400 && addr <= 0x07FF {
				data = ppu.nameTable[0][addr&0x03FF]
			} else if addr >= 0x0800 && addr <= 0x0BFF {
				data = ppu.nameTable[1][addr&0x03FF]
			} else if addr >= 0x0C00 && addr <= 0x0FFF {
				data = ppu.nameTable[1][addr&0x03FF]
			}
		default:
			break
		}
	} else if addr >= 0x3F00 && addr <= 0x3FFF {
		addr &= 0x001F
		switch addr {
		case 0x0010:
			addr = 0x0000
		case 0x0014:
			addr = 0x0004
		case 0x0018:
			addr = 0x0008
		case 0x001C:
			addr = 0x000C
		default:
			break
		}
		grayscale := uint8(0x3F)
		if ppu.mask.Grayscale != 0 {
			grayscale = 0x30
		}
		data = ppu.paletteTable[addr] & grayscale
	}
	return data
}

func (ppu *PPU) Write(addr uint16, data uint8) {
	addr &= 0x3FFF
	if ppu.rom.PPUWrite(addr, data) {
		// Write to the ROM or pass and write to PPU memory
	} else if addr >= 0x2000 && addr <= 0x3EFF {
		addr &= 0x0FFF
		switch ppu.rom.GetMirror() {
		case cartridge.MirrorVertical:
			if addr <= 0x03FF {
				ppu.nameTable[0][addr&0x03FF] = data
			} else if addr >= 0x0400 && addr <= 0x07FF {
				ppu.nameTable[1][addr&0x03FF] = data
			} else if addr >= 0x0800 && addr <= 0x0BFF {
				ppu.nameTable[0][addr&0x03FF] = data
			} else if addr >= 0x0C00 && addr <= 0x0FFF {
				ppu.nameTable[1][addr&0x03FF] = data
			}
		case cartridge.MirrorHorizontal:
			if addr <= 0x03FF {
				ppu.nameTable[0][addr&0x03FF] = data
			} else if addr >= 0x0400 && addr <= 0x07FF {
				ppu.nameTable[0][addr&0x03FF] = data
			} else if addr >= 0x0800 && addr <= 0x0BFF {
				ppu.nameTable[1][addr&0x03FF] = data
			} else if addr >= 0x0C00 && addr <= 0x0FFF {
				ppu.nameTable[1][addr&0x03FF] = data
			}
		default:
			break
		}
	} else if addr >= 0x3F00 && addr <= 0x3FFF {
		addr &= 0x001F
		switch addr {
		case 0x0010:
			addr = 0x0000
		case 0x0014:
			addr = 0x0004
		case 0x0018:
			addr = 0x0008
		case 0x001C:
			addr = 0x000C
		default:
			break
		}
		ppu.paletteTable[addr] = data
	}
}

func (ppu *PPU) BusRead(addr uint16, readOnly bool) uint8 {
	var data uint8 = 0x00

	switch addr {
	case 0x0000: // Control
		break
	case 0x0001: // Mask
		break
	case 0x0002: // Status
		ppu.status.VerticalBlank = 1
		ppu.status.Update()
		data = (ppu.status.Reg & 0xE0) | (ppu.dataBuffer & 0x1F)
		ppu.status.VerticalBlank = 0
		ppu.status.Update()
		ppu.addressLatch = 0
	case 0x0003: // OAM Address
		break
	case 0x0004: // OAM Data
		break
	case 0x0005: // Scroll
		break
	case 0x0006: // PPU Address
		break
	case 0x0007: // PPU Data
		data = ppu.dataBuffer
		ppu.dataBuffer = ppu.Read(ppu.address, false)
		if ppu.address > 0x3F00 {
			data = ppu.dataBuffer
		}
	default:
		break
	}

	return data
}

func (ppu *PPU) BusWrite(addr uint16, data uint8) {
	switch addr {
	case 0x0000: // Control
		break
	case 0x0001: // Mask
		break
	case 0x0002: // Status
		break
	case 0x0003: // OAM Address
		break
	case 0x0004: // OAM Data
		break
	case 0x0005: // Scroll
		break
	case 0x0006: // PPU Address
		if ppu.addressLatch == 0 {
			ppu.address = (ppu.address & 0xFF00) | uint16(data)
			ppu.addressLatch = 1
		} else {
			ppu.address = (ppu.address & 0x00FF) | uint16(data)<<8
			ppu.addressLatch = 0
		}
	case 0x0007: // PPU Data
		ppu.Write(addr, data)
	default:
		break
	}
}

func (ppu *PPU) ConnectCartridge(rom *cartridge.ROM) {
	ppu.rom = rom
}

func (ppu *PPU) ConnectRenderer(renderer *sdl.Renderer) {
	ppu.renderer = renderer
}

func (ppu *PPU) GetScreen() *sdl.Surface {
	return ppu.sprScreen
}

func (ppu *PPU) GetNameTable(i uint8) *sdl.Surface {
	return ppu.sprNameTable[i]
}

func (ppu *PPU) GetPatternTable(i uint8, palette uint8) *sdl.Surface {
	for x := range uint16(16) {
		for y := range uint16(16) {
			offset := y*256 + x*16
			for row := range uint16(8) {
				tileLSB := ppu.Read(uint16(i)*0x1000+offset+row, true)
				tileMSB := ppu.Read(uint16(i)*0x1000+offset+row+8, true)
				for col := range uint16(8) {
					pixel := (tileLSB & 0x01) + (tileMSB & 0x01)
					tileLSB >>= 1
					tileMSB >>= 1
					color := ppu.GetColorPallette(palette, pixel)
					argb := sdl.ARGB8888{A: color.R, R: color.G, G: color.B, B: color.A}
					ppu.sprPatternTable[i].Set(int(x*8+(7-col)), int(y*8+row), argb)
				}
			}
		}
	}
	return ppu.sprPatternTable[i]
}

func (ppu *PPU) GetColorPallette(palette uint8, pixel uint8) sdl.Color {
	return ppu.palScreen[ppu.Read(0x3F00+uint16((palette<<2)+pixel), true)]
}

func (ppu *PPU) Clock() {
	i := 0
	if rand.Int()%2 != 0 {
		i = 0x3F
	} else {
		i = 0x30
	}
	c := ppu.palScreen[i]
	ppu.renderer.SetDrawColor(c.R, c.G, c.B, c.A)
	ppu.renderer.DrawPoint(int32(ppu.cycle)-1, int32(ppu.scanLine))

	ppu.cycle++
	if ppu.cycle >= 341 {
		ppu.cycle = 0
		ppu.scanLine++
		if ppu.scanLine >= 261 {
			ppu.scanLine = 0xFFFF
			ppu.frameComplete = true
		}
	}

}

func (ppu *PPU) Reset() {
	ppu.frameComplete = false
}

func (ppu *PPU) FrameComplete() bool {
	return ppu.frameComplete
}
