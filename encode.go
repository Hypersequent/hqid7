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

const (
	base58Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

	// The encoder/decoder works two digits at a time. A pair has 58^2 possible
	// values; a 4-digit chunk has 58^4 possible values and fits comfortably in
	// uint32. Splitting 22 output digits into one 2-digit prefix plus five
	// 4-digit chunks keeps all intermediate division state in uint64.
	base58PairBase  = 58 * 58
	base58ChunkBase = base58PairBase * base58PairBase

	base58EncodedDigits = 22
	base58SeparatorAt   = 9
)

var base58Pairs = func() [base58PairBase]uint16 {
	var pairs [base58PairBase]uint16
	for i := range pairs {
		pairs[i] = uint16(base58Alphabet[i/58])<<8 | uint16(base58Alphabet[i%58])
	}
	return pairs
}()

func EncodeBase58(u UUID) string {
	// A 128-bit UUID fits in 22 Base58 digits. This function is the allocation-
	// free equivalent of:
	//
	//   s := base58.Encode(u[:])
	//   s = strings.Repeat("1", 22-len(s)) + s
	//   return s[:9] + "_" + s[9:]
	//
	// Instead of asking the generic Base58 package to divide a byte slice, the
	// UUID is treated as four big-endian uint32 limbs. Each unrolled block below
	// divides the whole 128-bit value by 58^4, stores the remainder as four
	// Base58 digits, and carries the quotient forward for the next more
	// significant chunk. Remainders come out least-significant first, so output
	// is filled from the right edge back toward the left.
	l0 := binary.BigEndian.Uint32(u[0:4])
	l1 := binary.BigEndian.Uint32(u[4:8])
	l2 := binary.BigEndian.Uint32(u[8:12])
	l3 := binary.BigEndian.Uint32(u[12:16])
	var out [23]byte
	out[base58SeparatorAt] = '_'

	cur := uint64(l0)
	l0 = uint32(cur / base58ChunkBase)
	rem := cur - uint64(l0)*base58ChunkBase
	cur = (rem << 32) | uint64(l1)
	l1 = uint32(cur / base58ChunkBase)
	rem = cur - uint64(l1)*base58ChunkBase
	cur = (rem << 32) | uint64(l2)
	l2 = uint32(cur / base58ChunkBase)
	rem = cur - uint64(l2)*base58ChunkBase
	cur = (rem << 32) | uint64(l3)
	l3 = uint32(cur / base58ChunkBase)
	rem = cur - uint64(l3)*base58ChunkBase
	loPair := base58Pairs[rem%base58PairBase]
	out[21] = byte(loPair >> 8)
	out[22] = byte(loPair)
	hiPair := base58Pairs[rem/base58PairBase]
	out[19] = byte(hiPair >> 8)
	out[20] = byte(hiPair)

	cur = uint64(l0)
	l0 = uint32(cur / base58ChunkBase)
	rem = cur - uint64(l0)*base58ChunkBase
	cur = (rem << 32) | uint64(l1)
	l1 = uint32(cur / base58ChunkBase)
	rem = cur - uint64(l1)*base58ChunkBase
	cur = (rem << 32) | uint64(l2)
	l2 = uint32(cur / base58ChunkBase)
	rem = cur - uint64(l2)*base58ChunkBase
	cur = (rem << 32) | uint64(l3)
	l3 = uint32(cur / base58ChunkBase)
	rem = cur - uint64(l3)*base58ChunkBase
	loPair = base58Pairs[rem%base58PairBase]
	out[17] = byte(loPair >> 8)
	out[18] = byte(loPair)
	hiPair = base58Pairs[rem/base58PairBase]
	out[15] = byte(hiPair >> 8)
	out[16] = byte(hiPair)

	cur = uint64(l0)
	l0 = uint32(cur / base58ChunkBase)
	rem = cur - uint64(l0)*base58ChunkBase
	cur = (rem << 32) | uint64(l1)
	l1 = uint32(cur / base58ChunkBase)
	rem = cur - uint64(l1)*base58ChunkBase
	cur = (rem << 32) | uint64(l2)
	l2 = uint32(cur / base58ChunkBase)
	rem = cur - uint64(l2)*base58ChunkBase
	cur = (rem << 32) | uint64(l3)
	l3 = uint32(cur / base58ChunkBase)
	rem = cur - uint64(l3)*base58ChunkBase
	loPair = base58Pairs[rem%base58PairBase]
	out[13] = byte(loPair >> 8)
	out[14] = byte(loPair)
	hiPair = base58Pairs[rem/base58PairBase]
	out[11] = byte(hiPair >> 8)
	out[12] = byte(hiPair)

	// This 4-digit chunk straddles the visual separator. The first three digits
	// land before out[9], the final digit lands after it.
	cur = uint64(l0)
	l0 = uint32(cur / base58ChunkBase)
	rem = cur - uint64(l0)*base58ChunkBase
	cur = (rem << 32) | uint64(l1)
	l1 = uint32(cur / base58ChunkBase)
	rem = cur - uint64(l1)*base58ChunkBase
	cur = (rem << 32) | uint64(l2)
	l2 = uint32(cur / base58ChunkBase)
	rem = cur - uint64(l2)*base58ChunkBase
	cur = (rem << 32) | uint64(l3)
	l3 = uint32(cur / base58ChunkBase)
	rem = cur - uint64(l3)*base58ChunkBase
	loPair = base58Pairs[rem%base58PairBase]
	out[8] = byte(loPair >> 8)
	out[10] = byte(loPair)
	hiPair = base58Pairs[rem/base58PairBase]
	out[6] = byte(hiPair >> 8)
	out[7] = byte(hiPair)

	cur = uint64(l0)
	l0 = uint32(cur / base58ChunkBase)
	rem = cur - uint64(l0)*base58ChunkBase
	cur = (rem << 32) | uint64(l1)
	l1 = uint32(cur / base58ChunkBase)
	rem = cur - uint64(l1)*base58ChunkBase
	cur = (rem << 32) | uint64(l2)
	l2 = uint32(cur / base58ChunkBase)
	rem = cur - uint64(l2)*base58ChunkBase
	cur = (rem << 32) | uint64(l3)
	l3 = uint32(cur / base58ChunkBase)
	rem = cur - uint64(l3)*base58ChunkBase
	loPair = base58Pairs[rem%base58PairBase]
	out[4] = byte(loPair >> 8)
	out[5] = byte(loPair)
	hiPair = base58Pairs[rem/base58PairBase]
	out[2] = byte(hiPair >> 8)
	out[3] = byte(hiPair)

	// After five 4-digit chunks have been removed the quotient is less than
	// 58^2, because max(uint128) < 58^22. Only the final two leading digits are
	// written; the discarded high pair would always be "11".
	cur = uint64(l0)
	l0 = uint32(cur / base58ChunkBase)
	rem = cur - uint64(l0)*base58ChunkBase
	cur = (rem << 32) | uint64(l1)
	l1 = uint32(cur / base58ChunkBase)
	rem = cur - uint64(l1)*base58ChunkBase
	cur = (rem << 32) | uint64(l2)
	l2 = uint32(cur / base58ChunkBase)
	rem = cur - uint64(l2)*base58ChunkBase
	cur = (rem << 32) | uint64(l3)
	l3 = uint32(cur / base58ChunkBase)
	rem = cur - uint64(l3)*base58ChunkBase
	loPair = base58Pairs[rem%base58PairBase]
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
	for i := range pairs {
		pairs[i] = 0xffff
	}
	// Index by two raw input bytes packed into uint16. Valid Base58 pairs store
	// their numeric value; every other byte combination keeps the sentinel.
	for c0, v0 := range base58Decode {
		if v0 == 0xff {
			continue
		}
		for c1, v1 := range base58Decode {
			if v1 != 0xff {
				pairs[c0<<8|c1] = uint16(v0)*58 + uint16(v1)
			}
		}
	}
	return pairs
}()

func DecodeBase58(s string) (UUID, error) {
	if len(s) != 23 {
		return UUID{}, errors.New("hqid7 base58: invalid length")
	}
	if s[base58SeparatorAt] != '_' {
		return UUID{}, errors.New("hqid7 base58: invalid separator")
	}

	// Decode the fixed 22 Base58 digits as one leading 2-digit chunk followed
	// by five 4-digit chunks. Each 4-digit chunk is formed from two lookup-table
	// pairs, then accumulated as:
	//
	//   value = value*58^4 + chunk
	//
	// The accumulator is held as hi:lo uint64 halves. bits.Mul64 gives the low
	// half of lo*58^4 plus the carry into hi; the explicit addition handles a
	// carry from adding the new chunk to lo. If the incoming string represents a
	// value above 128 bits, normal uint64 overflow preserves the same low 128
	// bits that the previous generic decoder kept after trimming leading bytes.
	p0 := base58DecodePairs[uint16(s[0])<<8|uint16(s[1])]
	if p0 == 0xffff {
		return UUID{}, invalidBase58Pair(s[0], s[1])
	}
	var hi uint64
	lo := uint64(p0)

	p0 = base58DecodePairs[uint16(s[2])<<8|uint16(s[3])]
	if p0 == 0xffff {
		return UUID{}, invalidBase58Pair(s[2], s[3])
	}
	p1 := base58DecodePairs[uint16(s[4])<<8|uint16(s[5])]
	if p1 == 0xffff {
		return UUID{}, invalidBase58Pair(s[4], s[5])
	}
	chunk := uint64(p0)*base58PairBase + uint64(p1)
	loCarry, loProd := bits.Mul64(lo, base58ChunkBase)
	lo = loProd + chunk
	if lo < loProd {
		loCarry++
	}
	hi = hi*base58ChunkBase + loCarry

	p0 = base58DecodePairs[uint16(s[6])<<8|uint16(s[7])]
	if p0 == 0xffff {
		return UUID{}, invalidBase58Pair(s[6], s[7])
	}
	p1 = base58DecodePairs[uint16(s[8])<<8|uint16(s[10])]
	if p1 == 0xffff {
		return UUID{}, invalidBase58Pair(s[8], s[10])
	}
	chunk = uint64(p0)*base58PairBase + uint64(p1)
	loCarry, loProd = bits.Mul64(lo, base58ChunkBase)
	lo = loProd + chunk
	if lo < loProd {
		loCarry++
	}
	hi = hi*base58ChunkBase + loCarry

	p0 = base58DecodePairs[uint16(s[11])<<8|uint16(s[12])]
	if p0 == 0xffff {
		return UUID{}, invalidBase58Pair(s[11], s[12])
	}
	p1 = base58DecodePairs[uint16(s[13])<<8|uint16(s[14])]
	if p1 == 0xffff {
		return UUID{}, invalidBase58Pair(s[13], s[14])
	}
	chunk = uint64(p0)*base58PairBase + uint64(p1)
	loCarry, loProd = bits.Mul64(lo, base58ChunkBase)
	lo = loProd + chunk
	if lo < loProd {
		loCarry++
	}
	hi = hi*base58ChunkBase + loCarry

	p0 = base58DecodePairs[uint16(s[15])<<8|uint16(s[16])]
	if p0 == 0xffff {
		return UUID{}, invalidBase58Pair(s[15], s[16])
	}
	p1 = base58DecodePairs[uint16(s[17])<<8|uint16(s[18])]
	if p1 == 0xffff {
		return UUID{}, invalidBase58Pair(s[17], s[18])
	}
	chunk = uint64(p0)*base58PairBase + uint64(p1)
	loCarry, loProd = bits.Mul64(lo, base58ChunkBase)
	lo = loProd + chunk
	if lo < loProd {
		loCarry++
	}
	hi = hi*base58ChunkBase + loCarry

	p0 = base58DecodePairs[uint16(s[19])<<8|uint16(s[20])]
	if p0 == 0xffff {
		return UUID{}, invalidBase58Pair(s[19], s[20])
	}
	p1 = base58DecodePairs[uint16(s[21])<<8|uint16(s[22])]
	if p1 == 0xffff {
		return UUID{}, invalidBase58Pair(s[21], s[22])
	}
	chunk = uint64(p0)*base58PairBase + uint64(p1)
	loCarry, loProd = bits.Mul64(lo, base58ChunkBase)
	lo = loProd + chunk
	if lo < loProd {
		loCarry++
	}
	hi = hi*base58ChunkBase + loCarry

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

func invalidBase58Digit(c byte) error {
	if c > 127 {
		return fmt.Errorf("High-bit set on invalid digit")
	}
	return fmt.Errorf("Invalid base58 digit (%q)", rune(c))
}
