# hqid7

This is a Go library for generating UUIDv7 based on [RFC 9562](https://www.rfc-editor.org/rfc/rfc9562.html) — the new UUID specification published in May 2024.

> **Note:** This package was previously known as **Hypersequent's UUID7**, but was renamed to **hqid7** to avoid confusion with the official UUID7 standard released in RFC 9562.

This library also implements a custom string encoding for UUIDv7, which is lexicographically sortable. This string
encoding is not defined in RFC 9562 and is based on the Base58 encoding used in Bitcoin.

To make string representation visually more distinguishable from other UUIDs, there is a dash `_` character 
inserted after the first 9 characters.

Example: 
```txt
1C3XR6Gzv_es6ViopPLabMW
1C3XR6Gzv_gnTYagGW7m6AU
1C3VGAJyH_iXkB2HfuhEusP
1C3Rttz29_K2U2o4AdhPF5b
```

## Binary format

```plain 
    0                   1                   2                   3
    0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |                           unix_ts_ms                          |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |          unix_ts_ms           |  ver  |       rand_a          |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |var|                        rand_b                             |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |                            rand_b                             |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
```

- **unix_ts_ms** is filled with Go's `time.Now().UnixNano() / 1e6`
- **ver** is `0b0111` for UUIDv7 (RFC 9562)
- **rand_a** is filled using "Replace Leftmost Random Bits with Increased Clock Precision" (Method 3 in RFC 9562)
- **var** is `0b10` for UUIDv7 (RFC 9562)
- **rand_b** is cryptographically random bits, generated using Go's `crypto/rand` package

## String Encoding 

The UUID is encoded using Base58 encoding using BTC alphabet, which is the same as the one used in [Bitcoin](https://en.bitcoinwiki.org/wiki/Base58). 
Bitcoin address checksum is not used. 

The encoded string is always 23 characters long (padded with leading "zero" digit `1` if needed).

To make string representation visually more distinguishable from other UUIDs, there is a dash `_` character
inserted after the first 9 characters.

String representation is sortable lexicographically, which a useful property when using as keys in databases.

## Library Usage

```go
package main

import "fmt"
import "github.com/hypersequent/hqid7"

func main() {
    uuid := hqid7.NewString() // returns a 23 character long string like "1C3Rttz29_K2U2o4AdhPF5b"
    fmt.Println(uuid)
}
```

## CLI Tool

The `hqid7` command-line tool can generate and parse hqid7 identifiers.

### Installation

```bash
go install github.com/hypersequent/hqid7/cmd/hqid7@latest
```

### Usage

**Generate a new hqid7:**
```bash
hqid7 new
```

**Parse an existing hqid7:**
```bash
hqid7 parse 1C3XR6Gzv_es6ViopPLabMW
# hqid7: 1C3XR6Gzv_es6ViopPLabMW
#
# Timestamp (UTC):   2023-09-22 11:48:35.074 UTC
# Timestamp (Local): 2023-09-22 15:48:35.074 +04
# ...
```

## Dependencies

- [github.com/mr-tron/base58](https://github.com/mr-tron/base58) - Base58 encoding/decoding package (MIT License)

## License
MIT 


