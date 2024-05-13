package gameboy

import (
	"log"
)

type CPU struct {
	// Registers
	AF, BC, DE, HL uint16

	// Stack pointer
	SP uint16

	// Program counter
	PC uint16

	// Memory
	mapper *Mapper

	// Special internal flags
	IME    bool
	halted bool

	// Debugging
	opDebug     string
	breakpoints []uint16
}

func NewCPU(mapper *Mapper) *CPU {
	// Initial state of the CPU for the classic GB
	// It represents the state of the CPU after the BIOS has run, as we skip that
	cpu := CPU{
		AF:     0x01B0,
		BC:     0x0013,
		DE:     0x00D8,
		HL:     0x014D,
		SP:     0xFFEE,
		PC:     0x0000,
		mapper: mapper,
	}

	cpu.setFlagZ(true)
	cpu.setFlagN(false)
	cpu.setFlagH(false)
	cpu.setFlagC(false)
	cpu.IME = false

	return &cpu
}

func (cpu *CPU) ExecuteNext(skipBreak bool) (cyclesSpent int) {
	if cpu.halted {
		// Even if halted, we still spend some cycles
		return 4
	}

	currentPC := cpu.PC

	// Fetch the next instruction, this will also increment the PC
	opcode := cpu.fetchPC()

	cpu.logMessage("%04X:%02X %s", currentPC, opcode, opcodeNames[opcode])

	// Check if we have hit a breakpoint
	for _, bp := range cpu.breakpoints {
		if bp == currentPC && !skipBreak {
			log.Printf(" !!! Breakpoint hit at 0x%04X\n", currentPC)
			cpu.PC = currentPC
			return -1
		}
	}

	// Check if the opcode is valid
	if opcodes[opcode] == nil {
		log.Printf(" !!! Unknown opcode: 0x%02X\n", opcode)
		cpu.PC = currentPC
		return -1
	}

	// Decode & execute the opcode
	opcodes[opcode](cpu)

	cycles := opcodeLengths[opcode]
	return cycles
}

func (cpu *CPU) handleInterrupt(interrupt byte) {
	// Disable interrupts
	cpu.IME = false

	// Clear the interrupt flag at the interrupt bit
	cpu.mapper.Write(IF, cpu.mapper.Read(IF)&^interrupt)

	// Push the current PC onto the stack
	cpu.pushStack(cpu.PC)

	// Jump to the interrupt handler
	switch interrupt {
	case 0x01:
		//log.Println(" !!! VBlank interrupt not implemented")
		cpu.PC = 0x0040
	case 0x02:
		cpu.PC = 0x0048
	case 0x04:
		cpu.PC = 0x0050
	case 0x08:
		cpu.PC = 0x0058
	case 0x10:
		//log.Println(" !!! Joypad interrupt not implemented")
		cpu.PC = 0x0060
	}
}

// =======================================
// Flag getters and setters
// =======================================

func (cpu *CPU) setFlagZ(value bool) {
	if value {
		cpu.AF |= 0x80
	} else {
		cpu.AF &^= 0x80
	}
}

func (cpu *CPU) setFlagN(value bool) {
	if value {
		cpu.AF |= 0x40
	} else {
		cpu.AF &^= 0x40
	}
}

func (cpu *CPU) setFlagH(value bool) {
	if value {
		cpu.AF |= 0x20
	} else {
		cpu.AF &^= 0x20
	}
}

func (cpu *CPU) setFlagC(value bool) {
	// In bit 4
	if value {
		cpu.AF |= 0x10
	} else {
		cpu.AF &^= 0x10
	}
}

func (cpu *CPU) getFlagZ() bool { return cpu.AF&0x80 != 0 }

func (cpu *CPU) getFlagN() bool { return cpu.AF&0x40 != 0 }

func (cpu *CPU) getFlagH() bool { return cpu.AF&0x20 != 0 }

func (cpu *CPU) getFlagC() bool { return cpu.AF&0x10 != 0 }

// =======================================
// Register getters and setters
// =======================================

func (cpu *CPU) setA(value byte) { setHighByte(&cpu.AF, value) }

func (cpu *CPU) setB(value byte) { setHighByte(&cpu.BC, value) }

func (cpu *CPU) setC(value byte) { setLowByte(&cpu.BC, value) }

func (cpu *CPU) setD(value byte) { setHighByte(&cpu.DE, value) }

func (cpu *CPU) setE(value byte) { setLowByte(&cpu.DE, value) }

func (cpu *CPU) setH(value byte) { setHighByte(&cpu.HL, value) }

func (cpu *CPU) setL(value byte) { setLowByte(&cpu.HL, value) }

func (cpu *CPU) A() byte { return getHighByte(cpu.AF) }

// Note there is no getter for the F register as it is not used directly

func (cpu *CPU) B() byte { return getHighByte(cpu.BC) }

func (cpu *CPU) C() byte { return getLowByte(cpu.BC) }

func (cpu *CPU) D() byte { return getHighByte(cpu.DE) }

func (cpu *CPU) E() byte { return getLowByte(cpu.DE) }

func (cpu *CPU) H() byte { return getHighByte(cpu.HL) }

func (cpu *CPU) L() byte { return getLowByte(cpu.HL) }

// =======================================
// Helpers
// =======================================

// Logs a message if logging is enabled
func (cpu *CPU) logMessage(s string, a ...any) {
	if logging {
		log.Printf(s, a...)
	}
}
