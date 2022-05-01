package api

import (
	"encoding/json"
	"net/http"

	"github.com/ItsNotGoodName/reciva-web-remote/core/dto"
	"github.com/ItsNotGoodName/reciva-web-remote/left/presenter"
)

func GetStates(app dto.App) presenter.Requester {
	return func(r *http.Request) presenter.Response {
		res, err := app.StateList(r.Context())
		if err != nil {
			return handleRadioError(err)
		}

		return presenter.Response{Code: http.StatusOK, Data: res.States}
	}
}

func GetState(app dto.App) UUIDRequester {
	return func(r *http.Request, uuid string) presenter.Response {
		res, err := app.StateGet(r.Context(), &dto.StateRequest{UUID: uuid})
		if err != nil {
			return handleRadioError(err)
		}

		return presenter.Response{Code: http.StatusOK, Data: res.State}
	}
}

func PatchState(app dto.App) UUIDRequester {
	return func(r *http.Request, uuid string) presenter.Response {
		var req dto.StatePatchRequest = dto.StatePatchRequest{UUID: uuid}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return presenter.Response{Code: http.StatusBadRequest, Error: err}
		}

		if err := app.StatePatch(r.Context(), &req); err != nil {
			return handleRadioError(err)
		}

		return presenter.Response{Code: http.StatusOK}
	}
}
