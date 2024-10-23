package httpclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"view_count/model"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

type Client struct {
	GetAllViews     endpoint.Endpoint
	GetTopVideos    endpoint.Endpoint
	Increment       endpoint.Endpoint
	GetView         endpoint.Endpoint
	GetRecentVideos endpoint.Endpoint
}

func NewClient(baseURL string) *Client {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		panic(err)
	}
	return &Client{
		GetView: httptransport.NewClient(
			"GET",
			parsedURL.ResolveReference(&url.URL{Path: "/views/{vid}"}),
			encodeGetViewRequest,
			decodeGetViewResponse,
		).Endpoint(),

		GetAllViews: httptransport.NewClient(
			"GET",
			parsedURL.ResolveReference(&url.URL{Path: ""}),
			encodeRequestAll,
			decodeGetAllViewsResponse,
		).Endpoint(),

		Increment: httptransport.NewClient(
			"POST",
			parsedURL.ResolveReference(&url.URL{Path: "/increment/{vid}"}),
			encodeIncrementRequest,
			decodeResponse,
		).Endpoint(),

		GetTopVideos: httptransport.NewClient(
			"GET",
			parsedURL.ResolveReference(&url.URL{Path: "/top/{n}"}),
			encodeGetTopVideosRequest,
			decodeGetTopVideosResponse,
		).Endpoint(),

		GetRecentVideos: httptransport.NewClient(
			"GET",
			parsedURL.ResolveReference(&url.URL{Path: "/recent/{n}"}),
			encodeGetRecentVideosRequest,
			decodeGetRecentVideosResponse,
		).Endpoint(),
	}
}

// encoding request 

func encodeRequestAll(_ context.Context, r *http.Request, request interface{}) error {
	return nil
}

func encodeGetViewRequest(_ context.Context, r *http.Request, request interface{}) error {
	req := request.(string)
	r.URL.Path = fmt.Sprintf("/views/%s", req)
	return nil
}

// func encodeGetAllRequest() -> not needed

func encodeIncrementRequest(ctx context.Context, r *http.Request, request interface{}) error {
	req := request.(string)
	r.URL.Path = fmt.Sprintf("/increment/%s", req)
	return nil
}

func encodeGetTopVideosRequest(ctx context.Context, r *http.Request, request interface{}) error {
	req := request.(int)
	r.URL.Path = fmt.Sprintf("/top/%d", req)
	return nil
}

func encodeGetRecentVideosRequest(ctx context.Context, r *http.Request, request interface{}) error {
	req := request.(int)
	r.URL.Path = fmt.Sprintf("/recent/%d", req)
	return nil
}

// decoding request

func decodeResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode == http.StatusNoContent {
		return nil, nil 
	}
	if r.StatusCode != http.StatusOK {
		return nil, errors.New("request failed")
	}
	return nil, nil
}

func decodeGetViewResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var response struct {
		Views int `json:"views"`
	}
	err := json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		return nil, err
	}
	return response.Views, nil
}

func decodeGetAllViewsResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var response struct {
		Videos []model.VideoInfo `json:"videos"`
	}
	err := json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		return nil, err
	}
	return response.Videos, nil
}

func decodeGetTopVideosResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var response struct {
		Videos []model.VideoInfo `json:"videos"`
	}
	err := json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		return nil, err
	}
	return response.Videos, nil
}

func decodeGetRecentVideosResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var response struct {
		Videos []model.VideoInfo `json:"videos"`
	}
	err := json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		return nil, err
	}
	return response.Videos, nil
}
