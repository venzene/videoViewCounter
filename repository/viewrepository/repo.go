package viewrepository

import (
	"context"
	"errors"
	"view_count/model"
)

var (
	ErrVideoIdNotFound = errors.New("video id not found")
)

// TODO: write expectation of result

type Repository interface {
	//returns all the listed videos along with thier respectivce views
	GetAllViews(ctx context.Context) (info []model.VideoInfo, err error)

	//takes videoId as the input and Increment the videoId by 1 & updates the heap
	Increment(ctx context.Context, videoId string) (err error)

	//takes the videoId as the input and returns the total views
	GetView(ctx context.Context, videoId string) (view int, err error)

	//returns top N videos on the basis of views
	GetTopVideos(ctx context.Context, n int) (info []model.VideoInfo, err error) // add n as param : Done

	//returns top N recently viewed videos
	GetRecentVideos(ctx context.Context, n int) (info []model.VideoInfo, err error) // n as param : Done
}

