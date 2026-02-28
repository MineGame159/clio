package library

import (
	"clio/rd"
	"net/http"
)

func (a *addon) handlePlay(res http.ResponseWriter, req *http.Request) {
	link := "https://real-debrid.com/d/" + req.PathValue("id")

	download, err := rd.GetDownloadLink(a.token, link)
	if err != nil {
		writeError(res, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(res, req, download, http.StatusFound)
}
