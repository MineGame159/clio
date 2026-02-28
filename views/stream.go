package views

import (
	"clio/core"
	"clio/stremio"
	"clio/ui"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

type Stream struct {
	Name string

	Url         string
	RedirectUrl bool
	CheckUrl    string

	Resolution string
	VideoCodec string
	AudioCodec string
	Seeders    int
	Size       core.ByteSize

	Cache CacheStatus
}

type CacheStatus uint8

const (
	Unknown CacheStatus = iota
	Waiting
	Uncached
	Cached
)

// Parsing

func ParseStream(stream stremio.Stream) Stream {
	var s Stream

	s.Name = parseStreamName(stream)
	s.Url = stream.Url
	s.RedirectUrl = stream.RedirectUrl
	s.CheckUrl = stream.CheckUrl

	names := []string{stream.TitleDescription(), s.Name, stream.Hints.Filename}

	s.Resolution = parseStreamResolution(names)
	s.VideoCodec = parseStreamVideoCodec(names)
	s.AudioCodec = parseStreamAudioCodec(names)
	s.Seeders = parseStreamSeeders(names)

	if stream.Hints.VideoSize != 0 {
		s.Size = core.ByteSize(stream.Hints.VideoSize)
	} else {
		s.Size = parseStreamSize(names)
	}

	return s
}

var filenameRegex = regexp.MustCompile("(?im:(?:📄|📁|Name:)\\s*(?:.*/)?(.*)$)")

var resolutionRegex = regexp.MustCompile("\\D(\\d{3,4}[pP])\\W")
var videoCodecRegex = regexp.MustCompile("(?i:(x264|h264|h\\.264|AVC|x265|h265|h\\.265|HEVC|XviD|DivX|AV1|VP9|MPEG2|MPEG-2|VC-1|VC1))")
var audioCodecRegex = regexp.MustCompile("(?i:(AAC|AAC2\\.0|AAC5\\.1|AC3|AC-3|DD5\\.1|DD2\\.0|EAC3|E-AC-3|DDP|DD\\+|TrueHD|DTS|DTS-HD|DTS-MA|DTSHD|FLAC|MP3|Opus|Vorbis))")
var seedersRegex = regexp.MustCompile("(?i:👥|👤|Seeders:)\\s*(\\d+)")
var sizeRegex = regexp.MustCompile("\\D(\\d+(?:\\.\\d+)? [a-zA-Z]{2})\\W")

func parseStreamName(stream stremio.Stream) string {
	filename := stream.Hints.Filename

	if filename == "" {
		if matches := filenameRegex.FindStringSubmatch(stream.TitleDescription()); len(matches) == 2 {
			filename = matches[1]
		}
	}

	index := strings.LastIndexByte(filename, '.')
	if index != -1 {
		return filename[:index]
	}

	return filename
}

func parseStreamResolution(names []string) string {
	for _, name := range names {
		if matches := resolutionRegex.FindStringSubmatch(name); len(matches) == 2 {
			return matches[1]
		}
	}

	return ""
}

func parseStreamVideoCodec(names []string) string {
	for _, name := range names {
		if matches := videoCodecRegex.FindStringSubmatch(name); len(matches) == 2 {
			switch strings.ToLower(matches[1]) {
			case "x264", "h264", "h.264", "avc":
				return "H264"
			case "x265", "h265", "h.265", "hevc":
				return "H265"
			case "xvid", "Divx":
				return "MPEG-4 Part 2"
			case "av1":
				return "AV1"
			case "vp9":
				return "VP9"
			case "mpeg2", "mpeg-2":
				return "MPEG2"
			case "vc1", "vc-1":
				return "VC1"
			}
		}
	}

	return ""
}

func parseStreamAudioCodec(names []string) string {
	for _, name := range names {
		if matches := audioCodecRegex.FindStringSubmatch(name); len(matches) == 2 {
			switch strings.ToLower(matches[1]) {
			case "aac", "aac2.0", "aac5.1":
				return "AAC"
			case "ac3", "ac-3", "dd5.1", "dd2.0":
				return "Dolby"
			case "eac3", "e-ac-3", "ddp", "dd+":
				return "Dolby+"
			case "truehd":
				return "TrueHD"
			case "dts":
				return "DTS"
			case "dts-hd", "dts-ma", "dtshd":
				return "DTS-HD"
			case "flac":
				return "FLAC"
			case "mp3":
				return "MPEG Layer 3"
			case "opus":
				return "Opus"
			case "Vorbis":
				return "Vorbis"
			}
		}
	}

	return ""
}

func parseStreamSeeders(names []string) int {
	for _, name := range names {
		if matches := seedersRegex.FindStringSubmatch(name); len(matches) == 2 {
			seeders, _ := strconv.Atoi(matches[1])
			return seeders
		}
	}

	return -1
}

func parseStreamSize(names []string) core.ByteSize {
	for _, name := range names {
		if matches := sizeRegex.FindStringSubmatch(name); len(matches) == 2 {
			if size, err := core.ParseByteSize(matches[1]); err == nil {
				return size
			}
		}
	}

	return core.ByteSize(0)
}

// Display

func StreamText(stream Stream) string {
	return stream.Name
}

func StreamWidget(stream Stream, selected bool) ui.Widget {
	style := tcell.StyleDefault
	if selected {
		style = ui.Fg(color.Lime)
	}

	spans := []ui.Span{{stream.Name + "\n", style}}

	addStreamMeta(&spans, stream.Resolution, color.Gray)
	addStreamMeta(&spans, stream.VideoCodec, color.Gray)
	addStreamMeta(&spans, stream.AudioCodec, color.Gray)
	if stream.Seeders >= 0 {
		addStreamMeta(&spans, fmt.Sprintf("%d seeders", stream.Seeders), color.Gray)
	}
	addStreamMeta(&spans, stream.Size.String(), color.Gray)

	switch stream.Cache {
	case Unknown:
	case Waiting:
		addStreamMeta(&spans, "...", color.Gray)
	case Uncached:
		addStreamMeta(&spans, "Uncached", color.Black)
	case Cached:
		addStreamMeta(&spans, "Cached", color.Olive)
	}

	return &ui.Paragraph{Spans: spans}
}

func addStreamMeta(spans *[]ui.Span, value string, fg color.Color) {
	if value == "" {
		return
	}

	if len(*spans) > 1 {
		*spans = append(*spans, ui.Span{Text: ", ", Style: ui.Fg(color.Silver)})
	}

	*spans = append(*spans, ui.Span{Text: value, Style: ui.Fg(fg)})
}
