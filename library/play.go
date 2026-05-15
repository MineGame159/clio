package library

import (
	"clio/core"
	"net/http"
	"strconv"
)

func (a *Addon) handlePlay(res http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")
	fileIdStr := req.PathValue("file_id")

	fileId, err := strconv.ParseUint(fileIdStr, 10, 32)
	if err != nil {
		core.WriteError(res, "invalid file_id", http.StatusBadRequest)
		return
	}

	download, err := a.client.GetDownload(req.Context(), id, uint(fileId))
	if err != nil {
		core.WriteError(res, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(res, req, download, http.StatusFound)
}
