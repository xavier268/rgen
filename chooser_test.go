package revregex

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
)

// compiler checks
var _ Chooser = new(rand.Rand)
var _ Chooser = new(bigChooser)

func TestInterRandIntn(t *testing.T) {

	const loop = 10
	it := NewRandChooserSeed(42)
	m := make(map[int]bool, loop)
	for i := 0; i < loop; i++ {
		n := it.Intn(loop * loop)
		fmt.Println(n)
		if m[n] {
			t.Fatal("duplicated value")
		}
	}

}

func TestInterRandExpVisual(t *testing.T) {
	const loop = 100_000
	it := NewRandChooserSeed(42)
	m := make(map[int]float64, loop)

	for i := 0; i < loop; i++ {
		n := exp(it)
		m[n] += 1
	}
	fmt.Println("Value :  Frequence")
	for i := 0; i < len(m); i++ {
		freq := 100 * m[i] / float64(loop)
		exp := 100. / math.Exp2(float64(i+1))
		fmt.Printf("%3d   :  %2.3f%% (expected : %2.3f%%)\n", i, freq, exp)
		if i < 8 && math.Abs((freq-exp)/exp) > 0.1 {
			t.Fatal("more than 10%% error frequency vs expected)")
		}
	}
}

func TestInterBytes0(t *testing.T) {

	if NewBytesChooser([]byte{}).Intn(7) != 0 {
		t.Fatal()
	}

	if NewBytesChooser([]byte{15, 35}).Intn(5) != 0 {
		t.Fatal()
	}

}

func TestInterBytes1(t *testing.T) {

	b := []byte{224} // start with 224
	it := NewBytesChooser(b)
	fmt.Println(it.(*bigChooser).big)
	if it.Intn(3) != 2 { // 224 -> 74
		fmt.Println(it.(*bigChooser).big)
		t.Fatal(it.(*bigChooser).big)
	}
	fmt.Println(it.(*bigChooser).big)
	if it.Intn(2) != 0 { // 74 -> 37
		fmt.Println(it.(*bigChooser).big)
		t.Fatal(it.(*bigChooser).big)
	}
	fmt.Println(it.(*bigChooser).big)
	if it.Intn(1) != 0 { // 37 -> 37
		fmt.Println(it.(*bigChooser).big)
		t.Fatal(it.(*bigChooser).big)
	}
	fmt.Println(it.(*bigChooser).big)
	if it.Intn(70) != 37 { // 37 -> 0
		fmt.Println(it.(*bigChooser).big)
		t.Fatal(it.(*bigChooser).big)
	}
	fmt.Println(it.(*bigChooser).big)
	if it.Intn(11) != 0 { // 0 -> 0
		fmt.Println(it.(*bigChooser).big)
		t.Fatal(it.(*bigChooser).big)
	}
	fmt.Println(it.(*bigChooser).big)
	for i := 0; i < 10; i++ {
		if it.Intn(i) != 0 {
			fmt.Println(it.(*bigChooser).big)
			t.Fatal(it.(*bigChooser).big)
		}
	}

}
