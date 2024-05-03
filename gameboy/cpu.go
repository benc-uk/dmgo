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

	opDebug string

	// Interrupts
	IME bool

	breakpoint uint16
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
		PC:     0x0100,
		mapper: mapper,
	}

	cpu.setFlagZ(true)
	cpu.setFlagN(false)
	cpu.setFlagH(false)
	cpu.setFlagC(false)
	cpu.IME = false

	return &cpu
}

func (cpu *CPU) ExecuteNext(ignoreBreak bool) bool {
	oldPC := cpu.PC

	// Fetch the next instruction, this will also increment the PC
	opcode := cpu.fetchPC()

	cpu.logMessage("%04X:%02X %s", oldPC, opcode, opcodeNames[opcode])
	cpu.opDebug = opcodeNames[opcode]

	if oldPC == cpu.breakpoint && !ignoreBreak {
		log.Printf(">>> Breakpoint hit at %04X\n", oldPC)
		cpu.PC--
		return false
	}

	// Check if the opcode is valid
	if opcodes[opcode] == nil {
		log.Printf(" !!! Unknown opcode: 0x%02X\n", opcode)
		cpu.PC--
		return false
	}

	// Decode & execute the opcode
	opcodes[opcode](cpu)

	return true
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

// Fetches the next byte from memory and increments the program counter
func (cpu *CPU) fetchPC() byte {
	v := cpu.mapper.Read(cpu.PC)
	cpu.PC += 1
	return v
}

// Fetches the next 16-bit word from memory and increments the program counter
func (cpu *CPU) fetchPC16() uint16 {
	lo := cpu.mapper.Read(cpu.PC)
	hi := cpu.mapper.Read(cpu.PC + 1)
	cpu.PC += 2
	return uint16(hi)<<8 | uint16(lo)
}

// Used to call a subroutine at the given address
func (cpu *CPU) callSub(addr uint16) {
	log.Printf(">>> Calling %04X from PC:%04X\n", addr, cpu.PC)
	cpu.pushStack(cpu.PC)
	cpu.PC = addr
}

// Pushes a 16-bit value onto the stack, often the PC but can be any value
func (cpu *CPU) pushStack(addr uint16) {
	log.Printf(">>>> Pushing %04X to stack at SP:%04X\n", addr, cpu.SP)

	cpu.mapper.Write(cpu.SP-1, byte(addr>>8))
	cpu.mapper.Write(cpu.SP-2, byte(addr&0xFF))
	cpu.SP -= 2
}

// Returns from a subroutine by popping the address from the stack
func (cpu *CPU) returnSub() {
	pc := cpu.popStack()
	log.Printf("Returning to %04X from SP:%04X\n", pc, cpu.SP)
	cpu.PC = pc
}

// Pops a 16-bit value from the stack
func (cpu *CPU) popStack() uint16 {
	sp := cpu.SP
	lo := cpu.mapper.Read(sp)
	hi := cpu.mapper.Read(sp + 1)
	cpu.SP += 2
	return uint16(hi)<<8 | uint16(lo)
}

// Performs an OR operation on two bytes and sets the flags accordingly & returns the result
func (cpu *CPU) byteOR(a, b byte) byte {
	cpu.logMessage("--- OR: a:%02X, b:%02X", a, b)
	result := a | b
	cpu.setFlagZ(result == 0)
	cpu.setFlagN(false)
	cpu.setFlagH(false)
	cpu.setFlagC(false)
	return a | b
}

// Performs an AND operation on two bytes and sets the flags accordingly & returns the result
func (cpu *CPU) byteAND(a, b byte) byte {
	cpu.logMessage("--- AND: a:%02X, b:%02X", a, b)
	result := a & b
	cpu.setFlagZ(result == 0)
	cpu.setFlagN(false)
	cpu.setFlagH(true)
	cpu.setFlagC(false)
	return result
}

// Performs an XOR operation on two bytes and sets the flags accordingly & returns the result
func (cpu *CPU) byteXOR(a, b byte) byte {
	cpu.logMessage("--- XOR: a:%02X, b:%02X", a, b)
	result := a ^ b
	cpu.setFlagZ(result == 0)
	cpu.setFlagN(false)
	cpu.setFlagH(false)
	cpu.setFlagC(false)
	return result
}

// Performs addition between two bytes and sets the flags accordingly
func (cpu *CPU) byteAdd(a, b byte) byte {
	cpu.logMessage("--- ADD: a:%02X, b:%02X", a, b)
	result := a + b
	cpu.setFlagZ(result == 0)
	cpu.setFlagN(false)
	cpu.setFlagH(halfCarryAdd(a, b))
	// Note that the carry flag is left unchanged
	return result
}

// Performs subtraction between two bytes and sets the flags accordingly
func (cpu *CPU) byteSub(a, b byte) byte {
	cpu.logMessage("--- SUB: a:%02X, b:%02X", a, b)
	result := a - b
	cpu.setFlagZ(result == 0)
	cpu.setFlagN(true)
	cpu.setFlagH(halfCarrySub(a, b))
	// Note that the carry flag is left unchanged
	return result
}

// Performs comparison between two bytes sets the flags accordingly
func (cpu *CPU) cmp(a, b byte) {
	cpu.logMessage("--- CMP: a:%02X, b:%02X", a, b)
	result := a - b
	cpu.setFlagZ(result == 0)
	cpu.setFlagN(true)
	cpu.setFlagH((a & 0xF) < (b & 0xF))
	cpu.setFlagC(a < b)
}
