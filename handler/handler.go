package handler

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/bukalapak/packen/metric"
	"github.com/bukalapak/packen/response"
	"github.com/julienschmidt/httprouter"

	gx "github.com/bukalapak/go-xample"
)

type Handler struct {
	Gx gx.GoXample
}

type Meta struct {
	HTTPStatus int `json:"http_status"`
}

func NewHandler(goXample gx.GoXample) Handler {
	return Handler{
		Gx: goXample,
	}
}

func (h *Handler) Healthz(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, "ok")
}

func (h *Handler) Metric(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	metric.Handler(w, r)
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// TODO: prepare context

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		res := response.BuildError([]error{err})
		response.Write(w, res)
		return
	}

	user, err := h.Gx.CreateUser(string(body))
	if err != nil {
		res := response.BuildError([]error{err})
		response.Write(w, res)
		return
	}

	res := response.BuildSuccess(user, Meta{HTTPStatus: 200})
	response.Write(w, res)
}
