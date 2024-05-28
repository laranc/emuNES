package mos6502

func (cpu *CPU) ADC() uint8 {
	cpu.fetch()
	temp := uint16(cpu.a + cpu.fetched + cpu.getFlag(C))
	cpu.setFlag(C, temp > 255)
	cpu.setFlag(Z, (temp&0x00FF) == 0)
	cpu.setFlag(N, (temp&0x0080) != 0)
	v := ^(uint16(cpu.a) ^ uint16(cpu.fetched)&(uint16(cpu.a)^temp)) & 0x0080
	cpu.setFlag(V, v != 0)
	cpu.a = uint8(temp & 0x00FF)
	return 1
}

func (cpu *CPU) AND() uint8 {
	cpu.fetch()
	cpu.a &= cpu.fetched
	cpu.setFlag(Z, cpu.a == 0x00)
	cpu.setFlag(N, (cpu.a&0x80) != 0)
	return 1
}

func (cpu *CPU) ASL() uint8 {
	cpu.fetch()
	temp := uint16(cpu.fetched) << 1
	cpu.setFlag(C, (temp&0xFF00) > 0)
	cpu.setFlag(Z, (temp&0x00FF) == 0x00)
	cpu.setFlag(N, (temp&0x0080) != 0)
	if cpu.isIMP() {
		cpu.a = uint8(temp & 0x00FF)
	} else {
		cpu.write(cpu.addrAbs, uint8(temp&0x00FF))
	}
	return 0
}

func (cpu *CPU) BCC() uint8 {
	if cpu.getFlag(C) == 0 {
		cpu.cycles++
		cpu.addrAbs = cpu.pc + cpu.addrRel
		if (cpu.addrAbs & 0xFF00) != (cpu.pc & 0xFF00) {
			cpu.cycles++
		}
		cpu.pc = cpu.addrAbs
	}
	return 0
}

func (cpu *CPU) BCS() uint8 {
	if cpu.getFlag(C) == 1 {
		cpu.cycles++
		cpu.addrAbs = cpu.pc + cpu.addrRel
		if (cpu.addrAbs & 0xFF00) != (cpu.pc & 0xFF00) {
			cpu.cycles++
		}
		cpu.pc = cpu.addrAbs
	}
	return 0
}

func (cpu *CPU) BEQ() uint8 {
	if cpu.getFlag(Z) == 1 {
		cpu.cycles++
		cpu.addrAbs = cpu.pc + cpu.addrRel
		if (cpu.addrAbs & 0xFF00) != (cpu.pc & 0xFF00) {
			cpu.cycles++
		}
		cpu.pc = cpu.addrAbs
	}
	return 0
}

func (cpu *CPU) BIT() uint8 {
	cpu.fetch()
	temp := cpu.a & cpu.fetched
	cpu.setFlag(Z, (temp&0x00FF == 0x00))
	cpu.setFlag(N, (cpu.fetched&(1<<7)) != 0)
	cpu.setFlag(V, (cpu.fetched&(1<<6)) != 0)
	return 0
}

func (cpu *CPU) BMI() uint8 {
	if cpu.getFlag(N) == 1 {
		cpu.cycles++
		cpu.addrAbs = cpu.pc + cpu.addrRel
		if (cpu.addrAbs & 0xFF00) != (cpu.pc & 0xFF00) {
			cpu.cycles++
		}
		cpu.pc = cpu.addrAbs
	}
	return 0
}

func (cpu *CPU) BNE() uint8 {
	if cpu.getFlag(Z) == 0 {
		cpu.cycles++
		cpu.addrAbs = cpu.pc + cpu.addrRel
		if (cpu.addrAbs & 0xFF00) != (cpu.pc & 0xFF00) {
			cpu.cycles++
		}
		cpu.pc = cpu.addrAbs
	}
	return 0
}

func (cpu *CPU) BPL() uint8 {
	if cpu.getFlag(N) == 0 {
		cpu.cycles++
		cpu.addrAbs = cpu.pc + cpu.addrRel
		if (cpu.addrAbs & 0xFF00) != (cpu.pc & 0xFF00) {
			cpu.cycles++
		}
		cpu.pc = cpu.addrAbs
	}
	return 0
}
func (cpu *CPU) BRK() uint8 {
	cpu.pc++
	cpu.setFlag(I, true)
	cpu.write(0x0100+uint16(cpu.sp), uint8((cpu.pc>>8)&0x00FF))
	cpu.sp++
	cpu.write(0x0100+uint16(cpu.sp), uint8(cpu.pc&0x00FF))
	cpu.sp++
	cpu.setFlag(B, true)
	cpu.write(0x0100+uint16(cpu.sp), cpu.status)
	cpu.sp--
	cpu.setFlag(B, false)
	cpu.pc = uint16(uint16(cpu.read(0xFFFE)) | uint16(cpu.read(0xFFFF))<<8)
	return 0
}

func (cpu *CPU) BVC() uint8 {
	if cpu.getFlag(V) == 0 {
		cpu.cycles++
		cpu.addrAbs = cpu.pc + cpu.addrRel
		if (cpu.addrAbs & 0xFF00) != (cpu.pc & 0xFF00) {
			cpu.cycles++
		}
		cpu.pc = cpu.addrAbs
	}
	return 0
}

func (cpu *CPU) BVS() uint8 {
	if cpu.getFlag(V) == 1 {
		cpu.cycles++
		cpu.addrAbs = cpu.pc + cpu.addrRel
		if (cpu.addrAbs & 0xFF00) != (cpu.pc & 0xFF00) {
			cpu.cycles++
		}
		cpu.pc = cpu.addrAbs
	}
	return 0
}

func (cpu *CPU) CLC() uint8 {
	cpu.setFlag(C, false)
	return 0
}

func (cpu *CPU) CLD() uint8 {
	cpu.setFlag(D, false)
	return 0
}

func (cpu *CPU) CLI() uint8 {
	cpu.setFlag(I, false)
	return 0
}

func (cpu *CPU) CLV() uint8 {
	cpu.setFlag(V, false)
	return 0
}

func (cpu *CPU) CMP() uint8 {
	cpu.fetch()
	temp := uint16(cpu.a) - uint16(cpu.fetched)
	cpu.setFlag(C, cpu.a >= cpu.fetched)
	cpu.setFlag(Z, (temp&0x00FF) == 0x0000)
	cpu.setFlag(N, (temp&0x0080) != 0)
	return 1
}

func (cpu *CPU) CPX() uint8 {
	cpu.fetch()
	temp := uint16(cpu.x) - uint16(cpu.fetched)
	cpu.setFlag(C, cpu.x >= cpu.fetched)
	cpu.setFlag(Z, (temp&0x00FF) == 0x0000)
	cpu.setFlag(N, (temp&0x0080) != 0)
	return 0
}

func (cpu *CPU) CPY() uint8 {
	cpu.fetch()
	temp := uint16(cpu.y) - uint16(cpu.fetched)
	cpu.setFlag(C, cpu.y >= cpu.fetched)
	cpu.setFlag(Z, (temp&0x00FF) == 0x0000)
	cpu.setFlag(N, (temp&0x0080) != 0)
	return 0
}

func (cpu *CPU) DEC() uint8 {
	cpu.fetch()
	temp := uint16(cpu.fetched) - 1
	cpu.write(cpu.addrAbs, uint8(temp&0x00FF))
	cpu.setFlag(Z, (temp&0x00FF) == 0x0000)
	cpu.setFlag(N, (temp&0x0080) != 0)
	return 0
}

func (cpu *CPU) DEX() uint8 {
	cpu.x--
	cpu.setFlag(Z, cpu.x == 0x00)
	cpu.setFlag(N, (cpu.x&0x80) != 0)
	return 0
}

func (cpu *CPU) DEY() uint8 {
	cpu.y--
	cpu.setFlag(Z, cpu.y == 0x00)
	cpu.setFlag(N, (cpu.y&0x80) != 0)
	return 0
}
func (cpu *CPU) EOR() uint8 {
	cpu.fetch()
	cpu.a = cpu.a ^ cpu.fetched
	cpu.setFlag(Z, cpu.a == 0x00)
	cpu.setFlag(N, (cpu.a&0x80) != 0)
	return 1
}

func (cpu *CPU) INC() uint8 {
	cpu.fetch()
	temp := uint16(cpu.fetched + 1)
	cpu.write(cpu.addrAbs, uint8(temp&0x00FF))
	cpu.setFlag(Z, (temp&0x00FF) == 0x0000)
	cpu.setFlag(N, (temp&0x0080) != 0)
	return 0
}

func (cpu *CPU) INX() uint8 {
	cpu.x++
	cpu.setFlag(Z, cpu.x == 0x00)
	cpu.setFlag(N, (cpu.x&0x80) != 0)
	return 0
}

func (cpu *CPU) INY() uint8 {
	cpu.y++
	cpu.setFlag(Z, cpu.y == 0x00)
	cpu.setFlag(N, (cpu.y&0x80) != 0)
	return 0
}

func (cpu *CPU) JMP() uint8 {
	cpu.pc = cpu.addrAbs
	return 0
}

func (cpu *CPU) JSR() uint8 {
	cpu.pc--
	cpu.write(0x0100+uint16(cpu.sp), uint8((cpu.pc>>8)&0x00FF))
	cpu.sp++
	cpu.write(0x0100+uint16(cpu.sp), uint8(cpu.pc&0x00FF))
	cpu.sp--
	cpu.pc = cpu.addrAbs
	return 0
}

func (cpu *CPU) LDA() uint8 {
	cpu.fetch()
	cpu.a = cpu.fetched
	cpu.setFlag(Z, cpu.a == 0x00)
	cpu.setFlag(N, (cpu.a&0x80) != 0)
	return 0
}

func (cpu *CPU) LDX() uint8 {
	cpu.fetch()
	cpu.x = cpu.fetched
	cpu.setFlag(Z, cpu.x == 0x00)
	cpu.setFlag(N, (cpu.x&0x80) != 0)
	return 1
}

func (cpu *CPU) LDY() uint8 {
	cpu.fetch()
	cpu.y = cpu.fetched
	cpu.setFlag(Z, cpu.y == 0x00)
	cpu.setFlag(N, (cpu.y&0x80) != 0)
	return 0
}

func (cpu *CPU) LSR() uint8 {
	cpu.fetch()
	cpu.setFlag(C, (cpu.fetched&0x0001) != 0)
	temp := uint16(cpu.fetched >> 1)
	cpu.setFlag(Z, (temp&0x00FF) == 0x0000)
	cpu.setFlag(N, (temp&0x0080) != 0)
	if cpu.isIMP() {
		cpu.a = uint8(temp & 0x00FF)
	} else {
		cpu.write(cpu.addrAbs, uint8(temp&0x00FF))
	}
	return 0
}

func (cpu *CPU) NOP() uint8 {
	switch cpu.opcode {
	case 0x1C:
	case 0x3C:
	case 0x5C:
	case 0x7C:
	case 0xDC:
	case 0xFC:
		return 1
	}
	return 0
}
func (cpu *CPU) ORA() uint8 {
	cpu.fetch()
	cpu.a = cpu.a | cpu.fetched
	cpu.setFlag(Z, cpu.a == 0x00)
	cpu.setFlag(N, (cpu.a&0x80) != 0)
	return 1
}

func (cpu *CPU) PHA() uint8 {
	cpu.write(0x0100+uint16(cpu.sp), cpu.a)
	cpu.sp--
	return 0
}

func (cpu *CPU) PHP() uint8 {
	cpu.write(0x0100+uint16(cpu.sp), cpu.status|B|U)
	cpu.setFlag(B, false)
	cpu.setFlag(U, false)
	cpu.sp--
	return 0
}

func (cpu *CPU) PLA() uint8 {
	cpu.sp++
	cpu.a = cpu.read(0x0100 + uint16(cpu.sp))
	cpu.setFlag(Z, cpu.a == 0x00)
	cpu.setFlag(N, (cpu.a&0x80) != 0)
	return 0
}

func (cpu *CPU) PLP() uint8 {
	cpu.sp++
	cpu.status = cpu.read(0x0100 + uint16(cpu.sp))
	cpu.setFlag(U, true)
	return 0
}

func (cpu *CPU) ROL() uint8 {
	cpu.fetch()
	temp := uint16(cpu.fetched<<1) | uint16(cpu.getFlag(C))
	cpu.setFlag(C, (temp&0xFF00) != 0)
	cpu.setFlag(Z, (temp&0x00FF) == 0x0000)
	cpu.setFlag(N, (temp&0x0080) != 0)
	if cpu.isIMP() {
		cpu.a = uint8(temp & 0x00FF)
	} else {
		cpu.write(cpu.addrAbs, uint8(temp&0x00FF))
	}
	return 0
}

func (cpu *CPU) ROR() uint8 {
	cpu.fetch()
	temp := uint16(cpu.fetched<<7) | uint16(cpu.fetched>>1)
	cpu.setFlag(C, (cpu.fetched&0x01) != 0)
	cpu.setFlag(Z, (temp&0x00FF) == 0x0000)
	cpu.setFlag(N, (temp&0x0080) != 0)
	if cpu.isIMP() {
		cpu.a = uint8(temp & 0x00FF)
	} else {
		cpu.write(cpu.addrAbs, uint8(temp&0x00FF))
	}
	return 0
}

func (cpu *CPU) RTI() uint8 {
	cpu.sp++
	cpu.status = cpu.read(0x0100 + uint16(cpu.sp))
	cpu.status &= ^B
	cpu.status &= ^U
	cpu.sp++
	cpu.pc = uint16(cpu.read(0x0100 + uint16(cpu.sp)))
	cpu.sp++
	cpu.pc |= uint16(cpu.read(0x0100+uint16(cpu.sp))) << 8
	return 0
}

func (cpu *CPU) RTS() uint8 {
	cpu.sp++
	cpu.pc = uint16(cpu.read(0x0100 + uint16(cpu.sp)))
	cpu.sp++
	cpu.pc |= uint16(cpu.read(0x0100+uint16(cpu.sp))) << 8
	cpu.pc++
	return 0
}

func (cpu *CPU) SBC() uint8 {
	cpu.fetch()
	value := uint16(cpu.fetched) ^ 0x00FF
	cpu.fetch()
	temp := uint16(cpu.a) + value + uint16(cpu.getFlag(C))
	cpu.setFlag(C, temp > 255)
	cpu.setFlag(Z, (temp&0x00FF) == 0)
	cpu.setFlag(N, (temp&0x0080) != 0)
	v := ^(uint16(cpu.a) ^ uint16(cpu.fetched)&(uint16(cpu.a)^temp)) & 0x0080
	cpu.setFlag(V, v != 0)
	cpu.a = uint8(temp & 0x00FF)
	return 1
}

func (cpu *CPU) SEC() uint8 {
	cpu.setFlag(C, true)
	return 0
}

func (cpu *CPU) SED() uint8 {
	cpu.setFlag(D, true)
	return 0
}

func (cpu *CPU) SEI() uint8 {
	cpu.setFlag(I, true)
	return 0
}

func (cpu *CPU) STA() uint8 {
	cpu.write(cpu.addrAbs, cpu.a)
	return 0
}

func (cpu *CPU) STX() uint8 {
	cpu.write(cpu.addrAbs, cpu.x)
	return 0
}

func (cpu *CPU) STY() uint8 {
	cpu.write(cpu.addrAbs, cpu.y)
	return 0
}

func (cpu *CPU) TAX() uint8 {
	cpu.x = cpu.a
	cpu.setFlag(Z, cpu.x == 0x00)
	cpu.setFlag(N, (cpu.x&0x80) != 0)
	return 0
}

func (cpu *CPU) TAY() uint8 {
	cpu.y = cpu.a
	cpu.setFlag(Z, cpu.y == 0x00)
	cpu.setFlag(N, (cpu.y&0x80) != 0)
	return 0
}

func (cpu *CPU) TSX() uint8 {
	cpu.x = cpu.sp
	cpu.setFlag(Z, cpu.x == 0x00)
	cpu.setFlag(N, (cpu.x&0x80) != 0)
	return 0
}

func (cpu *CPU) TXA() uint8 {
	cpu.a = cpu.x
	cpu.setFlag(Z, cpu.a == 0x00)
	cpu.setFlag(N, (cpu.a&0x80) != 0)
	return 0
}

func (cpu *CPU) TXS() uint8 {
	cpu.sp = cpu.x
	return 0
}

func (cpu *CPU) TYA() uint8 {
	cpu.a = cpu.y
	cpu.setFlag(Z, cpu.a == 0x00)
	cpu.setFlag(N, (cpu.a&0x80) != 0)
	return 0
}
func (cpu *CPU) XXX() uint8 { return 0 }
