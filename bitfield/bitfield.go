package bitfield

type Bitfield []byte

// HasPiece returns whether a bitfield has a particular index set to 1
func (bf Bitfield) HasPiece(i int) bool {
	byteIndex := i / 8
	offset := i % 8

	if byteIndex < 0 || byteIndex >= len(bf) {
		return false
	}
	// since offset is a signed integer, >> extends (preserves) the sign
	// as an arithmetic right shift, same as >> in Java

	// Basically, get the section of 8 bits which is bf[byteIndex]
	// Shift until the specific bit you are looking at (bit at index i)
	// is the least significant digit (right-most bit).
	// Do a logical AND with 00000001.
	// Will return 00000000 or 00000001 depending on whether the right-most bit is 1

	// eg 00000000 0100000, we want index = 9 where the offset is 1 (0 based indexing)
	// therefore 01000000 >> 6 which becomes 00000001
	return (bf[byteIndex] >> uint(7-offset)) & 1 == 1
}

// SetPiece sets a bit in the bitfield
func (bf Bitfield) SetPiece(i int) {
	byteIndex := i / 8
	offset := i % 8
	// bf[byteIndex] = bf[byeIndex] | (1 << (7-offset))
	// get the specific 8 bits that we are looking at with bf[byteIndex]
	// set the specific bit to 1 with 00000001 << 7-offset

	// index = 9 (0 based)
	// byteIndex = 1 which is the second byte
	// where you want 00000000 01000000 given 00000000 00000000
	// the offset is 1
	// if 1 << 6, the result is 01000000
	// logical OR will set the 7th bit of 00000000 to 1 (and leave the rest of the bits as is)
	bf[byteIndex] |= 1 << uint(7 - offset)
}
