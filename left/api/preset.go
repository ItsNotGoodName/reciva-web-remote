package api

import (
	"encoding/json"
	"net/http"

	"github.com/ItsNotGoodName/reciva-web-remote/core"
	"github.com/ItsNotGoodName/reciva-web-remote/core/dto"
	"github.com/ItsNotGoodName/reciva-web-remote/left/presenter"
	"github.com/go-chi/chi/v5"
)

func handlePresetError(err error) presenter.Response {
	code := http.StatusInternalServerError
	if err == core.ErrPresetNotFound {
		code = http.StatusNotFound
	}

	return presenter.Response{Code: code, Error: err}
}

func PresetGet(app dto.App) presenter.Requester {
	return func(r *http.Request) presenter.Response {
		url := chi.URLParam(r, "*")

		res, err := app.PresetGet(r.Context(), &dto.PresetGetRequest{URL: url})
		if err != nil {
			return handlePresetError(err)
		}

		return presenter.Response{Code: http.StatusOK, Data: res.Preset}
	}
}

func PresetList(app dto.App) presenter.Requester {
	return func(r *http.Request) presenter.Response {
		res, err := app.PresetList(r.Context())
		if err != nil {
			handlePresetError(err)
		}

		return presenter.Response{Code: http.StatusOK, Data: res.Presets}
	}
}

func PresetUpdate(app dto.App) presenter.Requester {
	return func(r *http.Request) presenter.Response {
		preset := dto.Preset{}
		err := json.NewDecoder(r.Body).Decode(&preset)
		if err != nil {
			return presenter.Response{Code: http.StatusBadRequest, Error: err}
		}

		if err := app.PresetUpdate(r.Context(), &dto.PresetUpdateRequest{Preset: preset}); err != nil {
			return handlePresetError(err)
		}

		return presenter.Response{Code: http.StatusOK}
	}
}

func PresetGetURLNew(app dto.App, url string) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		res, err := app.PresetGet(r.Context(), &dto.PresetGetRequest{URL: url})
		if err != nil {
			http.Error(rw, err.Error(), handlePresetError(err).Code)
			return
		}

		rw.Write([]byte(res.Preset.URLNew))
	}
}
