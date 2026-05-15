package library

import (
	"clio/core"
	"fmt"
	"net/http"
)

func (a *Addon) handleDelete(res http.ResponseWriter, req *http.Request) {
	if err := a.client.DeleteTorrent(req.Context(), req.PathValue("id")); err != nil {
		core.WriteError(res, err.Error(), http.StatusBadRequest)
	}

	res.Header().Set("Content-Type", "application/json")
	_, _ = fmt.Fprint(res, "{}")
}
