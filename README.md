# revregex


[![Go Reference](https://pkg.go.dev/badge/github.com/xavier268/revregex.svg)](https://pkg.go.dev/github.com/xavier268/revregex) [![Go Report Card](https://goreportcard.com/badge/github.com/xavier268/revregex)](https://goreportcard.com/report/github.com/xavier268/revregex)

A package you can use to generate all possible strings matching a given regular expression.


# How to use

```golang

import (
    "fmt"
    "github.com/xavier268/revregex"
    )

    // start from a POSIX regular expression pattern
    pattern := "a(b|c)a?"
    // create a generator for this pattern
	generator, err := NewGen(pattern)
    if err != nil { 
        ... 
        }
    // create a source of entropy
	chooser := NewRandChooser() 

	
	// Now, you can ask the generator to create a string using the chooser to make its decisions.
	result := generator.Next(chooser) // for instance, result will get "aca" or "ab".
		
    // if you don't trust the package and want to verify that the string actually matches ...
    err = generator.Verify(chooser)
    if err != nil {
        ...
        }
    

```


See full example file [here](./example_test.go)

# Regex syntax

Any **POSIX** format regexes valid in Golang is accepted. 

No flags are available. 
Grouping or named captures have no effect.
Word, line or text boundaries are meaningless here and should be avoided.

Unbounded metacharacteres such as * or + are accepted. 
They generate shorter strings first, giving an exponentially decreasing priority to the longer strings.

Beware of *negative* classes, such as [^a] or the dot "." operator, because they will likely generate a lot of strange unicode caracters ( up to 4 bytes characters ! ). Prefer to use *positive* classes such as [[:digit:]] or [a-zA-Z] to limit unprintable characters.

See reference [here](https://golang.org/s/re2syntax)

# Deterministic or random ?

To choose which string to generate among the many - possibly unlimited - strings matching the provided pattern, one needs to make choices.
These choices are driven by something that fullfils the Chooser interface.

```golang

type Chooser interface {
	// Intn provides a number between 0 and (n-1) included.
	Intn(n int) int
}

```

Two ways to construct a Chooser are provided :

* **NewRandChooser** to make decision randomly. If you computer has a good random generator, you most likely will endup generating all the shorter possible strings.

* **NewBytesChooser**, which takes a []byte as input, will use the *information* contained in this array to make its choices. There is a one-to-one relation between a given byte array and the sequence of strings generated. However, the information from the provided byte array is consumed as we generate strings and make choices, and, at some point, once there is no more information (underlying array is nil or contains only zeros), the defaults choices will always be made, generating from that point always the same default answer. 

# Use for Fuzzing

This can be useful for fuzzing if you want to fuzz a function whose argument matches a specific regex pattern. You could, of course, check if the pattern is matched when starting the test, and skip if not matching. However, even with medium complexity pattern, the likelyhood to find relevant pattern is very low, and the likelyhood to explore the pattern space is even lower.

With this library, instead of fuzzing based on the function argument itself, you can fuzz an int64 that will be used as the random seed to generate a few strings matching the required pattern. The benefit is that you would explore the valid entry set more thoroughly, and faster. You could still, however, fuzz the direct way as well to check how your function behaves hen the input pattern is invalid.

# Deduplicate generated strings

The Deduper interface can be used to detect duplicated string. Querying for uniqueness will both register the string and return a uniqueness response. Querying twice will always return false ...

It is currently available in two flavors :

* **NewDedupMap** is based on a map[string]bool. Is is always exact, but the memory footprint is not bounded.
* **NewDedupBloom** is a bloom-filter implementation. It has a bounded memory footprint and query time, but when a certain volume is reached, false positive (wrongly flaging original strings as duplicates) will start to occur. However, if a string is flaged as unique, it is really unique. 
The default size for the bloom filter is a reasonable value up to a few hundred thousands different strings.