package season

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

type transformer func(s string) string

func transform(input string, episodeNumberLength int, prexformers ...transformer) string {
	ext := filepath.Ext(input)
	fn := strings.TrimSuffix(input, ext)

	for _, xformer := range prexformers {
		fn = xformer(fn)
	}

	for _, xformer := range []transformer{
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
