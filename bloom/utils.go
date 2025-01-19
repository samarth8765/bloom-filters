package bloom

func setBit(buf []byte, pos uint64) {
	idx, offset := pos/8, pos%8
	buf[idx] |= 1 << offset
}

func getBit(buf []byte, pos uint64) byte {
	idx, offset := pos/8, pos%8
	return buf[idx] & (1 << offset)
}
