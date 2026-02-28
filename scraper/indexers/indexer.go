package indexers

import (
	"clio/core"
	"context"
	"fmt"
	"sync"
)

type Torrent struct {
	Name    string
	Hash    string
	Magnet  string
	Size    core.ByteSize
	Seeders uint

	Season  int
	Episode int
}

var imdbIndexers = []func(ctx context.Context, _ *sync.WaitGroup, id string, torrents chan Torrent){
	scrapeEzTv,
	//scrapeThePirateBay,
	scrapeTorrentGalaxy,
}

var nameIndexers = []func(ctx context.Context, _ *sync.WaitGroup, name string, torrents chan Torrent){
	scrapeKnaben,
}

func Scrape(ctx context.Context, kind, id string) chan Torrent {
	wg := sync.WaitGroup{}
	torrents := make(chan Torrent)

	// ID indexers
	for _, indexer := range imdbIndexers {
		wg.Go(func() {
			indexer(ctx, &wg, id, torrents)
		})
	}

	// Name indexers
	url := fmt.Sprintf("https://v3-cinemeta.strem.io/meta/%s/%s.json", kind, id)

	if res, err := core.GetJson[struct {
		Meta struct {
			Name string
		}
	}](url); err == nil {
		for _, indexer := range nameIndexers {
			wg.Go(func() {
				indexer(ctx, &wg, res.Meta.Name, torrents)
			})
		}
	}

	// Return
	go func() {
		wg.Wait()
		close(torrents)
	}()

	return torrents
}
