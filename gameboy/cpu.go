package gameboy

import (
	"log"
)

type CPU struct {
	// Registers
	AF, BC, DE, HL uint16

	// Stack pointer
	SP uint16

	// Program counter
	PC uint16

	// Memory
	mapper *Mapper

	opDebug string

	// Interrupts
	IME bool
}

func NewCPU(mapper *Mapper) *CPU {
	// Initial state of the CPU for the classic GB
	// It represents the state of the CPU after the BIOS has run, as we skip that
	cpu := CPU{
		AF:     0x01b0,
		BC:     0,
		DE:     0xff56,
		HL:     0x000d,
		SP:     0xfffe,
		PC:     0x100,
		mapper: mapper,
	}

	cpu.setFlagZ(true)
	return &cpu
}

func (cpu *CPU) ExecuteNext() bool {
	// Fetch the next instruction
	opcode := cpu.fetchPC()

	if opcodes[opcode] == nil {
		log.Printf("Unknown opcode: 0x%02X\n", opcode)
		cpu.PC--
		return false
	}

	// Decode & execute the opcode
	opcodes[opcode](cpu)

	cpu.logMessage(cpu.opDebug)

	return true
}

func (cpu *CPU) setFlagZ(value bool) {
	if value {
		cpu.AF |= 0x80
	} else {
		cpu.AF &^= 0x80
	}
}

func (cpu *CPU) setFlagN(value bool) {
	if value {
		cpu.AF |= 0x40
	} else {
		cpu.AF &^= 0x40
	}
}

func (cpu *CPU) setFlagH(value bool) {
	if value {
		cpu.AF |= 0x20
	} else {
		cpu.AF &^= 0x20
	}
}

func (cpu *CPU) setFlagC(value bool) {
	// In bit 4
	if value {
		cpu.AF |= 0x10
	} else {
		cpu.AF &^= 0x10
	}
}

func (cpu *CPU) getFlagZ() bool {
	return cpu.AF&0x80 != 0
}

func (cpu *CPU) getFlagN() bool {
	return cpu.AF&0x40 != 0
}

func (cpu *CPU) getFlagH() bool {
	return cpu.AF&0x20 != 0
}

func (cpu *CPU) getFlagC() bool {
	return cpu.AF&0x10 != 0
}

func (cpu *CPU) setRegA(value byte) {
	setHighByte(&cpu.AF, value)
}

func (cpu *CPU) setRegB(value byte) {
	setHighByte(&cpu.BC, value)
}

func (cpu *CPU) setRegC(value byte) {
	setLowByte(&cpu.BC, value)
}

func (cpu *CPU) setRegD(value byte) {
	setHighByte(&cpu.DE, value)
}

func (cpu *CPU) setRegE(value byte) {
	setLowByte(&cpu.DE, value)
}

func (cpu *CPU) setRegH(value byte) {
	setHighByte(&cpu.HL, value)
}

func (cpu *CPU) setRegL(value byte) {
	setLowByte(&cpu.HL, value)
}

func (cpu *CPU) getRegA() byte {
	return getHighByte(cpu.AF)
}

func (cpu *CPU) getRegB() byte {
	return getHighByte(cpu.BC)
}

func (cpu *CPU) getRegC() byte {
	return getLowByte(cpu.BC)
}

func (cpu *CPU) getRegD() byte {
	return getHighByte(cpu.DE)
}

func (cpu *CPU) getRegE() byte {
	return getLowByte(cpu.DE)
}

func (cpu *CPU) getRegH() byte {
	return getHighByte(cpu.HL)
}

func (cpu *CPU) getRegL() byte {
	return getLowByte(cpu.HL)
}

// =======================================
// Helpers
// =======================================

func setHighByte(reg *uint16, value byte) {
	*reg = uint16(value)<<8 | *reg&0xff
}

func setLowByte(reg *uint16, value byte) {
	*reg = uint16(value) | *reg&0xff00
}

func (cpu *CPU) logMessage(s string, a ...any) {
	if logging {
		log.Printf(s, a...)
	}
}

func getHighByte(reg uint16) byte {
	return byte(reg >> 8)
}

func getLowByte(reg uint16) byte {
	return byte(reg & 0xff)
}
