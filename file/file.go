// Copyright © 2020 Developer Network, LLC
//
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of this source code package.

package file

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Info is an interface expansion of os.FileInfo
// which includes a Path() method which returns
// the path of the absolute path of a file if one
// is set
type Info interface {
	os.FileInfo
	Path() string
}

type info struct {
	os.FileInfo
	path string
}

// Path returns the path value of the file if one is set
func (i info) Path() string {
	return i.path
}

// Recurse is a method for concurrently recursing through a directory
// structure and returning the file info over a channel to the caller
func Recurse(ctx context.Context, path string) (<-chan Info, error) {
	if path == "" {
		return nil, errors.New("empty path")
	}

	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	fmt.Println(path)

	// Pull the file stats
	dir, err := os.Stat(abs)
	if err != nil {
		return nil, err
	}

	if !dir.IsDir() {
		return nil, fmt.Errorf(
			"%s%s%s is not a directory",
			path,
			string(os.PathSeparator),
			dir,
		)
	}

	files := make(chan Info)

	// Start the internal routine
	go func() {
		defer close(files)
		//TODO: decide what to do about errors here
		_ = recurse(ctx, path, files)
	}()

	return files, nil
}

// recurse iterates through each file in a directory
// spawning new routines for additional directories
// and pushing files to the channel
func recurse(
	ctx context.Context,
	path string,
	out chan<- Info,
) error {

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {

			// Recurse the directory
			// TODO: do something with the errors
			_ = recurse(ctx, filepath.Join(path, file.Name()), out)
		}

		// Push the file to the files channel for processing
		select {
		case <-ctx.Done():
			return nil
		case out <- info{file, path}:
			// Pushed file to the files channel
		}

	}

	return nil
}
