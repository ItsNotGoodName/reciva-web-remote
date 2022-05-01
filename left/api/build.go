package api

import (
	"net/http"

	"github.com/ItsNotGoodName/reciva-web-remote/core/dto"
	"github.com/ItsNotGoodName/reciva-web-remote/left/presenter"
)

func GetBuild(app dto.App) presenter.Requester {
	return func(r *http.Request) presenter.Response {
		return presenter.Response{Code: http.StatusOK, Data: app.Build()}
	}
}
