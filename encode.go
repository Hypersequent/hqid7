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
	// A 128-bit UUID is at most 22 Base58 digits. Divide four uint32 limbs by 58
	// directly, then left-pad with BTC Base58 zeroes ('1') to hqid7's fixed width.
	l0 := binary.BigEndian.Uint32(u[0:4])
	l1 := binary.BigEndian.Uint32(u[4:8])
	l2 := binary.BigEndian.Uint32(u[8:12])
	l3 := binary.BigEndian.Uint32(u[12:16])

	var raw [22]byte
	for i := range raw {
		raw[i] = '1'
	}
	for i := len(raw) - 1; (l0 | l1 | l2 | l3) != 0; i -= 5 {
		cur := uint64(l0)
		l0 = uint32(cur / 656356768)
		rem := cur - uint64(l0)*656356768
		cur = (rem << 32) | uint64(l1)
		l1 = uint32(cur / 656356768)
		rem = cur - uint64(l1)*656356768
		cur = (rem << 32) | uint64(l2)
		l2 = uint32(cur / 656356768)
		rem = cur - uint64(l2)*656356768
		cur = (rem << 32) | uint64(l3)
		l3 = uint32(cur / 656356768)
		rem = cur - uint64(l3)*656356768
		loPair := base58Pairs[rem%3364]
		raw[i-1] = byte(loPair >> 8)
		raw[i] = byte(loPair)
		if i >= 4 {
			rem /= 3364
			midPair := base58Pairs[rem%3364]
			raw[i-3] = byte(midPair >> 8)
			raw[i-2] = byte(midPair)
			raw[i-4] = base58Alphabet[rem/3364]
		}
	}

	var out [23]byte
	copy(out[:9], raw[:9])
	out[9] = '_'
	copy(out[10:], raw[9:])
	return string(out[:])
}

var base58Decode = [256]byte{
	'1': 1, '2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7, '8': 8, '9': 9,
	'A': 10, 'B': 11, 'C': 12, 'D': 13, 'E': 14, 'F': 15, 'G': 16, 'H': 17,
	'J': 18, 'K': 19, 'L': 20, 'M': 21, 'N': 22, 'P': 23, 'Q': 24, 'R': 25,
	'S': 26, 'T': 27, 'U': 28, 'V': 29, 'W': 30, 'X': 31, 'Y': 32, 'Z': 33,
	'a': 34, 'b': 35, 'c': 36, 'd': 37, 'e': 38, 'f': 39, 'g': 40, 'h': 41,
	'i': 42, 'j': 43, 'k': 44, 'm': 45, 'n': 46, 'o': 47, 'p': 48, 'q': 49,
	'r': 50, 's': 51, 't': 52, 'u': 53, 'v': 54, 'w': 55, 'x': 56, 'y': 57, 'z': 58,
}

func DecodeBase58(s string) (UUID, error) {
	if len(s) != 23 {
		return UUID{}, errors.New("hqid7 base58: invalid length")
	}
	if s[9] != '_' {
		return UUID{}, errors.New("hqid7 base58: invalid separator")
	}

	// Decode the fixed 22 Base58 digits into two uint64 limbs. hqid7 keeps the
	// low 128 bits after removing leading Base58 zeroes, so arithmetic wraps
	// above the UUID width to match the original decoder.
	var hi, lo uint64
	for i := 0; i < 9; i++ {
		c := s[i]
		v := base58Decode[c]
		if v == 0 {
			if c > 127 {
				return UUID{}, fmt.Errorf("High-bit set on invalid digit")
			}
			return UUID{}, fmt.Errorf("Invalid base58 digit (%q)", rune(c))
		}

		loCarry, loProd := bits.Mul64(lo, 58)
		lo = loProd + uint64(v-1)
		if lo < loProd {
			loCarry++
		}
		hi = hi*58 + loCarry
	}
	for i := 10; i < 23; i++ {
		c := s[i]
		v := base58Decode[c]
		if v == 0 {
			if c > 127 {
				return UUID{}, fmt.Errorf("High-bit set on invalid digit")
			}
			return UUID{}, fmt.Errorf("Invalid base58 digit (%q)", rune(c))
		}

		loCarry, loProd := bits.Mul64(lo, 58)
		lo = loProd + uint64(v-1)
		if lo < loProd {
			loCarry++
		}
		hi = hi*58 + loCarry
	}

	var uuid UUID
	binary.BigEndian.PutUint64(uuid[0:8], hi)
	binary.BigEndian.PutUint64(uuid[8:16], lo)
	return uuid, nil
}
