package season

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// File is a single file to be renamed, with it's original file name and the
// modified version. The file names do not include the path.
type File struct {
	Original string
	Modified string
}

// Scan looks for all files in a given directory that have an appropriate
// extension and compiles them with transformed file names. The garbage
// parameter can be used to strip some known quantity from each file name prior
// to running the common normalization. This can be helpful when a set of files
// contains some bits of information that should not appear in the final
// filenames.
func Scan(base string, garbage string) (Files, error) {
	fs := Files{BasePath: base}

	allFiles, err := filepath.Glob(filepath.Join(base, "*"))
	if err != nil {
		return Files{}, errors.Wrap(err, "scan")
	}

	var files []string
	for _, file := range allFiles {
		fi, err := os.Stat(file)
		if err != nil {
			return Files{}, errors.Wrap(err, "scan")
		}
		if fi.IsDir() {
			continue
		}
		if !isVideo(file) {
			continue
		}
		files = append(files, file)
	}

	n := len(files)
	episodeNumberLength := len(strconv.Itoa(n))

	var prexformers []transformer
	if garbage != "" {
		prexformers = append(prexformers, removeGarbage(garbage))
	}
	for _, file := range files {
		dir, fn := filepath.Split(file)
		mod := transform(fn, episodeNumberLength, prexformers...)
		if dir != "" {
			currDir := filepath.Base(dir)
			currDir = regexp.MustCompile(" ").ReplaceAllString(currDir, "_")
			mod = currDir + "_" + mod
		}

		fs.Files = append(fs.Files, File{
			Original: fn,
			Modified: mod,
		})
	}

	return fs, nil
}

// Files is a compilation of a path and all of the files to be transformed
// within that path.
type Files struct {
	BasePath string
	Files    []File
}

// Display prints information about the path and the files to be renamed, along
// with their transformed name for review by the user.
func (fs Files) Display(w io.Writer) {
	fmt.Fprintf(w, " Path: %s\n", fs.BasePath)
	fmt.Fprintln(w, " The following files will be renamed:")
	fmt.Fprintln(w)

	longest := findLongest(fs.Files)
	for _, file := range fs.Files {
		fmt.Fprintf(w, "    %s\n", padToLength(file.Original, longest))
		fmt.Fprintf(w, "      -> %s\n\n", file.Modified)
	}
}

// Move moves each file from its original location to its modified location
// within the same base path.
func (fs Files) Move() []error {
	var errs []error
	for _, file := range fs.Files {
		oldpath := filepath.Join(fs.BasePath, file.Original)
		newpath := filepath.Join(fs.BasePath, file.Modified)
		if err := os.Rename(oldpath, newpath); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

func padToLength(s string, l int) string {
	if len(s) > l {
		panic(fmt.Errorf("invalid input to padToLength(%s, %d)", s, l))
	}
	needed := l - len(s)
	return s + strings.Repeat(" ", needed)
}

func findLongest(files []File) int {
	var longest int
	for _, file := range files {
		if len(file.Original) > longest {
			longest = len(file.Original)
		}
	}
	return longest
}

func isVideo(file string) bool {
	ext := strings.ToLower(filepath.Ext(file))
	switch ext {
	case
		".mp4",
		".mov":
		return true
	default:
		return false
	}
}
