package common

import (
	"testing"
	"time"
)

func TestPow(t *testing.T) {
	for i := 6; i < 12; i++ {
		now := time.Now()
		hash := Pow(i)
		t.Logf("target: %d, took %d ms, hash: %v\n", i, time.Since(now), hash)
	}
}
