package stremio

import (
	"clio/core"
	"context"
	"fmt"
	"strings"
)

type StreamProvider struct {
	Addon *Addon `json:"-"`

	Kinds      []MediaKind `json:"types"`
	IdPrefixes []string    `json:"idPrefixes"`
}

type Stream struct {
	Name        string        `json:"name"`
	Title       string        `json:"title,omitempty"`
	Description string        `json:"description,omitempty"`
	Url         string        `json:"url,omitempty"`
	RedirectUrl bool          `json:"redirectUrl,omitempty"`
	CheckUrl    string        `json:"checkUrl,omitempty"`
	Hints       BehaviorHints `json:"behaviorHints,omitempty"`
}

type BehaviorHints struct {
	BingeGroup string `json:"bingeGroup"`
	Filename   string `json:"filename,omitempty"`
	VideoSize  uint64 `json:"videoSize,omitempty"`
}

type StreamCheck struct {
	Cached bool `json:"cached"`
}

// StreamProvider

func (s *StreamProvider) SupportsKindId(kind MediaKind, id string) bool {
	supportsKind := false

	for _, providerKind := range s.Kinds {
		if providerKind == kind {
			supportsKind = true
			break
		}
	}

	if !supportsKind {
		return false
	}

	for _, prefix := range s.IdPrefixes {
		if strings.HasPrefix(id, prefix) {
			return true
		}
	}

	return false
}

func (s *StreamProvider) Search(ctx context.Context, kind MediaKind, id string) ([]Stream, error) {
	searchUrl := fmt.Sprintf("%s/stream/%s/%s.json", s.Addon.Url, kind, id)

	res, err := core.GetJsonCtx[struct{ Streams []Stream }](ctx, searchUrl)
	if err != nil {
		return nil, err
	}

	return res.Streams, nil
}

func (s *StreamProvider) SearchEpisode(ctx context.Context, kind MediaKind, id string, season uint, episode uint) ([]Stream, error) {
	searchUrl := fmt.Sprintf("%s/stream/%s/%s:%d:%d.json", s.Addon.Url, kind, id, season, episode)

	res, err := core.GetJsonCtx[struct{ Streams []Stream }](ctx, searchUrl)
	if err != nil {
		return nil, err
	}

	return res.Streams, nil
}

// Stream

func (s *Stream) TitleDescription() string {
	if s.Description != "" {
		return s.Description
	}
	return s.Title
}
