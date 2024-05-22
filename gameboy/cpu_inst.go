package gameboy

// Fetches the next byte from memory and increments the program counter
func (cpu *CPU) fetchPC() byte {
	v := cpu.mapper.read(cpu.pc)
	cpu.pc += 1
	return v
}

// Fetches the next 16-bit word from memory and increments the program counter
func (cpu *CPU) fetchPC16() uint16 {
	lo := cpu.mapper.read(cpu.pc)
	hi := cpu.mapper.read(cpu.pc + 1)
	cpu.pc += 2
	return uint16(hi)<<8 | uint16(lo)
}

// Used to call a subroutine at the given address
func (cpu *CPU) callSub(addr uint16) {
	cpu.pushStack(cpu.pc)
	cpu.pc = addr
}

// Pushes a 16-bit value onto the stack, often the PC but can be any value
func (cpu *CPU) pushStack(addr uint16) {
	sp := cpu.sp
	cpu.mapper.write(sp-1, byte(uint16(addr&0xFF00)>>8))
	cpu.mapper.write(sp-2, byte(addr&0xFF))
	cpu.sp -= 2
}

// Returns from a subroutine by popping the address from the stack
func (cpu *CPU) returnSub() {
	pc := cpu.popStack()
	cpu.pc = pc
}

// Pops a 16-bit value from the stack
func (cpu *CPU) popStack() uint16 {
	sp := cpu.sp
	lo := cpu.mapper.read(sp)
	hi := cpu.mapper.read(sp + 1)
	cpu.sp += 2
	return uint16(hi)<<8 | uint16(lo)
}

// Performs an OR operation on two bytes and sets the flags accordingly & returns the result
func (cpu *CPU) byteOR(a, b byte) byte {
	result := a | b
	cpu.setFlagZ(result == 0)
	cpu.setFlagN(false)
	cpu.setFlagH(false)
	cpu.setFlagC(false)
	return a | b
}

// Performs an AND operation on two bytes and sets the flags accordingly & returns the result
func (cpu *CPU) byteAND(a, b byte) byte {
	result := a & b
	cpu.setFlagZ(result == 0)
	cpu.setFlagN(false)
	cpu.setFlagH(true)
	cpu.setFlagC(false)
	return result
}

// Performs an XOR operation on two bytes and sets the flags accordingly & returns the result
func (cpu *CPU) byteXOR(a, b byte) byte {
	result := a ^ b
	cpu.setFlagZ(result == 0)
	cpu.setFlagN(false)
	cpu.setFlagH(false)
	cpu.setFlagC(false)
	return result
}

// Performs 8-bit addition with carry between two bytes and sets the flags accordingly
func (cpu *CPU) byteAdd(a, b byte, addCarry bool) byte {
	carry := int16(BoolToByte(cpu.getFlagC() && addCarry))
	result := int16(a) + int16(b) + carry

	cpu.setFlagZ(byte(result) == 0)
	cpu.setFlagN(false)
	cpu.setFlagH((b&0xF)+(a&0xF)+byte(carry) > 0xF)
	cpu.setFlagC(result > 0xFF) // If result is greater than 255

	return byte(result)
}

// Performs 16-bit addition with carry between two words and sets the flags accordingly
func (cpu *CPU) wordAdd(a, b uint16) uint16 {
	result := a + b
	cpu.setFlagN(false)
	cpu.setFlagH((a&0xFFF)+(b&0xFFF) > 0xFFF)
	cpu.setFlagC(result < a)

	return result
}

// Performs 8-bit subtraction with carry between two bytes and sets the flags accordingly
func (cpu *CPU) byteSub(a, b byte, addCarry bool) byte {
	carry := int16(BoolToByte(cpu.getFlagC() && addCarry))
	dirtySum := int16(a) - int16(b) - carry
	result := byte(dirtySum)

	cpu.setFlagZ(result == 0)
	cpu.setFlagN(true)
	cpu.setFlagH(int16(a&0x0f)-int16(b&0xF)-int16(carry) < 0)
	cpu.setFlagC(dirtySum < 0)

	return result
}

// Performs 8-bit increment on a byte and sets the flags accordingly
func (cpu *CPU) byteInc(a byte) byte {
	result := a + 1
	cpu.setFlagZ(result == 0)
	cpu.setFlagN(false)

	// Set if overflow from bit 3 (?)	(not sure about this one)
	cpu.setFlagH(a&0xF == 0xF)

	return result
}

// Performs 8-bit decrement on a byte and sets the flags accordingly
func (cpu *CPU) byteDec(a byte) byte {
	result := a - 1
	cpu.setFlagZ(result == 0)
	cpu.setFlagN(true)

	// Set if no borrow from bit 4 (?) (not sure about this one)
	cpu.setFlagH((a & 0xF) < 1)

	return result
}

// Performs comparison between two bytes sets the flags accordingly
func (cpu *CPU) cmp(a, b byte) {
	result := a - b
	cpu.setFlagZ(result == 0)
	cpu.setFlagN(true)
	cpu.setFlagH((a & 0xF) < (b & 0xF))
	cpu.setFlagC(a < b)
}

// Bit test on a register
func (cpu *CPU) bitTest(reg byte, bit uint) {
	cpu.setFlagZ(reg&(1<<bit) == 0)
	cpu.setFlagN(false)
	cpu.setFlagH(true)
}

// Rotate left through carry and sets the flags accordingly
func (cpu *CPU) rotLeft(val byte) byte {
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
	newCarry := val >> 7
	shifted := val << 1

	cpu.setFlagZ(shifted == 0)
	cpu.setFlagN(false)
	cpu.setFlagH(false)
	cpu.setFlagC(newCarry == 1)

	return shifted
}

func (cpu *CPU) rotRight(val byte) byte {
	newCarry := val & 1
	oldCarry := byte(BoolToInt(cpu.getFlagC()))
	rot := (val >> 1) | (oldCarry << 7)

	cpu.setFlagZ(rot == 0)
	cpu.setFlagN(false)
	cpu.setFlagH(false)
	cpu.setFlagC(newCarry == 1)

	return rot
}

func (cpu *CPU) rotRightCarry(val byte) byte {
	newCarry := val & 1
	rot := (val >> 1) | (val << 7)

	cpu.setFlagZ(rot == 0)
	cpu.setFlagN(false)
	cpu.setFlagH(false)
	cpu.setFlagC(newCarry == 1)

	return rot
}

func (cpu *CPU) shiftRightArithmetic(val byte) byte {
	newCarry := val & 1
	shifted := (val >> 1) | (val & 0x80)

	cpu.setFlagZ(shifted == 0)
	cpu.setFlagN(false)
	cpu.setFlagH(false)
	cpu.setFlagC(newCarry == 1)

	return shifted
}

func (cpu *CPU) shiftRightLogical(val byte) byte {
	newCarry := val & 1
	shifted := val >> 1

	cpu.setFlagZ(shifted == 0)
	cpu.setFlagN(false)
	cpu.setFlagH(false)
	cpu.setFlagC(newCarry == 1)

	return shifted
}

func (cpu *CPU) rotLeftCarry(val byte) byte {
	newCarry := val >> 7
	rot := (val << 1) | newCarry

	cpu.setFlagZ(rot == 0)
	cpu.setFlagN(false)
	cpu.setFlagH(false)
	cpu.setFlagC(newCarry == 1)

	return rot
}

// Swaps the nibbles of a byte
func (cpu *CPU) swapNibbles(val byte) byte {
	swapped := val<<4 | val>>4

	cpu.setFlagZ(swapped == 0)
	cpu.setFlagN(false)
	cpu.setFlagH(false)
	cpu.setFlagC(false)

	return swapped
}
