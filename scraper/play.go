package scraper

import (
	"clio/core"
	"clio/rd"
	"encoding/base64"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func (a *Addon) handlePlay(res http.ResponseWriter, req *http.Request) {
	// Read magnet link
	magnet, hash, err := readMagnet(req)
	if err != nil {
		writeError(res, err.Error(), http.StatusBadRequest)
		return
	}

	// Read season
	season, err := strconv.Atoi(req.PathValue("season"))
	if err != nil {
		writeError(res, err.Error(), http.StatusBadRequest)
		return
	}

	// Episode
	episode, err := strconv.Atoi(req.PathValue("episode"))
	if err != nil {
		writeError(res, err.Error(), http.StatusBadRequest)
		return
	}

	// Find existing torrent with the same hash
	torrents, err := rd.GetAllTorrents(a.token)
	if err != nil {
		writeError(res, err.Error(), http.StatusBadRequest)
		return
	}

	for _, torrent := range torrents {
		if torrent.Hash == hash {
			download, err := a.getDownloadFromTorrent(torrent.Id, season, episode)
			if err != nil {
				writeError(res, err.Error(), http.StatusBadRequest)
				return
			}

			if download != "" {
				http.Redirect(res, req, download, http.StatusFound)
				return
			}

			break
		}
	}

	// Add magnet to library
	id, err := rd.AddMagnet(a.token, magnet)
	if err != nil {
		writeError(res, err.Error(), http.StatusBadRequest)
		return
	}

	// Select files
	if err := a.selectFilesFromTorrent(id, season == -1 && episode == -1); err != nil {
		writeError(res, err.Error(), http.StatusBadRequest)
		return
	}

	time.Sleep(time.Millisecond * 250)

	// Redirect
	download, err := a.getDownloadFromTorrent(id, season, episode)
	if err != nil {
		writeError(res, err.Error(), http.StatusBadRequest)
		return
	}

	if download != "" {
		http.Redirect(res, req, download, http.StatusFound)
	}
}

func (a *Addon) selectFilesFromTorrent(id string, movie bool) error {
	_, files, err := rd.GetTorrent(a.token, id)
	if err != nil {
		return err
	}

	var fileIds []uint
	var biggestFile rd.File

	for _, file := range files {
		if !core.IsVideoFile(file.Path) {
			continue
		}

		if movie {
			if file.Size > biggestFile.Size {
				biggestFile = file
			}
		} else {
			name := file.Path
			if index := strings.LastIndexByte(name, '/'); index != -1 {
				name = name[index+1:]
			}

			info := core.ParseTorrentName(name)

			if info.Season != -1 && info.Episode != -1 {
				fileIds = append(fileIds, file.Id)
			}
		}
	}

	if movie && biggestFile.Path != "" {
		fileIds = append(fileIds, biggestFile.Id)
	}

	return rd.SelectFiles(a.token, id, fileIds)
}

func (a *Addon) getDownloadFromTorrent(id string, season, episode int) (string, error) {
	_, files, err := rd.GetTorrent(a.token, id)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		if fileMatches(file, season, episode) {
			download, err := rd.GetDownloadLink(a.token, file.Link)
			if err != nil {
				return "", err
			}

			return download, nil
		}
	}

	return "", nil
}

func fileMatches(file rd.File, season, episode int) bool {
	if file.Selected == 0 {
		return false
	}

	if season == -1 && episode == -1 {
		return true
	}

	name := file.Path
	if index := strings.LastIndexByte(name, '/'); index != -1 {
		name = name[index+1:]
	}

	info := core.ParseTorrentName(name)
	return info.Season == season && info.Episode == episode
}

func readMagnet(req *http.Request) (string, string, error) {
	reader := base64.NewDecoder(base64.URLEncoding, strings.NewReader(req.PathValue("magnet")))
	magnetBytes, err := io.ReadAll(reader)
	if err != nil {
		return "", "", err
	}

	magnet := string(magnetBytes)
	magnetQuery, _ := strings.CutPrefix(magnet, "magnet:?")

	values, err := url.ParseQuery(magnetQuery)
	if err != nil {
		return "", "", err
	}

	hash, _ := strings.CutPrefix(values.Get("xt"), "urn:btih:")
	return magnet, hash, nil
}
