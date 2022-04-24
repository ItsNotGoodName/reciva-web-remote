package presenter

import (
	"net/http"
)

type (
	Response struct {
		Code  int
		Data  interface{}
		Error error
	}

	Render func(http.ResponseWriter, Response)

	Requester func(*http.Request) Response

	Presenter func(Requester) http.HandlerFunc
)

func New(render Render) Presenter {
	return func(requester Requester) http.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) {
			render(rw, requester(r))
		}
	}
}
