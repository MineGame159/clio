package stremio

import (
	"clio/core"
	"fmt"
	"regexp"
	"strings"
)

type StreamProvider struct {
	Addon *Addon `json:"-"`

	Types      []string `json:"types"`
	IdPrefixes []string `json:"idPrefixes"`
}

type Stream struct {
	Name        string        `json:"name"`
	Title       string        `json:"title,omitempty"`
	Description string        `json:"description,omitempty"`
	Url         string        `json:"url,omitempty"`
	Hints       BehaviorHints `json:"behaviorHints,omitempty"`
}

type BehaviorHints struct {
	BingeGroup string `json:"bingeGroup"`
	Filename   string `json:"filename,omitempty"`
	VideoSize  uint64 `json:"videoSize,omitempty"`
}

// StreamProvider

func (s *StreamProvider) SupportsId(id string) bool {
	for _, prefix := range s.IdPrefixes {
		if strings.HasPrefix(id, prefix) {
			return true
		}
	}

	return false
}

func (s *StreamProvider) Search(kind string, id string) ([]Stream, error) {
	searchUrl := fmt.Sprintf("%s/stream/%s/%s.json", s.Addon.Url, kind, id)

	res, err := core.GetJson[struct{ Streams []Stream }](searchUrl)
	if err != nil {
		return nil, err
	}

	return res.Streams, nil
}

func (s *StreamProvider) SearchEpisode(kind string, id string, season uint, episode uint) ([]Stream, error) {
	searchUrl := fmt.Sprintf("%s/stream/%s/%s:%d:%d.json", s.Addon.Url, kind, id, season, episode)

	res, err := core.GetJson[struct{ Streams []Stream }](searchUrl)
	if err != nil {
		return nil, err
	}

	return res.Streams, nil
}

// Stream

var resolutionRegex = regexp.MustCompile("\\D(\\d{3,4}[pP])\\W")
var sizeRegex = regexp.MustCompile("\\D(\\d+(?:\\.\\d+)? [a-zA-Z]{2})\\W")

func (s *Stream) TitleDescription() string {
	if s.Description != "" {
		return s.Description
	}
	return s.Title
}

func (s *Stream) TorrentName() string {
	if s.Hints.Filename != "" {
		return stringUpToLast(s.Hints.Filename, '.')
	}

	if s.TitleDescription() != "" {
		return strings.TrimSpace(strings.TrimPrefix(stringUpToFirst(s.TitleDescription(), '\n'), "📄"))
	}

	return ""
}

func (s *Stream) Resolution() string {
	if submatches := resolutionRegex.FindStringSubmatch(s.TitleDescription()); len(submatches) == 2 {
		return submatches[1]
	}

	return ""
}

func (s *Stream) Size() core.ByteSize {
	if s.Hints.VideoSize != 0 {
		return core.ByteSize(s.Hints.VideoSize)
	}

	if submatches := sizeRegex.FindStringSubmatch(s.TitleDescription()); len(submatches) == 2 {
		if size, err := core.ParseByteSize(submatches[1]); err == nil {
			return size
		}
	}

	return 0
}

func stringUpToFirst(str string, char byte) string {
	index := strings.IndexByte(str, char)

	if index == -1 {
		return str
	}

	return str[:index]
}

func stringUpToLast(str string, char byte) string {
	index := strings.LastIndexByte(str, char)

	if index == -1 {
		return str
	}

	return str[:index]
}
