package hqid7

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/mr-tron/base58"
)

func NewString() string {
	return EncodeBase58(MustUUID7())
}

func encodeBase58Raw(u UUID) string {
	return base58.Encode(u[:])
}

const base58Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

func EncodeBase58(u UUID) string {
	// A 128-bit UUID is at most 22 Base58 digits. Convert directly into a fixed
	// stack buffer and left-pad with BTC Base58 zeroes ('1'), avoiding the
	// intermediate string allocation from base58.Encode plus concatenation.
	var digits [22]byte // little-endian Base58 digits
	length := 0
	for _, b := range u {
		carry := uint32(b)
		for i := 0; i < length; i++ {
			carry += uint32(digits[i]) << 8
			digits[i] = byte(carry % 58)
			carry /= 58
		}
		for carry > 0 {
			digits[length] = byte(carry % 58)
			length++
			carry /= 58
		}
	}

	var raw [22]byte
	for i := range raw {
		raw[i] = '1'
	}
	for i := 0; i < length; i++ {
		raw[len(raw)-1-i] = base58Alphabet[digits[i]]
	}

	var out [23]byte
	copy(out[:9], raw[:9])
	out[9] = '_'
	copy(out[10:], raw[9:])
	return string(out[:])
}

var base58Decode = [128]byte{
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

	// Decode the fixed 22 Base58 digits into six uint32 limbs, mirroring the
	// mr-tron/base58 fast decoder. hqid7 keeps the low 128 bits after removing
	// leading Base58 zeroes, so the UUID is the final four limbs.
	var out [6]uint32
	for i := 0; i < len(s); i++ {
		if i == 9 {
			continue
		}
		c := s[i]
		if c > 127 {
			return UUID{}, fmt.Errorf("High-bit set on invalid digit")
		}
		v := base58Decode[c]
		if v == 0 {
			return UUID{}, fmt.Errorf("Invalid base58 digit (%q)", rune(c))
		}

		carry := uint32(v - 1)
		for j := len(out) - 1; j >= 0; j-- {
			t := uint64(out[j])*58 + uint64(carry)
			carry = uint32(t>>32) & 0x3f
			out[j] = uint32(t)
		}
		if carry > 0 {
			return UUID{}, fmt.Errorf("Output number too big (carry to the next int32)")
		}
		if out[0]&0xffff0000 != 0 {
			return UUID{}, fmt.Errorf("Output number too big (last int32 filled too far)")
		}
	}

	var uuid UUID
	binary.BigEndian.PutUint32(uuid[0:4], out[2])
	binary.BigEndian.PutUint32(uuid[4:8], out[3])
	binary.BigEndian.PutUint32(uuid[8:12], out[4])
	binary.BigEndian.PutUint32(uuid[12:16], out[5])
	return uuid, nil
}
