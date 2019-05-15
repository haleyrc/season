# season

[`season`](https://github.com/haleyrc/season) is a utility for quickly renaming
sets of files according to an opinionated set of rules with the end result that
they are parseable by the [Plex Media Server](https://www.plex.tv/) as a TV show.

This utility was designed exclusively to help me with transforming a large
number of online course videos I had downloaded that were named in entirely
disparate ways, causing Plex to fail to correctly organize and present them.

There is no guarantee that this utility will work correctly for any other
material, so use it at your own risk. See below for instructions on how to
verify the intended output before actually going through with the renaming
process.

## Install

This is currently only available as a Go packge, so you will need Go version
1.11 or higher to install it.

`$ go install github.com/haleyrc/season/cmd/season`

## Usage

```
Usage: season [OPTION]...

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
  2  if help option provided
```
