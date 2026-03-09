package scraper

import (
	"clio/core"
	"clio/rd"
	"net/http"
)

type Addon struct {
	rd      *rd.Client
	baseUrl string
}

func Start(token string) (string, error) {
	// Initialize addon
	a := &Addon{
		rd: rd.NewClient(token),
	}

	// Routes
	mux := http.NewServeMux()

	mux.HandleFunc("GET /manifest.json", a.handleManifest)
	mux.HandleFunc("GET /stream/{kind}/{id}", a.handleStream)
	mux.HandleFunc("GET /check/{magnet}", a.handleCheck)
	mux.HandleFunc("GET /play/{magnet}/{season}/{episode}", a.handlePlay)

	// Listen
	var err error
	a.baseUrl, err = core.ListenAndServe(mux)
	if err != nil {
		return "", err
	}

	return a.baseUrl + "/manifest.json", nil
}
