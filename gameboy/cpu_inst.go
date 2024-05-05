package gameboy

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
	// log.Printf(">>> Calling %04X from PC:%04X\n", addr, cpu.PC)
	cpu.pushStack(cpu.PC)
	cpu.PC = addr
}

// Pushes a 16-bit value onto the stack, often the PC but can be any value
func (cpu *CPU) pushStack(addr uint16) {
	// log.Printf(">>>> Pushing %04X to stack at SP:%04X\n", addr, cpu.SP)
	sp := cpu.SP
	cpu.mapper.Write(sp-1, byte(uint16(addr&0xFF00)>>8))
	cpu.mapper.Write(sp-2, byte(addr&0xFF))
	cpu.SP -= 2
}

// Returns from a subroutine by popping the address from the stack
func (cpu *CPU) returnSub() {
	pc := cpu.popStack()
	// log.Printf("Returning to %04X from SP:%04X\n", pc, cpu.SP)
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

// Performs 8-bit addition between two bytes and sets the flags accordingly
func (cpu *CPU) byteAdd(a, b byte) byte {
	cpu.logMessage("--- ADD: a:%02X, b:%02X", a, b)
	result := a + b
	carry := (uint16(a) + uint16(b)) > 0xFF

	cpu.setFlagZ(result == 0)
	cpu.setFlagN(false)
	cpu.setFlagH((a&0xF)+(b&0xF) > 0xF)
	cpu.setFlagC(carry)

	return result
}

func (cpu *CPU) wordAdd(a, b uint16) uint16 {
	cpu.logMessage("--- WADD: a:%04X, b:%04X", a, b)
	result := a + b
	carry := (a + b) > 0xFFFF

	cpu.setFlagN(false)
	cpu.setFlagH((a&0xFFF)+(b&0xFFF) > 0xFFF)
	cpu.setFlagC(carry)

	return result
}

// Performs 8-bit subtraction between two bytes and sets the flags accordingly
func (cpu *CPU) byteSub(a, b byte) byte {
	cpu.logMessage("--- SUB: a:%02X, b:%02X", a, b)
	result := a - b
	carry := a < b

	cpu.setFlagZ(result == 0)
	cpu.setFlagN(true)
	cpu.setFlagH(a&0xF < b&0xF)
	cpu.setFlagC(carry)

	return result
}

// Performs 8-bit increment on a byte and sets the flags accordingly
func (cpu *CPU) byteInc(a byte) byte {
	cpu.logMessage("--- INC: a:%02X", a)
	result := a + 1
	cpu.setFlagZ(result == 0)
	cpu.setFlagN(false)

	// Set if overflow from bit 3 (?)	(not sure about this one)
	cpu.setFlagH(a&0xF == 0xF)

	return result
}

// Performs 8-bit decrement on a byte and sets the flags accordingly
func (cpu *CPU) byteDec(a byte) byte {
	cpu.logMessage("--- DEC: a:%02X", a)
	result := a - 1
	cpu.setFlagZ(result == 0)
	cpu.setFlagN(true)

	// Set if no borrow from bit 4 (?) (not sure about this one)
	cpu.setFlagH((a & 0xF) < 1)

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

// Bit test on a register
func (cpu *CPU) bitTest(reg byte, bit uint) {
	cpu.logMessage("--- BIT: %d, %02X", bit, reg)
	cpu.setFlagZ(reg&(1<<bit) == 0)
	cpu.setFlagN(false)
	cpu.setFlagH(true)
}

func (cpu *CPU) rotLeft(val byte) byte {
	cpu.logMessage("--- RL: %02X", val)
	newCarry := val >> 7
	oldCarry := byte(BoolToInt(cpu.getFlagC()))
	rot := (val<<1)&0xFF | oldCarry

	cpu.setFlagZ(rot == 0)
	cpu.setFlagN(false)
	cpu.setFlagH(false)
	cpu.setFlagC(newCarry == 1)

	return rot
}

func (cpu *CPU) shiftLeftArithmetic(val byte) byte {
	cpu.logMessage("--- SLA: %02X", val)
	newCarry := val >> 7
	shifted := val << 1

	cpu.setFlagZ(shifted == 0)
	cpu.setFlagN(false)
	cpu.setFlagH(false)
	cpu.setFlagC(newCarry == 1)

	return shifted
}
