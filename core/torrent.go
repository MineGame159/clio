package core

import (
	"math"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

type NameInfo struct {
	Name    string
	Season  int
	Episode int
}

var seasonRegex = regexp.MustCompile("(?:^|[.\\s\\d_-])[sS](\\d{1,4})")
var episodeRegex = regexp.MustCompile("(?:^|[.\\s\\d_-])[eE](\\d{1,4})")
var stopWordRegex = regexp.MustCompile("(?i:\\b(?:1080p|720p|2160p|4k|480p|bluray|web-?dl|webrip|bdremux|remux|hdtv|divx|xvid|x264|x265|h\\.264|h\\.265|av1|hdr|hevc|avc|10bit|aac|dts|truehd)\\b)")

var replacer = strings.NewReplacer(".", " ", "_", " ", "-", " ")

func ParseTorrentName(name string) NameInfo {
	name = strings.TrimSpace(name)

	// Remove starting brace pairs

	for {
		var trimmed bool
		name, trimmed = trimStartBracePair(name)

		if !trimmed {
			break
		}

		name = strings.TrimSpace(name)
	}

	// Season + Episode

	seasonLoc := seasonRegex.FindStringSubmatchIndex(name)
	episodeLoc := episodeRegex.FindStringSubmatchIndex(name)

	season := -1
	episode := -1

	if seasonLoc != nil {
		v, _ := strconv.ParseInt(name[seasonLoc[2]:seasonLoc[3]], 10, 32)
		season = int(v)
	}

	if episodeLoc != nil {
		v, _ := strconv.ParseInt(name[episodeLoc[2]:episodeLoc[3]], 10, 32)
		episode = int(v)
	}

	// Stop word

	stopWordLoc := stopWordRegex.FindStringIndex(name)

	// Trim up to season / episode / stop word

	if index := minIndex(seasonLoc, episodeLoc, stopWordLoc); index != -1 {
		name = name[:index]
	}

	// Return

	return NameInfo{
		Name:    replacer.Replace(name),
		Season:  season,
		Episode: episode,
	}
}

func minIndex(locs ...[]int) int {
	index := math.MaxInt

	for _, loc := range locs {
		if loc != nil {
			index = min(index, loc[0])
		}
	}

	if index == math.MaxInt {
		return -1
	}

	return index
}

func trimStartBracePair(value string) (string, bool) {
	depth := 0

	for i, ch := range value {
		if depth == 0 {
			switch ch {
			case '(', '[', '{':
				depth++
			default:
				return value, false
			}
		} else {
			switch ch {
			case '(', '[', '{':
				depth++
			case ')', ']', '}':
				depth--

				if depth == 0 {
					return value[i+utf8.RuneLen(ch):], true
				}
			}
		}
	}

	return value, false
}
