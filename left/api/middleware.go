package api

import (
	"net/http"

	"github.com/ItsNotGoodName/reciva-web-remote/left/presenter"
	"github.com/go-chi/chi/v5"
)

type UUIDRequester func(r *http.Request, uuid string) presenter.Response

func UUIDRequire(next UUIDRequester) presenter.Requester {
	return func(r *http.Request) presenter.Response {
		return next(r, chi.URLParam(r, "uuid"))
	}
}
