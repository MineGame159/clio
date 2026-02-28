package indexers

import (
	"clio/core"
	"fmt"
	"net/url"
	"strconv"
)

var trackers = []string{
	"udp://open.stealth.si:80/announce",
	"udp://tracker-udp.gbitt.info:80/announce",
	"udp://tracker.tryhackx.org:6969/announce",
	"udp://retracker.lanta.me:2710/announce",
	"https://shahidrazi.online:443/announce",
	"udp://tracker.t-1.org:6969/announce",
	"https://tr.nyacat.pw:443/announce",
	"udp://extracker.dahrkael.net:6969/announce",
	"udp://tracker.plx.im:6969/announce",
	"udp://tracker.1h.is:1337/announce",
	"https://tracker.iochimari.moe:443/announce",
	"udp://tracker.flatuslifir.is:6969/announce",
	"udp://tracker.opentorrent.top:6969/announce",
	"udp://ns575949.ip-51-222-82.net:6969/announce",
	"udp://tracker.bluefrog.pw:2710/announce",
	"udp://tracker.ixuexi.click:6969/announce",
	"udp://tracker.corpscorp.online:80/announce",
	"https://tracker.qingwapt.org:443/announce",
	"udp://tracker.riverarmy.xyz:6969/announce",
	"udp://tracker.opentrackr.org:1337/announce",
	"udp://tracker.torrent.eu.org:451/announce",
	"udp://martin-gebhardt.eu:25/announce",
	"udp://tr4ck3r.duckdns.org:6969/announce",
	"https://t.213891.xyz:443/announce",
	"udp://tracker.playground.ru:6969/announce",
	"udp://torrentvpn.club:6990/announce",
	"https://tracker.novy.vip:443/announce",
	"udp://open.demonii.com:1337/announce",
	"udp://evan.im:6969/announce",
	"udp://uabits.today:6990/announce",
	"udp://tracker.gmi.gd:6969/announce",
	"udp://tracker.fnix.net:6969/announce",
	"https://tracker.ghostchu-services.top:443/announce",
	"https://tracker.manager.v6.navy:443/announce",
	"udp://tracker.srv00.com:6969/announce",
	"https://tracker.zhuqiy.com:443/announce",
	"udp://tracker.dler.com:6969/announce",
	"udp://tracker.breizh.pm:6969/announce",
	"udp://tracker.qu.ax:6969/announce",
	"udp://tracker.therarbg.to:6969/announce",
	"udp://utracker.ghostchu-services.top:6969/announce",
	"udp://udp.tracker.projectk.org:23333/announce",
	"udp://tracker.tvunderground.org.ru:3218/announce",
	"udp://tracker.iperson.xyz:6969/announce",
	"udp://tracker.sharebro.in:6969/announce",
	"udp://wepzone.net:6969/announce",
	"udp://tracker.wepzone.net:6969/announce",
	"udp://tracker.theoks.net:6969/announce",
	"udp://tracker.filemail.com:6969/announce",
}

func parseSeasonEpisode(torrent *Torrent) {
	// Season
	seasonMatches := core.SeasonRegex.FindStringSubmatch(torrent.Name)

	if len(seasonMatches) == 2 {
		torrent.Season, _ = strconv.Atoi(seasonMatches[1])
	} else {
		torrent.Season = -1
	}

	// Episode
	episodeMatches := core.EpisodeRegex.FindStringSubmatch(torrent.Name)

	if len(episodeMatches) == 2 {
		torrent.Episode, _ = strconv.Atoi(episodeMatches[1])
	} else {
		torrent.Episode = -1
	}
}

func getMagnetLink(name, hash string) string {
	xtParam := "xt=urn:btih:" + hash

	values := url.Values{}
	values.Set("dn", name)

	for _, tracker := range trackers {
		values.Add("tr", tracker)
	}

	return fmt.Sprintf("magnet:?%s&%s", xtParam, values.Encode())
}
