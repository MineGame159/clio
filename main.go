package main

import (
	"clio/debrid"
	"clio/debrid/rd"
	"clio/library"
	"clio/scraper"
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
	clients := make(map[string]debrid.Client)

	for _, url := range config.Addons {
		if strings.HasPrefix(url, "<library:") && strings.HasSuffix(url, ">") {
			var err error
			url, err = library.Start(getClient(clients, url[9:len(url)-1]))

			if err != nil {
				panic(err.Error())
			}
		} else if strings.HasPrefix(url, "<scraper:") && strings.HasSuffix(url, ">") {
			var err error
			url, err = scraper.Start(getClient(clients, url[9:len(url)-1]))

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

func getClient(clients map[string]debrid.Client, str string) debrid.Client {
	if client, ok := clients[str]; ok {
		return client
	}

	if strings.HasPrefix(str, "RD:") {
		client := rd.NewClient(str[3:])
		clients[str] = client

		return client
	}

	panic("Invalid debrid token: " + str)
}
