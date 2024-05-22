package gameboy

func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func BoolToByte(b bool) byte {
	if b {
		return 1
	}
	return 0
}

func setHighByte(reg *uint16, value byte) {
	*reg = uint16(value)<<8 | *reg&0x00ff
}

func setLowByte(reg *uint16, value byte) {
	*reg = uint16(value) | *reg&0xff00
}

func getHighByte(reg uint16) byte {
	return byte(reg >> 8)
}

func getLowByte(reg uint16) byte {
	return byte(reg & 0xff)
}

// Set given bit in the byte to 0
func bitReset(b byte, bit uint) byte {
	return b &^ (1 << bit)
}

// Set given bit in the byte to 1
func bitSet(b byte, bit uint) byte {
	return b | (1 << bit)
}

func twoBitValue(b byte, bit uint) int {
	// get two bits at position bit and bit+1 and return them as a int from 0 to 3
	return int((b >> bit) & 0x3)
}

func checkBit(b byte, bit uint) bool {
	return b&(1<<bit) != 0
}
