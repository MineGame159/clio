package stremio

import (
	"clio/core"
	"fmt"
	"net/url"
)

type Catalog struct {
	Addon *Addon `json:"-"`

	Kind   MediaKind `json:"type"`
	Id     string    `json:"id"`
	Name   string    `json:"name"`
	Extras []Extra   `json:"extra"`
}

type MediaKind string

const (
	Movie  MediaKind = "movie"
	Series MediaKind = "series"
	Anime  MediaKind = "anime"
	Other  MediaKind = "other"
)

type Extra struct {
	Name string `json:"name"`
}

type SearchResult struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Poster string `json:"poster"`
}

// Catalog

func (c *Catalog) FullName() string {
	return fmt.Sprintf("%s | %s - %s", c.Addon.Name, c.Kind.Name(), c.Name)
}

func (c *Catalog) HasExtra(name string) bool {
	for _, extra := range c.Extras {
		if extra.Name == name {
			return true
		}
	}

	return false
}

func (c *Catalog) Search(query string) ([]SearchResult, error) {
	searchUrl := fmt.Sprintf("%s/catalog/%s/%s/search=%s.json", c.Addon.Url, c.Kind, c.Id, url.QueryEscape(query))

	res, err := core.GetJson[struct{ Metas []SearchResult }](searchUrl)
	if err != nil {
		return nil, err
	}

	return res.Metas, nil
}

// MediaKind

func (m MediaKind) HasEpisodes() bool {
	switch m {
	case Series, Anime:
		return true

	default:
		return false
	}
}

func (m MediaKind) Name() string {
	switch m {
	case Movie:
		return "Movie"
	case Series:
		return "Series"
	case Anime:
		return "Anime"
	case Other:
		return "Other"
	default:
		return core.Capitalize(string(m))
	}
}
