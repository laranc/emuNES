package cartridge

import (
	"encoding/binary"
	"io"
	"log"
	"os"

	"github.com/laranc/emuNES/mapper"
)

const (
	MirrorHorizontal uint8 = iota
	MirrorVertical
	MirrorOnescreenLow
	MirrorOnescreenHigh
)

type ROM struct {
	prg        []uint8
	chr        []uint8
	imageValid bool
	mapperID   uint8
	prgBanks   uint8
	chrBanks   uint8
	mirror     uint8
	mapper     mapper.Mapper
}

type Header struct {
	Name       [4]byte
	PrgChunks  uint8
	ChrChunks  uint8
	Mapper1    uint8
	Mapper2    uint8
	PrgRamSize uint8
	TVSystem1  uint8
	TVSystem2  uint8
	Unused     [4]byte
}

func NewROM(file string) *ROM {
	var rom *ROM = &ROM{}
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	h := &Header{}
	err = binary.Read(f, binary.NativeEndian, h)
	if err != nil {
		log.Fatal(err)
	}

	if (h.Mapper1 & 0x04) != 0 {
		f.Seek(512, io.SeekCurrent)
	}
	rom.mapperID = ((h.Mapper2 >> 4) << 4) | (h.Mapper1 >> 4)
	fileType := 1
	switch fileType {
	case 1:
		rom.prgBanks = h.PrgChunks
		rom.prg = make([]uint8, int(rom.prgBanks)*16384)
		_, err := f.Read(rom.prg)
		if err != nil {
			log.Fatal(err)
		}
		rom.chrBanks = h.ChrChunks
		rom.chr = make([]uint8, int(rom.chrBanks)*8192)
		_, err = f.Read(rom.chr)
		if err != nil {
			log.Fatal(err)
		}
	default:
		break
	}
	switch rom.mapperID {
	case 0:
		rom.mapper = mapper.MakeMapper000(rom.prgBanks, rom.chrBanks)
	default:
		break
	}
	rom.imageValid = true
	rom.mirror = MirrorHorizontal
	return rom
}

func (rom *ROM) CPUWrite(addr uint16, data uint8) bool {
	var mappedAddr uint32 = 0
	if rom.mapper.CPUMapWrite(addr, &mappedAddr) {
		rom.prg[mappedAddr] = data
		return true
	}
	return false
}

func (rom *ROM) CPURead(addr uint16, data *uint8) bool {
	var mappedAddr uint32 = 0
	if rom.mapper.CPUMapRead(addr, &mappedAddr) {
		*data = rom.prg[mappedAddr]
		return true
	}
	return false
}

func (rom *ROM) PPUWrite(addr uint16, data uint8) bool {
	var mappedAddr uint32 = 0
	if rom.mapper.PPUMapWrite(addr, &mappedAddr) {
		rom.chr[mappedAddr] = data
		return true
	}
	return false
}

func (rom *ROM) PPURead(addr uint16, data *uint8) bool {
	var mappedAddr uint32 = 0
	if rom.mapper.PPUMapRead(addr, &mappedAddr) {
		*data = rom.chr[mappedAddr]
		return true
	}
	return false
}

func (rom *ROM) ImageValid() bool {
	return rom.imageValid
}

func (rom *ROM) Reset() {
	if rom.mapper != nil {
		rom.mapper.Reset()
	}
}

func (rom *ROM) GetMirror() uint8 {
	return rom.mirror
}
