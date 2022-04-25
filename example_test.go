package revregex

import "fmt"

func ExampleNewGen() {

	pattern := "a(b|c)a?"
	generator, err := NewGen(pattern)
	if err != nil {
		fmt.Println(err)
	}
	entropy := NewRandChooserSeed(42) // or use NewRandInter() for real randomness ...

	// Generate 5 strings that match "a(b|c)a?"
	for i := 0; i < 5; i++ {
		result := generator.Next(entropy)
		fmt.Println(result)

		// Verify each generated string actually matches ?
		if err = generator.Verify(result); err != nil {
			fmt.Println("Verification failed with error : ", err)
		}
	}

	// Output:
	// aca
	// ab
	// aca
	// ac
	// aba

}
