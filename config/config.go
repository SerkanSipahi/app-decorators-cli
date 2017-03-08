package config

type Appdec struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
}

const (
	AppName     = "appdec"
	AuthorName  = "Serkan Sipahi"
	AuthorEmail = "serkan.sipahi@yahoo.de"
	AppVersion  = "0.8.204"
	NpmPackage  = "app-decorators"
)