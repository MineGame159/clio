package library

import (
	"clio/core"
	"net/http"
)

func (a *Addon) handlePlay(res http.ResponseWriter, req *http.Request) {
	link := "https://real-debrid.com/d/" + req.PathValue("id")

	download, err := a.rd.GetDownloadLink(link)
	if err != nil {
		core.WriteError(res, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(res, req, download, http.StatusFound)
}
