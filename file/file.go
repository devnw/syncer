package file

import "os"

// Recurse is a method for concurrently recursing through a directory
// structure and returning the file info over a channel to the caller
func Recurse(dir string) <-chan os.FileInfo {

	return nil
}
