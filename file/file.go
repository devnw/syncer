// Copyright © 2020 Developer Network, LLC
//
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of this source code package.

package file

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/devnw/alog"
)

// Info is an interface expansion of os.FileInfo
// which includes a Path() method which returns
// the path of the absolute path of a file if one
// is set
type Info interface {
	os.FileInfo
	Path() string
	Hash() string
}

type info struct {
	os.FileInfo
	path string
	hash string
}

// Path returns the path value of the file if one is set
func (i info) Path() string {
	return i.path
}

// Hash returns the sha256 hash of the file
func (i info) Hash() string {
	return i.hash
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
		select {
		case <-ctx.Done():
			return
		case files <- <-recurse(ctx, path):
		}
	}()

	// Hash the files before passing out
	out := H256(ctx, files)

	return out, nil
}

// recurse iterates through each file in a directory
// spawning new routines for additional directories
// and pushing files to the channel
func recurse(
	ctx context.Context,
	path string,
) <-chan Info {

	out := make(chan Info)

	go func() {
		//defer close(out)

		files, err := ioutil.ReadDir(path)
		if err != nil {
			return
		}

		for _, file := range files {
			fmt.Printf("file %s\n", file.Name())
			if file.IsDir() {

				// Recurse the directory
				// TODO: do something with the errors
				files := recurse(
					ctx,
					filepath.Join(
						path,
						file.Name(),
					),
				)

				go func(files <-chan Info) {

					for {
						select {
						case <-ctx.Done():
							return
						case out <- <-files:
						}
					}
				}(files)

			}

			// Push the file to the files channel for processing
			select {
			case <-ctx.Done():
				return
			case out <- info{file, path, ""}:
				// Pushed file to the files channel
				fmt.Printf("Pushed %s\n", file.Name())
			}

		}
	}()

	return out
}

// H256 takes in file information and hashes the file
func H256(ctx context.Context, files <-chan Info) <-chan Info {
	out := make(chan Info)

	// Start the internal routine
	go func() {
		defer close(out)
		h := sha256.New()

		select {
		case <-ctx.Done():
			return
		case file, ok := <-files:
			if ok {
				fmt.Printf("summing file %s\n", file.Name())
				f, err := os.Open(
					filepath.Join(
						file.Path(),
						file.Name(),
					),
				)

				if err != nil {
					alog.Fatal(err, "error loading file")
				}

				defer func() {
					_ = f.Close()
				}()

				if _, err := io.Copy(h, f); err != nil {
					alog.Fatal(err)
				}

				select {
				case <-ctx.Done():
					return
				case out <- info{
					file,
					file.Path(),
					base64.StdEncoding.EncodeToString(
						h.Sum(nil),
					),
				}:
				}

			} else {
				return
			}
		}
	}()

	return out
}
