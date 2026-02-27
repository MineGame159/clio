package library

import (
	"clio/core"
	"clio/rd"
	"clio/stremio"
	"cmp"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"sync"
)

func (a *addon) handleStream(res http.ResponseWriter, req *http.Request) {
	// Parse ID
	id, _ := strings.CutSuffix(req.PathValue("id"), ".json")

	if !strings.HasPrefix(id, "tt") {
		writeError(res, "Invalid ID prefix", http.StatusBadRequest)
		return
	}

	id, season, episode, err := parseId(id)
	if err != nil {
		writeError(res, err.Error(), http.StatusBadRequest)
		return
	}

	// Fetch torrent files matching the id
	wg := sync.WaitGroup{}
	files := make(chan rd.File)

	if info, ok := a.media[id]; ok {
		for _, id := range info.torrentIds {
			wg.Go(func() {
				if tFiles, err := rd.GetTorrentFiles(a.token, id); err == nil {
					for _, file := range tFiles {
						if file.Link != "" {
							files <- file
						}
					}
				}
			})
		}
	}

	go func() {
		wg.Wait()
		close(files)
	}()

	// Get streams
	var streams []stremio.Stream

	for file := range files {
		filename := file.Path
		if index := strings.LastIndexByte(filename, '/'); index != -1 {
			filename = filename[index+1:]
		}

		include := true
		if season != -1 || episode != -1 {
			info := core.ParseTorrentName(filename)
			include = info.Season == season && info.Episode == episode
		}

		if include {
			streams = append(streams, a.getStream(file, filename))
		}
	}

	slices.SortFunc(streams, func(a, b stremio.Stream) int {
		return cmp.Compare(b.Hints.VideoSize, a.Hints.VideoSize)
	})

	// Write response
	writeJson(res, struct {
		Streams []stremio.Stream
	}{streams})
}

func (a *addon) getStream(file rd.File, filename string) stremio.Stream {
	id := file.Link
	if index := strings.LastIndexByte(id, '/'); index != -1 {
		id = id[index+1:]
	}

	return stremio.Stream{
		Name:        "Library",
		Title:       "",
		Description: filename,
		Url:         fmt.Sprintf("%s/play/%s", a.baseUrl, id),
		RedirectUrl: true,
		Hints: stremio.BehaviorHints{
			BingeGroup: "",
			Filename:   filename,
			VideoSize:  uint64(file.Size),
		},
	}
}

func parseId(id string) (string, int, int, error) {
	ids := strings.Split(id, ":")

	season := -1
	episode := -1

	if len(ids) == 3 {
		seasonV, err := strconv.ParseInt(ids[1], 10, 32)
		if err != nil {
			return "", -1, -1, errors.New("invalid season number")
		}

		episodeV, err := strconv.ParseInt(ids[2], 10, 32)
		if err != nil {
			return "", -1, -1, errors.New("invalid episode number")
		}

		season = int(seasonV)
		episode = int(episodeV)
	}

	return ids[0], season, episode, nil
}
