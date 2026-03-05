package rd

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/sync/semaphore"
)

const base = "https://api.real-debrid.com/rest/1.0"

type Client struct {
	token string
	sem   *semaphore.Weighted
}

func NewClient(token string) *Client {
	return &Client{
		token: token,
		sem:   semaphore.NewWeighted(1),
	}
}

func (c *Client) GetDownloads(page uint) ([]Download, error) {
	url_ := fmt.Sprintf("%s/downloads?limit=100&page=%d", base, page)
	return get[[]Download](c, url_)
}

func (c *Client) GetAllDownloads() ([]Download, error) {
	var downloads []Download

	page := uint(1)

	for {
		pageDownloads, err := c.GetDownloads(page)
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

func (c *Client) Unrestrict(link string) (Download, error) {
	url_ := fmt.Sprintf("%s/unrestrict/link", base)
	values := url.Values{"link": {link}}
	return post[Download](c, url_, values)
}

func (c *Client) GetDownloadLink(link string) (string, error) {
	// Try to find an existing download link
	downloads, err := c.GetAllDownloads()
	if err != nil {
		return "", err
	}

	for _, download := range downloads {
		if download.Link == link {
			return download.Download, nil
		}
	}

	// Generate a download link
	download, err := c.Unrestrict(link)
	if err != nil {
		return "", err
	}

	return download.Download, nil
}

func (c *Client) GetTorrents(page uint) ([]Torrent, error) {
	url_ := fmt.Sprintf("%s/torrents?limit=100&page=%d", base, page)
	return get[[]Torrent](c, url_)
}

func (c *Client) GetAllTorrents() ([]Torrent, error) {
	var torrents []Torrent

	page := uint(1)

	for {
		pageTorrents, err := c.GetTorrents(page)
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

func (c *Client) GetTorrent(id string) (Torrent, []File, error) {
	url_ := fmt.Sprintf("%s/torrents/info/%s", base, id)

	body, err := get[struct {
		Torrent
		Files []File
		Links []string
	}](c, url_)
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

func (c *Client) AddMagnet(magnet string) (string, error) {
	url_ := fmt.Sprintf("%s/torrents/addMagnet", base)
	values := url.Values{"magnet": {magnet}}

	body, err := post[struct {
		Id string
	}](c, url_, values)
	if err != nil {
		return "", err
	}

	return body.Id, nil
}

func (c *Client) SelectFiles(id string, fileIds []uint) error {
	url_ := fmt.Sprintf("%s/torrents/selectFiles/%s", base, id)

	var strFileIds strings.Builder

	for i, fileId := range fileIds {
		if i > 0 {
			strFileIds.WriteRune(',')
		}

		_, _ = fmt.Fprintf(&strFileIds, "%d", fileId)
	}

	values := url.Values{"files": {strFileIds.String()}}

	_, err := post[struct{}](c, url_, values)
	return err
}

func (c *Client) DeleteTorrent(id string) error {
	url_ := fmt.Sprintf("%s/torrents/delete/%s", base, id)

	req, err := http.NewRequest("DELETE", url_, nil)
	if err != nil {
		return err
	}

	_, err = doRequest[struct{}](c, req)
	return err
}
