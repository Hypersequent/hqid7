package hqid7

import (
	"crypto/rand"
	"encoding/binary"
	"time"
)

type UUID [16]byte

func MustUUID7() UUID {
	uuid, err := UUID7()
	if err != nil {
		panic(err)
	}
	return uuid
}

func UUID7() (UUID, error) {
	return FromTime(time.Now())
}

func FromTime(uuidTime time.Time) (UUID, error) {
	unixNano := uuidTime.UnixNano()
	milliseconds := uint64(unixNano / int64(time.Millisecond))
	nanoseconds := uint64(unixNano % int64(time.Millisecond))

	// Calculate the 12-bit sub-millisecond precision time.
	precisionBitsValue := nanoseconds * 4096 / uint64(time.Millisecond)

	var uuidBytes UUID

	// Manually construct the first 64-bit field:
	// 48 bits for the timestamp, 4 bits for version, and 12 bits for sub-ms time.
	sixtyFourBitField := (milliseconds << 16) | (uint64(7) << 12) | precisionBitsValue
	binary.BigEndian.PutUint64(uuidBytes[0:], sixtyFourBitField)

	// Generate 62 random bits and set the UUID variant bits to 0b10.
	if _, err := rand.Read(uuidBytes[8:]); err != nil {
		return UUID{}, err
	}
	uuidBytes[8] = (uuidBytes[8] & 0x3f) | 0x80

	return uuidBytes, nil
}
