package main

import (
	"os"
	"testing"
)

var benchmarkCLIArgsSink []string

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
