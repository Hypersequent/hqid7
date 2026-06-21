package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/hypersequent/hqid7"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]

	switch command {
	case "new", "generate":
		generateUUID()
	case "parse":
		parseUUID()
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func generateUUID() {
	// Generate a new UUID7
	uuid, err := hqid7.UUID7()
	if err != nil {
		fmt.Printf("Error generating UUID: %v\n", err)
		os.Exit(1)
	}

	// Print only canonical hqid7 base58 encoded string
	base58String := hqid7.EncodeBase58(uuid)
	_, _ = os.Stdout.WriteString(base58String + "\n")
}

func parseUUID() {
	if len(os.Args) < 3 {
		fmt.Println("Error: parse command requires an hqid7 string")
		fmt.Println("Usage: hqid7 parse <hqid7-string>")
		os.Exit(1)
	}

	idString := os.Args[2]

	// Decode the hqid7 string
	uuid, err := hqid7.DecodeBase58(idString)
	if err != nil {
		fmt.Printf("Error decoding hqid7: %v\n", err)
		os.Exit(1)
	}

	// Extract timestamp (first 48 bits / 6 bytes)
	first64Bits := binary.BigEndian.Uint64(uuid[0:8])
	timestampMs := first64Bits >> 16
	version := (first64Bits >> 12) & 0xF
	subMsPrecision := first64Bits & 0xFFF

	// Convert timestamp to time
	timestamp := time.UnixMilli(int64(timestampMs))

	// Extract random bits (last 8 bytes)
	last64Bits := binary.BigEndian.Uint64(uuid[8:16])
	variant := (last64Bits >> 62) & 0x3
	randomBits := last64Bits & 0x3FFFFFFFFFFFFFFF

	// Display information
	var storage [320]byte
	buf := storage[:0]
	buf = append(buf, "hqid7: "...)
	buf = append(buf, idString...)
	buf = append(buf, "\n\nTimestamp (UTC):   "...)
	buf = appendToolTime(buf, timestamp.UTC())
	buf = append(buf, "\nTimestamp (Local): "...)
	buf = appendToolTime(buf, timestamp.Local())
	buf = append(buf, "\nUnix milliseconds: "...)
	buf = strconv.AppendUint(buf, timestampMs, 10)
	buf = append(buf, "\n\nVersion:           "...)
	buf = strconv.AppendUint(buf, version, 10)
	buf = append(buf, "\nVariant:           "...)
	buf = strconv.AppendUint(buf, variant, 10)
	buf = append(buf, " (binary: "...)
	buf = appendBits2(buf, variant)
	buf = append(buf, ")\nSub-ms precision:  "...)
	buf = strconv.AppendUint(buf, subMsPrecision, 10)
	buf = append(buf, " (binary: "...)
	buf = appendBits12(buf, subMsPrecision)
	buf = append(buf, ")\nRandom bits (62):  0x"...)
	buf = appendHex15(buf, randomBits)
	buf = append(buf, '\n')
	_, _ = os.Stdout.Write(buf)
}

func appendToolTime(buf []byte, t time.Time) []byte {
	year, month, day := t.Date()
	if year < 0 || year > 9999 {
		return t.AppendFormat(buf, timeFormat)
	}
	hour, minute, second := t.Clock()
	zoneName, zoneOffset := t.Zone()

	buf = append4Digits(buf, year)
	buf = append(buf, '-')
	buf = append2Digits(buf, int(month))
	buf = append(buf, '-')
	buf = append2Digits(buf, day)
	buf = append(buf, ' ')
	buf = append2Digits(buf, hour)
	buf = append(buf, ':')
	buf = append2Digits(buf, minute)
	buf = append(buf, ':')
	buf = append2Digits(buf, second)
	buf = append(buf, '.')
	buf = append3Digits(buf, t.Nanosecond()/int(time.Millisecond))
	buf = append(buf, ' ')
	if zoneName == "" {
		buf = appendZoneOffset(buf, zoneOffset)
	} else {
		buf = append(buf, zoneName...)
	}
	return buf
}

func append2Digits(buf []byte, v int) []byte {
	return append(buf, byte('0'+v/10), byte('0'+v%10))
}

func append3Digits(buf []byte, v int) []byte {
	return append(buf, byte('0'+v/100), byte('0'+(v/10)%10), byte('0'+v%10))
}

func append4Digits(buf []byte, v int) []byte {
	return append(buf, byte('0'+v/1000), byte('0'+(v/100)%10), byte('0'+(v/10)%10), byte('0'+v%10))
}

func appendZoneOffset(buf []byte, offset int) []byte {
	if offset < 0 {
		buf = append(buf, '-')
		offset = -offset
	} else {
		buf = append(buf, '+')
	}
	offset /= 60
	hours := offset / 60
	minutes := offset % 60
	buf = append2Digits(buf, hours)
	buf = append2Digits(buf, minutes)
	return buf
}

func appendBits2(buf []byte, v uint64) []byte {
	return append(buf, '0'+byte((v>>1)&1), '0'+byte(v&1))
}

func appendBits12(buf []byte, v uint64) []byte {
	return append(buf,
		'0'+byte((v>>11)&1), '0'+byte((v>>10)&1), '0'+byte((v>>9)&1), '0'+byte((v>>8)&1),
		'0'+byte((v>>7)&1), '0'+byte((v>>6)&1), '0'+byte((v>>5)&1), '0'+byte((v>>4)&1),
		'0'+byte((v>>3)&1), '0'+byte((v>>2)&1), '0'+byte((v>>1)&1), '0'+byte(v&1),
	)
}

func appendHex15(buf []byte, v uint64) []byte {
	const hex = "0123456789ABCDEF"
	return append(buf,
		hex[(v>>56)&0xF], hex[(v>>52)&0xF], hex[(v>>48)&0xF], hex[(v>>44)&0xF],
		hex[(v>>40)&0xF], hex[(v>>36)&0xF], hex[(v>>32)&0xF], hex[(v>>28)&0xF],
		hex[(v>>24)&0xF], hex[(v>>20)&0xF], hex[(v>>16)&0xF], hex[(v>>12)&0xF],
		hex[(v>>8)&0xF], hex[(v>>4)&0xF], hex[v&0xF],
	)
}

const timeFormat = "2006-01-02 15:04:05.000 MST"

const usageText = `Hypersequent hqid7 Tool

Usage:
  hqid7 <command>

Commands:
  new, generate       Generate and print a new hqid7
  parse <hqid7>       Parse an hqid7 and show timestamp and random parts
  help, -h, --help    Show this help message

Examples:
  hqid7 new
  hqid7 generate
  hqid7 parse 1C3XR6Gzv_es6ViopPLabMW

Installation:
  go install github.com/hypersequent/hqid7/cmd/hqid7@latest

Development:
  go run cmd/hqid7/tool.go <command>
`

func printUsage() {
	_, _ = os.Stdout.WriteString(usageText)
}
