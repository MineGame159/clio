package rd

import (
	"fmt"
	"net/url"
)

const base = "https://api.real-debrid.com/rest/1.0"

func GetDownloads(token string, page uint) ([]Download, error) {
	url_ := fmt.Sprintf("%s/downloads?limit=100&page=%d", base, page)
	return get[[]Download](token, url_)
}

func GetAllDownloads(token string) ([]Download, error) {
	var downloads []Download

	page := uint(1)

	for {
		pageDownloads, err := GetDownloads(token, page)
		if err != nil {
			return nil, err
		}

		downloads = append(downloads, pageDownloads...)

		if len(pageDownloads) < 100 {
			break
		}
	}

	return downloads, nil
}

func Unrestrict(token string, link string) (Download, error) {
	url_ := fmt.Sprintf("%s/unrestrict/link", base)
	values := url.Values{"link": {link}}
	return post[Download](token, url_, values)
}

func GetTorrents(token string, page uint) ([]Torrent, error) {
	url_ := fmt.Sprintf("%s/torrents?limit=100&page=%d", base, page)
	return get[[]Torrent](token, url_)
}

func GetAllTorrents(token string) ([]Torrent, error) {
	var torrents []Torrent

	page := uint(1)

	for {
		pageTorrents, err := GetTorrents(token, page)
		if err != nil {
			return nil, err
		}

		torrents = append(torrents, pageTorrents...)

		if len(pageTorrents) < 100 {
			break
		}
	}

	return torrents, nil
}

func GetTorrentFiles(token string, id string) ([]File, error) {
	url_ := fmt.Sprintf("%s/torrents/info/%s", base, id)

	body, err := get[struct {
		Files []File
		Links []string
	}](token, url_)
	if err != nil {
		return nil, err
	}

	linkIndex := 0

	for fileIndex := range len(body.Files) {
		file := &body.Files[fileIndex]

		if file.Selected != 0 {
			file.Link = body.Links[linkIndex]
			linkIndex++
		}
	}

	return body.Files, nil
}
