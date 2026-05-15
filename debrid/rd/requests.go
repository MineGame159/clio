package rd

import (
	"clio/debrid"
	"context"
	"fmt"
	"net/http"
	"net/url"
)

func getTorrent(ctx context.Context, c debrid.HttpClient, id string) (Torrent, []File, error) {
	url_ := fmt.Sprintf("%s/torrents/info/%s", base, id)

	body, err := debrid.Get[struct {
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

func getTorrents(ctx context.Context, c debrid.HttpClient, page uint) ([]Torrent, error) {
	url_ := fmt.Sprintf("%s/torrents?limit=500&page=%d", base, page)
	return debrid.Get[[]Torrent](c, ctx, url_)
}

func addMagnet(ctx context.Context, c debrid.HttpClient, magnet string) (string, error) {
	url_ := fmt.Sprintf("%s/torrents/addMagnet", base)
	values := url.Values{"magnet": {magnet}}

	body, err := debrid.Post[struct {
		Id string
	}](c, ctx, url_, values)

	if err != nil {
		return "", err
	}

	return body.Id, nil
}

func selectAllFiles(ctx context.Context, c debrid.HttpClient, id string) error {
	url_ := fmt.Sprintf("%s/torrents/selectFiles/%s", base, id)

	values := url.Values{"files": {"all"}}
	_, err := debrid.Post[struct{}](c, ctx, url_, values)

	return err
}

func deleteTorrent(ctx context.Context, c debrid.HttpClient, id string) error {
	url_ := fmt.Sprintf("%s/torrents/delete/%s", base, id)

	req, err := http.NewRequest("DELETE", url_, nil)
	if err != nil {
		return err
	}

	_, err = debrid.Do[struct{}](c, ctx, req)
	return err
}

func getDownloads(ctx context.Context, c debrid.HttpClient, page uint) ([]Download, error) {
	url_ := fmt.Sprintf("%s/downloads?limit=500&page=%d", base, page)
	return debrid.Get[[]Download](c, ctx, url_)
}

func unrestrict(ctx context.Context, c debrid.HttpClient, link string) (Download, error) {
	url_ := fmt.Sprintf("%s/unrestrict/link", base)
	values := url.Values{"link": {link}}
	return debrid.Post[Download](c, ctx, url_, values)
}

// Utils

func getAllDownloads(ctx context.Context, c debrid.HttpClient) ([]Download, error) {
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
