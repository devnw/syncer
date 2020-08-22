package main

// Mode is the operating mode which comes in from a syncer command
// line flag so that the application knows which logic to execute
type Mode string

const (
	//SYNC is the mode for synchronizing directories versus straight copy
	SYNC Mode = "SYNC"

	// COPY is the mode for copying files from the source directory to
	// to the destination directory disregarding directory differences
	// between the two
	COPY Mode = "COPY"
)

// Options are the command line options for syncer
type Options struct {

	// SrcDir is the flag for the source directory
	SrcDir string `short:"s" long:"src" description:"Source Directory" required:"true"`

	//DestDir is the flag for the destination directory
	DestDir string `short:"d" long:"dest" description:"Destination Directory" required:"true"`

	// Mode is the mode for configuring the sync logic
	Mode Mode `short:"m" long:"mode" choice:"SYNC" choice:"COPY" default:"SYNC"`
}
