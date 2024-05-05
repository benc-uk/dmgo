package gameboy

var cbOpcodes = [0x100]func(*CPU){

	0x11: func(cpu *CPU) { cpu.setC(cpu.rotLeft(cpu.C())) },

	0x27: func(cpu *CPU) { cpu.setA(cpu.shiftLeftArithmetic(cpu.A())) },

	0x37: func(cpu *CPU) { cpu.setA(swapNibbles(cpu.A())) },

	0x40: func(cpu *CPU) { cpu.bitTest(cpu.B(), 0) },

	0x50: func(cpu *CPU) { cpu.bitTest(cpu.B(), 2) },

	0x58: func(cpu *CPU) { cpu.bitTest(cpu.B(), 3) },

	0x60: func(cpu *CPU) { cpu.bitTest(cpu.B(), 4) },

	0x68: func(cpu *CPU) { cpu.bitTest(cpu.B(), 5) },

	0x55: func(cpu *CPU) { cpu.bitTest(cpu.L(), 2) },

	0x70: func(cpu *CPU) { cpu.bitTest(cpu.B(), 6) },

	0x7C: func(cpu *CPU) { cpu.bitTest(cpu.H(), 7) },

	0x7E: func(cpu *CPU) { cpu.bitTest(cpu.mapper.Read(cpu.HL), 7) },

	0x7F: func(cpu *CPU) { cpu.setA(bitSet(cpu.A(), 7)) },

	0x87: func(cpu *CPU) { cpu.setA(bitReset(cpu.A(), 0)) },

	0xDE: func(cpu *CPU) { cpu.mapper.Write(cpu.HL, bitSet(cpu.mapper.Read(cpu.HL), 3)) },
}
