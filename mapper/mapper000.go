package mapper

type Mapper000 struct {
	prgBanks uint8
	chrBanks uint8
}

func MakeMapper000(prgBanks uint8, chrBanks uint8) Mapper000 {
	return Mapper000{
		prgBanks: prgBanks,
		chrBanks: chrBanks,
	}
}

func (m Mapper000) CPUMapRead(addr uint16, mappedAddr *uint32) bool {
	if addr >= 0x8000 && addr <= 0xFFFF {
		a := uint16(0x3FFF)
		if m.prgBanks > 1 {
			a = 0x7FFF
		}
		*mappedAddr = uint32(addr & a)
		return true
	}
	return false
}

func (m Mapper000) CPUMapWrite(addr uint16, mappedAddr *uint32) bool {
	if addr >= 0x8000 && addr <= 0xFFFF {
		a := uint16(0x3FFF)
		if m.prgBanks > 1 {
			a = 0x7FFF
		}
		*mappedAddr = uint32(addr & a)
		return true
	}
	return false
}

func (m Mapper000) PPUMapRead(addr uint16, mappedAddr *uint32) bool {
	if addr <= 0x1FFF {
		*mappedAddr = uint32(addr)
		return true
	}
	return false
}

func (m Mapper000) PPUMapWrite(addr uint16, mappedAddr *uint32) bool {
	if addr <= 0x1FFF && m.chrBanks == 0 {
		*mappedAddr = uint32(addr)
		return true
	}
	return false
}

func (m Mapper000) Reset() {}
