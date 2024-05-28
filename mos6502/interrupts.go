package mos6502

func (cpu *CPU) Clock() {
	if cpu.cycles == 0 {
		cpu.opcode = cpu.read(cpu.pc)
		cpu.setFlag(U, true)
		cpu.pc++
		cpu.cycles = cpu.lookup[cpu.opcode].Cycles
		cycle1 := cpu.lookup[cpu.opcode].AddrMode()
		cycle2 := cpu.lookup[cpu.opcode].Operate()
		cpu.cycles += (cycle1 & cycle2)
		cpu.setFlag(U, true)
	}
	cpu.cycles--
}

func (cpu *CPU) Reset() {
	cpu.a = 0x00
	cpu.x = 0x00
	cpu.y = 0x00
	cpu.sp = 0xFD
	cpu.status = 0x00 | U
	cpu.addrAbs = 0xFFFC
	low := uint16(cpu.read(cpu.addrAbs + 0))
	high := uint16(cpu.read(cpu.addrAbs + 1))
	cpu.pc = (high << 8) | low
	cpu.addrRel = 0x0000
	cpu.addrAbs = 0x0000
	cpu.fetched = 0x00
	cpu.cycles = 8
}

func (cpu *CPU) IRQ() {
	if cpu.getFlag(I) == 0 {
		cpu.write(0x0100+uint16(cpu.sp), uint8((cpu.pc>>8)&0x00FF))
		cpu.sp--
		cpu.write(0x0100+uint16(cpu.sp), uint8(cpu.pc&0x00FF))
		cpu.sp--
		cpu.setFlag(B, false)
		cpu.setFlag(U, true)
		cpu.setFlag(I, true)
		cpu.write(0x0100+uint16(cpu.sp), cpu.status)
		cpu.sp--
		cpu.addrAbs = 0xFFFE
		low := uint16(cpu.read(cpu.addrAbs + 0))
		high := uint16(cpu.read(cpu.addrAbs + 1))
		cpu.pc = (high << 8) | low
		cpu.cycles = 7
	}
}
func (cpu *CPU) NMI() {
	cpu.write(0x0100+uint16(cpu.sp), uint8((cpu.pc>>8)&0x00FF))
	cpu.sp--
	cpu.write(0x0100+uint16(cpu.sp), uint8(cpu.pc&0x00FF))
	cpu.sp--
	cpu.setFlag(B, false)
	cpu.setFlag(U, true)
	cpu.setFlag(I, true)
	cpu.write(0x0100+uint16(cpu.sp), cpu.status)
	cpu.sp--
	cpu.addrAbs = 0xFFFE
	low := uint16(cpu.read(cpu.addrAbs + 0))
	high := uint16(cpu.read(cpu.addrAbs + 1))
	cpu.pc = (high << 8) | low
	cpu.cycles = 8
}
