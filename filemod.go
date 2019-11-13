package season

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
)

type FileMod struct {
	base   string
	subdir string
	from   string
	toName string
	toExt  string
}

func NewFileMods(base string) FileMods {
	return FileMods{
		base:    base,
		mods:    make([]*FileMod, 0),
		seasons: make(map[string]int),
	}
}

type FileMods struct {
	base    string
	mods    []*FileMod
	seasons map[string]int
}

// Display prints information about the path and the files to be renamed, along
// with their transformed name for review by the user.
func (mods FileMods) Display(w io.Writer) {
	fmt.Fprintf(w, " Path: %s\n", mods.base)
	fmt.Fprintln(w, " The following files will be renamed:")
	fmt.Fprintln(w)

	longest := findLongest(mods.mods)
	for _, mod := range mods.mods {
		from := filepath.Join(mod.subdir, mod.from)
		to := filepath.Join(mod.toName + mod.toExt)
		fmt.Fprintf(w, "    %s\n", padToLength(from, longest))
		fmt.Fprintf(w, "      -> %s\n\n", to)
	}
}

// Move moves each file from its original location to its modified location
// within the same base path.
func (mods FileMods) Move() []error {
	var errs []error
	for _, mod := range mods.mods {
		oldpath := filepath.Join(mod.base, mod.subdir, mod.from)
		newpath := filepath.Join(mod.base, mod.toName+mod.toExt)
		if err := os.Rename(oldpath, newpath); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

func ScanV2(base string, opts ...Option) (FileMods, error) {
	var cfg scanConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	var err error
	mods := NewFileMods(base)
	if cfg.nested {
		mods, err = getNestedFiles(base)
	} else {
		mods, err = getFiles(base)
	}
	if err != nil {
		return FileMods{}, errors.Wrap(err, "ScanV2")
	}
	if cfg.debug {
		spew.Dump(mods)
	}

	var prexformers []transformer
	if cfg.garbage != "" {
		prexformers = append(prexformers)
	}

	for _, mod := range mods.mods {
		// This will just be 0 if we're in flat mode
		numSeasons := len(mods.seasons)
		// Default to setting number of episodes to number of total files, which
		// is true in the flat case.
		numEpisodes := len(mods.mods)
		if cfg.debug {
			fmt.Printf("num seasons : %d\n", numSeasons)
			fmt.Printf("num episodes: %d\n", numEpisodes)
		}
		if cfg.nested {
			// If we're in the nested case and have seasons, set it to the
			// number of files for just that season.
			numEpisodes = mods.seasons[mod.subdir]
			if cfg.debug {
				fmt.Println("Using seasons...")
				fmt.Printf("num episodes: %d\n", numEpisodes)
			}
		}
		if cfg.debug {
			fmt.Println("Pre-transform:")
			spew.Dump(mod)
			fmt.Println()
		}
		for _, transformer := range []struct {
			name string
			f    transformer
		}{
			{"removing garbage", removeGarbage(cfg.garbage)},
			{"replacing escaped", replaceEscaped},
			{"removing nonalphanumeric", removeNonAlphaNumeric},
			{"trimming space 1", trimSpace},
			{"trim leading zeroes", trimZeroes},
			{"prepending episode", prependEpisode(numSeasons, numEpisodes)},
			{"trimming space 2", trimSpace},
			{"replacing spaces", replaceSpaces},
			{"removing dupe underscores", removeDupeUnderscores},
		} {
			if cfg.debug {
				fmt.Printf("%s:\n", transformer.name)
			}
			transformer.f(mod)
			if cfg.debug {
				spew.Dump(mod)
				fmt.Println()
			}
		}
	}

	return mods, nil
}

func getFiles(base string) (FileMods, error) {
	all, err := filepath.Glob(filepath.Join(base, "*"))
	if err != nil {
		return FileMods{}, errors.Wrap(err, "ScanV2")
	}

	mods := NewFileMods(base)
	for _, match := range all {
		fi, err := os.Stat(match)
		if err != nil {
			return FileMods{}, errors.Wrap(err, "ScanV2")
		}
		if fi.IsDir() {
			continue
		}
		ext := filepath.Ext(match)
		if !isVideo(ext) {
			continue
		}
		fn := strings.TrimSuffix(filepath.Base(match), ext)
		mods.mods = append(mods.mods, &FileMod{
			base:   base,
			subdir: "",
			from:   fn + ext,
			toName: fn,
			toExt:  ext,
		})
	}

	return mods, nil
}

func getNestedFiles(base string) (FileMods, error) {
	all, err := filepath.Glob(filepath.Join(base, "*"))
	if err != nil {
		return FileMods{}, errors.Wrap(err, "ScanV2")
	}

	mods := NewFileMods(base)
	for _, match := range all {
		fi, err := os.Stat(match)
		if err != nil {
			return FileMods{}, errors.Wrap(err, "ScanV2")
		}
		if !fi.IsDir() {
			continue
		}
		subfiles, err := filepath.Glob(filepath.Join(match, "*"))
		if err != nil {
			return FileMods{}, errors.Wrap(err, "ScanV2")
		}
		for _, submatch := range subfiles {
			ext := filepath.Ext(submatch)
			if !isVideo(ext) {
				continue
			}
			fn := strings.TrimSuffix(filepath.Base(submatch), ext)
			mods.mods = append(mods.mods, &FileMod{
				base:   base,
				subdir: fi.Name(),
				from:   filepath.Base(submatch),
				toName: fn,
				toExt:  ext,
			})
			mods.seasons[fi.Name()]++
		}
	}

	return mods, nil
}

func padToLength(s string, l int) string {
	if len(s) > l {
		panic(fmt.Errorf("invalid input to padToLength(%s, %d)", s, l))
	}
	needed := l - len(s)
	return s + strings.Repeat(" ", needed)
}

func findLongest(mods []*FileMod) int {
	var longest int
	for _, mod := range mods {
		cl := len(filepath.Join(mod.subdir, mod.from))
		if cl > longest {
			longest = cl
		}
	}
	return longest
}

func isVideo(ext string) bool {
	ext = strings.ToLower(ext)
	switch ext {
	case
		".mp4",
		".mov":
		return true
	default:
		return false
	}
}
