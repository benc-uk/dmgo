package gameboy

import (
	"log"
)

type CPU struct {
	// Registers
	af, bc, de, hl uint16

	// Stack pointer
	sp uint16

	// Program counter
	pc uint16

	// Memory
	mapper *Mapper

	// Special internal flags
	ime    bool // Interrupt Master Enable
	halted bool // Halt state

	// Debugging
	opDebug     []byte
	breakpoints []uint16
}

func NewCPU(mapper *Mapper) *CPU {
	// Initial state of the CPU for the classic GB
	// It represents the state of the CPU after the BIOS has run, as we skip that
	cpu := CPU{
		af:     0x01B0,
		bc:     0x0013,
		de:     0x00D8,
		hl:     0x014D,
		sp:     0xFFEE,
		pc:     0x0000,
		mapper: mapper,
	}

	cpu.setFlagZ(true)
	cpu.setFlagN(false)
	cpu.setFlagH(false)
	cpu.setFlagC(false)
	cpu.ime = false

	return &cpu
}

func (cpu *CPU) ExecuteNext(skipBreak bool) (cyclesSpent int) {
	if cpu.halted {
		// Even if halted, we still spend some cycles
		return 4
	}

	currentPC := cpu.pc

	// Fetch the next instruction, this will also increment the PC
	opcode := cpu.fetchPC()

	// Check if we have hit a breakpoint
	for _, bp := range cpu.breakpoints {
		if bp == currentPC && !skipBreak {
			log.Printf("!!! Breakpoint hit at 0x%04X\n", currentPC)
			cpu.pc = currentPC
			return -1
		}
	}

	// Check if the opcode is valid
	if opcodes[opcode] == nil {
		log.Printf("!!! Unknown opcode: 0x%02X at 0x%04X\n", opcode, currentPC)
		cpu.pc = currentPC
		return -1
	}

	// Debugging output
	for _, b := range cpu.opDebug {
		if b == opcode {
			log.Printf("0x%04X -> 0x%02X (%s)\n", currentPC, opcode, opcodeNames[opcode])
		}
	}

	// Decode & execute the opcode
	opcodes[opcode](cpu)

	cycles := opcodeLengths[opcode]
	return cycles
}

func (cpu *CPU) handleInterrupt(interrupt byte) {
	//log.Printf("Handling interrupt %08b\n", interrupt)

	// Disable interrupts
	cpu.ime = false

	// Clear the interrupt flag at the interrupt bit
	cpu.mapper.Write(IF, cpu.mapper.Read(IF)&^interrupt)

	// Push the current PC onto the stack
	cpu.pushStack(cpu.pc)

	// Jump to the interrupt handler
	switch interrupt {
	case 0x01:
		cpu.pc = 0x0040 // VBLANK interrupt handler address
	case 0x02:
		cpu.pc = 0x0048 // LCD interrupt handler address
	case 0x04:
		cpu.pc = 0x0050 // Timer interrupt handler address
	case 0x08:
		cpu.pc = 0x0058 // Serial interrupt handler address
	case 0x10:
		cpu.pc = 0x0060 // Joypad interrupt handler address
	}

}

// =======================================
// Flag getters and setters
// =======================================

func (cpu *CPU) setFlagZ(value bool) {
	if value {
		cpu.af |= 0x80
	} else {
		cpu.af &^= 0x80
	}
}

func (cpu *CPU) setFlagN(value bool) {
	if value {
		cpu.af |= 0x40
	} else {
		cpu.af &^= 0x40
	}
}

func (cpu *CPU) setFlagH(value bool) {
	if value {
		cpu.af |= 0x20
	} else {
		cpu.af &^= 0x20
	}
}

func (cpu *CPU) setFlagC(value bool) {
	// In bit 4
	if value {
		cpu.af |= 0x10
	} else {
		cpu.af &^= 0x10
	}
}

func (cpu *CPU) getFlagZ() bool { return cpu.af&0x80 != 0 }

func (cpu *CPU) getFlagN() bool { return cpu.af&0x40 != 0 }

func (cpu *CPU) getFlagH() bool { return cpu.af&0x20 != 0 }

func (cpu *CPU) getFlagC() bool { return cpu.af&0x10 != 0 }

// =======================================
// Register getters and setters
// =======================================

func (cpu *CPU) setA(value byte) { setHighByte(&cpu.af, value) }

func (cpu *CPU) setB(value byte) { setHighByte(&cpu.bc, value) }

func (cpu *CPU) setC(value byte) { setLowByte(&cpu.bc, value) }

func (cpu *CPU) setD(value byte) { setHighByte(&cpu.de, value) }

func (cpu *CPU) setE(value byte) { setLowByte(&cpu.de, value) }

func (cpu *CPU) setH(value byte) { setHighByte(&cpu.hl, value) }

func (cpu *CPU) setL(value byte) { setLowByte(&cpu.hl, value) }

func (cpu *CPU) A() byte { return getHighByte(cpu.af) }

// Note there is no getter for the F register as it is not used directly

func (cpu *CPU) B() byte { return getHighByte(cpu.bc) }

func (cpu *CPU) C() byte { return getLowByte(cpu.bc) }

func (cpu *CPU) D() byte { return getHighByte(cpu.de) }

func (cpu *CPU) E() byte { return getLowByte(cpu.de) }

func (cpu *CPU) H() byte { return getHighByte(cpu.hl) }

func (cpu *CPU) L() byte { return getLowByte(cpu.hl) }
