package dwg

import "fmt"

// decryptR2004Header reverses the fixed LCG-based XOR obfuscation AutoCAD
// applies to the R2004+ file header at offset 0x80. The keystream depends
// only on position, not on the source bytes, so it can be computed for any
// length independent of what's decrypted.
func decryptR2004Header(buf []byte) []byte {
	dst := make([]byte, len(buf))
	var seed uint32 = 1
	for i := range buf {
		seed = seed*0x343fd + 0x269ec3
		dst[i] = buf[i] ^ byte(seed>>0x10)
	}
	return dst
}

// decompressR2004Section implements AutoCAD's R2004+ LZ77 variant. It reads
// compressed bytes from src and writes decompressed output into dst starting
// at dst[dstOff:], stopping at the terminator opcode (0x11), when src is
// exhausted, or when dst fills up.
func decompressR2004Section(src []byte, dst []byte, dstOff int) error {
	sp := 0
	dp := dstOff

	readSrc := func() byte {
		if sp >= len(src) {
			return 0
		}
		b := src[sp]
		sp++
		return b
	}

	readLiteralLength := func(opcode byte) int {
		lowbits := int(opcode & 0xf)
		if lowbits == 0 {
			var lastbyte byte
			for {
				lastbyte = readSrc()
				if lastbyte != 0 || sp >= len(src) {
					break
				}
				lowbits += 0xFF
			}
			lowbits += 0xf + int(lastbyte)
		}
		return lowbits + 3
	}

	readCompressedBytes := func(opcode byte, bits int) int {
		compressedBytes := int(opcode) & bits
		if compressedBytes == 0 {
			var lastbyte byte
			for {
				lastbyte = readSrc()
				if lastbyte != 0 || sp >= len(src) {
					break
				}
				compressedBytes += 0xFF
			}
			compressedBytes += int(lastbyte) + bits
		}
		return compressedBytes + 2
	}

	twoByteOffset := func(plus int, offset *int) byte {
		first := readSrc()
		second := readSrc()
		*offset |= int(first) >> 2
		*offset |= int(second) << 6
		*offset += plus
		return first
	}

	copyBytes := func(n int) byte {
		for i := 0; i < n; i++ {
			b := readSrc()
			if dp < len(dst) {
				dst[dp] = b
			}
			dp++
		}
		return readSrc()
	}

	if sp > len(src) {
		return fmt.Errorf("dwg: invalid compressed section")
	}

	opcode1 := readSrc()
	if opcode1&0xF0 == 0 {
		opcode1 = copyBytes(readLiteralLength(opcode1))
	}

	for sp < len(src) && dp < len(dst) && opcode1 != 0x11 {
		compBytes := 0
		compOffset := 0
		switch {
		case opcode1 < 0x10 || opcode1 >= 0x40:
			compBytes = int(opcode1>>4) - 1
			opcode2 := readSrc()
			compOffset = (((int(opcode1) >> 2) & 3) | (int(opcode2) << 2)) + 1
		case opcode1 < 0x20:
			compBytes = readCompressedBytes(opcode1, 7)
			compOffset = (int(opcode1) & 8) << 11
			opcode1 = twoByteOffset(0x4000, &compOffset)
		default: // opcode1 >= 0x20
			compBytes = readCompressedBytes(opcode1, 0x1f)
			opcode1 = twoByteOffset(1, &compOffset)
		}

		pos := dp
		end := pos + compBytes
		if end > len(dst) || pos < compOffset || pos-compOffset >= len(dst) || compOffset > len(dst) {
			return fmt.Errorf("dwg: invalid decompression bytes %d, offset %d", compBytes, compOffset)
		}
		for ; pos < end; pos++ {
			dst[pos] = dst[pos-compOffset]
		}
		dp = end

		litLength := int(opcode1) & 3
		if litLength == 0 {
			opcode1 = readSrc()
			if opcode1&0xf0 == 0 {
				litLength = readLiteralLength(opcode1)
			}
		}
		if litLength != 0 && end+litLength <= len(dst) {
			opcode1 = copyBytes(litLength)
		} else if litLength != 0 {
			break
		}
	}
	return nil
}
