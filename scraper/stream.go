package scraper

import (
	"bytes"
	"clio/scraper/indexers"
	"clio/stremio"
	"cmp"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
)

func (a *Addon) handleStream(res http.ResponseWriter, req *http.Request) {
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

	// Scrape torrents using an IMDB ID
	var torrents []indexers.Torrent
	torrentSet := make(map[string]any)

	for torrent := range indexers.Scrape(req.Context(), req.PathValue("kind"), id) {
		if season != -1 && episode != -1 {
			if torrent.Season != season {
				continue
			}
			if torrent.Episode != -1 && torrent.Episode != episode {
				continue
			}
		}

		if _, ok := torrentSet[torrent.Hash]; !ok {
			torrents = append(torrents, torrent)
			torrentSet[torrent.Hash] = nil
		}
	}

	// Sort torrents
	totalSeeders := uint64(0)

	for _, torrent := range torrents {
		totalSeeders += uint64(torrent.Seeders)
	}

	seederGroup := uint(1)
	if len(torrents) > 0 {
		seederGroup = max(uint(totalSeeders/uint64(len(torrents))), 1)
	}
	halfSeederGroup := seederGroup / 2

	slices.SortFunc(torrents, func(a, b indexers.Torrent) int {
		// 1. prioritise complete seasons
		aComplete := a.Season > 0 && a.Episode == -1
		bComplete := b.Season > 0 && b.Episode == -1

		if aComplete != bComplete {
			if aComplete {
				return -1
			}

			return 1
		}

		// 2. sort by grouped seeder count
		aSeeders := uint(uint64(a.Seeders)+uint64(halfSeederGroup)) / seederGroup
		bSeeders := uint(uint64(b.Seeders)+uint64(halfSeederGroup)) / seederGroup

		if aSeeders != bSeeders {
			return cmp.Compare(bSeeders, aSeeders)
		}

		// 3 sort by size
		return cmp.Compare(b.Size, a.Size)
	})

	// Create streams
	streams := make([]stremio.Stream, len(torrents))

	for i, torrent := range torrents {
		var buf bytes.Buffer

		encoder := base64.NewEncoder(base64.URLEncoding, &buf)
		_, _ = encoder.Write([]byte(torrent.Magnet))
		_ = encoder.Close()

		magnet := buf.String()

		streams[i] = stremio.Stream{
			Name:        "Scraper",
			Title:       "",
			Description: fmt.Sprintf("%s\n👥 %d", torrent.Name, torrent.Seeders),
			Url:         fmt.Sprintf("%s/play/%s/%d/%d", a.baseUrl, magnet, season, episode),
			RedirectUrl: true,
			CheckUrl:    fmt.Sprintf("%s/check/%s/%d/%d", a.baseUrl, magnet, season, episode),
			Hints: stremio.BehaviorHints{
				BingeGroup: "",
				Filename:   torrent.Name,
				VideoSize:  uint64(torrent.Size),
			},
		}
	}

	// Write response
	writeJson(res, struct {
		Streams []stremio.Stream
	}{streams})
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
