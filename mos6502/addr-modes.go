package mos6502

func (cpu *CPU) IMP() uint8 {
	cpu.fetched = cpu.a
	return 0
}

func (cpu *CPU) ZP0() uint8 {
	cpu.addrAbs = uint16(cpu.read(cpu.pc))
	cpu.pc++
	cpu.addrAbs &= 0x00FF
	return 0
}

func (cpu *CPU) ZPY() uint8 {
	cpu.addrAbs = uint16(cpu.read(cpu.pc) + cpu.y)
	return 0
}

func (cpu *CPU) ABS() uint8 {
	low := uint16(cpu.read(cpu.pc))
	cpu.pc++
	high := uint16(cpu.read(cpu.pc))
	cpu.pc++
	cpu.addrAbs = (high << 8) | low
	return 0
}

func (cpu *CPU) ABY() uint8 {
	low := uint16(cpu.read(cpu.pc))
	cpu.pc++
	high := uint16(cpu.read(cpu.pc))
	cpu.pc++
	cpu.addrAbs = (high << 8) | low
	cpu.addrAbs += uint16(cpu.y)
	if (cpu.addrAbs & 0xFF00) != (high << 8) {
		return 1
	} else {
		return 0
	}
}

func (cpu *CPU) IZX() uint8 {
	t := uint16(cpu.read(cpu.pc))
	cpu.pc++
	low := uint16(cpu.read((t + uint16(cpu.x)) & 0x00FF))
	high := uint16(cpu.read(t+uint16(cpu.x)+1) & 0x00FF)
	cpu.addrAbs = (high << 8) | low
	return 0
}

func (cpu *CPU) IMM() uint8 {
	cpu.addrAbs = cpu.pc
	cpu.pc++
	return 0
}

func (cpu *CPU) ZPX() uint8 {
	cpu.addrAbs = uint16(cpu.read(cpu.pc) + cpu.x)
	return 0
}

func (cpu *CPU) REL() uint8 {
	cpu.addrRel = uint16(cpu.read(cpu.pc))
	cpu.pc++
	if (cpu.addrRel & 0x0080) != 1 {
		cpu.addrRel |= 0xFF00
	}
	return 0
}

func (cpu *CPU) ABX() uint8 {
	low := uint16(cpu.read(cpu.pc))
	cpu.pc++
	high := uint16(cpu.read(cpu.pc))
	cpu.pc++
	cpu.addrAbs = (high << 8) | low
	cpu.addrAbs += uint16(cpu.x)
	if (cpu.addrAbs & 0xFF00) != (high << 8) {
		return 1
	} else {
		return 0
	}
}

func (cpu *CPU) IND() uint8 {
	ptr_low := uint16(cpu.read(cpu.pc))
	cpu.pc++
	ptr_high := uint16(cpu.read(cpu.pc))
	cpu.pc++
	ptr := (ptr_high << 8) | ptr_low
	if ptr_low == 0x00FF {
		cpu.addrAbs = (uint16(cpu.read(ptr&0xFF00)) << 8) | uint16(cpu.read(ptr))
	} else {
		cpu.addrAbs = (uint16(cpu.read(ptr+1)) << 8) | uint16(cpu.read(ptr))
	}
	return 0
}

func (cpu *CPU) IZY() uint8 {
	t := uint16(cpu.pc)
	cpu.pc++
	low := uint16(cpu.read(t & 0x00FF))
	high := uint16(cpu.read((t + 1) & 0x00FF))
	cpu.addrAbs = (high << 8) | low
	cpu.addrAbs += uint16(cpu.y)
	if (cpu.addrAbs & 0xFF00) != (high << 8) {
		return 1
	} else {
		return 0
	}
}
