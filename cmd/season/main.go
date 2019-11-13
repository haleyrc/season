package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/haleyrc/season"
)

func main() {
	confirm := flag.Bool(
		"confirm",
		false,
		"Set to perform the rename instead of just displaying the intended changes")
	remove := flag.String(
		"remove",
		"",
		"Garbage stuff to remove preprocessing",
	)
	helpShort := flag.Bool(
		"h",
		false,
		"Print usage information",
	)
	multi := flag.Bool(
		"multi",
		false,
		"Directory contains multiple seasons in subdirectories",
	)
	debug := flag.Bool(
		"debug",
		false,
		"Output debugging information",
	)
	help := flag.Bool(
		"help",
		false, "Print usage information",
	)
	flag.Parse()

	if *help || *helpShort {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(2)
	}

	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: could not get current directory: %v\n", err)
		os.Exit(1)
	}

	mods, err := season.ScanV2(
		dir,
		season.WithGarbage(*remove),
		season.WithNested(*multi),
		season.WithDebug(*debug),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: could not scan files: %s: %v", dir, err)
		os.Exit(1)
	}

	fmt.Fprintln(os.Stderr)
	mods.Display(os.Stderr)

	if !*confirm {
		fmt.Fprintln(os.Stderr, " Run with --confirm to rename the files.")
		os.Exit(0)
	}

	if errs := mods.Move(); errs != nil {
		fmt.Fprintf(os.Stderr, " Errors:\n\n")
		fmt.Fprintln(os.Stderr, err)
	}
}

var usage = `Usage: season [OPTION]...

If no option is provided, season will process the files in the current directory
and print a list of the proposed renames to stderr.

The following options are available:

  -h, --help      Display this usage and exit. Overrides other options.
  --remove=EXP    Remove the expression given by EXP from each filename prior
                  to running the standard normalization routines.
  --confirm       Rename all files in the current directory according to the
                  normalization rules. Applies any additional preprocessing
                  given by the remove option prior to normalization.

Exit status:
0  if ok
1  if error
2  if help option provided`
