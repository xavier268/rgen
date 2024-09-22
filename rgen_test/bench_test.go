package rgen_test

import (
	"context"
	"testing"

	"github.com/xavier268/rgen"
)

const pattern = "a*b*(d|e)*" // same pattern for all benches !

// goos: windows
// goarch: amd64
// pkg: github.com/xavier268/rgen/rgen_test
// cpu: Intel(R) Core(TM) i7-10700 CPU @ 2.90GHz
// BenchmarkAllSynch2-2              166117              7148 ns/op            3736 B/op        153 allocs/op
// BenchmarkAllSynch2-8              165837              7461 ns/op            3736 B/op        153 allocs/op
// BenchmarkAllAsync2-2               99858             11713 ns/op            8128 B/op        218 allocs/op
// BenchmarkAllAsync2-8              101503             11615 ns/op            8128 B/op        218 allocs/op
// BenchmarkAllSynch4-2               26162             45608 ns/op           11728 B/op        922 allocs/op
// BenchmarkAllSynch4-8               23916             48037 ns/op           11728 B/op        922 allocs/op
// BenchmarkAllAsync4-2               26544             45138 ns/op           22384 B/op       1092 allocs/op
// BenchmarkAllAsync4-8               24645             48460 ns/op           22384 B/op       1092 allocs/op
// BenchmarkAllSynch8-2                1050           1175203 ns/op          307872 B/op      21932 allocs/op
// BenchmarkAllSynch8-8                1048           1198312 ns/op          307873 B/op      21932 allocs/op
// BenchmarkAllAsync8-2                1418            837494 ns/op          336721 B/op      22420 allocs/op
// BenchmarkAllAsync8-8                1186           1013251 ns/op          336728 B/op      22420 allocs/op
// BenchmarkAllSynch16-2                  2         521580050 ns/op        232448640 B/op   8750766 allocs/op
// BenchmarkAllSynch16-8                  2         530774900 ns/op        232449504 B/op   8750775 allocs/op
// BenchmarkAllAsync16-2                  3         471671167 ns/op        232536768 B/op   8752322 allocs/op
// BenchmarkAllAsync16-8                  3         401646933 ns/op        232538864 B/op   8752328 allocs/op
// BenchmarkAllSynch20-2                  1       12394618300 ns/op       5479923624 B/op 162826963 allocs/op
// BenchmarkAllSynch20-8                  1       12071839700 ns/op       5479930088 B/op 162827027 allocs/op
// BenchmarkAllAsync20-2                  1        8966445300 ns/op        5480086688 B/op 162829409 allocs/op
// BenchmarkAllAsync20-8                  1        8494588700 ns/op        5480084560 B/op 162829385 allocs/op

// Keytakeaways :
// As length increases, Async gets substantially faster.
// The number of available cpus does not have significant impact.

func BenchmarkAllSynch2(b *testing.B) {
	doBench(2, b)
}

func BenchmarkAllAsync2(b *testing.B) {
	doBenchAsync(2, b)
}
func BenchmarkAllSynch4(b *testing.B) {
	doBench(4, b)
}

func BenchmarkAllAsync4(b *testing.B) {
	doBenchAsync(4, b)
}

func BenchmarkAllSynch8(b *testing.B) {
	doBench(8, b)
}
func BenchmarkAllAsync8(b *testing.B) {
	doBenchAsync(8, b)
}

func BenchmarkAllSynch16(b *testing.B) {
	doBench(16, b)
}
func BenchmarkAllAsync16(b *testing.B) {
	doBenchAsync(16, b)
}

func BenchmarkAllSynch20(b *testing.B) {
	doBench(20, b)
}
func BenchmarkAllAsync20(b *testing.B) {
	doBenchAsync(20, b)
}

// ============================================

// synchroneous
func doBench(len int, b *testing.B) {

	ss := "" // prevent compiler eleiminating loop content !
	for i := 0; i < b.N; i++ {
		for s := range rgen.All(pattern, len) {
			ss = s
		}
	}
	_ = ss // keep compiler happy !
}

// asynchroneous
func doBenchAsync(len int, b *testing.B) {
	// blackhole will receive anything
	blackhole := make(chan string, 100_000)
	defer close(blackhole)

	// ensure blackhole is always emptied
	go func() {
		for range blackhole {
			// discard all content !
		}
	}()

	// actual bench starts here !
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := rgen.Generate(context.Background(), blackhole, pattern, len); err != nil {
			panic(err)
		}
	}

}
