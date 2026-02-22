package views

import (
	"clio/ui"
	"time"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

type View interface {
	Title() string
	Keys() []Key

	Widgets() []ui.Widget

	HandleEvent(event any)
}

type Key struct {
	Bind        string
	Description string
}

// Stack

type Stack struct {
	views  []stackedView
	events chan any

	running       bool
	width, height int
}

type stackedView struct {
	view   View
	widget *ui.Block
}

func NewStack() *Stack {
	return &Stack{
		views:  nil,
		events: make(chan any, 4),
	}
}

func (s *Stack) Push(view View) {
	s.views = append(s.views, stackedView{view, nil})
}

func (s *Stack) Post(event any) {
	s.events <- event
}

func (s *Stack) Stop() {
	s.running = false
}

func (s *Stack) Run() {
	// Create screen
	screen, err := tcell.NewScreen()
	if err != nil {
		panic(err.Error())
	}

	time.Sleep(100 * time.Millisecond)

	if err := screen.Init(); err != nil {
		panic(err.Error())
	}

	defer screen.Fini()

	screen.Clear()

	// Run application
	s.running = true

	for s.running {
		// Get view widget
		if len(s.views) == 0 {
			s.running = false
			break
		}

		view := &s.views[len(s.views)-1]

		if view.widget == nil {
			view.createWidget()

			if s.width > 0 && s.height > 0 {
				event := tcell.NewEventResize(s.width, s.height)

				view.widget.HandleEvent(event)
				view.view.HandleEvent(event)
			}
		} else {
			view.updateFooter()
		}

		root := view.widget

		// Draw entire application
		root.CalcRequiredSize()
		width, height := screen.Size()

		screen.Clear()
		screen.HideCursor()

		root.Draw(screen, ui.Rect{
			X:      0,
			Y:      0,
			Width:  width,
			Height: height,
		})

		screen.Show()

		// Wait for a single event
		select {
		case event := <-screen.EventQ():
			s.handleEvent(screen, root, event)
		case event := <-s.events:
			view.view.HandleEvent(event)
		}

		// Immediately process buffered events without waiting
	outer:
		for {
			select {
			case event := <-screen.EventQ():
				s.handleEvent(screen, root, event)
			case event := <-s.events:
				view.view.HandleEvent(event)
			default:
				break outer
			}
		}
	}
}

func (s *Stack) handleEvent(screen tcell.Screen, root ui.Widget, event any) {
	switch event := event.(type) {
	case *tcell.EventResize:
		screen.Clear()
		screen.Sync()

		s.width, s.height = event.Size()

	case *tcell.EventKey:
		switch event.Key() {
		case tcell.KeyCtrlC:
			s.running = false
		case tcell.KeyEscape:
			s.views = s.views[:len(s.views)-1]
		default:
		}
	}

	if len(s.views) > 0 {
		root.HandleEvent(event)
		s.views[len(s.views)-1].view.HandleEvent(event)
	}
}

func (s *stackedView) createWidget() {
	s.widget = &ui.Block{
		Container: ui.Container{
			PrimaryAlignment:   ui.Stretch,
			SecondaryAlignment: ui.Stretch,
			Children:           s.view.Widgets(),
		},
		Runes: ui.Rounded,
		Title: []ui.Span{{s.view.Title(), ui.Fg(color.White).Bold(true)}},
	}

	s.updateFooter()
}

func (s *stackedView) updateFooter() {
	var footer []ui.Span

	for i, key := range s.view.Keys() {
		bind := key.Bind
		if i > 0 {
			bind = "   " + bind
		}

		footer = append(footer, ui.Span{Text: bind, Style: ui.Fg(color.Lime)})
		footer = append(footer, ui.Span{Text: " " + key.Description, Style: ui.Fg(color.Gray)})
	}

	s.widget.Footer = footer
}
