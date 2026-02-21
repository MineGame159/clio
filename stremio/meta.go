package stremio

import (
	"clio/core"
	"fmt"
	"iter"
	"slices"
	"strings"
)

type MetaProvider struct {
	Addon *Addon

	IdPrefixes []string
}

type Meta struct {
	Id string `json:"id"`

	Name        string `json:"name"`
	Description string `json:"description"`

	Poster      string   `json:"poster"`
	Cast        []string `json:"cast"`
	Rating      string   `json:"imdbRation"`
	ReleaseInfo string   `json:"releaseInfo"`
	Awards      string   `json:"awards"`
	Genres      []string `json:"genre"`

	Runtime string
	Status  string
	Videos  []Video `json:"videos"`
}

type Video struct {
	Season  uint `json:"season"`
	Episode uint `json:"episode"`

	Title string `json:"title"`
	Name  string `json:"name"`
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
	metaUrl := fmt.Sprintf("%s/meta/%s/%s.json", m.Addon.Url, kind, id)

	res, err := core.GetJson[struct{ Meta Meta }](metaUrl)
	if err != nil {
		return Meta{}, err
	}

	return res.Meta, nil
}

// Meta

func (m *Meta) Seasons() []uint {
	var seasons []uint

	for _, video := range m.Videos {
		if !slices.Contains(seasons, video.Season) {
			seasons = append(seasons, video.Season)
		}
	}

	return seasons
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

// Video

func (v *Video) ActualTitle() string {
	if v.Title != "" {
		return v.Title
	}

	return v.Name
}
