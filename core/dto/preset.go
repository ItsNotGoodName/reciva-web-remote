package dto

type Preset struct {
	URL      string `json:"url"`
	TitleNew string `json:"title_new"`
	URLNew   string `json:"url_new"`
}

type PresetGetRequest struct {
	URL string `json:"url"`
}

type PresetGetResponse struct {
	Preset Preset `json:"preset"`
}

type PresetListResponse struct {
	Presets []Preset `json:"presets"`
}

type PresetUpdateRequest struct {
	Preset Preset `json:"preset"`
}
