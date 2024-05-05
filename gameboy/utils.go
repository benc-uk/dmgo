package gameboy

func BoolToInt(b bool) int {
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

func swapNibbles(b byte) byte {
	return b<<4 | b>>4
}

// Set given bit in the byte to 0
func bitReset(b byte, bit uint) byte {
	return b &^ (1 << bit)
}

// Set given bit in the byte to 1
func bitSet(b byte, bit uint) byte {
	return b | (1 << bit)
}
