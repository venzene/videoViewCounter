package viewservice

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

// TODO: write gokit client also

// Write unit test cases. Hint: use httptest package : done
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

func encodeResponse(_ context.Context, w http.ResponseWriter, response any) error {
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
	// Validatiaons should be at service level
	videoId := vars["id"]
	return getViewRequest{videoId: videoId}, nil
}

func decodeGetAllViewsRequest(_ context.Context, r *http.Request) (any, error) {
	return nil, nil
}

func decodeIncrementRequest(_ context.Context, r *http.Request) (any, error) {
	vars := mux.Vars(r)
	videoId := vars["id"]
	return incrementRequest{videoId: videoId}, nil
}

func decodeGetRecentVideosRequest(_ context.Context, r *http.Request) (any, error) {
	vars := mux.Vars(r)
	nStr := vars["n"]
	nInt, _ := strconv.Atoi(nStr)

	return getRecentVideosRequest{n: nInt}, nil
}

func decodeGetTopVideosRequest(_ context.Context, r *http.Request) (any, error) {
	vars := mux.Vars(r)
	nStr := vars["n"]
	nInt, _ := strconv.Atoi(nStr)

	return getTopVideosRequest{n: nInt}, nil
}
