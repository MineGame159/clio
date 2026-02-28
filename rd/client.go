package rd

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
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

func GetDownloadLink(token, link string) (string, error) {
	// Try to find an existing download link
	downloads, err := GetAllDownloads(token)
	if err != nil {
		return "", err
	}

	for _, download := range downloads {
		if download.Link == link {
			return download.Download, nil
		}
	}

	// Generate a download link
	download, err := Unrestrict(token, link)
	if err != nil {
		return "", err
	}

	return download.Download, nil
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

func GetTorrent(token string, id string) (Torrent, []File, error) {
	url_ := fmt.Sprintf("%s/torrents/info/%s", base, id)

	body, err := get[struct {
		Torrent
		Files []File
		Links []string
	}](token, url_)
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

func AddMagnet(token string, magnet string) (string, error) {
	url_ := fmt.Sprintf("%s/torrents/addMagnet", base)
	values := url.Values{"magnet": {magnet}}

	body, err := post[struct {
		Id string
	}](token, url_, values)
	if err != nil {
		return "", err
	}

	return body.Id, nil
}

func SelectFiles(token string, id string, fileIds []uint) error {
	url_ := fmt.Sprintf("%s/torrents/selectFiles/%s", base, id)

	var strFileIds strings.Builder

	for i, fileId := range fileIds {
		if i > 0 {
			strFileIds.WriteRune(',')
		}

		_, _ = fmt.Fprintf(&strFileIds, "%d", fileId)
	}

	values := url.Values{"files": {strFileIds.String()}}

	_, err := post[struct{}](token, url_, values)
	return err
}

func DeleteTorrent(token string, id string) error {
	url_ := fmt.Sprintf("%s/torrents/delete/%s", base, id)

	req, err := http.NewRequest("DELETE", url_, nil)
	if err != nil {
		return err
	}

	_, err = doRequest[struct{}](token, req)
	return err
}
