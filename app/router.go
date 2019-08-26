package app

import (
	"io"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func setRouter(r *mux.Router, a *APP, w io.Writer) {

	// r.Use(Monitor)
	r.HandleFunc("/version", a.Version)

	apiCNV1 := r.PathPrefix("/v1cn").Subrouter()
	apiENV1 := r.PathPrefix("/v1").Subrouter()

	apiCNV1.Use(Monitor)
	apiENV1.Use(Monitor)

	router := make([]*mux.Router, 0, 5)
	router = append(router, apiCNV1)
	router = append(router, apiENV1)

	for _, api := range router {

		// 排名
		api.Methods(http.MethodGet).Path("/ranking/{name:.+}.svg").Handler(
			handlers.CombinedLoggingHandler(w, http.HandlerFunc(a.Badge)),
		)

		// 通过的题目/问题总数
		api.Methods(http.MethodGet).Path("/solved/{name:.+}.svg").Handler(
			handlers.CombinedLoggingHandler(w, http.HandlerFunc(a.Badge)),
		)

		// 通过的题目/问题总数 rate
		api.Methods(http.MethodGet).Path("/solved-rate/{name:.+}.svg").Handler(
			handlers.CombinedLoggingHandler(w, http.HandlerFunc(a.Badge)),
		)

		// 通过提交/提交的总数
		api.Methods(http.MethodGet).Path("/accepted/{name:.+}.svg").Handler(
			handlers.CombinedLoggingHandler(w, http.HandlerFunc(a.Badge)),
		)

		// 通过提交/提交的总数 rate
		api.Methods(http.MethodGet).Path("/accepted-rate/{name:.+}.svg").Handler(
			handlers.CombinedLoggingHandler(w, http.HandlerFunc(a.Badge)),
		)

		api.Methods(http.MethodGet).Path("/{name:.+}.svg").Handler(
			handlers.CombinedLoggingHandler(w, http.HandlerFunc(a.Profile)),
		)
	}
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
}

// Monitor Monitor
func Monitor(h http.Handler) http.Handler {
	return promhttp.InstrumentHandlerCounter(requestsTotal,
		promhttp.InstrumentHandlerDuration(requestDurationHistogram, h))
}
