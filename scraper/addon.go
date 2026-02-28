package scraper

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
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

func writeJson(res http.ResponseWriter, data any) {
	res.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(res).Encode(data)
}

func writeError(res http.ResponseWriter, msg string, code int) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(code)

	_, _ = fmt.Fprintf(res, "{\"error\":\"%s\"}", msg)
}
