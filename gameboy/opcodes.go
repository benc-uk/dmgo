package gameboy

import "log"

var opcodes = [0x100]func(cpu *CPU){
	// NOP
	0x00: func(cpu *CPU) {},

	// LD BC, nn
	0x01: func(cpu *CPU) { cpu.bc = cpu.fetchPC16() },

	// LD (BC), A
	0x02: func(cpu *CPU) { cpu.mapper.write(cpu.bc, cpu.A()) },

	// INC BC
	0x03: func(cpu *CPU) { cpu.bc++ },

	// INC B
	0x04: func(cpu *CPU) { cpu.setB(cpu.byteInc(cpu.B())) },

	// DEC B
	0x05: func(cpu *CPU) { cpu.setB(cpu.byteDec(cpu.B())) },

	// LD B, n
	0x06: func(cpu *CPU) { cpu.setB(cpu.fetchPC()) },

	// RLCA
	0x07: func(cpu *CPU) {
		value := cpu.A()
		result := byte(value<<1) | byte(value>>7)

		cpu.setA(result)
		cpu.setFlagZ(false)
		cpu.setFlagN(false)
		cpu.setFlagH(false)
		cpu.setFlagC(value > 0x7F)
	},

	// LD (nn), SP
	0x08: func(cpu *CPU) {
		addr := cpu.fetchPC16()
		cpu.mapper.write(addr, byte(cpu.sp))
		cpu.mapper.write(addr+1, byte(cpu.sp>>8))
	},

	// LD HL, BC
	0x09: func(cpu *CPU) { cpu.hl = cpu.wordAdd(cpu.hl, cpu.bc) },

	// LD A, (BC)
	0x0A: func(cpu *CPU) { cpu.setA(cpu.mapper.read(cpu.bc)) },

	// DEC BC
	0x0B: func(cpu *CPU) { cpu.bc-- },

	// INC C
	0x0C: func(cpu *CPU) { cpu.setC(cpu.byteInc(cpu.C())) },

	// DEC C
	0x0D: func(cpu *CPU) { cpu.setC(cpu.byteDec(cpu.C())) },

	// LD C, n
	0x0E: func(cpu *CPU) { cpu.setC(cpu.fetchPC()) },

	// RRCA
	0x0F: func(cpu *CPU) {
		value := cpu.A()
		result := byte(value>>1) | byte(value<<7)

		cpu.setA(result)
		cpu.setFlagZ(false)
		cpu.setFlagN(false)
		cpu.setFlagH(false)
		cpu.setFlagC(value&1 == 1)
	},

	// STOP
	0x10: func(cpu *CPU) { log.Fatalf("STOP instruction not implemented") },

	// LD DE, nn
	0x11: func(cpu *CPU) { cpu.de = cpu.fetchPC16() },

	// LD (DE), A
	0x12: func(cpu *CPU) { cpu.mapper.write(cpu.de, cpu.A()) },

	// INC DE
	0x13: func(cpu *CPU) { cpu.de++ },

	// INC D
	0x14: func(cpu *CPU) { cpu.setD(cpu.byteInc(cpu.D())) },

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
	0x18: func(cpu *CPU) { cpu.pc += uint16(int8(cpu.fetchPC())) },

	// ADD HL, DE
	0x19: func(cpu *CPU) { cpu.hl = cpu.wordAdd(cpu.hl, cpu.de) },

	// LD A, (DE)
	0x1A: func(cpu *CPU) { cpu.setA(cpu.mapper.read(cpu.de)) },

	// DEC DE
	0x1B: func(cpu *CPU) { cpu.de-- },

	// INC E
	0x1C: func(cpu *CPU) { cpu.setE(cpu.byteInc(cpu.E())) },

	// DEC E
	0x1D: func(cpu *CPU) { cpu.setE(cpu.byteDec(cpu.E())) },

	// LD E, n
	0x1E: func(cpu *CPU) { cpu.setE(cpu.fetchPC()) },

	// RRA
	0x1F: func(cpu *CPU) {
		value := cpu.A()
		var carry byte
		if cpu.getFlagC() {
			carry = 0x80
		}

		result := byte(value>>1) + carry

		cpu.setA(result)
		cpu.setFlagZ(false)
		cpu.setFlagN(false)
		cpu.setFlagH(false)
		cpu.setFlagC(value&1 == 1)
	},

	// JR NZ,e
	0x20: func(cpu *CPU) {
		e := cpu.fetchPC() // Important fetch & inc PC before the condition!
		if !cpu.getFlagZ() {
			cpu.pc += uint16(int8(e))
		}
	},

	// LD HL, nn
	0x21: func(cpu *CPU) { cpu.hl = cpu.fetchPC16() },

	// LD (HL+), A
	0x22: func(cpu *CPU) {
		cpu.mapper.write(cpu.hl, cpu.A())
		cpu.hl++
	},

	// INC HL
	0x23: func(cpu *CPU) { cpu.hl++ },

	// INC H
	0x24: func(cpu *CPU) { cpu.setH(cpu.byteInc(cpu.H())) },

	// DEC H
	0x25: func(cpu *CPU) { cpu.setH(cpu.byteDec(cpu.H())) },

	// LD H, n
	0x26: func(cpu *CPU) { cpu.setH(cpu.fetchPC()) },

	// DAA
	0x27: func(cpu *CPU) {
		a := cpu.A()
		c := cpu.getFlagC()
		h := cpu.getFlagH()
		n := cpu.getFlagN()

		if !n {
			if c || a > 0x99 {
				a += 0x60
				cpu.setFlagC(true)
			}

			if h || a&0x0F > 0x09 {
				a += 0x06
			}
		} else {
			if c {
				a -= 0x60
			}

			if h {
				a -= 0x06
			}
		}

		cpu.setA(a)
		cpu.setFlagZ(a == 0)
		cpu.setFlagH(false)
	},

	// JR Z, e
	0x28: func(cpu *CPU) {
		e := cpu.fetchPC() // Important fetch & inc PC before the condition!
		if cpu.getFlagZ() {
			cpu.pc += uint16(int8(e))
		}
	},

	// ADD HL, HL
	0x29: func(cpu *CPU) { cpu.hl = cpu.wordAdd(cpu.hl, cpu.hl) },

	// LD A, (HL+)
	0x2A: func(cpu *CPU) {
		cpu.setA(cpu.mapper.read(cpu.hl))
		cpu.hl++
	},

	// DEC HL
	0x2B: func(cpu *CPU) { cpu.hl-- },

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

	// JR NC, e
	0x30: func(cpu *CPU) {
		e := cpu.fetchPC() // Important fetch & inc PC before the condition!
		if !cpu.getFlagC() {
			cpu.pc += uint16(int8(e))
		}
	},

	// LD SP, nn
	0x31: func(cpu *CPU) { cpu.sp = cpu.fetchPC16() },

	// LD [HL-], A
	0x32: func(cpu *CPU) {
		cpu.mapper.write(cpu.hl, cpu.A())
		cpu.hl--
	},

	// INC SP
	0x33: func(cpu *CPU) { cpu.sp++ },

	// INC (HL)
	0x34: func(cpu *CPU) {
		value := cpu.mapper.read(cpu.hl)
		cpu.mapper.write(cpu.hl, cpu.byteInc(value))
	},

	// DEC (HL)
	0x35: func(cpu *CPU) {
		value := cpu.mapper.read(cpu.hl)
		cpu.mapper.write(cpu.hl, cpu.byteDec(value))
	},

	// LD (HL), n
	0x36: func(cpu *CPU) { cpu.mapper.write(cpu.hl, cpu.fetchPC()) },

	// SCF
	0x37: func(cpu *CPU) {
		cpu.setFlagN(false)
		cpu.setFlagH(false)
		cpu.setFlagC(true)
	},

	// JR C, e
	0x38: func(cpu *CPU) {
		e := cpu.fetchPC() // Important fetch & inc PC before the condition!
		if cpu.getFlagC() {
			cpu.pc += uint16(int8(e))
		}
	},

	// ADD HL, SP
	0x39: func(cpu *CPU) { cpu.hl = cpu.wordAdd(cpu.hl, cpu.sp) },

	// LD A, (HL-)
	0x3A: func(cpu *CPU) {
		cpu.setA(cpu.mapper.read(cpu.hl))
		cpu.hl--
	},

	// DEC SP
	0x3B: func(cpu *CPU) { cpu.sp-- },

	// INC A
	0x3C: func(cpu *CPU) { cpu.setA(cpu.byteInc(cpu.A())) },

	// DEC A
	0x3D: func(cpu *CPU) { cpu.setA(cpu.byteDec(cpu.A())) },

	// LD A, n
	0x3E: func(cpu *CPU) { cpu.setA(cpu.fetchPC()) },

	// CCF
	0x3F: func(cpu *CPU) {
		cpu.setFlagN(false)
		cpu.setFlagH(false)
		cpu.setFlagC(!cpu.getFlagC())
	},

	// LD, B, B
	0x40: func(cpu *CPU) { cpu.setB(cpu.B()) },

	// LD B, C
	0x41: func(cpu *CPU) { cpu.setB(cpu.C()) },

	// LD B, D
	0x42: func(cpu *CPU) { cpu.setB(cpu.D()) },

	// LD B, E
	0x43: func(cpu *CPU) { cpu.setB(cpu.E()) },

	// LD B, H
	0x44: func(cpu *CPU) { cpu.setB(cpu.H()) },

	// LD B, L
	0x45: func(cpu *CPU) { cpu.setB(cpu.L()) },

	// LD B, (HL)
	0x46: func(cpu *CPU) { cpu.setB(cpu.mapper.read(cpu.hl)) },

	// LD B, A
	0x47: func(cpu *CPU) { cpu.setB(cpu.A()) },

	// LB C, B
	0x48: func(cpu *CPU) { cpu.setC(cpu.B()) },

	// LD C, C
	0x49: func(cpu *CPU) { cpu.setC(cpu.C()) },

	// LD C, D
	0x4A: func(cpu *CPU) { cpu.setC(cpu.D()) },

	// LD C, E
	0x4B: func(cpu *CPU) { cpu.setC(cpu.E()) },

	// LD C, H
	0x4C: func(cpu *CPU) { cpu.setC(cpu.H()) },

	// LD C, L
	0x4D: func(cpu *CPU) { cpu.setC(cpu.L()) },

	// LD C, (HL)
	0x4E: func(cpu *CPU) { cpu.setC(cpu.mapper.read(cpu.hl)) },

	// LD C, A
	0x4F: func(cpu *CPU) { cpu.setC(cpu.A()) },

	// LD D, B
	0x50: func(cpu *CPU) { cpu.setD(cpu.B()) },

	// LD D, C
	0x51: func(cpu *CPU) { cpu.setD(cpu.C()) },

	// LD D, D
	0x52: func(cpu *CPU) { cpu.setD(cpu.D()) },

	// LD D, E
	0x53: func(cpu *CPU) { cpu.setD(cpu.E()) },

	// LD D, H
	0x54: func(cpu *CPU) { cpu.setD(cpu.H()) },

	// LD D, L
	0x55: func(cpu *CPU) { cpu.setD(cpu.L()) },

	// LD D, (HL)
	0x56: func(cpu *CPU) { cpu.setD(cpu.mapper.read(cpu.hl)) },

	// LD D, A
	0x57: func(cpu *CPU) { cpu.setD(cpu.A()) },

	// LD E, B
	0x58: func(cpu *CPU) { cpu.setE(cpu.B()) },

	// LD E, C
	0x59: func(cpu *CPU) { cpu.setE(cpu.C()) },

	// LD E, D
	0x5A: func(cpu *CPU) { cpu.setE(cpu.D()) },

	// LD E, E
	0x5B: func(cpu *CPU) { cpu.setE(cpu.E()) },

	// LD E, H
	0x5C: func(cpu *CPU) { cpu.setE(cpu.H()) },

	// LD E, L
	0x5D: func(cpu *CPU) { cpu.setE(cpu.L()) },

	// LD E, (HL)
	0x5E: func(cpu *CPU) { cpu.setE(cpu.mapper.read(cpu.hl)) },

	// LD E, A
	0x5F: func(cpu *CPU) { cpu.setE(cpu.A()) },

	// LD H, B
	0x60: func(cpu *CPU) { cpu.setH(cpu.B()) },

	// LD H, C
	0x61: func(cpu *CPU) { cpu.setH(cpu.C()) },

	// LD H, D
	0x62: func(cpu *CPU) { cpu.setH(cpu.D()) },

	// LD H, E
	0x63: func(cpu *CPU) { cpu.setH(cpu.E()) },

	// LD H, H
	0x64: func(cpu *CPU) { cpu.setH(cpu.H()) },

	// LD H, L
	0x65: func(cpu *CPU) { cpu.setH(cpu.L()) },

	// LD H, (HL)
	0x66: func(cpu *CPU) { cpu.setH(cpu.mapper.read(cpu.hl)) },

	// LD H, A
	0x67: func(cpu *CPU) { cpu.setH(cpu.A()) },

	// LD L, B
	0x68: func(cpu *CPU) { cpu.setL(cpu.B()) },

	// LD L, C
	0x69: func(cpu *CPU) { cpu.setL(cpu.C()) },

	// LD L, D
	0x6A: func(cpu *CPU) { cpu.setL(cpu.D()) },

	// LD L, E
	0x6B: func(cpu *CPU) { cpu.setL(cpu.E()) },

	// LD L, H
	0x6C: func(cpu *CPU) { cpu.setL(cpu.H()) },

	// LD L, L
	0x6D: func(cpu *CPU) { cpu.setL(cpu.L()) },

	// LD L, (HL)
	0x6E: func(cpu *CPU) { cpu.setL(cpu.mapper.read(cpu.hl)) },

	// LD L, A
	0x6F: func(cpu *CPU) { cpu.setL(cpu.A()) },

	// LD (HL), B
	0x70: func(cpu *CPU) { cpu.mapper.write(cpu.hl, cpu.B()) },

	// LD (HL), C
	0x71: func(cpu *CPU) { cpu.mapper.write(cpu.hl, cpu.C()) },

	// LD (HL), D
	0x72: func(cpu *CPU) { cpu.mapper.write(cpu.hl, cpu.D()) },

	// LD (HL), E
	0x73: func(cpu *CPU) { cpu.mapper.write(cpu.hl, cpu.E()) },

	// LD (HL), H
	0x74: func(cpu *CPU) { cpu.mapper.write(cpu.hl, cpu.H()) },

	// LD (HL), L
	0x75: func(cpu *CPU) { cpu.mapper.write(cpu.hl, cpu.L()) },

	// HALT
	0x76: func(cpu *CPU) { cpu.halted = true },

	// LD (HL), A
	0x77: func(cpu *CPU) { cpu.mapper.write(cpu.hl, cpu.A()) },

	// LD A, B
	0x78: func(cpu *CPU) { cpu.setA(cpu.B()) },

	// LD A, C
	0x79: func(cpu *CPU) { cpu.setA(cpu.C()) },

	// LD A, D
	0x7A: func(cpu *CPU) { cpu.setA(cpu.D()) },

	// LD A, E
	0x7B: func(cpu *CPU) { cpu.setA(cpu.E()) },

	// LD A, H
	0x7C: func(cpu *CPU) { cpu.setA(cpu.H()) },

	// LD A, L
	0x7D: func(cpu *CPU) { cpu.setA(cpu.L()) },

	// LD A, (HL)
	0x7E: func(cpu *CPU) { cpu.setA(cpu.mapper.read(cpu.hl)) },

	// LD A, A
	0x7F: func(cpu *CPU) { cpu.setA(cpu.A()) },

	// ADD B
	0x80: func(cpu *CPU) { cpu.setA(cpu.byteAdd(cpu.A(), cpu.B())) },

	// ADD C
	0x81: func(cpu *CPU) { cpu.setA(cpu.byteAdd(cpu.A(), cpu.C())) },

	// ADD D
	0x82: func(cpu *CPU) { cpu.setA(cpu.byteAdd(cpu.A(), cpu.D())) },

	// ADD E
	0x83: func(cpu *CPU) { cpu.setA(cpu.byteAdd(cpu.A(), cpu.E())) },

	// ADD H
	0x84: func(cpu *CPU) { cpu.setA(cpu.byteAdd(cpu.A(), cpu.H())) },

	// ADD L
	0x85: func(cpu *CPU) { cpu.setA(cpu.byteAdd(cpu.A(), cpu.L())) },

	// ADD A, (HL)
	0x86: func(cpu *CPU) { cpu.setA(cpu.byteAdd(cpu.A(), cpu.mapper.read(cpu.hl))) },

	// ADD A
	0x87: func(cpu *CPU) { cpu.setA(cpu.byteAdd(cpu.A(), cpu.A())) },

	// ADC A, B
	0x88: func(cpu *CPU) { cpu.setA(cpu.byteAddCarry(cpu.A(), cpu.B())) },

	// ADC A, C
	0x89: func(cpu *CPU) { cpu.setA(cpu.byteAddCarry(cpu.A(), cpu.C())) },

	// ADC A, D
	0x8A: func(cpu *CPU) { cpu.setA(cpu.byteAddCarry(cpu.A(), cpu.D())) },

	// ADC A, E
	0x8B: func(cpu *CPU) { cpu.setA(cpu.byteAddCarry(cpu.A(), cpu.E())) },

	// ADC A, H
	0x8C: func(cpu *CPU) { cpu.setA(cpu.byteAddCarry(cpu.A(), cpu.H())) },

	// ADC A, L
	0x8D: func(cpu *CPU) { cpu.setA(cpu.byteAddCarry(cpu.A(), cpu.L())) },

	// ADC A, (HL)
	0x8E: func(cpu *CPU) { cpu.setA(cpu.byteAddCarry(cpu.A(), cpu.mapper.read(cpu.hl))) },

	// ADC A, A
	0x8F: func(cpu *CPU) { cpu.setA(cpu.byteAddCarry(cpu.A(), cpu.A())) },

	// SUB A, B
	0x90: func(cpu *CPU) { cpu.setA(cpu.byteSub(cpu.A(), cpu.B())) },

	// SUB A, C
	0x91: func(cpu *CPU) { cpu.setA(cpu.byteSub(cpu.A(), cpu.C())) },

	// SUB A, D
	0x92: func(cpu *CPU) { cpu.setA(cpu.byteSub(cpu.A(), cpu.D())) },

	// SUB A, E
	0x93: func(cpu *CPU) { cpu.setA(cpu.byteSub(cpu.A(), cpu.E())) },

	// SUB A, H
	0x94: func(cpu *CPU) { cpu.setA(cpu.byteSub(cpu.A(), cpu.H())) },

	// SUB A, L
	0x95: func(cpu *CPU) { cpu.setA(cpu.byteSub(cpu.A(), cpu.L())) },

	// SUB A, (HL)
	0x96: func(cpu *CPU) { cpu.setA(cpu.byteSub(cpu.A(), cpu.mapper.read(cpu.hl))) },

	// SUB A, A
	0x97: func(cpu *CPU) { cpu.setA(cpu.byteSub(cpu.A(), cpu.A())) },

	// SBC A, B
	0x98: func(cpu *CPU) { cpu.setA(cpu.byteSubCarry(cpu.A(), cpu.B())) },

	// SBC A, C
	0x99: func(cpu *CPU) { cpu.setA(cpu.byteSubCarry(cpu.A(), cpu.C())) },

	// SBC A, D
	0x9A: func(cpu *CPU) { cpu.setA(cpu.byteSubCarry(cpu.A(), cpu.D())) },

	// SBC A, E
	0x9B: func(cpu *CPU) { cpu.setA(cpu.byteSubCarry(cpu.A(), cpu.E())) },

	// SBC A, H
	0x9C: func(cpu *CPU) { cpu.setA(cpu.byteSubCarry(cpu.A(), cpu.H())) },

	// SBC A, L
	0x9D: func(cpu *CPU) { cpu.setA(cpu.byteSubCarry(cpu.A(), cpu.L())) },

	// SBC A, (HL)
	0x9E: func(cpu *CPU) { cpu.setA(cpu.byteSubCarry(cpu.A(), cpu.mapper.read(cpu.hl))) },

	// SBC A, A
	0x9F: func(cpu *CPU) { cpu.setA(cpu.byteSubCarry(cpu.A(), cpu.A())) },

	// AND B
	0xA0: func(cpu *CPU) { cpu.setA(cpu.byteAND(cpu.A(), cpu.B())) },

	// AND C
	0xA1: func(cpu *CPU) { cpu.setA(cpu.byteAND(cpu.A(), cpu.C())) },

	// AND D
	0xA2: func(cpu *CPU) { cpu.setA(cpu.byteAND(cpu.A(), cpu.D())) },

	// AND E
	0xA3: func(cpu *CPU) { cpu.setA(cpu.byteAND(cpu.A(), cpu.E())) },

	// AND H
	0xA4: func(cpu *CPU) { cpu.setA(cpu.byteAND(cpu.A(), cpu.H())) },

	// AND L
	0xA5: func(cpu *CPU) { cpu.setA(cpu.byteAND(cpu.A(), cpu.L())) },

	// AND (HL)
	0xA6: func(cpu *CPU) { cpu.setA(cpu.byteAND(cpu.A(), cpu.mapper.read(cpu.hl))) },

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

	// XOR A, (HL)
	0xAE: func(cpu *CPU) { cpu.setA(cpu.byteXOR(cpu.A(), cpu.mapper.read(cpu.hl))) },

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

	// OR A, (HL)
	0xB6: func(cpu *CPU) { cpu.setA(cpu.byteOR(cpu.A(), cpu.mapper.read(cpu.hl))) },

	// OR A, A
	0xB7: func(cpu *CPU) { cpu.setA(cpu.byteOR(cpu.A(), cpu.A())) },

	// CP B
	0xB8: func(cpu *CPU) { cpu.cmp(cpu.A(), cpu.B()) },

	// CP C
	0xB9: func(cpu *CPU) { cpu.cmp(cpu.A(), cpu.C()) },

	// CP D
	0xBA: func(cpu *CPU) { cpu.cmp(cpu.A(), cpu.D()) },

	// CP E
	0xBB: func(cpu *CPU) { cpu.cmp(cpu.A(), cpu.E()) },

	// CP H
	0xBC: func(cpu *CPU) { cpu.cmp(cpu.A(), cpu.H()) },

	// CP L
	0xBD: func(cpu *CPU) { cpu.cmp(cpu.A(), cpu.L()) },

	// CP A, [HL]
	0xBE: func(cpu *CPU) { cpu.cmp(cpu.A(), cpu.mapper.read(cpu.hl)) },

	// CP A
	0xBF: func(cpu *CPU) { cpu.cmp(cpu.A(), cpu.A()) },

	// RET NZ
	0xC0: func(cpu *CPU) {
		if !cpu.getFlagZ() {
			cpu.returnSub()
		}
	},

	// POP BC
	0xC1: func(cpu *CPU) { cpu.bc = cpu.popStack() },

	// JP NZ, nn
	0xC2: func(cpu *CPU) {
		nn := cpu.fetchPC16()
		if !cpu.getFlagZ() {
			cpu.pc = nn
		}
	},

	// JP nn
	0xC3: func(cpu *CPU) { cpu.pc = cpu.fetchPC16() },

	// CALL NZ, nn
	0xC4: func(cpu *CPU) {
		nn := cpu.fetchPC16()
		if !cpu.getFlagZ() {
			cpu.callSub(nn)
		}
	},

	// PUSH BC
	0xC5: func(cpu *CPU) { cpu.pushStack(cpu.bc) },

	// ADD A, n
	0xC6: func(cpu *CPU) { cpu.setA(cpu.byteAdd(cpu.A(), cpu.fetchPC())) },

	// RST 0x00
	0xC7: func(cpu *CPU) { cpu.callSub(0x0000) },

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
			cpu.pc = nn
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

	// CALL Z, nn
	0xCC: func(cpu *CPU) {
		nn := cpu.fetchPC16()
		if cpu.getFlagZ() {
			cpu.callSub(nn)
		}
	},

	// CALL nn
	0xCD: func(cpu *CPU) { cpu.callSub(cpu.fetchPC16()) },

	// ADD A, n
	0xCE: func(cpu *CPU) { cpu.setA(cpu.byteAdd(cpu.A(), cpu.fetchPC())) },

	// RST 0x8
	0xCF: func(cpu *CPU) { cpu.callSub(0x0008) },

	// RET NC
	0xD0: func(cpu *CPU) {
		if !cpu.getFlagC() {
			cpu.returnSub()
		}
	},

	// POP DE
	0xD1: func(cpu *CPU) { cpu.de = cpu.popStack() },

	// JP NC, nn
	0xD2: func(cpu *CPU) {
		if !cpu.getFlagC() {
			cpu.pc = cpu.fetchPC16()
		}
	},

	// CALL NC, nn
	0xD4: func(cpu *CPU) {
		nn := cpu.fetchPC16()
		if !cpu.getFlagC() {
			cpu.callSub(nn)
		}
	},

	// PUSH DE
	0xD5: func(cpu *CPU) { cpu.pushStack(cpu.de) },

	// SUB A, n
	0xD6: func(cpu *CPU) {
		result := cpu.byteSub(cpu.A(), cpu.fetchPC())
		cpu.setFlagC(result > cpu.A())
		cpu.setA(result)
	},

	// RST 0x10
	0xD7: func(cpu *CPU) { cpu.callSub(0x0010) },

	// RET C
	0xD8: func(cpu *CPU) {
		if cpu.getFlagC() {
			cpu.returnSub()
		}
	},

	// RETI
	0xD9: func(cpu *CPU) {
		cpu.returnSub()
		cpu.ime = true
	},

	// JP C, nn
	0xDA: func(cpu *CPU) {
		if cpu.getFlagC() {
			cpu.pc = cpu.fetchPC16()
		}
	},

	// CALL C, nn
	0xDC: func(cpu *CPU) {
		nn := cpu.fetchPC16()
		if cpu.getFlagC() {
			cpu.callSub(nn)
		}
	},

	// SBC A, n
	0xDE: func(cpu *CPU) {
		result := cpu.byteSubCarry(cpu.A(), cpu.fetchPC())
		cpu.setFlagC(result > cpu.A())
		cpu.setA(result)
	},

	// RST 0x18
	0xDF: func(cpu *CPU) { cpu.callSub(0x0018) },

	// LDH (n), A
	0xE0: func(cpu *CPU) {
		cpu.mapper.write(0xFF00+uint16(cpu.fetchPC()), cpu.A())
	},

	// POP HL
	0xE1: func(cpu *CPU) { cpu.hl = cpu.popStack() },

	// LD (C), A
	0xE2: func(cpu *CPU) { cpu.mapper.write(0xFF00+uint16(cpu.C()), cpu.A()) },

	// PUSH HL
	0xE5: func(cpu *CPU) { cpu.pushStack(cpu.hl) },

	// AND A, n
	0xE6: func(cpu *CPU) { cpu.setA(cpu.byteAND(cpu.A(), cpu.fetchPC())) },

	// RST 20H
	0xE7: func(cpu *CPU) { cpu.callSub(0x0020) },

	// ADD SP, n
	0xE8: func(cpu *CPU) {
		n := int8(cpu.fetchPC())
		result := uint16(int32(cpu.sp) + int32(n))
		cpu.setFlagZ(false)
		cpu.setFlagN(false)
		cpu.setFlagH((cpu.sp&0xF)+(uint16(n)&0xF) > 0xF)
		cpu.setFlagC((cpu.sp&0xFF)+(uint16(n)&0xFF) > 0xFF)
		cpu.sp = result
	},

	// JP HL
	0xE9: func(cpu *CPU) { cpu.pc = cpu.hl },

	// LD (nn), A
	0xEA: func(cpu *CPU) { cpu.mapper.write(cpu.fetchPC16(), cpu.A()) },

	// XOR A, n
	0xEE: func(cpu *CPU) { cpu.setA(cpu.byteXOR(cpu.A(), cpu.fetchPC())) },

	// RST 28H
	0xEF: func(cpu *CPU) { cpu.callSub(0x0028) },

	// LDH A, (n)
	0xF0: func(cpu *CPU) { cpu.setA(cpu.mapper.read(0xFF00 + uint16(cpu.fetchPC()))) },

	// POP AF
	0xF1: func(cpu *CPU) { cpu.af = cpu.popStack() },

	// LD A, (C)
	0xF2: func(cpu *CPU) { cpu.setA(cpu.mapper.read(0xFF00 + uint16(cpu.C()))) },

	// DI
	0xF3: func(cpu *CPU) { cpu.ime = false },

	// PUSH AF
	0xF5: func(cpu *CPU) { cpu.pushStack(cpu.af) },

	// OR A, n
	0xF6: func(cpu *CPU) { cpu.setA(cpu.byteOR(cpu.A(), cpu.fetchPC())) },

	// RST 30H
	0xF7: func(cpu *CPU) { cpu.callSub(0x0030) },

	// LD HL, SP+n
	0xF8: func(cpu *CPU) {
		n := int8(cpu.fetchPC())
		result := uint16(int32(cpu.sp) + int32(n))
		cpu.setFlagZ(false)
		cpu.setFlagN(false)
		cpu.setFlagH((cpu.sp&0xF)+(uint16(n)&0xF) > 0xF)
		cpu.setFlagC((cpu.sp&0xFF)+(uint16(n)&0xFF) > 0xFF)
		cpu.hl = result
	},

	// LD SP, HL
	0xF9: func(cpu *CPU) { cpu.sp = cpu.hl },

	// LD A, (nn)
	0xFA: func(cpu *CPU) { cpu.setA(cpu.mapper.read(cpu.fetchPC16())) },

	// EI
	0xFB: func(cpu *CPU) { cpu.ime = true },

	// CP
	0xFE: func(cpu *CPU) { cpu.cmp(cpu.A(), cpu.fetchPC()) },

	// RST 38H
	0xFF: func(cpu *CPU) { cpu.callSub(0x0038) },
}
