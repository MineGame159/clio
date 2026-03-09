package scraper

import (
	"clio/core"
	"clio/rd"
	"clio/stremio"
	"context"
	"net/http"
	"time"
)

func (a *Addon) handleCheck(res http.ResponseWriter, req *http.Request) {
	// Read magnet link
	magnet, _, err := readMagnet(req)
	if err != nil {
		core.WriteError(res, err.Error(), http.StatusBadRequest)
		return
	}

	// Add magnet to library
	id, err := a.rd.AddMagnet(req.Context(), magnet)
	if err != nil {
		core.WriteError(res, err.Error(), http.StatusBadRequest)
		return
	}

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		_ = a.rd.DeleteTorrent(ctx, id)
	}()

	// Select files
	_, tFiles, err := a.rd.GetTorrent(req.Context(), id)
	if err != nil {
		core.WriteError(res, err.Error(), http.StatusBadRequest)
		return
	}

	fileIds := make([]uint, len(tFiles))
	files := make([]stremio.File, len(tFiles))

	for i, file := range tFiles {
		fileIds[i] = file.Id
		files[i] = stremio.File{
			Path: file.Path,
			Size: file.Size,
		}
	}

	if err := a.rd.SelectFiles(req.Context(), id, fileIds); err != nil {
		core.WriteError(res, err.Error(), http.StatusBadRequest)
		return
	}

	// Check status
	cached := false

	for i := 0; i < 10; i++ {
		time.Sleep(time.Millisecond * 250)

		torrent, _, err := a.rd.GetTorrent(req.Context(), id)
		if err != nil {
			core.WriteError(res, err.Error(), http.StatusBadRequest)
			return
		}

		if torrent.Status == rd.Downloaded {
			cached = true
			break
		}

		if torrent.Status == rd.Downloading || torrent.Status.Failed() {
			break
		}
	}

	// Response
	core.WriteJson(res, stremio.StreamCheck{
		Cached: cached,
		Files:  files,
	})
}
