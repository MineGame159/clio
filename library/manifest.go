package library

import (
	_ "embed"
	"net/http"
)

//go:embed manifest.json
var manifest []byte

func (a *addon) handleManifest(res http.ResponseWriter, _ *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	_, _ = res.Write(manifest)
}
