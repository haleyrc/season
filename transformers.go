package season

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

type transformer func(m *FileMod)

func replaceEscaped(mod *FileMod) {
	mod.toName = regexp.MustCompile(`%26`).ReplaceAllString(mod.toName, "And")
}

func removeGarbage(pattern string) transformer {
	return func(mod *FileMod) {
		if pattern == "" {
			return
		}
		mod.toName = regexp.MustCompile(pattern).ReplaceAllString(mod.toName, "")
	}
}

func removeDupeUnderscores(mod *FileMod) {
	mod.toName = regexp.MustCompile(`_+`).ReplaceAllString(mod.toName, "_")
}

func replaceSpaces(mod *FileMod) {
	mod.toName = regexp.MustCompile(" ").ReplaceAllString(mod.toName, "_")
}

func trimZeroes(mod *FileMod) {
	mod.toName = regexp.MustCompile("^0*").ReplaceAllString(mod.toName, "")
}

func prependEpisode(numSeasons, numEpisodes int) transformer {
	return func(mod *FileMod) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Error during processing:\n")
				fmt.Printf("\tCan't prepend episode: %s\n", r)
				fmt.Printf("\tNumber of seasons:  %d\n", numSeasons)
				fmt.Printf("\tNumber of episodes: %d\n", numEpisodes)
				spew.Dump(mod)
				os.Exit(1)
			}
		}()
		episode := regexp.MustCompile(`^\d*`).FindString(mod.toName)
		orig := strings.TrimPrefix(mod.toName, episode)

		episode = pad(episode, 2, len(strconv.Itoa(numEpisodes)))

		season := "1"
		if mod.subdir != "" {
			season = regexp.MustCompile(`^\d*`).FindString(mod.subdir)
		}
		season = pad(season, 2, len(strconv.Itoa(numSeasons)))

		prefix := fmt.Sprintf("S%sE%s", season, episode)

		mod.toName = prefix + " " + orig
	}
}

func pad(s string, min, max int) string {
	l := len(s)
	padding := max - l
	if l < min && padding < 1 {
		padding = 1
	}
	return strings.Repeat("0", padding) + s
}

func removeNonAlphaNumeric(mod *FileMod) {
	mod.toName = regexp.MustCompile(`[^a-zA-Z0-9_ ]`).ReplaceAllString(mod.toName, " ")
}

func trimSpace(mod *FileMod) {
	mod.toName = strings.TrimSpace(mod.toName)
}
