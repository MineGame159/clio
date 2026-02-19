package stremio

import (
	"encoding/json"
	"fmt"
	"iter"
	"net/http"
	"strings"
)

type MetaProvider struct {
	Addon *Addon

	IdPrefixes []string
}

type Meta struct {
	Videos []Video `json:"videos"`
}

type Video struct {
	Season uint   `json:"season"`
	Number uint   `json:"number"`
	Name   string `json:"name"`
}

// MetaProvider

func (m *MetaProvider) SupportsId(id string) bool {
	for _, prefix := range m.IdPrefixes {
		if strings.HasPrefix(id, prefix) {
			return true
		}
	}

	return false
}

func (m *MetaProvider) Get(kind string, id string) (Meta, error) {
	res, err := http.Get(fmt.Sprintf("%s/meta/%s/%s.json", m.Addon.Url, kind, id))
	if err != nil {
		return Meta{}, err
	}
	defer res.Body.Close()

	var body struct{ Meta Meta }
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return Meta{}, err
	}

	return body.Meta, nil
}

// Meta

func (m *Meta) Seasons() uint {
	count := uint(0)

	for _, video := range m.Videos {
		count = max(count, video.Season+1)
	}

	return count
}

func (m *Meta) Episodes(season uint) iter.Seq[Video] {
	return func(yield func(Video) bool) {
		for _, video := range m.Videos {
			if video.Season == season && !yield(video) {
				return
			}
		}
	}
}
