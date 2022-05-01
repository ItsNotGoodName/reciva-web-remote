package api

import (
	"net/http"
	"strconv"

	"github.com/ItsNotGoodName/reciva-web-remote/core"
	"github.com/ItsNotGoodName/reciva-web-remote/core/dto"
	"github.com/ItsNotGoodName/reciva-web-remote/left/presenter"
)

func handleRadioError(err error) presenter.Response {
	code := http.StatusInternalServerError
	if err == core.ErrRadioClosed {
		code = http.StatusGone
	} else if err == core.ErrRadioNotFound {
		code = http.StatusNotFound
	} else if err == core.ErrHubDiscovering {
		code = http.StatusConflict
	} else if err == core.ErrHubServiceClosed {
		code = http.StatusServiceUnavailable
	}

	return presenter.Response{Code: code, Error: err}
}

func GetRadios(app dto.App) presenter.Requester {
	return func(r *http.Request) presenter.Response {
		res, err := app.RadioList()
		if err != nil {
			return handleRadioError(err)
		}

		return presenter.Response{Code: http.StatusOK, Data: res.Radios}
	}
}

func PostRadios(app dto.App) presenter.Requester {
	return func(r *http.Request) presenter.Response {
		force, _ := strconv.ParseBool(r.URL.Query().Get("force"))

		res, err := app.RadioDiscover(&dto.RadioDiscoverRequest{Force: force})
		if err != nil {
			return handleRadioError(err)
		}

		return presenter.Response{Code: http.StatusOK, Data: res.Count}
	}
}

func GetRadio(app dto.App) UUIDRequester {
	return func(r *http.Request, uuid string) presenter.Response {
		res, err := app.RadioGet(&dto.RadioRequest{UUID: uuid})
		if err != nil {
			return handleRadioError(err)
		}

		return presenter.Response{Code: http.StatusOK, Data: res.Radio}
	}
}

func PostRadioSubscription(app dto.App) UUIDRequester {
	return func(r *http.Request, uuid string) presenter.Response {
		err := app.RadioRefreshSubscription(r.Context(), &dto.RadioRequest{UUID: uuid})
		if err != nil {
			return handleRadioError(err)
		}

		return presenter.Response{Code: http.StatusOK}
	}
}

func PostRadioVolume(app dto.App) UUIDRequester {
	return func(r *http.Request, uuid string) presenter.Response {
		err := app.RadioRefreshVolume(r.Context(), &dto.RadioRequest{UUID: uuid})
		if err != nil {
			return handleRadioError(err)
		}

		return presenter.Response{Code: http.StatusOK}
	}
}
