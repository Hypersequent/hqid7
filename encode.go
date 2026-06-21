package hqid7

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math/bits"

	"github.com/mr-tron/base58"
)

func NewString() string {
	return EncodeBase58(MustUUID7())
}

func encodeBase58Raw(u UUID) string {
	return base58.Encode(u[:])
}

const base58Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

var base58Pairs = func() [3364]uint16 {
	var pairs [3364]uint16
	for i := range pairs {
		pairs[i] = uint16(base58Alphabet[i/58])<<8 | uint16(base58Alphabet[i%58])
	}
	return pairs
}()

func EncodeBase58(u UUID) string {
	// A 128-bit UUID is always formatted as 22 Base58 digits. Divide four uint32
	// limbs by 58^4 chunks and write the fixed output positions directly.
	l0 := binary.BigEndian.Uint32(u[0:4])
	l1 := binary.BigEndian.Uint32(u[4:8])
	l2 := binary.BigEndian.Uint32(u[8:12])
	l3 := binary.BigEndian.Uint32(u[12:16])
	var out [23]byte
	out[9] = '_'

	cur := uint64(l0)
	l0 = uint32(cur / 11316496)
	rem := cur - uint64(l0)*11316496
	cur = (rem << 32) | uint64(l1)
	l1 = uint32(cur / 11316496)
	rem = cur - uint64(l1)*11316496
	cur = (rem << 32) | uint64(l2)
	l2 = uint32(cur / 11316496)
	rem = cur - uint64(l2)*11316496
	cur = (rem << 32) | uint64(l3)
	l3 = uint32(cur / 11316496)
	rem = cur - uint64(l3)*11316496
	loPair := base58Pairs[rem%3364]
	out[21] = byte(loPair >> 8)
	out[22] = byte(loPair)
	hiPair := base58Pairs[rem/3364]
	out[19] = byte(hiPair >> 8)
	out[20] = byte(hiPair)

	cur = uint64(l0)
	l0 = uint32(cur / 11316496)
	rem = cur - uint64(l0)*11316496
	cur = (rem << 32) | uint64(l1)
	l1 = uint32(cur / 11316496)
	rem = cur - uint64(l1)*11316496
	cur = (rem << 32) | uint64(l2)
	l2 = uint32(cur / 11316496)
	rem = cur - uint64(l2)*11316496
	cur = (rem << 32) | uint64(l3)
	l3 = uint32(cur / 11316496)
	rem = cur - uint64(l3)*11316496
	loPair = base58Pairs[rem%3364]
	out[17] = byte(loPair >> 8)
	out[18] = byte(loPair)
	hiPair = base58Pairs[rem/3364]
	out[15] = byte(hiPair >> 8)
	out[16] = byte(hiPair)

	cur = uint64(l0)
	l0 = uint32(cur / 11316496)
	rem = cur - uint64(l0)*11316496
	cur = (rem << 32) | uint64(l1)
	l1 = uint32(cur / 11316496)
	rem = cur - uint64(l1)*11316496
	cur = (rem << 32) | uint64(l2)
	l2 = uint32(cur / 11316496)
	rem = cur - uint64(l2)*11316496
	cur = (rem << 32) | uint64(l3)
	l3 = uint32(cur / 11316496)
	rem = cur - uint64(l3)*11316496
	loPair = base58Pairs[rem%3364]
	out[13] = byte(loPair >> 8)
	out[14] = byte(loPair)
	hiPair = base58Pairs[rem/3364]
	out[11] = byte(hiPair >> 8)
	out[12] = byte(hiPair)

	cur = uint64(l0)
	l0 = uint32(cur / 11316496)
	rem = cur - uint64(l0)*11316496
	cur = (rem << 32) | uint64(l1)
	l1 = uint32(cur / 11316496)
	rem = cur - uint64(l1)*11316496
	cur = (rem << 32) | uint64(l2)
	l2 = uint32(cur / 11316496)
	rem = cur - uint64(l2)*11316496
	cur = (rem << 32) | uint64(l3)
	l3 = uint32(cur / 11316496)
	rem = cur - uint64(l3)*11316496
	loPair = base58Pairs[rem%3364]
	out[8] = byte(loPair >> 8)
	out[10] = byte(loPair)
	hiPair = base58Pairs[rem/3364]
	out[6] = byte(hiPair >> 8)
	out[7] = byte(hiPair)

	cur = uint64(l0)
	l0 = uint32(cur / 11316496)
	rem = cur - uint64(l0)*11316496
	cur = (rem << 32) | uint64(l1)
	l1 = uint32(cur / 11316496)
	rem = cur - uint64(l1)*11316496
	cur = (rem << 32) | uint64(l2)
	l2 = uint32(cur / 11316496)
	rem = cur - uint64(l2)*11316496
	cur = (rem << 32) | uint64(l3)
	l3 = uint32(cur / 11316496)
	rem = cur - uint64(l3)*11316496
	loPair = base58Pairs[rem%3364]
	out[4] = byte(loPair >> 8)
	out[5] = byte(loPair)
	hiPair = base58Pairs[rem/3364]
	out[2] = byte(hiPair >> 8)
	out[3] = byte(hiPair)

	cur = uint64(l0)
	l0 = uint32(cur / 11316496)
	rem = cur - uint64(l0)*11316496
	cur = (rem << 32) | uint64(l1)
	l1 = uint32(cur / 11316496)
	rem = cur - uint64(l1)*11316496
	cur = (rem << 32) | uint64(l2)
	l2 = uint32(cur / 11316496)
	rem = cur - uint64(l2)*11316496
	cur = (rem << 32) | uint64(l3)
	l3 = uint32(cur / 11316496)
	rem = cur - uint64(l3)*11316496
	loPair = base58Pairs[rem%3364]
	out[0] = byte(loPair >> 8)
	out[1] = byte(loPair)

	return string(out[:])
}

var base58Decode = func() [256]byte {
	var decode [256]byte
	for i := range decode {
		decode[i] = 0xff
	}
	for i := 0; i < len(base58Alphabet); i++ {
		decode[base58Alphabet[i]] = byte(i)
	}
	return decode
}()

var base58DecodePairs = func() [65536]uint16 {
	var pairs [65536]uint16
	for c0, v0 := range base58Decode {
		if v0 == 0xff {
			continue
		}
		for c1, v1 := range base58Decode {
			if v1 != 0xff {
				pairs[c0<<8|c1] = uint16(v0)*58 + uint16(v1) + 1
			}
		}
	}
	return pairs
}()

func DecodeBase58(s string) (UUID, error) {
	if len(s) != 23 {
		return UUID{}, errors.New("hqid7 base58: invalid length")
	}
	if s[9] != '_' {
		return UUID{}, errors.New("hqid7 base58: invalid separator")
	}

	// Decode the fixed 22 Base58 digits as one leading 2-digit chunk followed by
	// five 4-digit chunks. hqid7 keeps the low 128 bits after removing leading
	// Base58 zeroes, so arithmetic wraps above the UUID width to match the
	// original decoder.
	p0 := base58DecodePairs[uint16(s[0])<<8|uint16(s[1])]
	if p0 == 0 {
		return UUID{}, invalidBase58Pair(s[0], s[1])
	}
	var hi uint64
	lo := uint64(p0 - 1)

	p0 = base58DecodePairs[uint16(s[2])<<8|uint16(s[3])]
	if p0 == 0 {
		return UUID{}, invalidBase58Pair(s[2], s[3])
	}
	p1 := base58DecodePairs[uint16(s[4])<<8|uint16(s[5])]
	if p1 == 0 {
		return UUID{}, invalidBase58Pair(s[4], s[5])
	}
	chunk := uint64(p0-1)*3364 + uint64(p1-1)
	loCarry, loProd := bits.Mul64(lo, 11316496)
	lo = loProd + chunk
	if lo < loProd {
		loCarry++
	}
	hi = hi*11316496 + loCarry

	p0 = base58DecodePairs[uint16(s[6])<<8|uint16(s[7])]
	if p0 == 0 {
		return UUID{}, invalidBase58Pair(s[6], s[7])
	}
	p1 = base58DecodePairs[uint16(s[8])<<8|uint16(s[10])]
	if p1 == 0 {
		return UUID{}, invalidBase58Pair(s[8], s[10])
	}
	chunk = uint64(p0-1)*3364 + uint64(p1-1)
	loCarry, loProd = bits.Mul64(lo, 11316496)
	lo = loProd + chunk
	if lo < loProd {
		loCarry++
	}
	hi = hi*11316496 + loCarry

	p0 = base58DecodePairs[uint16(s[11])<<8|uint16(s[12])]
	if p0 == 0 {
		return UUID{}, invalidBase58Pair(s[11], s[12])
	}
	p1 = base58DecodePairs[uint16(s[13])<<8|uint16(s[14])]
	if p1 == 0 {
		return UUID{}, invalidBase58Pair(s[13], s[14])
	}
	chunk = uint64(p0-1)*3364 + uint64(p1-1)
	loCarry, loProd = bits.Mul64(lo, 11316496)
	lo = loProd + chunk
	if lo < loProd {
		loCarry++
	}
	hi = hi*11316496 + loCarry

	p0 = base58DecodePairs[uint16(s[15])<<8|uint16(s[16])]
	if p0 == 0 {
		return UUID{}, invalidBase58Pair(s[15], s[16])
	}
	p1 = base58DecodePairs[uint16(s[17])<<8|uint16(s[18])]
	if p1 == 0 {
		return UUID{}, invalidBase58Pair(s[17], s[18])
	}
	chunk = uint64(p0-1)*3364 + uint64(p1-1)
	loCarry, loProd = bits.Mul64(lo, 11316496)
	lo = loProd + chunk
	if lo < loProd {
		loCarry++
	}
	hi = hi*11316496 + loCarry

	p0 = base58DecodePairs[uint16(s[19])<<8|uint16(s[20])]
	if p0 == 0 {
		return UUID{}, invalidBase58Pair(s[19], s[20])
	}
	p1 = base58DecodePairs[uint16(s[21])<<8|uint16(s[22])]
	if p1 == 0 {
		return UUID{}, invalidBase58Pair(s[21], s[22])
	}
	chunk = uint64(p0-1)*3364 + uint64(p1-1)
	loCarry, loProd = bits.Mul64(lo, 11316496)
	lo = loProd + chunk
	if lo < loProd {
		loCarry++
	}
	hi = hi*11316496 + loCarry

	var uuid UUID
	binary.BigEndian.PutUint64(uuid[0:8], hi)
	binary.BigEndian.PutUint64(uuid[8:16], lo)
	return uuid, nil
}

func invalidBase58Pair(c0, c1 byte) error {
	if base58Decode[c0] == 0xff {
		return invalidBase58Digit(c0)
	}
	return invalidBase58Digit(c1)
}

func invalidBase58Quad(c0, c1, c2, c3 byte) error {
	if base58Decode[c0] == 0xff {
		return invalidBase58Digit(c0)
	}
	if base58Decode[c1] == 0xff {
		return invalidBase58Digit(c1)
	}
	if base58Decode[c2] == 0xff {
		return invalidBase58Digit(c2)
	}
	return invalidBase58Digit(c3)
}

func invalidBase58Digit(c byte) error {
	if c > 127 {
		return fmt.Errorf("High-bit set on invalid digit")
	}
	return fmt.Errorf("Invalid base58 digit (%q)", rune(c))
}
