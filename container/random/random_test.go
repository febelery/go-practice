package main

import "testing"

const n = 16

func BenchmarkRandStringRunes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandStringRunes(n)
	}
}

func BenchmarkRandStringBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandStringBytes(n)
	}
}

func BenchmarkRandStringBytesRmndr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandStringBytesRmndr(n)
	}
}

func BenchmarkRandStringBytesMask(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandStringBytesMask(n)
	}
}

func BenchmarkRandStringBytesMaskImpr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandStringBytesMask(n)
	}
}

func BenchmarkRandStringBytesMaskImprSrc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandStringBytesMaskImprSrc(n)
	}
}
