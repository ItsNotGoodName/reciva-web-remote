package model

type Build struct {
	BuiltBy    string `json:"built_by" validate:"required"`
	Commit     string `json:"commit" validate:"required"`
	Date       string `json:"date" validate:"required"`
	ReleaseURL string `json:"release_url" validate:"required"`
	Summary    string `json:"summary" validate:"required"`
	Version    string `json:"version" validate:"required"`
}
