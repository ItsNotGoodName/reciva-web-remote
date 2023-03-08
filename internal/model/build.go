package model

type Build struct {
	BuiltBy    string `json:"built_by"`
	Commit     string `json:"commit"`
	Date       string `json:"date"`
	ReleaseURL string `json:"release_url"`
	Summary    string `json:"summary"`
	Version    string `json:"version"`
}
