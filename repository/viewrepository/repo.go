package viewrepository

import (
	"context"
	"errors"
	"view_count/model"
)

var (
	ErrVideoIdNotFound = errors.New("video id not found")
)

type Repository interface {
	GetAllViews(ctx context.Context) (info []model.VideoInfo, err error)

	Increment(ctx context.Context, videoId string) (err error)

	// TODO: write expectation of result
	GetView(ctx context.Context, videoId string) (view int, err error)

	GetTopVideos(ctx context.Context, n int) (info []model.VideoInfo, err error) // add n as param

	GetRecentVideos(ctx context.Context, n int) (info []model.VideoInfo, err error) // n as param
}
