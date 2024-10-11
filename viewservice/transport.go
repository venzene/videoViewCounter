package viewservice

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func MakeHandler(endpoints Endpoints, logger kitlog.Logger) http.Handler {
	r := mux.NewRouter()

	r.Handle("/", kithttp.NewServer(
		endpoints.GetAllViews,
		decodeGetAllViewsRequest,
		encodeResponse,
	)).Methods("GET")

	r.Handle("/views/{id}", kithttp.NewServer(
		endpoints.GetView,
		decodeGetViewRequest,
		encodeResponse,
	)).Methods("GET")

	r.Handle("/increment/{id}", kithttp.NewServer(
		endpoints.Increment,
		decodeIncrementRequest,
		encodeResponse,
	)).Methods("POST")

	r.Handle("/top/{n}", kithttp.NewServer(
		endpoints.GetTopVideos,
		decodeGetTopVideosRequest,
		encodeResponse,
	)).Methods("GET")

	r.Handle("/recent/{n}", kithttp.NewServer(
		endpoints.GetRecentVideos,
		decodeGetRecentVideosRequest,
		encodeResponse,
	)).Methods("GET")

	return r
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response any) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err, ok := response.(error); ok && err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
	}
	return json.NewEncoder(w).Encode(response)
}

func decodeGetViewRequest(_ context.Context, r *http.Request) (any, error) {
	vars := mux.Vars(r)
	videoId, ok := vars["id"]
	if !ok {
		return nil, errors.New("missing video id")
	}
	return getViewRequest{videoId: videoId}, nil
}

func decodeGetAllViewsRequest(_ context.Context, r *http.Request) (any, error) {
	return nil, nil
}

func decodeIncrementRequest(_ context.Context, r *http.Request) (any, error) {
	vars := mux.Vars(r)
	videoId, ok := vars["id"]
	if !ok {
		return nil, errors.New("missing video id")
	}
	return incrementRequest{videoId: videoId}, nil
}

func decodeGetRecentVideosRequest(_ context.Context, r *http.Request) (any, error) {
	vars := mux.Vars(r)
	nStr, ok := vars["n"]
	if !ok {
		return nil, errors.New("missing parameter 'n'")
	}

	nInt, err := strconv.Atoi(nStr)
	if err != nil {
		return nil, errors.New("invalid parameter 'n', must be an integer")
	}
	return getRecentVideosRequest{n: nInt}, nil
}

func decodeGetTopVideosRequest(_ context.Context, r *http.Request) (any, error) {
	vars := mux.Vars(r)
	nStr, ok := vars["n"]
	if !ok {
		return nil, errors.New("missing parameter 'n'")
	}

	nInt, err := strconv.Atoi(nStr)
	if err != nil {
		return nil, errors.New("invalid parameter 'n', must be an integer")
	}
	return getTopVideosRequest{n: nInt}, nil
}
