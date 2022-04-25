package api

import (
	"encoding/json"
	"net/http"

	"github.com/ItsNotGoodName/reciva-web-remote/core"
	"github.com/ItsNotGoodName/reciva-web-remote/core/dto"
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

		// Get preset
		p, err := presetStore.Get(r.Context(), url)
		if err != nil {
			return handlePresetError(err)
		}

		return presenter.Response{
			Code: http.StatusOK,
			Data: dto.NewPreset(p),
		}
	}
}

func GetPresets(presetStore preset.PresetStore) presenter.Requester {
	return func(r *http.Request) presenter.Response {
		// List presets
		p, err := presetStore.List(r.Context())
		if err != nil {
			handlePresetError(err)
		}

		return presenter.Response{
			Code: http.StatusOK,
			Data: dto.NewPresets(p),
		}
	}
}

func PostPreset(presetStore preset.PresetStore) presenter.Requester {
	return func(r *http.Request) presenter.Response {
		dtoPreset := &dto.Preset{}
		err := json.NewDecoder(r.Body).Decode(dtoPreset)
		if err != nil {
			return presenter.Response{
				Code:  http.StatusBadRequest,
				Error: err,
			}
		}

		p, err := dto.ConvertPreset(dtoPreset)
		if err != nil {
			return presenter.Response{
				Code:  http.StatusBadRequest,
				Error: err,
			}
		}

		// Update preset
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
		// Get preset
		dtoPreset, err := presetStore.Get(r.Context(), url)
		if err != nil {
			http.Error(rw, err.Error(), handlePresetError(err).Code)
			return
		}

		rw.Write([]byte(dto.NewPreset(dtoPreset).URLNew))
	}
}
