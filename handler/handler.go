package handler

import (
	"fmt"
	"net/http"

	"github.com/bukalapak/packen/metric"
)

func Healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, "ok")
}

func Metric(w http.ResponseWriter, r *http.Request) {
	metric.Handler(w, r)
}
