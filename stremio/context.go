package stremio

import "iter"

type Context struct {
	Addons []*Addon
}

func (c *Context) Catalogs() iter.Seq[*Catalog] {
	return func(yield func(*Catalog) bool) {
		for _, addon := range c.Addons {
			for _, catalog := range addon.Catalogs {
				if !yield(catalog) {
					return
				}
			}
		}
	}
}

func (c *Context) StreamProviders() iter.Seq[*StreamProvider] {
	return func(yield func(*StreamProvider) bool) {
		for _, addon := range c.Addons {
			for _, streamProvider := range addon.StreamProviders {
				if !yield(streamProvider) {
					return
				}
			}
		}
	}
}

func (c *Context) StreamProviderForKindId(kind MediaKind, id string) *StreamProvider {
	for streamProvider := range c.StreamProviders() {
		if streamProvider.SupportsKindId(kind, id) {
			return streamProvider
		}
	}

	return nil
}

func (c *Context) MetaProviders() iter.Seq[*MetaProvider] {
	return func(yield func(*MetaProvider) bool) {
		for _, addon := range c.Addons {
			for _, metaProvider := range addon.MetaProviders {
				if !yield(metaProvider) {
					return
				}
			}
		}
	}
}

func (c *Context) MetaProviderForKindId(kind MediaKind, id string) *MetaProvider {
	for metaProvider := range c.MetaProviders() {
		if metaProvider.SupportsKindId(kind, id) {
			return metaProvider
		}
	}

	return nil
}
