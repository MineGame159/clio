package main

import (
	"clio/core"
	"clio/stremio"
	"clio/views"
	"encoding/json"
	"os"
	"path"

	tea "github.com/charmbracelet/bubbletea"
)

type Config struct {
	Addons []string
}

func main() {
	// Read config
	configDir, err := os.UserConfigDir()
	if err != nil {
		panic(err.Error())
	}

	configFile, err := os.Open(path.Join(configDir, "clio.json"))
	if err != nil {
		panic(err.Error())
	}
	defer configFile.Close()

	var config Config
	if err := json.NewDecoder(configFile).Decode(&config); err != nil {
		panic(err.Error())
	}

	// Create app
	app := core.NewApp()

	// Load addons
	for _, url := range config.Addons {
		addon, err := stremio.Load(url)
		if err != nil {
			println("Failed to load addon:", err.Error())
			continue
		}

		app.Addons = append(app.Addons, addon)
	}

	// Push catalogs view
	app.Push(&views.Catalogs{App: app})

	// Run program
	p := tea.NewProgram(app, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		panic(err.Error())
	}
}
