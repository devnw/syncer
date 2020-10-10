// Copyright © 2020 Developer Network, LLC
//
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of this source code package.

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/davecgh/go-spew/spew"
	"github.com/devnw/syncer/file"
	"github.com/jessevdk/go-flags"
)

var options Options
var parser = flags.NewParser(&options, flags.Default)

func main() {
	ctx := initContext()

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

	files, err := file.Recurse(ctx, options.SrcDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	printFiles(ctx, files)
}

func printFiles(ctx context.Context, in <-chan file.Info) {

	for {
		select {
		case <-ctx.Done():
			return
		case f, ok := <-in:
			if ok {
				fmt.Println(
					filepath.Join(
						f.Path(),
						f.Name(),
					),
				)

				fmt.Println(f.Hash())
			} else {
				return
			}
		}
	}
}

func initContext() context.Context {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())

	// Setup interrupt monitoring for the agent
	go func() {
		defer cancel()

		select {
		case <-ctx.Done():
			return
		case <-sigs:
			fmt.Println("exiting syncer")
			os.Exit(1)
		}
	}()

	return ctx
}
