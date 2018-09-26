// Package handler manages the data flow from client to appropriate service.
package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"

	gx "github.com/bukalapak/go-xample"
)

// Handler controls request flow from client to service
type Handler struct {
	goXample *gx.GoXample
}

// Meta is used to consolidate all meta statuses
type Meta struct {
	HTTPStatus int `json:"http_status"`
}

// NewHandler returns a pointer of Handler instance
func NewHandler(goXample *gx.GoXample) *Handler {
	return &Handler{
		goXample: goXample,
	}
}

// Healthz is used to control the flow of GET /healthz endpoint
func (h *Handler) Healthz(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, "ok")
}

// CreateUser is used to control the flow of POST /users endpoint
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) error {
	ctx := r.Context()
	select {
	case <-ctx.Done():
		return errors.New("Timeout")
	default:
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeError(w, err)
		return err
	}

	var user gx.User

	err = json.Unmarshal(body, &user)
	if err != nil {
		return err
	}

	user, err = h.goXample.CreateUser(ctx, user)
	if err != nil {
		writeError(w, err)
		return err
	}

	writeSuccess(w, user)
	return nil
}

// GetUser is used to control the flow of GET /users endpoint
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) error {
	ctx := r.Context()
	select {
	case <-ctx.Done():
		return errors.New("Timeout")
	default:
	}

	userID, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		writeError(w, err)
		return err
	}

	user, err := h.goXample.GetUserByID(ctx, userID)
	if err != nil {
		writeError(w, err)
		return err
	}

	writeSuccess(w, user)
	return nil
}

// Login is used to control the flow of POST /login endpoint
func (h *Handler) Login(w http.ResponseWriter, r *http.Request, params httprouter.Params) error {
	ctx := r.Context()
	select {
	case <-ctx.Done():
		return errors.New("Timeout")
	default:
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeError(w, err)
		return err
	}

	var user gx.User

	err = json.Unmarshal(body, &user)
	if err != nil {
		return err
	}

	user, err = h.goXample.GetUserByCredential(ctx, user)
	if err != nil {
		writeError(w, err)
		return err
	}

	writeSuccess(w, user)
	return nil
}

// Logout is used to control the flow of GET /logout endpoint
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request, params httprouter.Params) error {
	ctx := r.Context()
	select {
	case <-ctx.Done():
		return errors.New("Timeout")
	default:
	}

	writeSuccess(w, "You've been successfully logout!")
	return nil
}

func writeError(w http.ResponseWriter, err error) {
	res := fmt.Sprintf("{\"ERROR\":%s}", err.Error())
	respWrite(w, []byte(res), 400)
}

func writeSuccess(w http.ResponseWriter, data interface{}) {
	res, _ := json.Marshal(data)
	respWrite(w, res, 200)
}

func respWrite(w http.ResponseWriter, res []byte, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(res)
}
