package rgen

import (
	_ "embed"
)

// v0.4.6   refactored packages boundaries, obsolete synchroneous API moved to internal package -
// v0.4.5   clean up, minor release -
// v0.4.4   added iterator API, bump to go 1.23, using task for convenience -
// v0.4.3	typos -
// v0.4.2	dedup added (bloom filter) -
// v0.4.1 	complete redesign from 0.3.x -
const VERSION = "0.4.6"

const COPYRIGHT = "(c) xavier gandillot 2022-2024"

var BUILDDATE string

//go:embed LICENSE
var license string

// Return licence text
func Licence() string {
	return license
}
