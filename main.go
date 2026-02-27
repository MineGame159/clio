package main

import (
	"clio/library"
	"clio/stremio"
	"clio/views"
	"encoding/json"
	"os"
	"path"
	"strings"
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
	stack := views.NewStack()

	// Create context
	ctx := &stremio.Context{}

	// Load addons
	for _, url := range config.Addons {
		if strings.HasPrefix(url, "<library:") && strings.HasSuffix(url, ">") {
			var err error
			url, err = library.Start(url[9 : len(url)-1])

			if err != nil {
				panic(err.Error())
			}
		}

		addon, err := stremio.Load(url)
		if err != nil {
			println("Failed to load addon:", err.Error())
			continue
		}

		ctx.Addons = append(ctx.Addons, addon)
	}

	// Push catalogs view
	stack.Push(&views.Catalogs{
		Stack: stack,
		Ctx:   ctx,
	})

	// Run application
	stack.Run()
}
