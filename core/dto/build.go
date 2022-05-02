package dto

type Build struct {
	Version    string `json:"version"`
	Commit     string `json:"commit"`
	Date       string `json:"date"`
	BuiltBy    string `json:"built_by"`
	ReleaseURL string `json:"release_url"`
	Summary    string `json:"summary"`
}
