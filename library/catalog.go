package library

import (
	"clio/core"
	"clio/stremio"
	"cmp"
	"net/http"
	"slices"
)

func (a *Addon) handleCatalog(res http.ResponseWriter, req *http.Request) {
	kind := stremio.MediaKind(req.PathValue("kind"))

	// Fetch results based on requested media kind
	var results []stremio.SearchResult

	for _, info := range a.media {
		if info.kind == kind {
			results = append(results, stremio.SearchResult{
				Id:     info.id,
				Name:   info.name,
				Poster: info.poster,
			})
		}
	}

	slices.SortFunc(results, func(a, b stremio.SearchResult) int {
		return cmp.Compare(a.Name, b.Name)
	})

	// Write response
	core.WriteJson(res, struct {
		Metas []stremio.SearchResult
	}{results})
}
