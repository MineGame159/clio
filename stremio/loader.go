package stremio

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type manifest struct {
	Name            string            `json:"name"`
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
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var man manifest
	if err := json.NewDecoder(res.Body).Decode(&man); err != nil {
		return nil, err
	}

	// Parse resources
	addon := &Addon{
		Name: man.Name,
		Url:  addonUrl,
	}

	for _, rawResource := range man.Resources {
		var resource Extra
		if err := json.Unmarshal(rawResource, &resource); err != nil {
			continue
		}

		switch resource.Name {
		case "catalog":
			var catalog Catalog
			if err := json.Unmarshal(rawResource, &catalog); err == nil {
				addon.Catalogs = append(addon.Catalogs, &catalog)
			}
		case "stream":
			var streamProvider StreamProvider
			if err := json.Unmarshal(rawResource, &streamProvider); err == nil {
				addon.StreamProviders = append(addon.StreamProviders, &streamProvider)
			}
		}
	}

	for _, catalog := range man.Catalogs {
		addon.Catalogs = append(addon.Catalogs, &catalog)
	}

	for _, streamProvider := range man.StreamProviders {
		addon.StreamProviders = append(addon.StreamProviders, &streamProvider)
	}

	// Store reference to addon in resources
	for _, catalog := range addon.Catalogs {
		catalog.Addon = addon
	}

	for _, streamProvider := range addon.StreamProviders {
		streamProvider.Addon = addon
	}

	return addon, nil
}
