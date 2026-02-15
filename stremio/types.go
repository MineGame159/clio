package stremio

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Addon struct {
	Name string

	Catalogs        []*Catalog
	StreamProviders []*StreamProvider

	Url string
}

type Catalog struct {
	Addon *Addon `json:"-"`

	Type   string  `json:"type"`
	Id     string  `json:"id"`
	Name   string  `json:"name"`
	Extras []Extra `json:"extra"`
}

type StreamProvider struct {
	Addon *Addon `json:"-"`

	Types      []string `json:"types"`
	IdPrefixes []string `json:"idPrefixes"`
}

type Extra struct {
	Name string `json:"name"`
}

type MetaBasic struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Stream struct {
	Name        string        `json:"name"`
	Title       string        `json:"title,omitempty"`
	Description string        `json:"description,omitempty"`
	Url         string        `json:"url,omitempty"`
	Hints       BehaviorHints `json:"behaviorHints,omitempty"`
}

type BehaviorHints struct {
	Filename  string `json:"filename,omitempty"`
	VideoSize uint64 `json:"videoSize,omitempty"`
}

// Catalog

func (c *Catalog) HasExtra(name string) bool {
	for _, extra := range c.Extras {
		if extra.Name == name {
			return true
		}
	}

	return false
}

func (c *Catalog) Search(query string) ([]MetaBasic, error) {
	res, err := http.Get(fmt.Sprintf("%s/catalog/%s/%s/search=%s.json", c.Addon.Url, c.Type, c.Id, url.QueryEscape(query)))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var body struct{ Metas []MetaBasic }
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return nil, err
	}

	return body.Metas, nil
}

func (c *Catalog) FilterValue() string {
	return fmt.Sprintf("%s | %s - %s", c.Addon.Name, c.Type, c.Name)
}

func (c *Catalog) Text() string {
	return fmt.Sprintf("%s | %s - %s", c.Addon.Name, c.Type, c.Name)
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
	res, err := http.Get(fmt.Sprintf("%s/stream/%s/%s.json", s.Addon.Url, kind, id))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var body struct{ Streams []Stream }
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return nil, err
	}

	return body.Streams, nil
}

// MetaBasic

func (m MetaBasic) FilterValue() string {
	return m.Name
}

func (m MetaBasic) Text() string {
	return m.Name
}

// Stream

func (s *Stream) TorrentName() string {
	if s.Title != "" {
		return strings.TrimSpace(strings.TrimPrefix(stringUpToFirst(s.Title, '\n'), "📄"))
	}

	if s.Description != "" {
		return strings.TrimSpace(strings.TrimPrefix(stringUpToFirst(s.Description, '\n'), "📄"))
	}

	if s.Hints.Filename != "" {
		return stringUpToLast(s.Hints.Filename, '.')
	}

	return ""
}

func (s *Stream) FilterValue() string {
	return s.TorrentName()
}

func (s *Stream) Text() string {
	return s.TorrentName()
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
