package gameboy

var cbOpcodes = [0x100]func(*CPU){
	// RL C
	0x11: func(cpu *CPU) { cpu.setC(cpu.rotLeft(cpu.C())) },

	0x37: func(cpu *CPU) { cpu.setA(swapNibbles(cpu.A())) },

	0x55: func(cpu *CPU) { cpu.bitTest(cpu.L(), 2) },

	0x7C: func(cpu *CPU) { cpu.bitTest(cpu.H(), 7) },

	// RES 0, A
	0x87: func(cpu *CPU) { cpu.setA(bitReset(cpu.A(), 0)) },
}
