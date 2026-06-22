package main

import (
	"io"
	"os"
	"testing"
	"time"
)

var benchmarkCLIArgsSink []string

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	read, write, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	oldStdout := os.Stdout
	os.Stdout = write
	defer func() {
		os.Stdout = oldStdout
	}()

	fn()

	if err := write.Close(); err != nil {
		t.Fatal(err)
	}
	out, err := io.ReadAll(read)
	if err != nil {
		t.Fatal(err)
	}
	if err := read.Close(); err != nil {
		t.Fatal(err)
	}

	return string(out)
}

func TestParseUUIDOutput(t *testing.T) {
	oldArgs := os.Args
	oldLocal := time.Local
	defer func() {
		os.Args = oldArgs
		time.Local = oldLocal
	}()

	os.Args = []string{"hqid7", "parse", "1C3XR6Gzv_es6ViopPLabMW"}
	time.Local = time.UTC

	got := captureStdout(t, parseUUID)
	want := `hqid7: 1C3XR6Gzv_es6ViopPLabMW

Timestamp (UTC):   2023-09-22 11:48:35.074 UTC
Timestamp (Local): 2023-09-22 11:48:35.074 UTC
Unix milliseconds: 1695383315074

Version:           7
Variant:           2 (binary: 10)
Sub-ms precision:  2875 (binary: 101100111011)
Random bits (62):  0x2E787F1B96267445
`
	if got != want {
		t.Fatalf("parse output mismatch\nwant:\n%s\ngot:\n%s", want, got)
	}
}

func benchmarkCLI(b *testing.B, args []string, fn func()) {
	b.Helper()

	devNull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil {
		b.Fatal(err)
	}
	defer devNull.Close()

	oldArgs := os.Args
	oldStdout := os.Stdout
	defer func() {
		os.Args = oldArgs
		os.Stdout = oldStdout
	}()

	os.Args = args
	os.Stdout = devNull
	benchmarkCLIArgsSink = args

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		fn()
	}
}

func BenchmarkMainGenerate(b *testing.B) {
	benchmarkCLI(b, []string{"hqid7", "new"}, main)
}

func BenchmarkGenerateUUID(b *testing.B) {
	benchmarkCLI(b, []string{"hqid7", "new"}, generateUUID)
}

func BenchmarkMainParse(b *testing.B) {
	benchmarkCLI(b, []string{"hqid7", "parse", "1C3XR6Gzv_es6ViopPLabMW"}, main)
}

func BenchmarkParseUUID(b *testing.B) {
	benchmarkCLI(b, []string{"hqid7", "parse", "1C3XR6Gzv_es6ViopPLabMW"}, parseUUID)
}

func BenchmarkMainHelp(b *testing.B) {
	benchmarkCLI(b, []string{"hqid7", "help"}, main)
}
