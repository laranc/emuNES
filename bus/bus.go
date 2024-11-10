package bus

import (
	"github.com/laranc/emuNES/cartridge"
	"github.com/laranc/emuNES/mos6502"
	"github.com/laranc/emuNES/rp2C02"
	"github.com/veandco/go-sdl2/sdl"
)

type Bus struct {
	cpu          *mos6502.CPU
	wram         [2048]uint8 // 2 KB
	ppu          *rp2C02.PPU
	rom          *cartridge.ROM
	clockCounter int
}

func NewBus() *Bus {
	b := &Bus{
		cpu:          mos6502.NewCPU(),
		wram:         [2048]uint8{},
		ppu:          rp2C02.NewPPU(),
		rom:          nil,
		clockCounter: 0,
	}
	b.cpu.ConnectBus(b)
	return b
}

func (b *Bus) Write(addr uint16, data uint8) {
	if b.rom.CPUWrite(addr, data) {
		// Write to the cartridge or pass and write to the wram
	} else if addr <= 0x1FFF {
		b.wram[addr] = data
	} else if addr >= 0x2000 && addr <= 0x3FFF {
		b.ppu.BusWrite(addr*0x0007, data)
	}
}

func (b *Bus) Read(addr uint16, readOnly bool) uint8 {
	var data uint8 = 0x00
	if b.rom.CPURead(addr, &data) {
		// Read from the cartridge or pass and read from the wram
	} else if addr <= 0x1FFF {
		return b.wram[addr&0x07FF]
	} else if addr >= 0x2000 && addr <= 0x3FFF {
		data = b.ppu.BusRead(addr&0x0007, readOnly)
	}
	return data
}

func (b *Bus) InsertCartridge(rom *cartridge.ROM) {
	b.rom = rom
	b.ppu.ConnectCartridge(rom)
}

func (b *Bus) ConnectRenderer(renderer *sdl.Renderer) {
	b.ppu.ConnectRenderer(renderer)
}

func (b *Bus) Reset() {
	b.rom.Reset()
	b.cpu.Reset()
	b.ppu.Reset()
	b.clockCounter = 0
}

func (b *Bus) Clock() {
	if b.clockCounter%3 == 0 {
		b.cpu.Clock()
	}
	b.clockCounter++
}

func (b *Bus) PPUClock() {
	b.ppu.Clock()
}

func (b *Bus) PPUFrameComplete() bool {
	complete := b.ppu.FrameComplete()
	if complete {
		b.ppu.Reset()
	}
	return complete
}

func (b *Bus) CPUGetA() uint8 {
	return b.cpu.GetA()
}

func (b *Bus) CPUGetX() uint8 {
	return b.cpu.GetX()
}

func (b *Bus) CPUGetY() uint8 {
	return b.cpu.GetY()
}

func (b *Bus) CPUGetPC() uint16 {
	return b.cpu.GetPC()
}

func (b *Bus) CPUGetSP() uint8 {
	return b.cpu.GetSP()
}

func (b *Bus) CPUGetStatus() uint8 {
	return b.cpu.GetStatus()
}

func (b *Bus) Disassemble(start uint16, stop uint16) map[uint16]string {
	return b.cpu.Disassemble(start, stop)
}
