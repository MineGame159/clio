package tb

import "clio/core"

type Torrent struct {
	Id    uint   `json:"id"`
	Name  string `json:"name"`
	Hash  string `json:"hash"`
	Files []File `json:"files"`
}

type File struct {
	Id   uint          `json:"id"`
	Path string        `json:"name"`
	Size core.ByteSize `json:"size"`
}
