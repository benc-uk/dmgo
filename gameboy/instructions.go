package gameboy

import (
	"fmt"
	"log"
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

	// INC B
	0x04: func(cpu *CPU) {
		cpu.opDebug = "INC B"
		cpu.setRegB(cpu.getRegB() + 1)
	},

	// DEC B
	0x05: func(cpu *CPU) {
		cpu.opDebug = "DEC B"
		cpu.setRegB(cpu.getRegB() - 1)
	},

	// LD B, n
	0x06: func(cpu *CPU) {
		n := cpu.fetchPC()
		cpu.opDebug = fmt.Sprintf("LD B,x%X", n)
		cpu.setRegB(n)
	},

	// DEC BC
	0x0B: func(cpu *CPU) {
		cpu.opDebug = "DEC BC"
		cpu.BC--
	},

	// INC C
	0x0C: func(cpu *CPU) {
		cpu.opDebug = "INC C"
		cpu.setRegC(cpu.getRegC() + 1)
	},

	// DEC C
	0x0D: func(cpu *CPU) {
		cpu.opDebug = "DEC C"
		cpu.setRegC(cpu.getRegC() - 1)
	},

	// LD C, n
	0x0E: func(cpu *CPU) {
		n := cpu.fetchPC()
		cpu.opDebug = fmt.Sprintf("LD C,x%X", n)
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

	// JR e
	0x18: func(cpu *CPU) {
		n := int8(cpu.fetchPC())
		cpu.opDebug = fmt.Sprintf("JR x%X", n)
		cpu.PC += uint16(n)
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

	// JR NZ,e
	0x20: func(cpu *CPU) {
		n := int8(cpu.fetchPC())
		cpu.opDebug = fmt.Sprintf("JR NZ,x%X", n)
		if !cpu.getFlagZ() {
			cpu.PC += uint16(n)
		}
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

	// LD A, (HL+)
	0x2A: func(cpu *CPU) {
		cpu.opDebug = "LD A,(HL+)"
		cpu.setRegA(cpu.mapper.Read(cpu.HL))
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

	// LD SP, nn
	0x31: func(cpu *CPU) {
		cpu.SP = cpu.fetchPC16()
		cpu.opDebug = fmt.Sprintf("LD SP,%04X", cpu.SP)
	},

	// LD [HL-], A
	0x32: func(cpu *CPU) {
		cpu.opDebug = "LD (HL-),A"
		cpu.mapper.Write(cpu.HL, cpu.getRegA())
		cpu.HL--
	},

	// INC SP
	0x33: func(cpu *CPU) {
		cpu.opDebug = "INC SP"
		cpu.SP++
	},

	// LD (HL), n
	0x36: func(cpu *CPU) {
		n := cpu.fetchPC()
		cpu.opDebug = fmt.Sprintf("LD (HL),x%X", n)
		cpu.mapper.Write(cpu.HL, n)
	},

	// DEC SP
	0x3B: func(cpu *CPU) {
		cpu.opDebug = "DEC SP"
		cpu.SP--
	},

	// LD A, n
	0x3E: func(cpu *CPU) {
		n := cpu.fetchPC()
		cpu.opDebug = fmt.Sprintf("LD A,x%X", n)
		cpu.setRegA(n)
	},

	// LB C, B
	0x48: func(cpu *CPU) {
		cpu.opDebug = "LD C,B"
		cpu.setRegC(cpu.getRegB())
	},

	// LD C, H
	0x4C: func(cpu *CPU) {
		cpu.opDebug = "LD C,H"
		cpu.setRegC(cpu.getRegH())
	},

	// LD C, L
	0x4D: func(cpu *CPU) {
		cpu.opDebug = "LD C,L"
		cpu.setRegC(cpu.getRegL())
	},

	// LD E, B
	0x58: func(cpu *CPU) {
		cpu.opDebug = "LD E,B"
		cpu.setRegE(cpu.getRegB())
	},

	// LD E, H
	0x5C: func(cpu *CPU) {
		cpu.opDebug = "LD E,H"
		cpu.setRegE(cpu.getRegH())
	},

	// LD E, L
	0x5D: func(cpu *CPU) {
		cpu.opDebug = "LD E,L"
		cpu.setRegE(cpu.getRegL())
	},

	// LD L, B
	0x68: func(cpu *CPU) {
		cpu.opDebug = "LD L,B"
		cpu.setRegL(cpu.getRegB())
	},

	// LD L, H
	0x6C: func(cpu *CPU) {
		cpu.opDebug = "LD L,H"
		cpu.setRegL(cpu.getRegH())
	},

	0x6D: func(cpu *CPU) {
		cpu.opDebug = "LD L,L"
		cpu.setRegL(cpu.getRegL())
	},

	// LD A, B
	0x78: func(cpu *CPU) {
		cpu.opDebug = "LD A,B"
		cpu.setRegA(cpu.getRegB())
	},

	// LD A, H
	0x7C: func(cpu *CPU) {
		cpu.opDebug = "LD A,H"
		cpu.setRegA(cpu.getRegH())
	},

	// LD A, L
	0x7D: func(cpu *CPU) {
		cpu.opDebug = "LD A,L"
		cpu.setRegA(cpu.getRegL())
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

	// c9 is the opcode for RET
	0xC9: func(cpu *CPU) {
		cpu.opDebug = "RET"
		cpu.returnSub()
	},

	// CALL nn
	0xCD: func(cpu *CPU) {
		addr := cpu.fetchPC16()
		log.Printf(">> current PC is %04X", cpu.PC)
		cpu.opDebug = fmt.Sprintf("CALL (%04X)", addr)
		cpu.callSub(addr)
	},

	// JP C, nn
	0xDA: func(cpu *CPU) {
		addr := cpu.fetchPC16()
		cpu.opDebug = fmt.Sprintf("JP C,(%04X)", addr)
		if cpu.getFlagC() {
			cpu.PC = addr
		}
	},

	// LDH (n), A
	0xE0: func(cpu *CPU) {
		addr := 0xFF00 + uint16(cpu.fetchPC())
		cpu.opDebug = fmt.Sprintf("LDH (%02X),A", addr)
		cpu.mapper.Write(addr, cpu.getRegA())
	},

	// LD (C), A
	0xE2: func(cpu *CPU) {
		cpu.opDebug = "LD (C),A"
		cpu.mapper.Write(0xFF00+uint16(cpu.getRegC()), cpu.getRegA())
	},

	// LD (nn), A
	0xEA: func(cpu *CPU) {
		addr := cpu.fetchPC16()
		cpu.opDebug = fmt.Sprintf("LD (%04X),A", addr)
		cpu.mapper.Write(addr, cpu.getRegA())
	},

	// LDH A, (n)
	0xF0: func(cpu *CPU) {
		addr := 0xFF00 + uint16(cpu.fetchPC())
		cpu.opDebug = fmt.Sprintf("LDH A,(%02X)", addr)
		cpu.setRegA(cpu.mapper.Read(addr))
	},

	// DI
	0xF3: func(cpu *CPU) {
		cpu.opDebug = "DI"
		cpu.IME = false
	},

	// LD A, (nn)
	0xFA: func(cpu *CPU) {
		addr := cpu.fetchPC16()
		cpu.opDebug = fmt.Sprintf("LD A,(%04X)", addr)

		cpu.setRegA(cpu.mapper.Read(addr))
	},

	// CP n
	0xFE: func(cpu *CPU) {
		n := cpu.fetchPC()
		cpu.opDebug = fmt.Sprintf("CP A,x%X", n)
		cpu.cmp(cpu.getRegA(), n)
	},
}
