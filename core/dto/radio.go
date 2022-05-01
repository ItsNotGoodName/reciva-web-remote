package dto

type Radio struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

type RadioRequest struct {
	UUID string `json:"uuid"`
}

type RadioGetResponse struct {
	Radio Radio `json:"radio"`
}

type RadioListResponse struct {
	Radios []Radio `json:"radios"`
}

type RadioDiscoverRequest struct {
	Force bool `json:"force"`
}

type RadioDiscoverResponse struct {
	Count int `json:"count"`
}
