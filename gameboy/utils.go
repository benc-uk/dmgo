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

func halfCarryAdd(a, b byte) bool {
	return (a&0xf)+(b&0xf) > 0xf
}

func halfCarrySub(a, b byte) bool {
	return (a & 0xf) < (b & 0xf)
}
