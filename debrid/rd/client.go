package rd

import (
	"clio/debrid"
	"context"
	"errors"
	"time"

	"golang.org/x/sync/semaphore"
)

type Client struct {
	token string
	sem   *semaphore.Weighted
}

func NewClient(token string) *Client {
	return &Client{
		token: token,
		sem:   semaphore.NewWeighted(2),
	}
}

// debrid.Client

func (c *Client) GetTorrent(ctx context.Context, id string) (debrid.Torrent, []debrid.File, error) {
	torrent, files, err := getTorrent(ctx, c, id)
	if err != nil {
		return debrid.Torrent{}, nil, err
	}

	dTorrent := getDebridTorrent(torrent)
	dFiles := make([]debrid.File, len(files))

	for i, file := range files {
		dFiles[i] = getDebridFile(file)
	}

	return dTorrent, dFiles, nil
}

func (c *Client) GetTorrents(ctx context.Context, page uint) ([]debrid.Torrent, error) {
	torrents, err := getTorrents(ctx, c, page)
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
	// Add magnet to library
	id, err := addMagnet(ctx, c, magnet)
	if err != nil {
		return false, nil, err
	}

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		_ = deleteTorrent(ctx, c, id)
	}()

	// Check status
	cached := false

	var torrent Torrent
	var files []File

	for i := 0; i < 10; i++ {
		time.Sleep(time.Millisecond * 250)

		torrent, files, err = getTorrent(ctx, c, id)
		if err != nil {
			return false, nil, err
		}

		if torrent.Status == Downloaded {
			cached = true
			break
		}

		if torrent.Status == Downloading || torrent.Status.Failed() {
			break
		}
	}

	// Convert files
	dFiles := make([]debrid.File, len(files))

	for i, file := range files {
		dFiles[i] = getDebridFile(file)
	}

	return cached, dFiles, nil
}

func (c *Client) AddMagnet(ctx context.Context, magnet string) (string, error) {
	id, err := addMagnet(ctx, c, magnet)
	if err != nil {
		return "", err
	}

	if err = selectAllFiles(ctx, c, id); err != nil {
		_ = deleteTorrent(ctx, c, id)
		return "", err
	}

	return id, nil
}

func (c *Client) DeleteTorrent(ctx context.Context, id string) error {
	return deleteTorrent(ctx, c, id)
}

func (c *Client) GetDownload(ctx context.Context, id string, fileId uint) (string, error) {
	// Get torrent file link
	_, files, err := getTorrent(ctx, c, id)
	if err != nil {
		return "", err
	}

	link := ""

	for _, file := range files {
		if file.Id == fileId {
			link = file.Link
			break
		}
	}

	if link == "" {
		return "", errors.New("invalid file id")
	}

	// Try to find an existing download link
	downloads, err := getAllDownloads(ctx, c)
	if err != nil {
		return "", err
	}

	for _, download := range downloads {
		if download.Link == link {
			return download.Download, nil
		}
	}

	// Generate a download link
	download, err := unrestrict(ctx, c, link)
	if err != nil {
		return "", err
	}

	return download.Download, nil
}

// Type conversion

func getDebridTorrent(torrent Torrent) debrid.Torrent {
	return debrid.Torrent{
		Id:   torrent.Id,
		Name: torrent.Filename,
		Hash: torrent.Hash,
	}
}

func getDebridFile(file File) debrid.File {
	return debrid.File{
		Id:       file.Id,
		Path:     file.Path,
		Size:     file.Size,
		Selected: file.Selected != 0,
	}
}
