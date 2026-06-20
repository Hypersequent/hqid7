package hqid7

import (
	"errors"

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

func DecodeBase58(s string) (UUID, error) {
	if len(s) != 23 {
		return UUID{}, errors.New("hqid7 base58: invalid length")
	}
	if s[9] != '_' {
		return UUID{}, errors.New("hqid7 base58: invalid separator")
	}
	s = s[0:9] + s[10:]
	d, err := base58.Decode(s)
	if err != nil {
		return UUID{}, err
	}
	if len(d) > 16 {
		d = d[len(d)-16:] // remove leading "zeroes" (1 in BTC base58)
	}
	return UUID(d), nil
}
