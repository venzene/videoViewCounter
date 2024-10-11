package main

import (
	"view_count/middleware"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func routeIntialiser(h handler) *mux.Router {

	r := mux.NewRouter()
	r.Use(middleware.LoggingMiddleware)

	r.HandleFunc("/", h.handleIndex)
	r.HandleFunc("/increment/{vID}", h.handleIncrement)
	r.HandleFunc("/views/{vID}", h.handleViews)
	r.HandleFunc("/top/{n}", h.handleTopVideos)
	r.HandleFunc("/recent/{n}", h.handleRecentVideos)

	r.Handle("/metrics", promhttp.Handler())

	// TODO: add handler which returns top 10 view video ids : Done
	// TODO: add handler which gives me 10 recent incrment video ids : Done

	return r
}
