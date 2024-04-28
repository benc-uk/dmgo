package gameboy

import (
	"fmt"
)

var opcodes = [0x100]func(cpu *CPU){
	// NOP
	0x00: func(cpu *CPU) {
		cpu.opDebug = "NOP"
	},

	// LD BC, nn
	0x01: func(cpu *CPU) {
		nn := cpu.fetchPC16()
		cpu.opDebug = fmt.Sprintf("LD BC,%04X", nn)
		cpu.BC = nn
	},

	// INC BC
	0x03: func(cpu *CPU) {
		cpu.opDebug = "INC BC"
		cpu.BC++
	},

	// DEC BC
	0x0B: func(cpu *CPU) {
		cpu.opDebug = "DEC BC"
		cpu.BC--
	},

	// LD C, n
	0x0E: func(cpu *CPU) {
		n := cpu.fetchPC()
		cpu.opDebug = fmt.Sprintf("LD C,%X", n)
		cpu.setRegC(n)
	},

	// LD DE, nn
	0x11: func(cpu *CPU) {
		nn := cpu.fetchPC16()
		cpu.opDebug = fmt.Sprintf("LD DE,%04X", nn)
		cpu.DE = nn
	},

	// INC DE
	0x13: func(cpu *CPU) {
		cpu.opDebug = "INC DE"
		cpu.DE++
	},

	// LD A, (DE)
	0x1A: func(cpu *CPU) {
		cpu.opDebug = fmt.Sprintf("LD A,(DE:%04X)", cpu.DE)
		cpu.setRegA(cpu.mapper.Read(cpu.DE))
	},

	// DEC DE
	0x1B: func(cpu *CPU) {
		cpu.opDebug = "DEC DE"
		cpu.DE--
	},

	// LD E, n
	0x1E: func(cpu *CPU) {
		cpu.opDebug = "LD E, n"
		cpu.setRegE(cpu.fetchPC())
	},

	// LD HL, nn
	0x21: func(cpu *CPU) {
		nn := cpu.fetchPC16()
		cpu.opDebug = fmt.Sprintf("LD HL,%04X", nn)
		cpu.HL = nn
	},

	// LD (HL+), A
	0x22: func(cpu *CPU) {
		cpu.opDebug = "LD (HL+),A"
		cpu.mapper.Write(cpu.HL, cpu.getRegA())
		cpu.HL++
	},

	// INC HL
	0x23: func(cpu *CPU) {
		cpu.opDebug = "INC HL"
		cpu.HL++
	},

	// DEC HL
	0x2B: func(cpu *CPU) {
		cpu.opDebug = "DEC HL"
		cpu.HL--
	},

	// LD L, n
	0x2E: func(cpu *CPU) {
		cpu.opDebug = "LD L, n"
		cpu.setRegL(cpu.fetchPC())
	},

	// INC SP
	0x33: func(cpu *CPU) {
		cpu.opDebug = "INC SP"
		cpu.SP++
	},

	// DEC SP
	0x3B: func(cpu *CPU) {
		cpu.opDebug = "DEC SP"
		cpu.SP--
	},

	// LD A, n
	0x3E: func(cpu *CPU) {
		n := cpu.fetchPC()
		cpu.opDebug = fmt.Sprintf("LD A,%X", n)
		cpu.setRegA(n)
	},

	// LB C, B
	0x48: func(cpu *CPU) {
		cpu.opDebug = "LD C,B"
		cpu.setRegC(cpu.getRegB())
	},

	// LD E, B
	0x58: func(cpu *CPU) {
		cpu.opDebug = "LD E,B"
		cpu.setRegE(cpu.getRegB())
	},

	// LD L, B
	0x68: func(cpu *CPU) {
		cpu.opDebug = "LD L,B"
		cpu.setRegL(cpu.getRegB())
	},

	// LD A, B
	0x78: func(cpu *CPU) {
		cpu.opDebug = "LD A,B"
		cpu.setRegA(cpu.getRegB())
	},

	// XOR A, A
	0xAF: func(cpu *CPU) {
		cpu.opDebug = "XOR A,A"
		cpu.setRegA(0)
		cpu.setFlagZ(true)
		cpu.setFlagN(false)
		cpu.setFlagH(false)
		cpu.setFlagC(false)
	},

	// OR A, B
	0xB0: func(cpu *CPU) {
		cpu.opDebug = "OR A,B"
		cpu.setRegA(cpu.byteOR(cpu.getRegA(), cpu.getRegB()))
	},

	// OR A, C
	0xB1: func(cpu *CPU) {
		cpu.opDebug = "OR A,C"
		cpu.setRegA(cpu.byteOR(cpu.getRegA(), cpu.getRegC()))
	},

	// OR A, D
	0xB2: func(cpu *CPU) {
		cpu.opDebug = "OR A,D"
		cpu.setRegA(cpu.byteOR(cpu.getRegA(), cpu.getRegD()))
	},

	// OR A, E
	0xB3: func(cpu *CPU) {
		cpu.opDebug = "OR A,E"
		cpu.setRegA(cpu.byteOR(cpu.getRegA(), cpu.getRegE()))
	},

	// OR A, H
	0xB4: func(cpu *CPU) {
		cpu.opDebug = "OR A,H"
		cpu.setRegA(cpu.byteOR(cpu.getRegA(), cpu.getRegH()))
	},

	// OR A, L
	0xB5: func(cpu *CPU) {
		cpu.opDebug = "OR A,L"

		cpu.setRegA(cpu.byteOR(cpu.getRegA(), cpu.getRegL()))
	},

	// JP NZ, nn
	0xC2: func(cpu *CPU) {
		addr := cpu.fetchPC16()
		cpu.opDebug = fmt.Sprintf("JP NZ,(%04X)", addr)
		if !cpu.getFlagZ() {
			cpu.PC = addr
		}
	},

	// JP nn
	0xC3: func(cpu *CPU) {
		addr := cpu.fetchPC16()
		cpu.opDebug = fmt.Sprintf("JP (%04X)", addr)
		cpu.PC = addr
	},

	// JP C, nn
	0xDA: func(cpu *CPU) {
		addr := cpu.fetchPC16()
		cpu.opDebug = fmt.Sprintf("JP C,(%04X)", addr)
		if cpu.getFlagC() {
			cpu.PC = addr
		}
	},

	// LD (nn), A
	0xEA: func(cpu *CPU) {
		addr := cpu.fetchPC16()
		cpu.opDebug = fmt.Sprintf("LD (%04X),A", addr)
		cpu.mapper.Write(addr, cpu.getRegA())
	},

	// LD A, (nn)
	0xFA: func(cpu *CPU) {
		addr := cpu.fetchPC16()
		cpu.opDebug = fmt.Sprintf("LD A,(%04X)", addr)

		cpu.setRegA(cpu.mapper.Read(addr))
	},

	// CP n
	0xFE: func(cpu *CPU) {
		v := cpu.fetchPC()
		cpu.opDebug = fmt.Sprintf("CP A,%X", v)
		cpu.cmp(cpu.getRegA(), v)
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

func (cpu *CPU) byteOR(a, b byte) byte {
	cpu.setFlagZ(a|b == 0)
	cpu.setFlagN(false)
	cpu.setFlagH(false)
	cpu.setFlagC(false)
	return a | b
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
