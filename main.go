// // Copyright © 2020 Developer Network, LLC
// //
// // This file is subject to the terms and conditions defined in
// // file 'LICENSE', which is part of this source code package.

package main

import (
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/jessevdk/go-flags"
)

var options Options
var parser = flags.NewParser(&options, flags.Default)

func main() {
	if _, err := parser.Parse(); err != nil {
		switch flagsErr := err.(type) {
		case flags.ErrorType:
			if flagsErr == flags.ErrHelp {
				os.Exit(0)
			}
			os.Exit(1)
		default:
			os.Exit(1)
		}
	}

	spew.Dump(options)

	switch options.Mode {
	case SYNC:
	case COPY:
	}
}
