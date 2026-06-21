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
	v0 := base58Decode[s[0]]
	if v0 == 0xff {
		c := s[0]
		if c > 127 {
			return UUID{}, fmt.Errorf("High-bit set on invalid digit")
		}
		return UUID{}, fmt.Errorf("Invalid base58 digit (%q)", rune(c))
	}
	v1 := base58Decode[s[1]]
	if v1 == 0xff {
		c := s[1]
		if c > 127 {
			return UUID{}, fmt.Errorf("High-bit set on invalid digit")
		}
		return UUID{}, fmt.Errorf("Invalid base58 digit (%q)", rune(c))
	}
	var hi uint64
	lo := uint64(v0)*58 + uint64(v1)

	v0 = base58Decode[s[2]]
	if v0 == 0xff {
		c := s[2]
		if c > 127 {
			return UUID{}, fmt.Errorf("High-bit set on invalid digit")
		}
		return UUID{}, fmt.Errorf("Invalid base58 digit (%q)", rune(c))
	}
	v1 = base58Decode[s[3]]
	if v1 == 0xff {
		c := s[3]
		if c > 127 {
			return UUID{}, fmt.Errorf("High-bit set on invalid digit")
		}
		return UUID{}, fmt.Errorf("Invalid base58 digit (%q)", rune(c))
	}
	v2 := base58Decode[s[4]]
	if v2 == 0xff {
		c := s[4]
		if c > 127 {
			return UUID{}, fmt.Errorf("High-bit set on invalid digit")
		}
		return UUID{}, fmt.Errorf("Invalid base58 digit (%q)", rune(c))
	}
	v3 := base58Decode[s[5]]
	if v3 == 0xff {
		c := s[5]
		if c > 127 {
			return UUID{}, fmt.Errorf("High-bit set on invalid digit")
		}
		return UUID{}, fmt.Errorf("Invalid base58 digit (%q)", rune(c))
	}
	chunk := (((uint64(v0)*58)+uint64(v1))*58+uint64(v2))*58 + uint64(v3)
	loCarry, loProd := bits.Mul64(lo, 11316496)
	lo = loProd + chunk
	if lo < loProd {
		loCarry++
	}
	hi = hi*11316496 + loCarry

	v0 = base58Decode[s[6]]
	if v0 == 0xff {
		c := s[6]
		if c > 127 {
			return UUID{}, fmt.Errorf("High-bit set on invalid digit")
		}
		return UUID{}, fmt.Errorf("Invalid base58 digit (%q)", rune(c))
	}
	v1 = base58Decode[s[7]]
	if v1 == 0xff {
		c := s[7]
		if c > 127 {
			return UUID{}, fmt.Errorf("High-bit set on invalid digit")
		}
		return UUID{}, fmt.Errorf("Invalid base58 digit (%q)", rune(c))
	}
	v2 = base58Decode[s[8]]
	if v2 == 0xff {
		c := s[8]
		if c > 127 {
			return UUID{}, fmt.Errorf("High-bit set on invalid digit")
		}
		return UUID{}, fmt.Errorf("Invalid base58 digit (%q)", rune(c))
	}
	v3 = base58Decode[s[10]]
	if v3 == 0xff {
		c := s[10]
		if c > 127 {
			return UUID{}, fmt.Errorf("High-bit set on invalid digit")
		}
		return UUID{}, fmt.Errorf("Invalid base58 digit (%q)", rune(c))
	}
	chunk = (((uint64(v0)*58)+uint64(v1))*58+uint64(v2))*58 + uint64(v3)
	loCarry, loProd = bits.Mul64(lo, 11316496)
	lo = loProd + chunk
	if lo < loProd {
		loCarry++
	}
	hi = hi*11316496 + loCarry

	v0 = base58Decode[s[11]]
	if v0 == 0xff {
		c := s[11]
		if c > 127 {
			return UUID{}, fmt.Errorf("High-bit set on invalid digit")
		}
		return UUID{}, fmt.Errorf("Invalid base58 digit (%q)", rune(c))
	}
	v1 = base58Decode[s[12]]
	if v1 == 0xff {
		c := s[12]
		if c > 127 {
			return UUID{}, fmt.Errorf("High-bit set on invalid digit")
		}
		return UUID{}, fmt.Errorf("Invalid base58 digit (%q)", rune(c))
	}
	v2 = base58Decode[s[13]]
	if v2 == 0xff {
		c := s[13]
		if c > 127 {
			return UUID{}, fmt.Errorf("High-bit set on invalid digit")
		}
		return UUID{}, fmt.Errorf("Invalid base58 digit (%q)", rune(c))
	}
	v3 = base58Decode[s[14]]
	if v3 == 0xff {
		c := s[14]
		if c > 127 {
			return UUID{}, fmt.Errorf("High-bit set on invalid digit")
		}
		return UUID{}, fmt.Errorf("Invalid base58 digit (%q)", rune(c))
	}
	chunk = (((uint64(v0)*58)+uint64(v1))*58+uint64(v2))*58 + uint64(v3)
	loCarry, loProd = bits.Mul64(lo, 11316496)
	lo = loProd + chunk
	if lo < loProd {
		loCarry++
	}
	hi = hi*11316496 + loCarry

	v0 = base58Decode[s[15]]
	if v0 == 0xff {
		c := s[15]
		if c > 127 {
			return UUID{}, fmt.Errorf("High-bit set on invalid digit")
		}
		return UUID{}, fmt.Errorf("Invalid base58 digit (%q)", rune(c))
	}
	v1 = base58Decode[s[16]]
	if v1 == 0xff {
		c := s[16]
		if c > 127 {
			return UUID{}, fmt.Errorf("High-bit set on invalid digit")
		}
		return UUID{}, fmt.Errorf("Invalid base58 digit (%q)", rune(c))
	}
	v2 = base58Decode[s[17]]
	if v2 == 0xff {
		c := s[17]
		if c > 127 {
			return UUID{}, fmt.Errorf("High-bit set on invalid digit")
		}
		return UUID{}, fmt.Errorf("Invalid base58 digit (%q)", rune(c))
	}
	v3 = base58Decode[s[18]]
	if v3 == 0xff {
		c := s[18]
		if c > 127 {
			return UUID{}, fmt.Errorf("High-bit set on invalid digit")
		}
		return UUID{}, fmt.Errorf("Invalid base58 digit (%q)", rune(c))
	}
	chunk = (((uint64(v0)*58)+uint64(v1))*58+uint64(v2))*58 + uint64(v3)
	loCarry, loProd = bits.Mul64(lo, 11316496)
	lo = loProd + chunk
	if lo < loProd {
		loCarry++
	}
	hi = hi*11316496 + loCarry

	v0 = base58Decode[s[19]]
	if v0 == 0xff {
		c := s[19]
		if c > 127 {
			return UUID{}, fmt.Errorf("High-bit set on invalid digit")
		}
		return UUID{}, fmt.Errorf("Invalid base58 digit (%q)", rune(c))
	}
	v1 = base58Decode[s[20]]
	if v1 == 0xff {
		c := s[20]
		if c > 127 {
			return UUID{}, fmt.Errorf("High-bit set on invalid digit")
		}
		return UUID{}, fmt.Errorf("Invalid base58 digit (%q)", rune(c))
	}
	v2 = base58Decode[s[21]]
	if v2 == 0xff {
		c := s[21]
		if c > 127 {
			return UUID{}, fmt.Errorf("High-bit set on invalid digit")
		}
		return UUID{}, fmt.Errorf("Invalid base58 digit (%q)", rune(c))
	}
	v3 = base58Decode[s[22]]
	if v3 == 0xff {
		c := s[22]
		if c > 127 {
			return UUID{}, fmt.Errorf("High-bit set on invalid digit")
		}
		return UUID{}, fmt.Errorf("Invalid base58 digit (%q)", rune(c))
	}
	chunk = (((uint64(v0)*58)+uint64(v1))*58+uint64(v2))*58 + uint64(v3)
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
