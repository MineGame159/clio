package scraper

import (
	"clio/core"
	"clio/rd"
	"clio/stremio"
	"net/http"
	"strconv"
	"time"
)

func (a *Addon) handleCheck(res http.ResponseWriter, req *http.Request) {
	// Read magnet link
	magnet, _, err := readMagnet(req)
	if err != nil {
		core.WriteError(res, err.Error(), http.StatusBadRequest)
		return
	}

	// Read season
	season, err := strconv.Atoi(req.PathValue("season"))
	if err != nil {
		core.WriteError(res, err.Error(), http.StatusBadRequest)
		return
	}

	// Episode
	episode, err := strconv.Atoi(req.PathValue("episode"))
	if err != nil {
		core.WriteError(res, err.Error(), http.StatusBadRequest)
		return
	}

	// Add magnet to library
	id, err := a.rd.AddMagnet(magnet)
	if err != nil {
		core.WriteError(res, err.Error(), http.StatusBadRequest)
		return
	}

	// Select files
	if err := a.selectFilesFromTorrent(id, season == -1 && episode == -1); err != nil {
		_ = a.rd.DeleteTorrent(id)

		core.WriteError(res, err.Error(), http.StatusBadRequest)
		return
	}

	// Check status
	cached := false

	for i := 0; i < 10; i++ {
		time.Sleep(time.Millisecond * 250)

		torrent, _, err := a.rd.GetTorrent(id)
		if err != nil {
			_ = a.rd.DeleteTorrent(id)

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

	// Delete torrent
	if err := a.rd.DeleteTorrent(id); err != nil {
		core.WriteError(res, err.Error(), http.StatusBadRequest)
		return
	}

	// Response
	core.WriteJson(res, stremio.StreamCheck{
		Cached: cached,
	})
}
