package debrid

import (
	"clio/core"
	"context"
)

type Torrent struct {
	Id   string
	Name string
	Hash string
}

type File struct {
	Id       uint
	Path     string
	Size     core.ByteSize
	Selected bool
}

type Client interface {
	GetTorrent(ctx context.Context, id string) (Torrent, []File, error)
	GetTorrents(ctx context.Context, page uint) ([]Torrent, error)

	CheckMagnet(ctx context.Context, magnet string) (bool, []File, error)
	AddMagnet(ctx context.Context, magnet string) (string, error)

	DeleteTorrent(ctx context.Context, id string) error

	GetDownload(ctx context.Context, id string, fileId uint) (string, error)
}

func GetAllTorrents(ctx context.Context, client Client) ([]Torrent, error) {
	var torrents []Torrent

	page := uint(1)

	for {
		pageTorrents, err := client.GetTorrents(ctx, page)
		if err != nil {
			return nil, err
		}

		torrents = append(torrents, pageTorrents...)

		if len(pageTorrents) < 500 {
			break
		}

		page++
	}

	return torrents, nil
}
