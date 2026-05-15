package rd

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

const base = "https://api.real-debrid.com/rest/1.0"

func getTorrent(ctx context.Context, c *Client, id string) (Torrent, []File, error) {
	url_ := fmt.Sprintf("%s/torrents/info/%s", base, id)

	body, err := get[struct {
		Torrent
		Files []File
		Links []string
	}](c, ctx, url_)
	if err != nil {
		return Torrent{}, nil, err
	}

	linkIndex := 0

	for fileIndex := range len(body.Files) {
		file := &body.Files[fileIndex]

		if file.Selected != 0 && linkIndex < len(body.Links) {
			file.Link = body.Links[linkIndex]
			linkIndex++
		}
	}

	return body.Torrent, body.Files, nil
}

func getTorrents(ctx context.Context, c *Client, page uint) ([]Torrent, error) {
	url_ := fmt.Sprintf("%s/torrents?limit=500&page=%d", base, page)
	return get[[]Torrent](c, ctx, url_)
}

func addMagnet(ctx context.Context, c *Client, magnet string) (string, error) {
	url_ := fmt.Sprintf("%s/torrents/addMagnet", base)
	values := url.Values{"magnet": {magnet}}

	body, err := post[struct {
		Id string
	}](c, ctx, url_, values)

	if err != nil {
		return "", err
	}

	return body.Id, nil
}

func selectAllFiles(ctx context.Context, c *Client, id string) error {
	url_ := fmt.Sprintf("%s/torrents/selectFiles/%s", base, id)

	values := url.Values{"files": {"all"}}
	_, err := post[struct{}](c, ctx, url_, values)

	return err
}

func deleteTorrent(ctx context.Context, c *Client, id string) error {
	url_ := fmt.Sprintf("%s/torrents/delete/%s", base, id)

	req, err := http.NewRequest("DELETE", url_, nil)
	if err != nil {
		return err
	}

	_, err = doRequest[struct{}](c, ctx, req)
	return err
}

func getDownloads(ctx context.Context, c *Client, page uint) ([]Download, error) {
	url_ := fmt.Sprintf("%s/downloads?limit=500&page=%d", base, page)
	return get[[]Download](c, ctx, url_)
}

func unrestrict(ctx context.Context, c *Client, link string) (Download, error) {
	url_ := fmt.Sprintf("%s/unrestrict/link", base)
	values := url.Values{"link": {link}}
	return post[Download](c, ctx, url_, values)
}

// Utils

func getAllDownloads(ctx context.Context, c *Client) ([]Download, error) {
	var downloads []Download

	page := uint(1)

	for {
		pageDownloads, err := getDownloads(ctx, c, page)
		if err != nil {
			return nil, err
		}

		downloads = append(downloads, pageDownloads...)

		if len(pageDownloads) < 500 {
			break
		}

		page++
	}

	return downloads, nil
}
