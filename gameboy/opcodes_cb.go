package gameboy

var cbOpcodes = [0x100]func(*CPU){}

func init() {
	// ==== ROTATE ===================================

	// 0x40 ~ 0x70
	for i := 0; i <= 3; i++ {
		cbOpcodes[0x40+0x10*i] = func(cpu *CPU) { cpu.bitTest(cpu.B(), uint(i*2)) }
	}

	// 0x41 ~ 0x71
	for i := 0; i <= 3; i++ {
		cbOpcodes[0x41+0x10*i] = func(cpu *CPU) { cpu.bitTest(cpu.C(), uint(i*2)) }
	}

	// 0x42 ~ 0x72
	for i := 0; i <= 3; i++ {
		cbOpcodes[0x42+0x10*i] = func(cpu *CPU) { cpu.bitTest(cpu.D(), uint(i*2)) }
	}

	// 0x43 ~ 0x73
	for i := 0; i <= 3; i++ {
		cbOpcodes[0x43+0x10*i] = func(cpu *CPU) { cpu.bitTest(cpu.E(), uint(i*2)) }
	}

	// 0x44 ~ 0x74
	for i := 0; i <= 3; i++ {
		cbOpcodes[0x44+0x10*i] = func(cpu *CPU) { cpu.bitTest(cpu.H(), uint(i*2)) }
	}

	// 0x45 ~ 0x75
	for i := 0; i <= 3; i++ {
		cbOpcodes[0x45+0x10*i] = func(cpu *CPU) { cpu.bitTest(cpu.L(), uint(i*2)) }
	}

	// 0x46 ~ 0x76
	for i := 0; i <= 3; i++ {
		cbOpcodes[0x46+0x10*i] = func(cpu *CPU) { cpu.bitTest(cpu.mapper.Read(cpu.HL), uint(i*2)) }
	}

	// 0x47 ~ 0x77
	for i := 0; i <= 3; i++ {
		cbOpcodes[0x47+0x10*i] = func(cpu *CPU) { cpu.bitTest(cpu.A(), uint(i*2)) }
	}

	// 0x48 ~ 0x78
	for i := 0; i <= 3; i++ {
		cbOpcodes[0x48+0x10*i] = func(cpu *CPU) { cpu.bitTest(cpu.B(), uint(i*2)+1) }
	}

	// 0x49 ~ 0x79
	for i := 0; i <= 3; i++ {
		cbOpcodes[0x49+0x10*i] = func(cpu *CPU) { cpu.bitTest(cpu.C(), uint(i*2)+1) }
	}

	// 0x4A ~ 0x7A
	for i := 0; i <= 3; i++ {
		cbOpcodes[0x4A+0x10*i] = func(cpu *CPU) { cpu.bitTest(cpu.D(), uint(i*2)+1) }
	}

	// 0x4B ~ 0x7B
	for i := 0; i <= 3; i++ {
		cbOpcodes[0x4B+0x10*i] = func(cpu *CPU) { cpu.bitTest(cpu.E(), uint(i*2)+1) }
	}

	// 0x4C ~ 0x7C
	for i := 0; i <= 3; i++ {
		cbOpcodes[0x4C+0x10*i] = func(cpu *CPU) { cpu.bitTest(cpu.H(), uint(i*2)+1) }
	}

	// 0x4D ~ 0x7D
	for i := 0; i <= 3; i++ {
		cbOpcodes[0x4D+0x10*i] = func(cpu *CPU) { cpu.bitTest(cpu.L(), uint(i*2)+1) }
	}

	// 0x4E ~ 0x7E
	for i := 0; i <= 3; i++ {
		cbOpcodes[0x4E+0x10*i] = func(cpu *CPU) { cpu.bitTest(cpu.mapper.Read(cpu.HL), uint(i*2)+1) }
	}

	// 0x4F ~ 0x7F
	for i := 0; i <= 3; i++ {
		cbOpcodes[0x4F+0x10*i] = func(cpu *CPU) { cpu.bitTest(cpu.A(), uint(i*2)+1) }
	}

	// ==== BIT RESET ===================================

	// 0x80 ~ 0xB0
	for i := 0; i <= 3; i++ {
		cbOpcodes[0x80+0x10*i] = func(cpu *CPU) { cpu.setB(bitReset(cpu.B(), uint(i*2))) }
	}

	// 0x81 ~ 0xB1
	for i := 0; i <= 3; i++ {
		cbOpcodes[0x81+0x10*i] = func(cpu *CPU) { cpu.setC(bitReset(cpu.C(), uint(i*2))) }
	}

	// 0x82 ~ 0xB2
	for i := 0; i <= 3; i++ {
		cbOpcodes[0x82+0x10*i] = func(cpu *CPU) { cpu.setD(bitReset(cpu.D(), uint(i*2))) }
	}

	// 0x83 ~ 0xB3
	for i := 0; i <= 3; i++ {
		cbOpcodes[0x83+0x10*i] = func(cpu *CPU) { cpu.setE(bitReset(cpu.E(), uint(i*2))) }
	}

	// 0x84 ~ 0xB4
	for i := 0; i <= 3; i++ {
		cbOpcodes[0x84+0x10*i] = func(cpu *CPU) { cpu.setH(bitReset(cpu.H(), uint(i*2))) }
	}

	// 0x85 ~ 0xB5
	for i := 0; i <= 3; i++ {
		cbOpcodes[0x85+0x10*i] = func(cpu *CPU) { cpu.setL(bitReset(cpu.L(), uint(i*2))) }
	}

	// 0x86 ~ 0xB6
	for i := 0; i <= 3; i++ {
		cbOpcodes[0x86+0x10*i] = func(cpu *CPU) { cpu.mapper.Write(cpu.HL, bitReset(cpu.mapper.Read(cpu.HL), uint(i*2))) }
	}

	// 0x87 ~ 0xB7
	for i := 0; i <= 3; i++ {
		cbOpcodes[0x87+0x10*i] = func(cpu *CPU) { cpu.setA(bitReset(cpu.A(), uint(i*2))) }
	}

	// 0x88 ~ 0xB8
	for i := 0; i <= 3; i++ {
		cbOpcodes[0x88+0x10*i] = func(cpu *CPU) { cpu.setB(bitReset(cpu.B(), uint(i*2)+1)) }
	}

	// 0x89 ~ 0xB9
	for i := 0; i <= 3; i++ {
		cbOpcodes[0x89+0x10*i] = func(cpu *CPU) { cpu.setC(bitReset(cpu.C(), uint(i*2)+1)) }
	}

	// 0x8A ~ 0xBA
	for i := 0; i <= 3; i++ {
		cbOpcodes[0x8A+0x10*i] = func(cpu *CPU) { cpu.setD(bitReset(cpu.D(), uint(i*2)+1)) }
	}

	// 0x8B ~ 0xBB
	for i := 0; i <= 3; i++ {
		cbOpcodes[0x8B+0x10*i] = func(cpu *CPU) { cpu.setE(bitReset(cpu.E(), uint(i*2)+1)) }
	}

	// 0x8C ~ 0xBC
	for i := 0; i <= 3; i++ {
		cbOpcodes[0x8C+0x10*i] = func(cpu *CPU) { cpu.setH(bitReset(cpu.H(), uint(i*2)+1)) }
	}

	// 0x8D ~ 0xBD
	for i := 0; i <= 3; i++ {
		cbOpcodes[0x8D+0x10*i] = func(cpu *CPU) { cpu.setL(bitReset(cpu.L(), uint(i*2)+1)) }
	}

	// 0x8E ~ 0xBE
	for i := 0; i <= 3; i++ {
		cbOpcodes[0x8E+0x10*i] = func(cpu *CPU) { cpu.mapper.Write(cpu.HL, bitReset(cpu.mapper.Read(cpu.HL), uint(i*2)+1)) }
	}

	// 0x8F ~ 0xBF
	for i := 0; i <= 3; i++ {
		cbOpcodes[0x8F+0x10*i] = func(cpu *CPU) { cpu.setA(bitReset(cpu.A(), uint(i*2)+1)) }
	}

	// ==== BIT SET ===================================

	// 0xC0 ~ 0xF0
	for i := 0; i <= 3; i++ {
		cbOpcodes[0xC0+0x10*i] = func(cpu *CPU) { cpu.setB(bitSet(cpu.B(), uint(i*2))) }
	}

	// 0xC1 ~ 0xF1
	for i := 0; i <= 3; i++ {
		cbOpcodes[0xC1+0x10*i] = func(cpu *CPU) { cpu.setC(bitSet(cpu.C(), uint(i*2))) }
	}

	// 0xC2 ~ 0xF2
	for i := 0; i <= 3; i++ {
		cbOpcodes[0xC2+0x10*i] = func(cpu *CPU) { cpu.setD(bitSet(cpu.D(), uint(i*2))) }
	}

	// 0xC3 ~ 0xF3
	for i := 0; i <= 3; i++ {
		cbOpcodes[0xC3+0x10*i] = func(cpu *CPU) { cpu.setE(bitSet(cpu.E(), uint(i*2))) }
	}

	// 0xC4 ~ 0xF4
	for i := 0; i <= 3; i++ {
		cbOpcodes[0xC4+0x10*i] = func(cpu *CPU) { cpu.setH(bitSet(cpu.H(), uint(i*2))) }
	}

	// 0xC5 ~ 0xF5
	for i := 0; i <= 3; i++ {
		cbOpcodes[0xC5+0x10*i] = func(cpu *CPU) { cpu.setL(bitSet(cpu.L(), uint(i*2))) }
	}

	// 0xC6 ~ 0xF6
	for i := 0; i <= 3; i++ {
		cbOpcodes[0xC6+0x10*i] = func(cpu *CPU) { cpu.mapper.Write(cpu.HL, bitSet(cpu.mapper.Read(cpu.HL), uint(i*2))) }
	}

	// 0xC7 ~ 0xF7
	for i := 0; i <= 3; i++ {
		cbOpcodes[0xC7+0x10*i] = func(cpu *CPU) { cpu.setA(bitSet(cpu.A(), uint(i*2))) }
	}

	// 0xC8 ~ 0xF8
	for i := 0; i <= 3; i++ {
		cbOpcodes[0xC8+0x10*i] = func(cpu *CPU) { cpu.setB(bitSet(cpu.B(), uint(i*2)+1)) }
	}

	// 0xC9 ~ 0xF9
	for i := 0; i <= 3; i++ {
		cbOpcodes[0xC9+0x10*i] = func(cpu *CPU) { cpu.setC(bitSet(cpu.C(), uint(i*2)+1)) }
	}

	// 0xCA ~ 0xFA
	for i := 0; i <= 3; i++ {
		cbOpcodes[0xCA+0x10*i] = func(cpu *CPU) { cpu.setD(bitSet(cpu.D(), uint(i*2)+1)) }
	}

	// 0xCB ~ 0xFB
	for i := 0; i <= 3; i++ {
		cbOpcodes[0xCB+0x10*i] = func(cpu *CPU) { cpu.setE(bitSet(cpu.E(), uint(i*2)+1)) }
	}

	// 0xCC ~ 0xFC
	for i := 0; i <= 3; i++ {
		cbOpcodes[0xCC+0x10*i] = func(cpu *CPU) { cpu.setH(bitSet(cpu.H(), uint(i*2)+1)) }
	}

	// 0xCD ~ 0xFD
	for i := 0; i <= 3; i++ {
		cbOpcodes[0xCD+0x10*i] = func(cpu *CPU) { cpu.setL(bitSet(cpu.L(), uint(i*2)+1)) }
	}

	// 0xCE ~ 0xFE
	for i := 0; i <= 3; i++ {
		cbOpcodes[0xCE+0x10*i] = func(cpu *CPU) { cpu.mapper.Write(cpu.HL, bitSet(cpu.mapper.Read(cpu.HL), uint(i*2)+1)) }
	}

	// 0xCF ~ 0xFF
	for i := 0; i <= 3; i++ {
		cbOpcodes[0xCF+0x10*i] = func(cpu *CPU) { cpu.setA(bitSet(cpu.A(), uint(i*2)+1)) }
	}

	// ==== ROTATE ===================================

	// 0x00
	cbOpcodes[0x00] = func(cpu *CPU) { cpu.setB(cpu.rotLeftCarry(cpu.B())) }
	cbOpcodes[0x01] = func(cpu *CPU) { cpu.setC(cpu.rotLeftCarry(cpu.C())) }
	cbOpcodes[0x02] = func(cpu *CPU) { cpu.setD(cpu.rotLeftCarry(cpu.D())) }
	cbOpcodes[0x03] = func(cpu *CPU) { cpu.setE(cpu.rotLeftCarry(cpu.E())) }
	cbOpcodes[0x04] = func(cpu *CPU) { cpu.setH(cpu.rotLeftCarry(cpu.H())) }
	cbOpcodes[0x05] = func(cpu *CPU) { cpu.setL(cpu.rotLeftCarry(cpu.L())) }
	cbOpcodes[0x06] = func(cpu *CPU) { cpu.mapper.Write(cpu.HL, cpu.rotLeftCarry(cpu.mapper.Read(cpu.HL))) }
	cbOpcodes[0x07] = func(cpu *CPU) { cpu.setA(cpu.rotLeftCarry(cpu.A())) }
	cbOpcodes[0x08] = func(cpu *CPU) { cpu.setB(cpu.rotRightCarry(cpu.B())) }
	cbOpcodes[0x09] = func(cpu *CPU) { cpu.setC(cpu.rotRightCarry(cpu.C())) }
	cbOpcodes[0x0A] = func(cpu *CPU) { cpu.setD(cpu.rotRightCarry(cpu.D())) }
	cbOpcodes[0x0B] = func(cpu *CPU) { cpu.setE(cpu.rotRightCarry(cpu.E())) }
	cbOpcodes[0x0C] = func(cpu *CPU) { cpu.setH(cpu.rotRightCarry(cpu.H())) }
	cbOpcodes[0x0D] = func(cpu *CPU) { cpu.setL(cpu.rotRightCarry(cpu.L())) }
	cbOpcodes[0x0E] = func(cpu *CPU) { cpu.mapper.Write(cpu.HL, cpu.rotRightCarry(cpu.mapper.Read(cpu.HL))) }
	cbOpcodes[0x0F] = func(cpu *CPU) { cpu.setA(cpu.rotRightCarry(cpu.A())) }

	// 0x10
	cbOpcodes[0x10] = func(cpu *CPU) { cpu.setB(cpu.rotLeft(cpu.B())) }
	cbOpcodes[0x11] = func(cpu *CPU) { cpu.setC(cpu.rotLeft(cpu.C())) }
	cbOpcodes[0x12] = func(cpu *CPU) { cpu.setD(cpu.rotLeft(cpu.D())) }
	cbOpcodes[0x13] = func(cpu *CPU) { cpu.setE(cpu.rotLeft(cpu.E())) }
	cbOpcodes[0x14] = func(cpu *CPU) { cpu.setH(cpu.rotLeft(cpu.H())) }
	cbOpcodes[0x15] = func(cpu *CPU) { cpu.setL(cpu.rotLeft(cpu.L())) }
	cbOpcodes[0x16] = func(cpu *CPU) { cpu.mapper.Write(cpu.HL, cpu.rotLeft(cpu.mapper.Read(cpu.HL))) }
	cbOpcodes[0x17] = func(cpu *CPU) { cpu.setA(cpu.rotLeft(cpu.A())) }
	cbOpcodes[0x18] = func(cpu *CPU) { cpu.setB(cpu.rotRight(cpu.B())) }
	cbOpcodes[0x19] = func(cpu *CPU) { cpu.setC(cpu.rotRight(cpu.C())) }
	cbOpcodes[0x1A] = func(cpu *CPU) { cpu.setD(cpu.rotRight(cpu.D())) }
	cbOpcodes[0x1B] = func(cpu *CPU) { cpu.setE(cpu.rotRight(cpu.E())) }
	cbOpcodes[0x1C] = func(cpu *CPU) { cpu.setH(cpu.rotRight(cpu.H())) }
	cbOpcodes[0x1D] = func(cpu *CPU) { cpu.setL(cpu.rotRight(cpu.L())) }
	cbOpcodes[0x1E] = func(cpu *CPU) { cpu.mapper.Write(cpu.HL, cpu.rotRight(cpu.mapper.Read(cpu.HL))) }
	cbOpcodes[0x1F] = func(cpu *CPU) { cpu.setA(cpu.rotRight(cpu.A())) }

	// 0x20
	cbOpcodes[0x20] = func(cpu *CPU) { cpu.setB(cpu.shiftLeftArithmetic(cpu.B())) }
	cbOpcodes[0x21] = func(cpu *CPU) { cpu.setC(cpu.shiftLeftArithmetic(cpu.C())) }
	cbOpcodes[0x22] = func(cpu *CPU) { cpu.setD(cpu.shiftLeftArithmetic(cpu.D())) }
	cbOpcodes[0x23] = func(cpu *CPU) { cpu.setE(cpu.shiftLeftArithmetic(cpu.E())) }
	cbOpcodes[0x24] = func(cpu *CPU) { cpu.setH(cpu.shiftLeftArithmetic(cpu.H())) }
	cbOpcodes[0x25] = func(cpu *CPU) { cpu.setL(cpu.shiftLeftArithmetic(cpu.L())) }
	cbOpcodes[0x26] = func(cpu *CPU) { cpu.mapper.Write(cpu.HL, cpu.shiftLeftArithmetic(cpu.mapper.Read(cpu.HL))) }
	cbOpcodes[0x27] = func(cpu *CPU) { cpu.setA(cpu.shiftLeftArithmetic(cpu.A())) }
	cbOpcodes[0x28] = func(cpu *CPU) { cpu.setB(cpu.shiftRightArithmetic(cpu.B())) }
	cbOpcodes[0x29] = func(cpu *CPU) { cpu.setC(cpu.shiftRightArithmetic(cpu.C())) }
	cbOpcodes[0x2A] = func(cpu *CPU) { cpu.setD(cpu.shiftRightArithmetic(cpu.D())) }
	cbOpcodes[0x2B] = func(cpu *CPU) { cpu.setE(cpu.shiftRightArithmetic(cpu.E())) }
	cbOpcodes[0x2C] = func(cpu *CPU) { cpu.setH(cpu.shiftRightArithmetic(cpu.H())) }
	cbOpcodes[0x2D] = func(cpu *CPU) { cpu.setL(cpu.shiftRightArithmetic(cpu.L())) }
	cbOpcodes[0x2E] = func(cpu *CPU) { cpu.mapper.Write(cpu.HL, cpu.shiftRightArithmetic(cpu.mapper.Read(cpu.HL))) }
	cbOpcodes[0x2F] = func(cpu *CPU) { cpu.setA(cpu.shiftRightArithmetic(cpu.A())) }

	// 0x30
	cbOpcodes[0x30] = func(cpu *CPU) { cpu.setB(cpu.swapNibbles(cpu.B())) }
	cbOpcodes[0x31] = func(cpu *CPU) { cpu.setC(cpu.swapNibbles(cpu.C())) }
	cbOpcodes[0x32] = func(cpu *CPU) { cpu.setD(cpu.swapNibbles(cpu.D())) }
	cbOpcodes[0x33] = func(cpu *CPU) { cpu.setE(cpu.swapNibbles(cpu.E())) }
	cbOpcodes[0x34] = func(cpu *CPU) { cpu.setH(cpu.swapNibbles(cpu.H())) }
	cbOpcodes[0x35] = func(cpu *CPU) { cpu.setL(cpu.swapNibbles(cpu.L())) }
	cbOpcodes[0x36] = func(cpu *CPU) { cpu.mapper.Write(cpu.HL, cpu.swapNibbles(cpu.mapper.Read(cpu.HL))) }
	cbOpcodes[0x37] = func(cpu *CPU) { cpu.setA(cpu.swapNibbles(cpu.A())) }
	cbOpcodes[0x38] = func(cpu *CPU) { cpu.setB(cpu.shiftRightLogical(cpu.B())) }
	cbOpcodes[0x39] = func(cpu *CPU) { cpu.setC(cpu.shiftRightLogical(cpu.C())) }
	cbOpcodes[0x3A] = func(cpu *CPU) { cpu.setD(cpu.shiftRightLogical(cpu.D())) }
	cbOpcodes[0x3B] = func(cpu *CPU) { cpu.setE(cpu.shiftRightLogical(cpu.E())) }
	cbOpcodes[0x3C] = func(cpu *CPU) { cpu.setH(cpu.shiftRightLogical(cpu.H())) }
	cbOpcodes[0x3D] = func(cpu *CPU) { cpu.setL(cpu.shiftRightLogical(cpu.L())) }
	cbOpcodes[0x3E] = func(cpu *CPU) { cpu.mapper.Write(cpu.HL, cpu.shiftRightLogical(cpu.mapper.Read(cpu.HL))) }
	cbOpcodes[0x3F] = func(cpu *CPU) { cpu.setA(cpu.shiftRightLogical(cpu.A())) }
}
