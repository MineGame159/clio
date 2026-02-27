package library

import (
	"clio/rd"
	"net/http"
)

func (a *addon) handlePlay(res http.ResponseWriter, req *http.Request) {
	link := "https://real-debrid.com/d/" + req.PathValue("id")

	// Try to find an existing download link
	downloads, err := rd.GetAllDownloads(a.token)
	if err != nil {
		writeError(res, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, download := range downloads {
		if download.Link == link {
			http.Redirect(res, req, download.Download, http.StatusFound)
			return
		}
	}

	// Generate a download link
	download, err := rd.Unrestrict(a.token, link)
	if err != nil {
		writeError(res, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(res, req, download.Download, http.StatusFound)
}
