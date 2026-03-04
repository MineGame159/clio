package library

import (
	"clio/core"
	"clio/stremio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"sync"
)

var noResult = errors.New("no result")

type imdbCache struct {
	Movie  map[string]stremio.SearchResult `json:"movie"`
	Series map[string]stremio.SearchResult `json:"series"`
	Anime  map[string]stremio.SearchResult `json:"anime"`
	Other  map[string]stremio.SearchResult `json:"other"`

	mutex sync.Mutex
}

func newImdbCache() *imdbCache {
	cache := &imdbCache{
		Movie:  make(map[string]stremio.SearchResult),
		Series: make(map[string]stremio.SearchResult),
		Anime:  make(map[string]stremio.SearchResult),
		Other:  make(map[string]stremio.SearchResult),
	}

	if err := cache.Load(); err != nil {
		panic(err.Error())
	}

	return cache
}

func (i *imdbCache) Get(kind stremio.MediaKind, name string) (stremio.SearchResult, error) {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	var cache map[string]stremio.SearchResult

	switch kind {
	case stremio.Movie:
		cache = i.Movie
	case stremio.Series:
		cache = i.Series
	case stremio.Anime:
		cache = i.Anime
	case stremio.Other:
		cache = i.Other
	default:
		panic("imdbIdCache.Get() - Invalid media kind")
	}

	if id, ok := cache[name]; ok {
		return id, nil
	}

	url := fmt.Sprintf("https://v3-cinemeta.strem.io/catalog/%s/top/search=%s.json", kind, name)

	body, err := core.GetJson[struct {
		Metas []stremio.SearchResult
	}](url)

	if err != nil {
		return stremio.SearchResult{}, err
	}
	if len(body.Metas) == 0 {
		return stremio.SearchResult{}, noResult
	}

	cache[name] = body.Metas[0]

	return body.Metas[0], nil
}

func (i *imdbCache) Load() error {
	base, err := os.UserCacheDir()
	if err != nil {
		return err
	}

	filename := path.Join(base, "clio", "imdb_cache.json")

	file, err := os.Open(filename)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}

	//goland:noinspection GoUnhandledErrorResult
	defer file.Close()

	if err := json.NewDecoder(file).Decode(i); err != nil {
		return err
	}

	return nil
}

func (i *imdbCache) Save() error {
	base, err := os.UserCacheDir()
	if err != nil {
		return err
	}

	base = path.Join(base, "clio")

	if err := os.MkdirAll(base, 0750); err != nil {
		return nil
	}

	filename := path.Join(base, "imdb_cache.json")

	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	//goland:noinspection GoUnhandledErrorResult
	defer file.Close()

	if err := json.NewEncoder(file).Encode(i); err != nil {
		return err
	}

	return nil
}
