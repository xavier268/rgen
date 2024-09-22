package rgen_test

import (
	"testing"

	"github.com/xavier268/rgen"
)

func BenchmarkAll2(b *testing.B) {
	doBench(2, b)
}

func BenchmarkAll4(b *testing.B) {
	doBench(4, b)
}
func BenchmarkAll8(b *testing.B) {
	doBench(8, b)
}
func BenchmarkAll16(b *testing.B) {
	doBench(16, b)
}

// ============================================

func doBench(len int, b *testing.B) {
	pattern := "a*b*(d|e)*"
	ss := "" // prevent compiler eleiminating loop content !
	for i := 0; i < b.N; i++ {
		for s := range rgen.All(pattern, len) {
			ss = s
		}
	}
	_ = ss // keep compiler happy !
}
