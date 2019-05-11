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

type File struct {
	Original string
	Modified string
}

func Scan(base string) (Files, error) {
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

	for _, file := range files {
		dir, fn := filepath.Split(file)
		mod := transform(fn, episodeNumberLength)
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

type Files struct {
	BasePath string
	Files    []File
}

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

type transformer func(s string) string

// TODO (RCH): Take flags like --remove "Some garbage$" and create xformers from them
func transform(input string, episodeNumberLength int) string {
	ext := filepath.Ext(input)
	fn := strings.TrimSuffix(input, ext)

	for _, xformer := range []transformer{
		removeGarbage(`- Meteor For Everyone$`),
		removeGarbage(`Mastering Figma `),
		removeGarbage(`^Level 2 Meteor 1\.4 \+ React `),
		removeGarbage(`Level 1 Apollo Client with React$`),
		removeGarbage(`Full-stack GraphQL with Apollo, Meteor & React-.*$`),
		removeGarbage(`Better JavaScript$`),
		removeGarbage(`VueJS For Everyone-.*$`),
		replaceEscaped,
		removeNonAlphaNumeric,
		prependEpisode(episodeNumberLength),
		strings.TrimSpace,
		replaceSpaces,
		removeDupeUnderscores,
	} {
		fn = xformer(fn)
	}

	return fn + ext
}

func replaceEscaped(s string) string {
	return regexp.MustCompile(`%26`).ReplaceAllString(s, "And")
}

func removeGarbage(pattern string) transformer {
	return func(s string) string {
		return regexp.MustCompile(pattern).ReplaceAllString(s, "")
	}
}

func removeDupeUnderscores(s string) string {
	return regexp.MustCompile(`_+`).ReplaceAllString(s, "_")
}

func replaceSpaces(s string) string {
	return regexp.MustCompile(" ").ReplaceAllString(s, "_")
}

func prependEpisode(l int) transformer {
	return func(s string) string {
		episode := regexp.MustCompile(`^\d*`).FindString(s)
		s = strings.TrimPrefix(s, episode)

		// We always want at least 2 digits, more only if there are
		// more than 99 "episodes"
		if len(episode) < 2 {
			episode = "0" + episode
		}
		episode = strings.Repeat("0", l-len(episode)) + episode
		prefix := fmt.Sprintf("S01E%s", episode)

		return prefix + " " + s
	}
}

func removeNonAlphaNumeric(s string) string {
	return regexp.MustCompile(`[^a-zA-Z0-9_ ]`).ReplaceAllString(s, "")
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
