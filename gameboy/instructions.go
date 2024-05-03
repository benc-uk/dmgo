package gameboy

var opcodes = [0x100]func(cpu *CPU){
	// NOP
	0x00: func(cpu *CPU) { cpu.opDebug = "NOP" },

	// LD BC, nn
	0x01: func(cpu *CPU) { cpu.BC = cpu.fetchPC16() },

	// INC BC
	0x03: func(cpu *CPU) { cpu.BC++ },

	// INC B
	0x04: func(cpu *CPU) { cpu.setB(cpu.byteAdd(cpu.B(), 1)) },

	// DEC B
	0x05: func(cpu *CPU) { cpu.setB(cpu.byteSub(cpu.B(), 1)) },

	// LD B, n
	0x06: func(cpu *CPU) { cpu.setB(cpu.fetchPC()) },

	// DEC BC
	0x0B: func(cpu *CPU) { cpu.BC-- },

	// INC C
	0x0C: func(cpu *CPU) { cpu.setC(cpu.byteAdd(cpu.C(), 1)) },

	// DEC C
	0x0D: func(cpu *CPU) { cpu.setC(cpu.byteSub(cpu.C(), 1)) },

	// LD C, n
	0x0E: func(cpu *CPU) { cpu.setC(cpu.fetchPC()) },

	// LD DE, nn
	0x11: func(cpu *CPU) { cpu.DE = cpu.fetchPC16() },

	// INC DE
	0x13: func(cpu *CPU) { cpu.DE++ },

	// JR e
	0x18: func(cpu *CPU) {
		cpu.PC += uint16(int8(cpu.fetchPC()))
	},

	// LD A, (DE)
	0x1A: func(cpu *CPU) { cpu.setA(cpu.mapper.Read(cpu.DE)) },

	// DEC DE
	0x1B: func(cpu *CPU) { cpu.DE-- },

	// INC E
	0x1C: func(cpu *CPU) { cpu.setE(cpu.byteAdd(cpu.E(), 1)) },

	// DEC E
	0x1D: func(cpu *CPU) { cpu.setE(cpu.byteSub(cpu.E(), 1)) },

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
	0x24: func(cpu *CPU) {
		cpu.setH(cpu.byteAdd(cpu.H(), 1))
	},

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
	0x2C: func(cpu *CPU) { cpu.setL(cpu.byteAdd(cpu.L(), 1)) },

	// DEC L
	0x2D: func(cpu *CPU) { cpu.setL(cpu.byteSub(cpu.L(), 1)) },

	// LD L, n
	0x2E: func(cpu *CPU) { cpu.setL(cpu.fetchPC()) },

	// LD SP, nn
	0x31: func(cpu *CPU) { cpu.SP = cpu.fetchPC16() },

	// LD [HL-], A
	0x32: func(cpu *CPU) {
		cpu.mapper.Write(cpu.HL, cpu.A())
		cpu.HL--
	},

	// INC SP
	0x33: func(cpu *CPU) { cpu.SP++ },

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

	// LD A, n
	0x3E: func(cpu *CPU) { cpu.setA(cpu.fetchPC()) },

	// LD B, H
	0x44: func(cpu *CPU) { cpu.setB(cpu.H()) },

	// LD B, L
	0x45: func(cpu *CPU) { cpu.setB(cpu.L()) },

	// LD B, (HL)
	0x46: func(cpu *CPU) { cpu.setB(cpu.mapper.Read(cpu.HL)) },

	// LB C, B
	0x48: func(cpu *CPU) { cpu.setC(cpu.B()) },

	// LD C, H
	0x4C: func(cpu *CPU) { cpu.setC(cpu.H()) },

	// LD C, L
	0x4D: func(cpu *CPU) { cpu.setC(cpu.L()) },

	// LD C, (HL)
	0x4E: func(cpu *CPU) { cpu.setC(cpu.mapper.Read(cpu.HL)) },

	// LD D, L
	0x55: func(cpu *CPU) { cpu.setD(cpu.L()) },

	// LD D, (HL)
	0x56: func(cpu *CPU) { cpu.setD(cpu.mapper.Read(cpu.HL)) },

	// LD E, B
	0x58: func(cpu *CPU) { cpu.setE(cpu.B()) },

	// LD E, H
	0x5C: func(cpu *CPU) { cpu.setE(cpu.H()) },

	// LD E, L
	0x5D: func(cpu *CPU) { cpu.setE(cpu.L()) },

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

	// LD A, H
	0x7C: func(cpu *CPU) { cpu.setA(cpu.H()) },

	// LD A, L
	0x7D: func(cpu *CPU) { cpu.setA(cpu.L()) },

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
	0xC6: func(cpu *CPU) {
		result := cpu.byteAdd(cpu.A(), cpu.fetchPC())
		cpu.setA(result)
		cpu.setFlagC(result < cpu.A())
	},

	// RET
	0xC9: func(cpu *CPU) { cpu.returnSub() },

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

	// LD (nn), A
	0xEA: func(cpu *CPU) { cpu.mapper.Write(cpu.fetchPC16(), cpu.A()) },

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

	// CP
	0xFE: func(cpu *CPU) { cpu.cmp(cpu.A(), cpu.fetchPC()) },
}
