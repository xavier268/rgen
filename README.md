# revregex


[![Go Reference](https://pkg.go.dev/badge/github.com/xavier268/revregex.svg)](https://pkg.go.dev/github.com/xavier268/revregex)


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
	generator := NewGen(pattern)
    // create a source of entropy
	chooser := NewRandChooser() 

	
	// Now, you can ask the generator to create a string using the chooser to make its decisions.
	result := generator.Next(chooser) // for instance, result will get "aca" or "ab".
		
    // if you don't trust the package and want to verify that the string actually matches ...
    err := generator.Verify(chooser)
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
They generate shorter strings first, giving an exponentially decreasing priority to the longuer strings.

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

* **NewBytesChooser**, which takes a []byte as input, will use the *information* contained in this array to make its choices. There is a one-to-one relation between a given byte array and the sequence of strings generated. However, the information from the provided byte array is consumed as we generate strings and make choices, and, at some point, once there is no more information (underlying array is nil or contains only zeros), the defaults choices will always be made, generating from tat point always the same defaul answer. 
