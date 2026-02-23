package stremio

import (
	"clio/core"
	"encoding/json"
	"errors"
	"strings"
)

type Addon struct {
	Name string

	Catalogs        []*Catalog
	MetaProviders   []*MetaProvider
	StreamProviders []*StreamProvider

	Url string
}

// Load

type manifest struct {
	Name            string            `json:"name"`
	Types           []string          `json:"types"`
	IdPrefixes      []string          `json:"idPrefixes"`
	Resources       []json.RawMessage `json:"resources,omitempty"`
	Catalogs        []Catalog         `json:"catalogs,omitempty"`
	StreamProviders []StreamProvider  `json:"streams,omitempty"`
}

func Load(url string) (*Addon, error) {
	addonUrl, ok := strings.CutSuffix(url, "/manifest.json")
	if !ok {
		return nil, errors.New("addon url needs to end with /manifest.json")
	}

	// Parse manifest
	man, err := core.GetJson[manifest](url)
	if err != nil {
		return nil, err
	}

	// Parse resources
	addon := &Addon{
		Name: man.Name,
		Url:  addonUrl,
	}

	for _, rawResource := range man.Resources {
		if string(rawResource) == "\"meta\"" {
			addon.MetaProviders = append(addon.MetaProviders, &MetaProvider{
				Addon:      addon,
				Types:      man.Types,
				IdPrefixes: man.IdPrefixes,
			})

			continue
		}

		var resource Extra
		if err := json.Unmarshal(rawResource, &resource); err != nil {
			continue
		}

		switch resource.Name {
		case "catalog":
			var catalog Catalog

			if err := json.Unmarshal(rawResource, &catalog); err == nil {
				catalog.Addon = addon
				addon.Catalogs = append(addon.Catalogs, &catalog)
			}

		case "meta":
			var metaProvider MetaProvider

			if err := json.Unmarshal(rawResource, &metaProvider); err == nil {
				metaProvider.Addon = addon

				if len(metaProvider.Types) == 0 {
					metaProvider.Types = man.Types
				}
				if len(metaProvider.IdPrefixes) == 0 {
					metaProvider.IdPrefixes = man.IdPrefixes
				}

				addon.MetaProviders = append(addon.MetaProviders, &metaProvider)
			}

		case "stream":
			var streamProvider StreamProvider

			if err := json.Unmarshal(rawResource, &streamProvider); err == nil {
				streamProvider.Addon = addon

				if len(streamProvider.Types) == 0 {
					streamProvider.Types = man.Types
				}
				if len(streamProvider.IdPrefixes) == 0 {
					streamProvider.IdPrefixes = man.IdPrefixes
				}

				addon.StreamProviders = append(addon.StreamProviders, &streamProvider)
			}
		}
	}

	for _, catalog := range man.Catalogs {
		catalog.Addon = addon
		addon.Catalogs = append(addon.Catalogs, &catalog)
	}

	for _, streamProvider := range man.StreamProviders {
		streamProvider.Addon = addon

		if len(streamProvider.IdPrefixes) == 0 {
			streamProvider.IdPrefixes = man.IdPrefixes
		}

		addon.StreamProviders = append(addon.StreamProviders, &streamProvider)
	}

	return addon, nil
}
