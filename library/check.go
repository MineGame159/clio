package library

import (
	"clio/core"
	"clio/stremio"
	"net/http"
)

func (a *Addon) handleCheck(res http.ResponseWriter, req *http.Request) {
	// Get torrent files
	_, tFiles, err := a.rd.GetTorrent(req.Context(), req.PathValue("id"))
	if err != nil {
		core.WriteError(res, err.Error(), http.StatusBadRequest)
		return
	}

	// Construct an array of stremio.File s
	files := make([]stremio.File, 0, len(tFiles))

	for _, file := range tFiles {
		if file.Selected != 0 {
			files = append(files, stremio.File{
				Path: file.Path,
				Size: file.Size,
			})
		}
	}

	// Write result
	core.WriteJson(res, stremio.StreamCheck{
		Cached: true,
		Files:  files,
	})
}
