
# rgen

[![Go Reference](https://pkg.go.dev/badge/github.com/xavier268/rgen.svg)](https://pkg.go.dev/github.com/xavier268/rgen) [![Go Report Card](https://goreportcard.com/badge/github.com/xavier268/rgen)](https://goreportcard.com/report/github.com/xavier268/rgen)

This package generate all the strings with a given length matching a provided regexp pattern.

Three different API can used :
  
1. **All(regex string, maxlen int) iter.Seq[string]** provides a synchroneous iterator, compliant with the new go 1.23 syntax.
   * This is the prefered method for synchroneous operation
   * AllExact is a variant that iterate on strings with exactly the provided length
   * The resulting iterator can be deduplicated, to generate only unique strings if the pattern is ambiguous (see Dedup and Deduper).

2. **Generate()** provides an asynchroneous generation model, to a channel.


## How to use 

See examples & test files.

## Supported operations :

* Parenthesis for grouping, without capture
* Zero-or-more(*)
* One-or-more(+)
* Zero-or-one (?)
* Repeat {n,m}
* Alternatives (|)
* Concatenation
* Character classes [a-z] or [abc] or [0-3 8-9]

The provided context allows for timeout and cancelation management.
Operations are threadsafe accross generators.

## Unsupported operations :

The following operations are not supported, because they make little sense in this context.

* dot(.) operator
* flags
* capture
* boundaries, start/end of word, text, line ...
