package json

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ItsNotGoodName/reciva-web-remote/left/presenter"
)

type Response struct {
	OK    bool        `json:"ok"`
	Error *Error      `json:"error,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

type Error struct {
	Message string `json:"message"`
}

func Render(rw http.ResponseWriter, r presenter.Response) {
	if r.Error != nil {
		renderError(rw, r.Code, r.Error)
		return
	}
	renderJSON(rw, r.Code, r.Data)
}

func renderError(rw http.ResponseWriter, code int, err error) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(code)
	if err := json.NewEncoder(rw).Encode(Response{
		OK: false,
		Error: &Error{
			Message: err.Error(),
		},
	}); err != nil {
		log.Println("json.renderError:", err)
	}
}

func renderJSON(rw http.ResponseWriter, code int, data interface{}) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(code)
	if err := json.NewEncoder(rw).Encode(Response{
		OK:   true,
		Data: data,
	}); err != nil {
		log.Println("json.renderJSON:", err)
	}
}
