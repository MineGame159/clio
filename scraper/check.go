package scraper

import (
	"clio/core"
	"clio/stremio"
	"net/http"
)

func (a *Addon) handleCheck(res http.ResponseWriter, req *http.Request) {
	// Read magnet link
	magnet, _, err := readMagnet(req)
	if err != nil {
		core.WriteError(res, err.Error(), http.StatusBadRequest)
		return
	}

	// Check magnet
	cached, dFiles, err := a.client.CheckMagnet(req.Context(), magnet)
	if err != nil {
		core.WriteError(res, err.Error(), http.StatusBadRequest)
		return
	}

	files := make([]stremio.File, 0, len(dFiles))

	for _, file := range dFiles {
		files = append(files, stremio.File{
			Path: file.Path,
			Size: file.Size,
		})
	}

	// Response
	core.WriteJson(res, stremio.StreamCheck{
		Cached: cached,
		Files:  files,
	})
}
