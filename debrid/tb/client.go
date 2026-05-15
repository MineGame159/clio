package tb

import (
	"clio/core"
	"clio/debrid"
	"context"
	"strconv"

	"golang.org/x/sync/semaphore"
)

type Client struct {
	token string
	sem   *semaphore.Weighted
}

func NewClient(token string) *Client {
	return &Client{
		token: token,
		sem:   semaphore.NewWeighted(3),
	}
}

// debrid.Client

func (c *Client) GetTorrent(ctx context.Context, id string) (debrid.Torrent, []debrid.File, error) {
	id_, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return debrid.Torrent{}, nil, err
	}

	torrent, err := getTorrent(ctx, c, uint(id_))
	if err != nil {
		return debrid.Torrent{}, nil, err
	}

	dTorrent := getDebridTorrent(torrent)
	dFiles := make([]debrid.File, len(torrent.Files))

	for i, file := range torrent.Files {
		dFiles[i] = getDebridFile(file)
	}

	return dTorrent, dFiles, nil
}

func (c *Client) GetTorrents(ctx context.Context, page uint) ([]debrid.Torrent, error) {
	torrents, err := getTorrentList(ctx, c, page)
	if err != nil {
		return nil, err
	}

	dTorrents := make([]debrid.Torrent, len(torrents))

	for i, torrent := range torrents {
		dTorrents[i] = getDebridTorrent(torrent)
	}

	return dTorrents, nil
}

func (c *Client) CheckMagnet(ctx context.Context, magnet string) (bool, []debrid.File, error) {
	hash, err := core.ParseMagnet(magnet)
	if err != nil {
		return false, nil, err
	}

	torrents, err := getTorrentCachedAvailability(ctx, c, hash)
	if err != nil {
		return false, nil, err
	}

	if len(torrents) == 0 {
		return false, nil, nil
	}

	dFiles := make([]debrid.File, len(torrents[0].Files))

	for i, file := range torrents[0].Files {
		dFiles[i] = getDebridFile(file)
	}

	return true, dFiles, nil
}

func (c *Client) AddMagnet(ctx context.Context, magnet string) (string, error) {
	id, err := createTorrent(ctx, c, magnet)
	if err != nil {
		return "", err
	}

	return strconv.FormatUint(uint64(id), 10), nil
}

func (c *Client) DeleteTorrent(ctx context.Context, id string) error {
	id_, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return err
	}

	return controlTorrent(ctx, c, uint(id_), opDelete)
}

func (c *Client) GetDownload(ctx context.Context, id string, fileId uint) (string, error) {
	id_, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return "", err
	}

	return requestDownloadLink(ctx, c, uint(id_), fileId)
}

// Type conversion

func getDebridTorrent(torrent Torrent) debrid.Torrent {
	return debrid.Torrent{
		Id:   strconv.FormatUint(uint64(torrent.Id), 10),
		Name: torrent.Name,
		Hash: torrent.Hash,
	}
}

func getDebridFile(file File) debrid.File {
	return debrid.File{
		Id:       file.Id,
		Path:     file.Path,
		Size:     file.Size,
		Selected: true,
	}
}
