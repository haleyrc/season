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
	flag.Parse()

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	files, err := season.Scan(dir)
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(os.Stderr)
	files.Display(os.Stderr)

	if !*confirm {
		fmt.Fprintln(os.Stderr, " Run with --confirm to rename the files.")
		os.Exit(0)
	}

	if errs := files.Move(); errs != nil {
		fmt.Fprintf(os.Stderr, " Errors:\n\n")
		fmt.Fprintln(os.Stderr, err)
	}
}
