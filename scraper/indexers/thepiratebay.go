package indexers

import (
	"clio/core"
	"context"
	"fmt"
	"strconv"
	"sync"
)

type tpbTorrent struct {
	Name    string `json:"name"`
	Hash    string `json:"info_hash"`
	Size    string `json:"size"`
	Seeders string `json:"seeders"`
}

func scrapeThePirateBay(ctx context.Context, _ *sync.WaitGroup, id string, torrents chan Torrent) {
	url := fmt.Sprintf("https://apibay.org/q.php?q=%s&cat=200", id)

	res, err := core.GetJsonCtx[[]tpbTorrent](ctx, url)
	if err != nil {
		return
	}

	for _, torrent := range res {
		seeders, err := strconv.ParseUint(torrent.Seeders, 10, 32)
		if err != nil || seeders == 0 {
			continue
		}

		size, err := strconv.ParseUint(torrent.Size, 10, 64)
		if err != nil {
			continue
		}

		t := Torrent{
			Name:    torrent.Name,
			Hash:    torrent.Hash,
			Magnet:  getMagnetLink(torrent.Name, torrent.Hash),
			Size:    core.ByteSize(size),
			Seeders: uint(seeders),
		}
		parseSeasonEpisode(&t)

		torrents <- t
	}
}
