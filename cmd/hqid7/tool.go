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
	var storage [512]byte
	buf := storage[:0]
	buf = append(buf, "hqid7: "...)
	buf = append(buf, idString...)
	buf = append(buf, "\n\nTimestamp (UTC):   "...)
	buf = timestamp.UTC().AppendFormat(buf, timeFormat)
	buf = append(buf, "\nTimestamp (Local): "...)
	buf = timestamp.Local().AppendFormat(buf, timeFormat)
	buf = append(buf, "\nUnix milliseconds: "...)
	buf = strconv.AppendUint(buf, timestampMs, 10)
	buf = append(buf, "\n\nVersion:           "...)
	buf = strconv.AppendUint(buf, version, 10)
	buf = append(buf, "\nVariant:           "...)
	buf = strconv.AppendUint(buf, variant, 10)
	buf = append(buf, " (binary: "...)
	buf = appendFixedBase(buf, variant, 2, "01")
	buf = append(buf, ")\nSub-ms precision:  "...)
	buf = strconv.AppendUint(buf, subMsPrecision, 10)
	buf = append(buf, " (binary: "...)
	buf = appendFixedBase(buf, subMsPrecision, 12, "01")
	buf = append(buf, ")\nRandom bits (62):  0x"...)
	buf = appendFixedBase(buf, randomBits, 15, "0123456789ABCDEF")
	buf = append(buf, '\n')
	_, _ = os.Stdout.Write(buf)
}

func appendFixedBase(buf []byte, v uint64, width int, alphabet string) []byte {
	start := len(buf)
	for i := 0; i < width; i++ {
		buf = append(buf, 0)
	}
	base := uint64(len(alphabet))
	for i := width - 1; i >= 0; i-- {
		buf[start+i] = alphabet[v%base]
		v /= base
	}
	return buf
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
