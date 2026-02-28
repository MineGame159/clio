package indexers

import (
	"clio/core"
	"context"
	"fmt"
	"strconv"
	"sync"
)

type eztvResponse struct {
	Count    int           `json:"torrents_count"`
	Torrents []eztvTorrent `json:"torrents"`
}

type eztvTorrent struct {
	Name    string `json:"filename"`
	Title   string `json:"title"`
	Hash    string `json:"hash"`
	Magnet  string `json:"magnet_url"`
	Size    string `json:"size_bytes"`
	Seeders uint   `json:"seeds"`
}

func scrapeEzTv(ctx context.Context, _ *sync.WaitGroup, id string, torrents chan Torrent) {
	page := 1

	for {
		url := fmt.Sprintf("https://eztvx.to/api/get-torrents?imdb_id=%s&limit=100&page=%d", id[2:], page)

		res, err := core.GetJsonCtx[eztvResponse](ctx, url)
		if err != nil {
			return
		}

		for _, torrent := range res.Torrents {
			if torrent.Seeders == 0 {
				continue
			}

			name := torrent.Name
			if name == "" {
				name = torrent.Title
			}

			size, err := strconv.ParseUint(torrent.Size, 10, 64)
			if err != nil {
				continue
			}

			t := Torrent{
				Name:    name,
				Hash:    torrent.Hash,
				Magnet:  torrent.Magnet,
				Size:    core.ByteSize(size),
				Seeders: torrent.Seeders,
			}
			parseSeasonEpisode(&t)

			torrents <- t
		}

		if (page-1)*100+len(res.Torrents) >= res.Count {
			break
		}

		page++
	}
}
