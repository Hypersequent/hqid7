package hqid7

import (
	"testing"
	"time"
)

var (
	benchmarkUUID = UUID{
		0x01, 0x8a, 0xa5, 0x87,
		0x2d, 0x82, 0x70, 0x00,
		0x9a, 0x2b, 0x3c, 0x4d,
		0x5e, 0x6f, 0x70, 0x81,
	}
	benchmarkTime          = time.Unix(1_700_000_000, 123_456_789)
	benchmarkEncodedString = EncodeBase58(benchmarkUUID)

	benchmarkUUIDSink   UUID
	benchmarkStringSink string
	benchmarkErrorSink  error
)

func BenchmarkNewString(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		benchmarkStringSink = NewString()
	}
}

func BenchmarkUUID7(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		benchmarkUUIDSink, benchmarkErrorSink = UUID7()
	}
}

func BenchmarkFromTime(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		benchmarkUUIDSink, benchmarkErrorSink = FromTime(benchmarkTime)
	}
}

func BenchmarkEncodeBase58(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		benchmarkStringSink = EncodeBase58(benchmarkUUID)
	}
}

func BenchmarkDecodeBase58(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		benchmarkUUIDSink, benchmarkErrorSink = DecodeBase58(benchmarkEncodedString)
	}
}

func BenchmarkEncodeDecodeBase58(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		encoded := EncodeBase58(benchmarkUUID)
		benchmarkUUIDSink, benchmarkErrorSink = DecodeBase58(encoded)
	}
}
