package scraper

import (
	"clio/core"
	"clio/debrid"
	"net/http"
)

type Addon struct {
	client  debrid.Client
	baseUrl string
}

func Start(client debrid.Client) (string, error) {
	// Initialize addon
	a := &Addon{
		client: client,
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
