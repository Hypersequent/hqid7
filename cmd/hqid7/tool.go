package main

import (
	"encoding/binary"
	"fmt"
	"os"
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
	fmt.Printf("%s\n", base58String)
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
	fmt.Printf("hqid7: %s\n", idString)
	fmt.Println()
	fmt.Printf("Timestamp (UTC):   %s\n", timestamp.UTC().Format("2006-01-02 15:04:05.000 MST"))
	fmt.Printf("Timestamp (Local): %s\n", timestamp.Local().Format("2006-01-02 15:04:05.000 MST"))
	fmt.Printf("Unix milliseconds: %d\n", timestampMs)
	fmt.Println()
	fmt.Printf("Version:           %d\n", version)
	fmt.Printf("Variant:           %d (binary: %02b)\n", variant, variant)
	fmt.Printf("Sub-ms precision:  %d (binary: %012b)\n", subMsPrecision, subMsPrecision)
	fmt.Printf("Random bits (62):  0x%015X\n", randomBits)
}

func printUsage() {
	fmt.Println("Hypersequent hqid7 Tool")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  go run cmd/hqid7/tool.go <command>")
	fmt.Println("  # or build and run:")
	fmt.Println("  go build -o hqid7 cmd/hqid7/tool.go")
	fmt.Println("  ./hqid7 <command>")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  new, generate       Generate and print a new hqid7")
	fmt.Println("  parse <hqid7>       Parse an hqid7 and show timestamp and random parts")
	fmt.Println("  help, -h, --help    Show this help message")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  go run cmd/hqid7/tool.go new")
	fmt.Println("  go run cmd/hqid7/tool.go generate")
	fmt.Println("  go run cmd/hqid7/tool.go parse 1C3XR6Gzv_es6ViopPLabMW")
}
