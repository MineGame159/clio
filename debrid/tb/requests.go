package tb

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

const base = "https://api.torbox.app/v1/api"

type torrentOperation = string

const (
	opDelete torrentOperation = "delete"
)

func createTorrent(ctx context.Context, c *Client, magnet string) (uint, error) {
	url_ := fmt.Sprintf("%s/torrents/createtorrent", base)

	values := url.Values{"magnet": {magnet}}

	res, err := post[struct {
		Id uint `json:"torrent_id"`
	}](c, ctx, url_, values)
	if err != nil {
		return 0, err
	}

	return res.Id, nil
}

func controlTorrent(ctx context.Context, c *Client, torrentId uint, operation torrentOperation) error {
	url_ := fmt.Sprintf("%s/torrents/controltorrent", base)

	values := url.Values{
		"torrent_id": {strconv.FormatUint(uint64(torrentId), 10)},
		"operation":  {operation},
	}

	_, err := post[struct{}](c, ctx, url_, values)
	return err
}

func requestDownloadLink(ctx context.Context, c *Client, torrentId uint, fileId uint) (string, error) {
	url_ := fmt.Sprintf("%s/torrents/requestdl?token=%s&torrent_id=%d&file_id=%d", base, c.token, torrentId, fileId)

	return get[string](c, ctx, url_)
}

func getTorrentList(ctx context.Context, c *Client, page uint) ([]Torrent, error) {
	url_ := fmt.Sprintf("%s/torrents/mylist?limit=500&offset=%d&bypass_cache=true", base, (page-1)*500)

	return get[[]Torrent](c, ctx, url_)
}

func getTorrent(ctx context.Context, c *Client, torrentId uint) (Torrent, error) {
	url_ := fmt.Sprintf("%s/torrents/mylist?id=%d&bypass_cache=true", base, torrentId)

	return get[Torrent](c, ctx, url_)
}

func getTorrentCachedAvailability(ctx context.Context, c *Client, hash string) ([]Torrent, error) {
	url_ := fmt.Sprintf("%s/torrents/checkcached?hash=%s&list_files=true", base, hash)

	map_, err := get[map[string]Torrent](c, ctx, url_)
	if err != nil {
		return nil, err
	}

	torrents := make([]Torrent, 0, len(map_))

	for _, torrent := range map_ {
		torrents = append(torrents, Torrent{
			Id:    0,
			Name:  torrent.Name,
			Hash:  torrent.Hash,
			Files: torrent.Files,
		})
	}

	return torrents, nil
}
