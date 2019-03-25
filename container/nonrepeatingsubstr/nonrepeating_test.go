package main

import "testing"

func BenchmarkSubstr(b *testing.B) {
	s := "1233123123123"
	ans := 3

	for i := 0; i < b.N; i++ {
		actual := lengthOfNonRepeatingSubStr(s)
		if actual != ans {
			b.Errorf("got %d for input %s; expected %d", actual, s, ans)
		}
	}
}
