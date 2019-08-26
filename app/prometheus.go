package app

import (
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/haozibi/zlog"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Subsystem: "http",
		Name:      "lc_requests_total",
		Help:      "total HTTP requests processed",
	}, []string{"code", "method"})

	requestDurationHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Subsystem: "http",
		Name:      "lc_request_duration_seconds",
		Help:      "Seconds spent serving HTTP requests.",
		Buckets:   prometheus.DefBuckets,
	}, []string{"code", "method"})
)

func (a *APP) runMonitor() error {

	r := mux.NewRouter()

	addr := ":2112"

	r.Handle("/metrics", promhttp.Handler())

	srv := &http.Server{
		Addr:         addr,
		WriteTimeout: 120 * time.Second,
		ReadTimeout:  120 * time.Second,
		Handler:      handlers.RecoveryHandler()(r),
	}

	zlog.ZInfo().Str("Addr", addr).Msg("[metrics]")

	if err := srv.ListenAndServe(); err != nil {
		return errors.Wrap(err, "metrics run error")
	}
	return nil
}
