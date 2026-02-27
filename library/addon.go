package library

import (
	"clio/core"
	"clio/rd"
	"clio/stremio"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"
)

type addon struct {
	token   string
	baseUrl string

	media map[string]mediaInfo
}

type mediaInfo struct {
	id string

	name   string
	kind   stremio.MediaKind
	poster string

	torrentIds []string
}

func Start(token string) (string, error) {
	// Initialize addon
	a := &addon{
		token: token,
		media: make(map[string]mediaInfo),
	}

	if err := a.fetchMedia(); err != nil {
		return "", err
	}

	// Routes
	mux := http.NewServeMux()

	mux.HandleFunc("GET /manifest.json", a.handleManifest)
	mux.HandleFunc("GET /catalog/{kind}/{id}", a.handleCatalog)
	mux.HandleFunc("GET /stream/{kind}/{id}", a.handleStream)
	mux.HandleFunc("GET /play/{id}", a.handlePlay)

	// Listen
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return "", err
	}

	//goland:noinspection GoUnhandledErrorResult
	go http.Serve(listener, mux)

	a.baseUrl = "http://localhost:" + strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)
	address := a.baseUrl + "/manifest.json"

	return address, nil
}

func (a *addon) fetchMedia() error {
	torrents, err := rd.GetAllTorrents(a.token)
	if err != nil {
		return err
	}

	// Group torrents by the cleaned up name
	mediaNameTorrents := make(map[string][]rd.Torrent)

	for _, torrent := range torrents {
		info := core.ParseTorrentName(torrent.Filename)

		var torrents []rd.Torrent

		if existing, ok := mediaNameTorrents[info.Name]; ok {
			torrents = append(existing, torrent)
		} else {
			torrents = []rd.Torrent{torrent}
		}

		mediaNameTorrents[info.Name] = torrents
	}

	// Fetch IMDB IDs for grouped media
	wg := sync.WaitGroup{}
	media := make(chan mediaInfo)

	for _, torrents := range mediaNameTorrents {
		wg.Go(func() {
			info := core.ParseTorrentName(torrents[0].Filename)

			kind := stremio.Movie
			if info.Season != -1 || info.Episode != -1 {
				kind = stremio.Series
			}

			url := fmt.Sprintf("https://v3-cinemeta.strem.io/catalog/%s/top/search=%s.json", kind, info.Name)

			body, err := core.GetJson[struct {
				Metas []stremio.SearchResult
			}](url)

			if err == nil && len(body.Metas) > 0 {
				meta := body.Metas[0]

				for _, torrent := range torrents {
					media <- mediaInfo{
						id:         meta.Id,
						name:       meta.Name,
						kind:       kind,
						poster:     meta.Poster,
						torrentIds: []string{torrent.Id},
					}
				}
			}
		})
	}

	go func() {
		wg.Wait()
		close(media)
	}()

	// Add torrents into addon media map
	for info := range media {
		if existing, ok := a.media[info.id]; ok {
			info.torrentIds = append(existing.torrentIds, info.torrentIds...)
		}

		a.media[info.id] = info
	}

	return nil
}

func writeJson(res http.ResponseWriter, data any) {
	res.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(res).Encode(data)
}

func writeError(res http.ResponseWriter, msg string, code int) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(code)

	_, _ = fmt.Fprintf(res, "{\"error\":\"%s\"}", msg)
}
