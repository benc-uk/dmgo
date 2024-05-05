package gameboy

import "log"

var opcodes = [0x100]func(cpu *CPU){
	// NOP
	0x00: func(cpu *CPU) { cpu.opDebug = "NOP" },

	// LD BC, nn
	0x01: func(cpu *CPU) { cpu.BC = cpu.fetchPC16() },

	// INC BC
	0x03: func(cpu *CPU) { cpu.BC++ },

	// INC B
	0x04: func(cpu *CPU) { cpu.setB(cpu.byteInc(cpu.B())) },

	// DEC B
	0x05: func(cpu *CPU) { cpu.setB(cpu.byteDec(cpu.B())) },

	// LD B, n
	0x06: func(cpu *CPU) { cpu.setB(cpu.fetchPC()) },

	// DEC BC
	0x0B: func(cpu *CPU) { cpu.BC-- },

	// INC C
	0x0C: func(cpu *CPU) { cpu.setC(cpu.byteInc(cpu.C())) },

	// DEC C
	0x0D: func(cpu *CPU) { cpu.setC(cpu.byteDec(cpu.C())) },

	// LD C, n
	0x0E: func(cpu *CPU) { cpu.setC(cpu.fetchPC()) },

	// LD DE, nn
	0x11: func(cpu *CPU) { cpu.DE = cpu.fetchPC16() },

	// LD (DE), A
	0x12: func(cpu *CPU) { cpu.mapper.Write(cpu.DE, cpu.A()) },

	// INC DE
	0x13: func(cpu *CPU) { cpu.DE++ },

	// DEC D
	0x15: func(cpu *CPU) { cpu.setD(cpu.byteDec(cpu.D())) },

	// LD D, n
	0x16: func(cpu *CPU) { cpu.setD(cpu.fetchPC()) },

	// RLA
	0x17: func(cpu *CPU) {
		value := cpu.A()
		var carry byte
		if cpu.getFlagC() {
			carry = 1
		}

		result := byte(value<<1) + carry

		cpu.setA(result)
		cpu.setFlagZ(false)
		cpu.setFlagN(false)
		cpu.setFlagH(false)
		cpu.setFlagC(value > 0x7F)
	},

	// JR e
	0x18: func(cpu *CPU) { cpu.PC += uint16(int8(cpu.fetchPC())) },

	// ADD HL, DE
	0x19: func(cpu *CPU) { cpu.HL = cpu.wordAdd(cpu.HL, cpu.DE) },

	// LD A, (DE)
	0x1A: func(cpu *CPU) { cpu.setA(cpu.mapper.Read(cpu.DE)) },

	// DEC DE
	0x1B: func(cpu *CPU) { cpu.DE-- },

	// INC E
	0x1C: func(cpu *CPU) { cpu.setE(cpu.byteInc(cpu.E())) },

	// DEC E
	0x1D: func(cpu *CPU) { cpu.setE(cpu.byteDec(cpu.E())) },

	// LD E, n
	0x1E: func(cpu *CPU) { cpu.setE(cpu.fetchPC()) },

	// JR NZ,e
	0x20: func(cpu *CPU) {
		e := cpu.fetchPC() // Important fetch & inc PC before the condition!
		if !cpu.getFlagZ() {
			cpu.PC += uint16(int8(e))
		}
	},

	// LD HL, nn
	0x21: func(cpu *CPU) { cpu.HL = cpu.fetchPC16() },

	// LD (HL+), A
	0x22: func(cpu *CPU) {
		cpu.mapper.Write(cpu.HL, cpu.A())
		cpu.HL++
	},

	// INC HL
	0x23: func(cpu *CPU) { cpu.HL++ },

	// INC H
	0x24: func(cpu *CPU) { cpu.setH(cpu.byteInc(cpu.H())) },

	// JR Z, e
	0x28: func(cpu *CPU) {
		e := cpu.fetchPC() // Important fetch & inc PC before the condition!
		if cpu.getFlagZ() {
			cpu.PC += uint16(int8(e))
		}
	},

	// LD A, (HL+)
	0x2A: func(cpu *CPU) {
		cpu.setA(cpu.mapper.Read(cpu.HL))
		cpu.HL++
	},

	// DEC HL
	0x2B: func(cpu *CPU) { cpu.HL-- },

	// INC L
	0x2C: func(cpu *CPU) { cpu.setL(cpu.byteInc(cpu.L())) },

	// DEC L
	0x2D: func(cpu *CPU) { cpu.setL(cpu.byteDec(cpu.L())) },

	// LD L, n
	0x2E: func(cpu *CPU) { cpu.setL(cpu.fetchPC()) },

	// CPL
	0x2F: func(cpu *CPU) {
		cpu.setA(^cpu.A())
		cpu.setFlagN(true)
		cpu.setFlagH(true)
	},

	// LD SP, nn
	0x31: func(cpu *CPU) { cpu.SP = cpu.fetchPC16() },

	// LD [HL-], A
	0x32: func(cpu *CPU) {
		cpu.mapper.Write(cpu.HL, cpu.A())
		cpu.HL--
	},

	// INC SP
	0x33: func(cpu *CPU) { cpu.SP++ },

	// INC (HL)
	0x34: func(cpu *CPU) {
		value := cpu.mapper.Read(cpu.HL)
		cpu.mapper.Write(cpu.HL, cpu.byteInc(value))
	},

	// LD (HL), n
	0x36: func(cpu *CPU) { cpu.mapper.Write(cpu.HL, cpu.fetchPC()) },

	// JR C, e
	0x38: func(cpu *CPU) {
		e := cpu.fetchPC() // Important fetch & inc PC before the condition!
		if cpu.getFlagC() {
			cpu.PC += uint16(int8(e))
		}
	},

	// DEC SP
	0x3B: func(cpu *CPU) { cpu.SP-- },

	// INC A
	0x3C: func(cpu *CPU) { cpu.setA(cpu.byteInc(cpu.A())) },

	// DEC A
	0x3D: func(cpu *CPU) { cpu.setA(cpu.byteDec(cpu.A())) },

	// LD A, n
	0x3E: func(cpu *CPU) { cpu.setA(cpu.fetchPC()) },

	// LD B, H
	0x44: func(cpu *CPU) { cpu.setB(cpu.H()) },

	// LD B, L
	0x45: func(cpu *CPU) { cpu.setB(cpu.L()) },

	// LD B, (HL)
	0x46: func(cpu *CPU) { cpu.setB(cpu.mapper.Read(cpu.HL)) },

	// LD B, A
	0x47: func(cpu *CPU) { cpu.setB(cpu.A()) },

	// LB C, B
	0x48: func(cpu *CPU) { cpu.setC(cpu.B()) },

	// LD C, H
	0x4C: func(cpu *CPU) { cpu.setC(cpu.H()) },

	// LD C, L
	0x4D: func(cpu *CPU) { cpu.setC(cpu.L()) },

	// LD C, (HL)
	0x4E: func(cpu *CPU) { cpu.setC(cpu.mapper.Read(cpu.HL)) },

	// LD C, A
	0x4F: func(cpu *CPU) { cpu.setC(cpu.A()) },

	// LD D, L
	0x55: func(cpu *CPU) { cpu.setD(cpu.L()) },

	// LD D, (HL)
	0x56: func(cpu *CPU) { cpu.setD(cpu.mapper.Read(cpu.HL)) },

	// LD D, A
	0x57: func(cpu *CPU) { cpu.setD(cpu.A()) },

	// LD E, B
	0x58: func(cpu *CPU) { cpu.setE(cpu.B()) },

	// LD E, H
	0x5C: func(cpu *CPU) { cpu.setE(cpu.H()) },

	// LD E, L
	0x5D: func(cpu *CPU) { cpu.setE(cpu.L()) },

	// LD E, (HL)
	0x5E: func(cpu *CPU) { cpu.setE(cpu.mapper.Read(cpu.HL)) },

	// LD E, A
	0x5F: func(cpu *CPU) { cpu.setE(cpu.A()) },

	// LD H, A
	0x67: func(cpu *CPU) { cpu.setH(cpu.A()) },

	// LD L, E
	0x6B: func(cpu *CPU) { cpu.setL(cpu.E()) },

	// LD L, B
	0x68: func(cpu *CPU) { cpu.setL(cpu.B()) },

	// LD L, H
	0x6C: func(cpu *CPU) { cpu.setL(cpu.H()) },

	// LD L, L
	0x6D: func(cpu *CPU) { cpu.setL(cpu.L()) },

	// LD (HL), A
	0x77: func(cpu *CPU) { cpu.mapper.Write(cpu.HL, cpu.A()) },

	// LD A, B
	0x78: func(cpu *CPU) { cpu.setA(cpu.B()) },

	// LD A, C
	0x79: func(cpu *CPU) { cpu.setA(cpu.C()) },

	// LD A, E
	0x7B: func(cpu *CPU) { cpu.setA(cpu.E()) },

	// LD A, H
	0x7C: func(cpu *CPU) { cpu.setA(cpu.H()) },

	// LD A, L
	0x7D: func(cpu *CPU) { cpu.setA(cpu.L()) },

	// LD A, (HL)
	0x7E: func(cpu *CPU) { cpu.setA(cpu.mapper.Read(cpu.HL)) },

	// ADD A, (HL)
	0x86: func(cpu *CPU) { cpu.setA(cpu.byteAdd(cpu.A(), cpu.mapper.Read(cpu.HL))) },

	// ADD A
	0x87: func(cpu *CPU) { cpu.setA(cpu.byteAdd(cpu.A(), cpu.A())) },

	// SUB A, B
	0x90: func(cpu *CPU) { cpu.setA(cpu.byteSub(cpu.A(), cpu.B())) },

	// AND C
	0xA1: func(cpu *CPU) { cpu.setA(cpu.byteAND(cpu.A(), cpu.C())) },

	// AND A
	0xA7: func(cpu *CPU) { cpu.setA(cpu.byteAND(cpu.A(), cpu.A())) },

	// XOR A, B
	0xA8: func(cpu *CPU) {
		cpu.setA(cpu.byteXOR(cpu.A(), cpu.B()))
	},

	// XOR A, C
	0xA9: func(cpu *CPU) { cpu.setA(cpu.byteXOR(cpu.A(), cpu.C())) },

	// XOR A, D
	0xAA: func(cpu *CPU) { cpu.setA(cpu.byteXOR(cpu.A(), cpu.D())) },

	// XOR A, E
	0xAB: func(cpu *CPU) { cpu.setA(cpu.byteXOR(cpu.A(), cpu.E())) },

	// XOR A, H
	0xAC: func(cpu *CPU) { cpu.setA(cpu.byteXOR(cpu.A(), cpu.H())) },

	// XOR A, L
	0xAD: func(cpu *CPU) { cpu.setA(cpu.byteXOR(cpu.A(), cpu.L())) },

	// XOR A, A
	0xAF: func(cpu *CPU) { cpu.setA(cpu.byteXOR(cpu.A(), cpu.A())) },

	// OR A, B
	0xB0: func(cpu *CPU) { cpu.setA(cpu.byteOR(cpu.A(), cpu.B())) },

	// OR A, C
	0xB1: func(cpu *CPU) { cpu.setA(cpu.byteOR(cpu.A(), cpu.C())) },

	// OR A, D
	0xB2: func(cpu *CPU) { cpu.setA(cpu.byteOR(cpu.A(), cpu.D())) },

	// OR A, E
	0xB3: func(cpu *CPU) { cpu.setA(cpu.byteOR(cpu.A(), cpu.E())) },

	// OR A, H
	0xB4: func(cpu *CPU) { cpu.setA(cpu.byteOR(cpu.A(), cpu.H())) },

	// OR A, L
	0xB5: func(cpu *CPU) { cpu.setA(cpu.byteOR(cpu.A(), cpu.L())) },

	// OR A, A
	0xB7: func(cpu *CPU) { cpu.setA(cpu.byteOR(cpu.A(), cpu.A())) },

	// CP A, [HL]
	0xBE: func(cpu *CPU) { cpu.cmp(cpu.A(), cpu.mapper.Read(cpu.HL)) },

	// RET NZ
	0xC0: func(cpu *CPU) {
		if !cpu.getFlagZ() {
			cpu.returnSub()
		}
	},

	// POP BC
	0xC1: func(cpu *CPU) { cpu.BC = cpu.popStack() },

	// JP NZ, nn
	0xC2: func(cpu *CPU) {
		if !cpu.getFlagZ() {
			cpu.PC = cpu.fetchPC16()
		}
	},

	// JP nn
	0xC3: func(cpu *CPU) { cpu.PC = cpu.fetchPC16() },

	// CALL NZ, nn
	0xC4: func(cpu *CPU) {
		nn := cpu.fetchPC16()
		if !cpu.getFlagZ() {
			cpu.callSub(nn)
		}
	},

	// PUSH BC
	0xC5: func(cpu *CPU) { cpu.pushStack(cpu.BC) },

	// ADD A, n
	0xC6: func(cpu *CPU) { cpu.setA(cpu.byteAdd(cpu.A(), cpu.fetchPC())) },

	// RET Z
	0xC8: func(cpu *CPU) {
		if cpu.getFlagZ() {
			cpu.returnSub()
		}
	},

	// RET
	0xC9: func(cpu *CPU) { cpu.returnSub() },

	// JP Z, nn
	0xCA: func(cpu *CPU) {
		nn := cpu.fetchPC16()
		if cpu.getFlagZ() {
			cpu.PC = nn
		}
	},

	// 0xCB
	0xCB: func(cpu *CPU) {
		opcode := cpu.fetchPC()

		if cbOpcodes[opcode] == nil {
			log.Fatalf(" !!! Unknown CB opcode: %02X", opcode)
			return
		}

		cbOpcodes[opcode](cpu)
	},

	// CALL nn
	0xCD: func(cpu *CPU) { cpu.callSub(cpu.fetchPC16()) },

	// POP DE
	0xD1: func(cpu *CPU) { cpu.DE = cpu.popStack() },

	// PUSH DE
	0xD5: func(cpu *CPU) { cpu.pushStack(cpu.DE) },

	// SUB A, n
	0xD6: func(cpu *CPU) {
		result := cpu.byteSub(cpu.A(), cpu.fetchPC())
		cpu.setA(result)
		cpu.setFlagC(result > cpu.A())
	},

	// RETI
	0xD9: func(cpu *CPU) {
		cpu.returnSub()
		cpu.IME = true
	},

	// JP C, nn
	0xDA: func(cpu *CPU) {
		if cpu.getFlagC() {
			cpu.PC = cpu.fetchPC16()
		}
	},

	// LDH (n), A
	0xE0: func(cpu *CPU) {
		cpu.mapper.Write(0xFF00+uint16(cpu.fetchPC()), cpu.A())
	},

	// POP HL
	0xE1: func(cpu *CPU) { cpu.HL = cpu.popStack() },

	// LD (C), A
	0xE2: func(cpu *CPU) { cpu.mapper.Write(0xFF00+uint16(cpu.C()), cpu.A()) },

	// PUSH HL
	0xE5: func(cpu *CPU) { cpu.pushStack(cpu.HL) },

	// AND A, n
	0xE6: func(cpu *CPU) { cpu.setA(cpu.byteAND(cpu.A(), cpu.fetchPC())) },

	// JP HL
	0xE9: func(cpu *CPU) { cpu.PC = cpu.HL },

	// LD (nn), A
	0xEA: func(cpu *CPU) { cpu.mapper.Write(cpu.fetchPC16(), cpu.A()) },

	// RST 28H
	0xEF: func(cpu *CPU) { cpu.callSub(0x0028) },

	// LDH A, (n)
	0xF0: func(cpu *CPU) { cpu.setA(cpu.mapper.Read(0xFF00 + uint16(cpu.fetchPC()))) },

	// POP AF
	0xF1: func(cpu *CPU) { cpu.AF = cpu.popStack() },

	// DI
	// TODO: Needs more implementation??
	0xF3: func(cpu *CPU) {
		cpu.IME = false
	},

	// PUSH AF
	0xF5: func(cpu *CPU) { cpu.pushStack(cpu.AF) },

	// LD A, (nn)
	0xFA: func(cpu *CPU) { cpu.setA(cpu.mapper.Read(cpu.fetchPC16())) },

	// EI
	0xFB: func(cpu *CPU) {
		// TODO: Needs more implementation??
		cpu.IME = true
	},

	// CP
	0xFE: func(cpu *CPU) { cpu.cmp(cpu.A(), cpu.fetchPC()) },
}
