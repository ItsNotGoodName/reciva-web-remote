package api

import (
	"encoding/json"
	"net/http"

	"github.com/ItsNotGoodName/reciva-web-remote/core"
	"github.com/ItsNotGoodName/reciva-web-remote/core/preset"
	"github.com/ItsNotGoodName/reciva-web-remote/left/presenter"
)

func handlePresetError(err error) presenter.Response {
	code := http.StatusInternalServerError
	if err == core.ErrPresetNotFound {
		code = http.StatusNotFound
	}

	return presenter.Response{
		Code:  code,
		Error: err,
	}
}

func GetPreset(presetStore preset.PresetStore) presenter.Requester {
	return func(r *http.Request) presenter.Response {
		url := r.URL.Query().Get("url")

		p, err := presetStore.Get(r.Context(), url)
		if err != nil {
			return handlePresetError(err)
		}

		return presenter.Response{
			Code: http.StatusOK,
			Data: p,
		}
	}
}

func GetPresets(presetStore preset.PresetStore) presenter.Requester {
	return func(r *http.Request) presenter.Response {
		ps, err := presetStore.List(r.Context())
		if err != nil {
			handlePresetError(err)
		}

		return presenter.Response{
			Code: http.StatusOK,
			Data: ps,
		}
	}
}

func PostPreset(presetStore preset.PresetStore) presenter.Requester {
	type request struct {
		URL      string `json:"url"`
		TitleNew string `json:"title_new"`
		URLNew   string `json:"url_new"`
	}

	return func(r *http.Request) presenter.Response {
		req := request{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			return presenter.Response{
				Code:  http.StatusBadRequest,
				Error: err,
			}
		}

		p, err := preset.ParsePreset(req.URL, req.TitleNew, req.URLNew)
		if err != nil {
			return presenter.Response{
				Code:  http.StatusBadRequest,
				Error: err,
			}
		}

		if err := presetStore.Update(r.Context(), p); err != nil {
			return handlePresetError(err)
		}

		return presenter.Response{
			Code: http.StatusOK,
		}
	}
}

func GetPresetURL(presetStore preset.PresetStore, url string) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		preset, err := presetStore.Get(r.Context(), url)
		if err != nil {
			http.Error(rw, err.Error(), handlePresetError(err).Code)
			return
		}

		rw.Write([]byte(preset.URLNew))
	}
}
