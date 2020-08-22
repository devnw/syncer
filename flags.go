package main

// Mode is the operating mode which comes in from a syncer command
// line flag so that the application knows which logic to execute
type Mode int8

const (
	//SYNC is the mode for synchronizing directories versus straight copy
	SYNC Mode = iota

	// COPY is the mode for copying files from the source directory to
	// to the destination directory disregarding directory differences
	// between the two
	COPY Mode = iota
)

var options struct {
	srcDir  string ``
	destDir string ``
	mode    string ``
}
