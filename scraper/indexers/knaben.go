package indexers

import (
	"bytes"
	"clio/core"
	"context"
	"encoding/json"
	"sync"
)

/*
{
  "search_type": "100%",
  "search_field": "title",
  "query": "dubioza kolektiv",
  "order_by": "peers",
  "order_direction": "desc",
  "categories": [
    1004000, 1001000
  ],
  "from": 0,
  "size": 150,
  "hide_unsafe": true,
  "hide_xxx": true,
  "seconds_since_last_seen": 86400
}
*/

type knabenRequest struct {
	SearchType           string `json:"search_type"`
	SearchField          string `json:"search_field,omitempty"`
	Query                string `json:"query,omitempty"`
	OrderBy              string `json:"order_by,omitempty"`
	OrderDirection       string `json:"order_direction"`
	Categories           []uint `json:"categories,omitempty"`
	From                 uint   `json:"from,omitempty"`
	Size                 uint   `json:"size,omitempty"`
	HideUnsafe           bool   `json:"hide_unsafe"`
	HideXXX              bool   `json:"hide_xxx,omitempty"`
	SecondsSinceLastSeen uint   `json:"seconds_since_last_seen,omitempty"`
}

type knabenTorrent struct {
	Title   string        `json:"title"`
	Hash    string        `json:"hash"`
	Magnet  string        `json:"magnetUrl"`
	Size    core.ByteSize `json:"bytes"`
	Seeders uint          `json:"seeders"`
}

func scrapeKnaben(ctx context.Context, _ *sync.WaitGroup, name string, torrents chan Torrent) {
	req := knabenRequest{
		SearchType:           "100%",
		SearchField:          "title",
		Query:                name,
		OrderBy:              "seeders",
		OrderDirection:       "desc",
		Categories:           nil,
		From:                 0,
		Size:                 300,
		HideUnsafe:           true,
		HideXXX:              false,
		SecondsSinceLastSeen: 0,
	}

	var buf bytes.Buffer

	for {
		_ = json.NewEncoder(&buf).Encode(&req)

		res, err := core.DoReqJsonCtx[struct {
			Hits []knabenTorrent
		}](ctx, "POST", "https://api.knaben.org/v1", "application/json", &buf)

		if err != nil {
			break
		}

		for _, torrent := range res.Hits {
			if torrent.Seeders != 0 {
				t := Torrent{
					Name:    torrent.Title,
					Hash:    torrent.Hash,
					Magnet:  torrent.Magnet,
					Size:    torrent.Size,
					Seeders: torrent.Seeders,
				}
				parseSeasonEpisode(&t)

				torrents <- t
			}
		}

		if len(res.Hits) < 300 {
			break
		}

		req.From += 300
	}
}
