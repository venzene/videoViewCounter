package viewservice

import (
	"context"
	"errors"
	"view_count/model"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	GetView         endpoint.Endpoint
	GetAllViews     endpoint.Endpoint
	Increment       endpoint.Endpoint
	GetTopVideos    endpoint.Endpoint
	GetRecentVideos endpoint.Endpoint
}

func MakeEndpoints(svc Service) Endpoints {
	return Endpoints{
		GetView:         MakeGetViewEndpoint(svc),
		GetAllViews:     MakeGetAllViewsEndpoint(svc),
		Increment:       MakeIncrementEndpoint(svc),
		GetTopVideos:    MakeGetTopVideosEndpoint(svc),
		GetRecentVideos: MakeGetRecentVideosEndpoint(svc),
	}
}

type getViewRequest struct {
	videoId string
}

type getViewResponse struct {
	Views int `json:"views"`
}

func MakeGetViewEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		req := request.(getViewRequest)
		views, err := svc.GetView(ctx, req.videoId)
		if err != nil {
			return nil, err
		}
		return getViewResponse{Views: views}, nil
	}
}

type getAllViewsResponse struct {
	Videos []model.VideoInfo `json:"videos"`
}

func MakeGetAllViewsEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		videos, err := svc.GetAllViews(ctx)
		if err != nil {
			return nil, err
		}
		return getAllViewsResponse{Videos: videos}, nil
	}
}

type incrementRequest struct {
	videoId string
}

type incrementResponse struct {
	Err error `json:"error,omitempty"`
}

func MakeIncrementEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		req, ok := request.(incrementRequest)
		if !ok {
			return nil, errors.New("invalid request")
		}
		err := svc.Increment(ctx, req.videoId)
		if err != nil {
			return nil, err
		}
		return incrementResponse{Err: nil}, nil
	}
}

type getRecentVideosRequest struct {
	n int
}

type getRecentVideosResponse struct {
	Videos []model.VideoInfo `json:"videos"`
}

func MakeGetRecentVideosEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		req, ok := request.(getRecentVideosRequest)
		if !ok {
			return nil, errors.New("invalid request")
		}
		videos, err := svc.GetRecentVideos(ctx, req.n)
		if err != nil {
			return nil, err
		}
		return getRecentVideosResponse{Videos: videos}, nil
	}
}

type getTopVideosRequest struct {
	n int
}

type getTopVideosResponse struct {
	Videos []model.VideoInfo `json:"videos"`
}

func MakeGetTopVideosEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		req, ok := request.(getTopVideosRequest)
		if !ok {
			return nil, errors.New("invalid request")
		}
		videos, err := svc.GetTopVideos(ctx, req.n)
		if err != nil {
			return nil, err
		}
		return getTopVideosResponse{Videos: videos}, nil
	}
}
