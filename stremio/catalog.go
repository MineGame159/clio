package stremio

import (
	"encoding/json"
	"fmt"
	"net/http"
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
	res, err := http.Get(fmt.Sprintf("%s/catalog/%s/%s/search=%s.json", c.Addon.Url, c.Type, c.Id, url.QueryEscape(query)))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var body struct{ Metas []SearchResult }
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return nil, err
	}

	return body.Metas, nil
}
