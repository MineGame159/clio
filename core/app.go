package core

import (
	"clio/stremio"
	"iter"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type View struct {
	model       tea.Model
	initialized bool
}

type App struct {
	views []View

	width  int
	height int

	Addons []*stremio.Addon
}

func NewApp() *App {
	return &App{}
}

func (a *App) Catalogs() iter.Seq[*stremio.Catalog] {
	return func(yield func(*stremio.Catalog) bool) {
		for _, addon := range a.Addons {
			for _, catalog := range addon.Catalogs {
				if !yield(catalog) {
					return
				}
			}
		}
	}
}

func (a *App) StreamProviders() iter.Seq[*stremio.StreamProvider] {
	return func(yield func(*stremio.StreamProvider) bool) {
		for _, addon := range a.Addons {
			for _, streamProvider := range addon.StreamProviders {
				if !yield(streamProvider) {
					return
				}
			}
		}
	}
}

func (a *App) Push(model tea.Model) {
	a.views = append(a.views, View{model, false})
}

func (a *App) Pop() tea.Cmd {
	a.views = a.views[:len(a.views)-1]
	return nil
}

// Model

var appStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true).BorderForeground(White)

func (a *App) Init() tea.Cmd {
	if len(a.views) == 0 {
		return tea.Quit
	}

	view := &a.views[len(a.views)-1]
	view.initialized = true

	return view.model.Init()
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var updateCmd tea.Cmd
	var initCmd tea.Cmd
	var windowSizeCmd tea.Cmd

	// Handle window size
	switch m := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = m.Width - appStyle.GetHorizontalFrameSize()
		a.height = m.Height - appStyle.GetVerticalFrameSize()

		msg = tea.WindowSizeMsg{Width: a.width, Height: a.height}
	}

	// Update current view
	{
		index := len(a.views) - 1

		var model tea.Model
		model, updateCmd = a.views[index].model.Update(msg)

		if index < len(a.views) {
			a.views[index].model = model
		}
	}

	// Quit app or initialize new view
	if len(a.views) == 0 {
		initCmd = tea.Quit
	} else {
		view := &a.views[len(a.views)-1]

		if !view.initialized {
			view.initialized = true
			initCmd = view.model.Init()

			view.model, windowSizeCmd = view.model.Update(tea.WindowSizeMsg{Width: a.width, Height: a.height})
		}
	}

	return a, tea.Sequence(updateCmd, initCmd, windowSizeCmd)
}

func (a *App) View() string {
	if len(a.views) == 0 {
		return ""
	}

	contents := a.views[len(a.views)-1].model.View()
	return appStyle.Width(a.width).Height(a.height).Render(contents)
}
