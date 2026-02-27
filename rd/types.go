package rd

import "clio/core"

type TorrentStatus string

const (
	MagnetStatus          TorrentStatus = "magnet_error"
	MagnetConversion      TorrentStatus = "magnet_conversion"
	WaitingFilesSelection TorrentStatus = "waiting_files_selection"
	Queued                TorrentStatus = "queued"
	Downloading           TorrentStatus = "downloading"
	Downloaded            TorrentStatus = "downloaded"
	Error                 TorrentStatus = "error"
	Virus                 TorrentStatus = "virus"
	Compressing           TorrentStatus = "compressing"
	Uploading             TorrentStatus = "uploading"
	Dead                  TorrentStatus = "dead"
)

type Torrent struct {
	Id       string        `json:"id"`
	Filename string        `json:"filename"`
	Hash     string        `json:"hash"`
	Size     core.ByteSize `json:"bytes"`
	Host     string        `json:"host"`
	Split    uint          `json:"split"`
	Progress uint          `json:"progress"`
	Status   TorrentStatus `json:"status"`
	Added    string        `json:"added"`
	Ended    string        `json:"ended"`
	Speed    uint          `json:"speed"`
	Seeders  uint          `json:"seeders"`
}

type File struct {
	Id       uint          `json:"id"`
	Path     string        `json:"path"`
	Size     core.ByteSize `json:"bytes"`
	Selected uint          `json:"selected"`
	Link     string        `json:"-"`
}

type Download struct {
	Id        string        `json:"id"`
	Filename  string        `json:"filename"`
	MineType  string        `json:"minetype"`
	Size      core.ByteSize `json:"filesize"`
	Link      string        `json:"link"`
	Host      string        `json:"host"`
	Chunks    uint          `json:"chunks"`
	Download  string        `json:"download"`
	Generated string        `json:"generated"`
}
