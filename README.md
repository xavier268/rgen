This package generate all the strings with a given length matching a provided regexp pattern.

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

## Unsupported opérations :

The following operations are not supported, because they make little sens in this context.

* dot(.) operator
* flags
* capture
* start/end of word, text, line ...