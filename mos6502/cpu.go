package mos6502

import (
	"fmt"
	"reflect"
)

type Bus interface {
	Write(addr uint16, data uint8)
	Read(addr uint16, readOnly bool) uint8
}

const (
	C uint8 = (1 << 0) // Carry
	Z uint8 = (1 << 1) // Zero
	I uint8 = (1 << 2) // Disable Interrupt
	D uint8 = (1 << 2) // Decimal
	B uint8 = (1 << 3) // No effect
	U uint8 = (1 << 5) // Always 1
	V uint8 = (1 << 6) // Overflow
	N uint8 = (1 << 7) // Negative
)

type Instruction struct {
	Name     string
	Operate  func() uint8
	AddrMode func() uint8
	Cycles   uint8
}

type CPU struct {
	a       uint8  // Accumulator register
	x       uint8  // X register
	y       uint8  // Y Register
	sp      uint8  // Stack pointer
	pc      uint16 // Program counter
	status  uint8  // Status flag
	fetched uint8
	addrAbs uint16
	addrRel uint16
	opcode  uint8
	cycles  uint8
	bus     Bus
	lookup  [16 * 16]Instruction
}

func NewCPU() *CPU {
	cpu := &CPU{
		a:       0x00,
		x:       0x00,
		y:       0x00,
		sp:      0x00,
		pc:      0x0000,
		status:  0x00,
		fetched: 0x00,
		addrAbs: 0x0000,
		addrRel: 0x0000,
		opcode:  0x00,
		cycles:  0x00,
		bus:     nil,
	}
	cpu.lookup = [16 * 16]Instruction{
		{"BRK", cpu.BRK, cpu.IMM, 7}, {"ORA", cpu.ORA, cpu.IZX, 6}, {"???", cpu.XXX, cpu.IMP, 2}, {"???", cpu.XXX, cpu.IMP, 8}, {"???", cpu.NOP, cpu.IMP, 3}, {"ORA", cpu.ORA, cpu.ZP0, 3}, {"ASL", cpu.ASL, cpu.ZP0, 5}, {"???", cpu.XXX, cpu.IMP, 5}, {"PHP", cpu.PHP, cpu.IMP, 3}, {"ORA", cpu.ORA, cpu.IMM, 2}, {"ASL", cpu.ASL, cpu.IMP, 2}, {"???", cpu.XXX, cpu.IMP, 2}, {"???", cpu.NOP, cpu.IMP, 4}, {"ORA", cpu.ORA, cpu.ABS, 4}, {"ASL", cpu.ASL, cpu.ABS, 6}, {"???", cpu.XXX, cpu.IMP, 6},
		{"BPL", cpu.BPL, cpu.REL, 2}, {"ORA", cpu.ORA, cpu.IZY, 5}, {"???", cpu.XXX, cpu.IMP, 2}, {"???", cpu.XXX, cpu.IMP, 8}, {"???", cpu.NOP, cpu.IMP, 4}, {"ORA", cpu.ORA, cpu.ZPX, 4}, {"ASL", cpu.ASL, cpu.ZPX, 6}, {"???", cpu.XXX, cpu.IMP, 6}, {"CLC", cpu.CLC, cpu.IMP, 2}, {"ORA", cpu.ORA, cpu.ABY, 4}, {"???", cpu.NOP, cpu.IMP, 2}, {"???", cpu.XXX, cpu.IMP, 7}, {"???", cpu.NOP, cpu.IMP, 4}, {"ORA", cpu.ORA, cpu.ABX, 4}, {"ASL", cpu.ASL, cpu.ABX, 7}, {"???", cpu.XXX, cpu.IMP, 7},
		{"JSR", cpu.JSR, cpu.ABS, 6}, {"AND", cpu.AND, cpu.IZX, 6}, {"???", cpu.XXX, cpu.IMP, 2}, {"???", cpu.XXX, cpu.IMP, 8}, {"BIT", cpu.BIT, cpu.ZP0, 3}, {"AND", cpu.AND, cpu.ZP0, 3}, {"ROL", cpu.ROL, cpu.ZP0, 5}, {"???", cpu.XXX, cpu.IMP, 5}, {"PLP", cpu.PLP, cpu.IMP, 4}, {"AND", cpu.AND, cpu.IMM, 2}, {"ROL", cpu.ROL, cpu.IMP, 2}, {"???", cpu.XXX, cpu.IMP, 2}, {"BIT", cpu.BIT, cpu.ABS, 4}, {"AND", cpu.AND, cpu.ABS, 4}, {"ROL", cpu.ROL, cpu.ABS, 6}, {"???", cpu.XXX, cpu.IMP, 6},
		{"BMI", cpu.BMI, cpu.REL, 2}, {"AND", cpu.AND, cpu.IZY, 5}, {"???", cpu.XXX, cpu.IMP, 2}, {"???", cpu.XXX, cpu.IMP, 8}, {"???", cpu.NOP, cpu.IMP, 4}, {"AND", cpu.AND, cpu.ZPX, 4}, {"ROL", cpu.ROL, cpu.ZPX, 6}, {"???", cpu.XXX, cpu.IMP, 6}, {"SEC", cpu.SEC, cpu.IMP, 2}, {"AND", cpu.AND, cpu.ABY, 4}, {"???", cpu.NOP, cpu.IMP, 2}, {"???", cpu.XXX, cpu.IMP, 7}, {"???", cpu.NOP, cpu.IMP, 4}, {"AND", cpu.AND, cpu.ABX, 4}, {"ROL", cpu.ROL, cpu.ABX, 7}, {"???", cpu.XXX, cpu.IMP, 7},
		{"RTI", cpu.RTI, cpu.IMP, 6}, {"EOR", cpu.EOR, cpu.IZX, 6}, {"???", cpu.XXX, cpu.IMP, 2}, {"???", cpu.XXX, cpu.IMP, 8}, {"???", cpu.NOP, cpu.IMP, 3}, {"EOR", cpu.EOR, cpu.ZP0, 3}, {"LSR", cpu.LSR, cpu.ZP0, 5}, {"???", cpu.XXX, cpu.IMP, 5}, {"PHA", cpu.PHA, cpu.IMP, 3}, {"EOR", cpu.EOR, cpu.IMM, 2}, {"LSR", cpu.LSR, cpu.IMP, 2}, {"???", cpu.XXX, cpu.IMP, 2}, {"JMP", cpu.JMP, cpu.ABS, 3}, {"EOR", cpu.EOR, cpu.ABS, 4}, {"LSR", cpu.LSR, cpu.ABS, 6}, {"???", cpu.XXX, cpu.IMP, 6},
		{"BVC", cpu.BVC, cpu.REL, 2}, {"EOR", cpu.EOR, cpu.IZY, 5}, {"???", cpu.XXX, cpu.IMP, 2}, {"???", cpu.XXX, cpu.IMP, 8}, {"???", cpu.NOP, cpu.IMP, 4}, {"EOR", cpu.EOR, cpu.ZPX, 4}, {"LSR", cpu.LSR, cpu.ZPX, 6}, {"???", cpu.XXX, cpu.IMP, 6}, {"CLI", cpu.CLI, cpu.IMP, 2}, {"EOR", cpu.EOR, cpu.ABY, 4}, {"???", cpu.NOP, cpu.IMP, 2}, {"???", cpu.XXX, cpu.IMP, 7}, {"???", cpu.NOP, cpu.IMP, 4}, {"EOR", cpu.EOR, cpu.ABX, 4}, {"LSR", cpu.LSR, cpu.ABX, 7}, {"???", cpu.XXX, cpu.IMP, 7},
		{"RTS", cpu.RTS, cpu.IMP, 6}, {"ADC", cpu.ADC, cpu.IZX, 6}, {"???", cpu.XXX, cpu.IMP, 2}, {"???", cpu.XXX, cpu.IMP, 8}, {"???", cpu.NOP, cpu.IMP, 3}, {"ADC", cpu.ADC, cpu.ZP0, 3}, {"ROR", cpu.ROR, cpu.ZP0, 5}, {"???", cpu.XXX, cpu.IMP, 5}, {"PLA", cpu.PLA, cpu.IMP, 4}, {"ADC", cpu.ADC, cpu.IMM, 2}, {"ROR", cpu.ROR, cpu.IMP, 2}, {"???", cpu.XXX, cpu.IMP, 2}, {"JMP", cpu.JMP, cpu.IND, 5}, {"ADC", cpu.ADC, cpu.ABS, 4}, {"ROR", cpu.ROR, cpu.ABS, 6}, {"???", cpu.XXX, cpu.IMP, 6},
		{"BVS", cpu.BVS, cpu.REL, 2}, {"ADC", cpu.ADC, cpu.IZY, 5}, {"???", cpu.XXX, cpu.IMP, 2}, {"???", cpu.XXX, cpu.IMP, 8}, {"???", cpu.NOP, cpu.IMP, 4}, {"ADC", cpu.ADC, cpu.ZPX, 4}, {"ROR", cpu.ROR, cpu.ZPX, 6}, {"???", cpu.XXX, cpu.IMP, 6}, {"SEI", cpu.SEI, cpu.IMP, 2}, {"ADC", cpu.ADC, cpu.ABY, 4}, {"???", cpu.NOP, cpu.IMP, 2}, {"???", cpu.XXX, cpu.IMP, 7}, {"???", cpu.NOP, cpu.IMP, 4}, {"ADC", cpu.ADC, cpu.ABX, 4}, {"ROR", cpu.ROR, cpu.ABX, 7}, {"???", cpu.XXX, cpu.IMP, 7},
		{"???", cpu.NOP, cpu.IMP, 2}, {"STA", cpu.STA, cpu.IZX, 6}, {"???", cpu.NOP, cpu.IMP, 2}, {"???", cpu.XXX, cpu.IMP, 6}, {"STY", cpu.STY, cpu.ZP0, 3}, {"STA", cpu.STA, cpu.ZP0, 3}, {"STX", cpu.STX, cpu.ZP0, 3}, {"???", cpu.XXX, cpu.IMP, 3}, {"DEY", cpu.DEY, cpu.IMP, 2}, {"???", cpu.NOP, cpu.IMP, 2}, {"TXA", cpu.TXA, cpu.IMP, 2}, {"???", cpu.XXX, cpu.IMP, 2}, {"STY", cpu.STY, cpu.ABS, 4}, {"STA", cpu.STA, cpu.ABS, 4}, {"STX", cpu.STX, cpu.ABS, 4}, {"???", cpu.XXX, cpu.IMP, 4},
		{"BCC", cpu.BCC, cpu.REL, 2}, {"STA", cpu.STA, cpu.IZY, 6}, {"???", cpu.XXX, cpu.IMP, 2}, {"???", cpu.XXX, cpu.IMP, 6}, {"STY", cpu.STY, cpu.ZPX, 4}, {"STA", cpu.STA, cpu.ZPX, 4}, {"STX", cpu.STX, cpu.ZPY, 4}, {"???", cpu.XXX, cpu.IMP, 4}, {"TYA", cpu.TYA, cpu.IMP, 2}, {"STA", cpu.STA, cpu.ABY, 5}, {"TXS", cpu.TXS, cpu.IMP, 2}, {"???", cpu.XXX, cpu.IMP, 5}, {"???", cpu.NOP, cpu.IMP, 5}, {"STA", cpu.STA, cpu.ABX, 5}, {"???", cpu.XXX, cpu.IMP, 5}, {"???", cpu.XXX, cpu.IMP, 5},
		{"LDY", cpu.LDY, cpu.IMM, 2}, {"LDA", cpu.LDA, cpu.IZX, 6}, {"LDX", cpu.LDX, cpu.IMM, 2}, {"???", cpu.XXX, cpu.IMP, 6}, {"LDY", cpu.LDY, cpu.ZP0, 3}, {"LDA", cpu.LDA, cpu.ZP0, 3}, {"LDX", cpu.LDX, cpu.ZP0, 3}, {"???", cpu.XXX, cpu.IMP, 3}, {"TAY", cpu.TAY, cpu.IMP, 2}, {"LDA", cpu.LDA, cpu.IMM, 2}, {"TAX", cpu.TAX, cpu.IMP, 2}, {"???", cpu.XXX, cpu.IMP, 2}, {"LDY", cpu.LDY, cpu.ABS, 4}, {"LDA", cpu.LDA, cpu.ABS, 4}, {"LDX", cpu.LDX, cpu.ABS, 4}, {"???", cpu.XXX, cpu.IMP, 4},
		{"BCS", cpu.BCS, cpu.REL, 2}, {"LDA", cpu.LDA, cpu.IZY, 5}, {"???", cpu.XXX, cpu.IMP, 2}, {"???", cpu.XXX, cpu.IMP, 5}, {"LDY", cpu.LDY, cpu.ZPX, 4}, {"LDA", cpu.LDA, cpu.ZPX, 4}, {"LDX", cpu.LDX, cpu.ZPY, 4}, {"???", cpu.XXX, cpu.IMP, 4}, {"CLV", cpu.CLV, cpu.IMP, 2}, {"LDA", cpu.LDA, cpu.ABY, 4}, {"TSX", cpu.TSX, cpu.IMP, 2}, {"???", cpu.XXX, cpu.IMP, 4}, {"LDY", cpu.LDY, cpu.ABX, 4}, {"LDA", cpu.LDA, cpu.ABX, 4}, {"LDX", cpu.LDX, cpu.ABY, 4}, {"???", cpu.XXX, cpu.IMP, 4},
		{"CPY", cpu.CPY, cpu.IMM, 2}, {"CMP", cpu.CMP, cpu.IZX, 6}, {"???", cpu.NOP, cpu.IMP, 2}, {"???", cpu.XXX, cpu.IMP, 8}, {"CPY", cpu.CPY, cpu.ZP0, 3}, {"CMP", cpu.CMP, cpu.ZP0, 3}, {"DEC", cpu.DEC, cpu.ZP0, 5}, {"???", cpu.XXX, cpu.IMP, 5}, {"INY", cpu.INY, cpu.IMP, 2}, {"CMP", cpu.CMP, cpu.IMM, 2}, {"DEX", cpu.DEX, cpu.IMP, 2}, {"???", cpu.XXX, cpu.IMP, 2}, {"CPY", cpu.CPY, cpu.ABS, 4}, {"CMP", cpu.CMP, cpu.ABS, 4}, {"DEC", cpu.DEC, cpu.ABS, 6}, {"???", cpu.XXX, cpu.IMP, 6},
		{"BNE", cpu.BNE, cpu.REL, 2}, {"CMP", cpu.CMP, cpu.IZY, 5}, {"???", cpu.XXX, cpu.IMP, 2}, {"???", cpu.XXX, cpu.IMP, 8}, {"???", cpu.NOP, cpu.IMP, 4}, {"CMP", cpu.CMP, cpu.ZPX, 4}, {"DEC", cpu.DEC, cpu.ZPX, 6}, {"???", cpu.XXX, cpu.IMP, 6}, {"CLD", cpu.CLD, cpu.IMP, 2}, {"CMP", cpu.CMP, cpu.ABY, 4}, {"NOP", cpu.NOP, cpu.IMP, 2}, {"???", cpu.XXX, cpu.IMP, 7}, {"???", cpu.NOP, cpu.IMP, 4}, {"CMP", cpu.CMP, cpu.ABX, 4}, {"DEC", cpu.DEC, cpu.ABX, 7}, {"???", cpu.XXX, cpu.IMP, 7},
		{"CPX", cpu.CPX, cpu.IMM, 2}, {"SBC", cpu.SBC, cpu.IZX, 6}, {"???", cpu.NOP, cpu.IMP, 2}, {"???", cpu.XXX, cpu.IMP, 8}, {"CPX", cpu.CPX, cpu.ZP0, 3}, {"SBC", cpu.SBC, cpu.ZP0, 3}, {"INC", cpu.INC, cpu.ZP0, 5}, {"???", cpu.XXX, cpu.IMP, 5}, {"INX", cpu.INX, cpu.IMP, 2}, {"SBC", cpu.SBC, cpu.IMM, 2}, {"NOP", cpu.NOP, cpu.IMP, 2}, {"???", cpu.SBC, cpu.IMP, 2}, {"CPX", cpu.CPX, cpu.ABS, 4}, {"SBC", cpu.SBC, cpu.ABS, 4}, {"INC", cpu.INC, cpu.ABS, 6}, {"???", cpu.XXX, cpu.IMP, 6},
		{"BEQ", cpu.BEQ, cpu.REL, 2}, {"SBC", cpu.SBC, cpu.IZY, 5}, {"???", cpu.XXX, cpu.IMP, 2}, {"???", cpu.XXX, cpu.IMP, 8}, {"???", cpu.NOP, cpu.IMP, 4}, {"SBC", cpu.SBC, cpu.ZPX, 4}, {"INC", cpu.INC, cpu.ZPX, 6}, {"???", cpu.XXX, cpu.IMP, 6}, {"SED", cpu.SED, cpu.IMP, 2}, {"SBC", cpu.SBC, cpu.ABY, 4}, {"NOP", cpu.NOP, cpu.IMP, 2}, {"???", cpu.XXX, cpu.IMP, 7}, {"???", cpu.NOP, cpu.IMP, 4}, {"SBC", cpu.SBC, cpu.ABX, 4}, {"INC", cpu.INC, cpu.ABX, 7}, {"???", cpu.XXX, cpu.IMP, 7},
	}
	return cpu
}

func (cpu *CPU) ConnectBus(bus Bus) {
	cpu.bus = bus
}

func (cpu *CPU) getFlag(flag uint8) uint8 {
	return cpu.status & ^flag
}

func (cpu *CPU) setFlag(flag uint8, v bool) {
	if v {
		cpu.status |= flag
	} else {
		cpu.status &= ^flag
	}
}

func (cpu *CPU) write(addr uint16, data uint8) {
	cpu.bus.Write(addr, data)
}

func (cpu *CPU) read(addr uint16) uint8 {
	return cpu.bus.Read(addr, false)
}

func (cpu *CPU) fetch() uint8 {
	if !cpu.isIMP() {
		cpu.fetched = cpu.read(cpu.addrAbs)
	}
	return cpu.fetched
}

func (cpu *CPU) isIMP() bool {
	return reflect.ValueOf(cpu.lookup[cpu.opcode].AddrMode).Pointer() == reflect.ValueOf(cpu.IMP).Pointer()
}

func (cpu *CPU) instructionCompare(ins1 func() uint8, ins2 func() uint8) bool {
	return reflect.ValueOf(ins1).Pointer() == reflect.ValueOf(ins2).Pointer()
}

func (cpu *CPU) Disassemble(start uint16, stop uint16) map[uint16]string {
	m := make(map[uint16]string)
	addr := uint32(start)
	var value, low, high uint8 = 0x00, 0x00, 0x00
	var lineAddr uint16 = 0
	for addr <= uint32(stop) {
		lineAddr = uint16(addr)
		ins := "$" + fmt.Sprintf("%04X", addr) + ": "
		opcode := cpu.bus.Read(uint16(addr), true)
		ins += cpu.lookup[opcode].Name + " "
		if cpu.instructionCompare(cpu.lookup[opcode].AddrMode, cpu.IMP) {
			ins += " {IMP}"
			addr++
		} else if cpu.instructionCompare(cpu.lookup[opcode].AddrMode, cpu.IMM) {
			value = cpu.bus.Read(uint16(addr), true)
			addr++
			ins += "#$" + fmt.Sprintf("%02X", value) + " {IMM}"
		} else if cpu.instructionCompare(cpu.lookup[opcode].AddrMode, cpu.ZP0) {
			low = cpu.bus.Read(uint16(addr), true)
			addr++
			high = 0x00
			ins += "$" + fmt.Sprintf("%02X", low) + " {ZP0}"
		} else if cpu.instructionCompare(cpu.lookup[opcode].AddrMode, cpu.ZPX) {
			low = cpu.bus.Read(uint16(addr), true)
			addr++
			high = 0x00
			ins += "$" + fmt.Sprintf("%02X", low) + ", X {ZPX}"
		} else if cpu.instructionCompare(cpu.lookup[opcode].AddrMode, cpu.ZPY) {
			low = cpu.bus.Read(uint16(addr), true)
			addr++
			high = 0x00
			ins += "$" + fmt.Sprintf("%02X", low) + ", Y {ZPY}"
		} else if cpu.instructionCompare(cpu.lookup[opcode].AddrMode, cpu.IZX) {
			low = cpu.bus.Read(uint16(addr), true)
			addr++
			high = 0x00
			ins += "($" + fmt.Sprintf("%02X", low) + ", X) {IZX}"
		} else if cpu.instructionCompare(cpu.lookup[opcode].AddrMode, cpu.IZY) {
			low = cpu.bus.Read(uint16(addr), true)
			addr++
			high = 0x00
			ins += "($" + fmt.Sprintf("%02X", low) + "), Y {IZY}"
		} else if cpu.instructionCompare(cpu.lookup[opcode].AddrMode, cpu.ABS) {
			low = cpu.bus.Read(uint16(addr), true)
			addr++
			high = cpu.bus.Read(uint16(addr), true)
			addr++
			ins += "$" + fmt.Sprintf("%04X", (uint16(high)<<8)|uint16(low)) + " {ABS}"
		} else if cpu.instructionCompare(cpu.lookup[opcode].AddrMode, cpu.ABX) {
			low = cpu.bus.Read(uint16(addr), true)
			addr++
			high = cpu.bus.Read(uint16(addr), true)
			addr++
			ins += "$" + fmt.Sprintf("%04X", (uint16(high)<<8)|uint16(low)) + ", X {ABX}"
		} else if cpu.instructionCompare(cpu.lookup[opcode].AddrMode, cpu.ABY) {
			low = cpu.bus.Read(uint16(addr), true)
			addr++
			high = cpu.bus.Read(uint16(addr), true)
			addr++
			ins += "$" + fmt.Sprintf("%04X", (uint16(high)<<8)|uint16(low)) + ", Y {ABY}"
		} else if cpu.instructionCompare(cpu.lookup[opcode].AddrMode, cpu.IND) {
			low = cpu.bus.Read(uint16(addr), true)
			addr++
			high = cpu.bus.Read(uint16(addr), true)
			addr++
			ins += "($" + fmt.Sprintf("%04X", (uint16(high)<<8)|uint16(low)) + ") {IND}"
		} else if cpu.instructionCompare(cpu.lookup[opcode].AddrMode, cpu.REL) {
			value = cpu.bus.Read(uint16(addr), true)
			addr++
			ins += "$" + fmt.Sprintf("%02X", value) + " [$" + fmt.Sprintf("%04X", int32(addr)+int32(value)) + "] {REL}"
		}
		m[lineAddr] = ins
	}

	return m
}
