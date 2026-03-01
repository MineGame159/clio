package scraper

import (
	"clio/core"
	"net/http"
)

type Addon struct {
	token   string
	baseUrl string
}

func Start(token string) (string, error) {
	// Initialize addon
	a := &Addon{
		token: token,
	}

	// Routes
	mux := http.NewServeMux()

	mux.HandleFunc("GET /manifest.json", a.handleManifest)
	mux.HandleFunc("GET /stream/{kind}/{id}", a.handleStream)
	mux.HandleFunc("GET /check/{magnet}/{season}/{episode}", a.handleCheck)
	mux.HandleFunc("GET /play/{magnet}/{season}/{episode}", a.handlePlay)

	// Listen
	var err error
	a.baseUrl, err = core.ListenAndServe(mux)
	if err != nil {
		return "", err
	}

	return a.baseUrl + "/manifest.json", nil
}
