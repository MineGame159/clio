package views

import (
	"clio/stremio"
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var metadataColorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
var normalMetadataStyle = normalItemStyle.Foreground(lipgloss.Color("8"))
var selectedMetadataStyle = selectedItemStyle.Foreground(lipgloss.Color("8"))

type Stream struct {
	Name            string
	Resolution      string
	Size            stremio.ByteSize
	VideosInTorrent int

	Url string
}

func (s Stream) FilterValue() string {
	return s.Name
}

type StreamDelegate struct{}

func (s StreamDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	stream := item.(*Stream)
	selected := index == m.Index()

	// Name
	nameStyle := normalItemStyle

	if selected {
		_, _ = fmt.Fprint(w, selectedItemStyle.Render("│"))
		nameStyle = selectedItemStyle
	}
	_, _ = fmt.Fprintln(w, nameStyle.Render(stream.Name))

	// Metadata
	metadataStyle := normalMetadataStyle

	if selected {
		_, _ = fmt.Fprint(w, selectedItemStyle.Render("│"))
		metadataStyle = selectedMetadataStyle
	}

	if stream.Resolution != "" {
		_, _ = fmt.Fprint(w, metadataStyle.Render(stream.Resolution), ", ")
		metadataStyle = metadataColorStyle
	}

	_, _ = fmt.Fprint(w, metadataStyle.Render(stream.Size.String()))
	//metadataStyle = metadataColorStyle

	//_, _ = fmt.Fprint(w, metadataStyle.Render(strconv.Itoa(stream.VideosInTorrent)))
}

func (s StreamDelegate) Height() int {
	return 2
}

func (s StreamDelegate) Spacing() int {
	return 0
}

func (s StreamDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}
