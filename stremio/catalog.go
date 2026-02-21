package stremio

import (
	"clio/core"
	"fmt"
	"net/url"
)

type Catalog struct {
	Addon *Addon `json:"-"`

	Type   string  `json:"type"`
	Id     string  `json:"id"`
	Name   string  `json:"name"`
	Extras []Extra `json:"extra"`
}

type Extra struct {
	Name string `json:"name"`
}

type SearchResult struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// Catalog

func (c *Catalog) FullName() string {
	return fmt.Sprintf("%s | %s - %s", c.Addon.Name, c.Type, c.Name)
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
	searchUrl := fmt.Sprintf("%s/catalog/%s/%s/search=%s.json", c.Addon.Url, c.Type, c.Id, url.QueryEscape(query))

	res, err := core.GetJson[struct{ Metas []SearchResult }](searchUrl)
	if err != nil {
		return nil, err
	}

	return res.Metas, nil
}
