package mos6502

func (cpu *CPU) GetA() uint8 {
	return cpu.a
}

func (cpu *CPU) GetX() uint8 {
	return cpu.x
}

func (cpu *CPU) GetY() uint8 {
	return cpu.y
}

func (cpu *CPU) GetPC() uint16 {
	return cpu.pc
}

func (cpu *CPU) GetSP() uint8 {
	return cpu.sp
}

func (cpu *CPU) GetStatus() uint8 {
	return cpu.status
}
