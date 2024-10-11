package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"view_count/viewservice"

	"github.com/gorilla/mux"
)

type handler struct {
	viewService viewservice.Service
}

func NewHandler(viewservice viewservice.Service) *handler {
	return &handler{
		viewService: viewservice,
	}
}

// Job of transport Routing, Encoding, Decoding : Done
func (h *handler) handleIndex(w http.ResponseWriter, r *http.Request) {

	acceptHeader := r.Header.Get("Accept")
	videos, err := h.viewService.GetAllViews(r.Context())
	if err != nil {
		http.Error(w, "server error.", http.StatusInternalServerError)
		return
	}

	if acceptHeader == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(videos)
		return
	}

	templ, err := template.ParseFiles("templates/index.gohtml")
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
	err = templ.Execute(w, videos)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		fmt.Println("Template execution error:", err)
		return
	}

}

func (h *handler) handleViews(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	videoID := vars["vID"]

	views, err := h.viewService.GetView(r.Context(), videoID)
	// TODO. add switch case for invalidarugment error and reutrn StatusBadRequest
	switch err {
	case nil:
	case viewservice.ErrInvalidArgument:
		http.Error(w, "VideoID is Required.", http.StatusBadRequest)
		return
	default:
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Number of views for video#%s: %d", videoID, views)
}

func (h *handler) handleIncrement(w http.ResponseWriter, r *http.Request) {
	// TODO: use gorilla mux : done
	vars := mux.Vars(r)
	videoID := vars["vID"]

	err := h.viewService.Increment(r.Context(), videoID)
	switch err {
	case nil:
	case viewservice.ErrInvalidArgument:
		http.Error(w, "VideoID is Required.", http.StatusBadRequest)
		return
	default:
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Success#%s", videoID)
}

func (h *handler) handleTopVideos(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	nStr := vars["n"]
	n, err := strconv.Atoi(nStr)
	if err != nil {
		http.Error(w, "Invalid n parameter", http.StatusBadRequest)
	}

	acceptHeader := r.Header.Get("Accept")
	videos, err := h.viewService.GetTopVideos(r.Context(), n)
	if err != nil {
		http.Error(w, "server error.", http.StatusInternalServerError)
		return
	}

	if acceptHeader == "application/json" {
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode(videos)
		return
	}

	templ, err := template.ParseFiles("templates/index.gohtml")
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
	err = templ.Execute(w, videos)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		fmt.Println("Template execution error:", err)
		return
	}
}

func (h *handler) handleRecentVideos(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nStr := vars["n"]
	n, err := strconv.Atoi(nStr)
	if err != nil {
		http.Error(w, "Invalid n parameter", http.StatusBadRequest)
	}

	videos, err := h.viewService.GetRecentVideos(r.Context(), n)
	if err != nil {
		http.Error(w, "server error.", http.StatusInternalServerError)
		return
	}

	acceptHeader := r.Header.Get("Accept")
	if acceptHeader == "application/json" {
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode(videos)
		return
	}

	templ, err := template.ParseFiles("templates/index.gohtml")
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
	templ.Execute(w, videos)
}
