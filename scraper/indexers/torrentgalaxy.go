package indexers

import (
	"clio/core"
	"context"
	"fmt"
	"math"
	"sync"
)

type tgResponse struct {
	PageSize uint        `json:"page_size"`
	Count    uint        `json:"count"`
	Total    uint        `json:"total"`
	Torrents []tgTorrent `json:"results"`
}

type tgTorrent struct {
	Name    string        `json:"n"`
	Hash    string        `json:"h"`
	Size    core.ByteSize `json:"s"`
	Seeders uint          `json:"se"`
}

func scrapeTorrentGalaxy(ctx context.Context, wg *sync.WaitGroup, id string, torrents chan Torrent) {
	// Fetch first page
	url := fmt.Sprintf("https://torrentgalaxy.space/get-posts/keywords:%s:format:json?page=1", id)

	res, err := core.GetJsonCtx[tgResponse](ctx, url)
	if err != nil {
		return
	}

	for _, torrent := range res.Torrents {
		processTorrentGalaxyTorrent(torrents, torrent)
	}

	// Fetch remaining pages
	pages := uint(math.Ceil(float64(res.Total) / float64(res.PageSize)))

	for page := uint(2); page <= pages; page++ {
		wg.Go(func() {
			url := fmt.Sprintf("https://torrentgalaxy.space/get-posts/keywords:%s:format:json?page=%d", id, page)

			res, err := core.GetJsonCtx[tgResponse](ctx, url)
			if err != nil {
				return
			}

			for _, torrent := range res.Torrents {
				processTorrentGalaxyTorrent(torrents, torrent)
			}
		})
	}
}

func processTorrentGalaxyTorrent(torrents chan Torrent, torrent tgTorrent) {
	if torrent.Seeders == 0 {
		return
	}

	t := Torrent{
		Name:    torrent.Name,
		Hash:    torrent.Hash,
		Magnet:  getMagnetLink(torrent.Name, torrent.Hash),
		Size:    torrent.Size,
		Seeders: torrent.Seeders,
	}
	parseSeasonEpisode(&t)

	torrents <- t
}
