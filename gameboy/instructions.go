package gameboy

import (
	"fmt"
)

var opcodes = [0x100]func(cpu *CPU){
	// NOP
	0x00: func(cpu *CPU) {
		cpu.opDebug = "NOP"
	},

	// JP nn
	0xC3: func(cpu *CPU) {
		cpu.opDebug = "JP nn"
		cpu.PC = cpu.fetchPC16()
	},

	// LD (nn), A
	0xEA: func(cpu *CPU) {
		addr := cpu.fetchPC16()
		cpu.opDebug = fmt.Sprintf("LD *%x, A", addr)
		cpu.mapper.Write(addr, cpu.getRegA())
	},

	// LD A, (nn)
	0xFA: func(cpu *CPU) {
		addr := cpu.fetchPC16()
		cpu.opDebug = fmt.Sprintf("LD A, *%x", addr)

		cpu.setRegA(cpu.mapper.Read(addr))
	},

	// CP n
	0xFE: func(cpu *CPU) {
		cpu.opDebug = "CP n"
		cpu.cmp(cpu.getRegA(), cpu.fetchPC())
	},

	// LD C, n
	0x0E: func(cpu *CPU) {
		cpu.opDebug = "LD C, n"
		cpu.setRegC(cpu.fetchPC())
	},

	// LD E, n
	0x1E: func(cpu *CPU) {
		cpu.opDebug = "LD E, n"
		cpu.setRegE(cpu.fetchPC())
	},

	// LD L, n
	0x2E: func(cpu *CPU) {
		cpu.opDebug = "LD L, n"
		cpu.setRegL(cpu.fetchPC())
	},

	// LD A, n
	0x3E: func(cpu *CPU) {
		n := cpu.fetchPC()
		cpu.opDebug = fmt.Sprintf("LD A, %x", n)
		cpu.setRegA(n)
	},

	// JP C, nn
	0xDA: func(cpu *CPU) {
		addr := cpu.fetchPC16()
		cpu.opDebug = fmt.Sprintf("JP C, *%04X", addr)
		if cpu.getFlagC() {
			cpu.PC = addr
		}
	},
}

// =======================================
// Helpers
// =======================================

func (cpu *CPU) fetchPC() byte {
	v := cpu.mapper.Read(cpu.PC)
	cpu.PC += 1
	return v
}

func (cpu *CPU) fetchPC16() uint16 {
	lo := cpu.mapper.Read(cpu.PC)
	hi := cpu.mapper.Read(cpu.PC + 1)
	cpu.PC += 2
	return uint16(hi)<<8 | uint16(lo)
}

func (cpu *CPU) cmp(a, b byte) {
	cpu.setFlagZ(a == b)
	cpu.setFlagN(true)
	cpu.setFlagH((a & 0xf) < (b & 0xf))
	cpu.setFlagC(a < b)
}
