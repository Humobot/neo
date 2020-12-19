package neo

import (


	. "humobot/neo/providers"
)

type Neo struct {
	Provider Provider
}

var providers = make(map[string]Provider, 0)

func NewNeo() *Neo {
	return &Neo{}
}

func AddProvider(name string, provider Provider) {
	providers[name] = provider
}

func (n *Neo) Use(providerName string) *Neo {
	n.Provider = providers[providerName]

	return n
}

func (n *Neo) Speak(word string, kind string, file string) error {
	if kind == "online" {
		return n.Provider.TtsOnline(word, file)
	}

	return nil

}

func (n *Neo) Listen() {
	
} 
