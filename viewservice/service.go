package viewservice

import (
	"context"
	"errors"
	"view_count/model"
	"view_count/repository/viewrepository"
)

var (
	ErrInvalidArgument = errors.New("Invalid Argument")
)

// TODO add middleware of Service for logging and instrumenting : done
// Validataion, Coordinator

// TODO write unit testcases using gomock
type Service interface {
	GetAllViews(ctx context.Context) (info []model.VideoInfo, err error)

	// Increment will increment view count of given videoId.
	// it will return ErrInvalidArgument if videoId is empty
	Increment(ctx context.Context, videoId string) (err error)

	GetView(ctx context.Context, videoId string) (view int, err error)

	GetTopVideos(ctx context.Context, n int) (info []model.VideoInfo, err error)

	GetRecentVideos(ctx context.Context, n int) (info []model.VideoInfo, err error)
}

type service struct {
	viewRepo viewrepository.Repository
}

func NewService(viewRepo viewrepository.Repository) *service {
	return &service{
		viewRepo: viewRepo,
	}
}

func (svc *service) GetAllViews(ctx context.Context) (info []model.VideoInfo, err error) {
	return svc.viewRepo.GetAllViews(ctx)
}

// TODO: read about context https://www.youtube.com/watch?v=LSzR0VEraWw : done
// TODO: do debug for this code
func (svc *service) Increment(ctx context.Context, videoId string) (err error) {
	if len(videoId) < 1 {
		return ErrInvalidArgument
	}

	return svc.viewRepo.Increment(ctx, videoId)
}

func (svc *service) GetView(ctx context.Context, videoId string) (view int, err error) {
	if len(videoId) < 1 {
		return 0, ErrInvalidArgument
	}

	return svc.viewRepo.GetView(ctx, videoId)
}

func (svc *service) GetTopVideos(ctx context.Context, n int) (info []model.VideoInfo, err error) {

	return svc.viewRepo.GetTopVideos(ctx, n)

}

func (svc *service) GetRecentVideos(ctx context.Context, n int) ([]model.VideoInfo, error) {

	return svc.viewRepo.GetRecentVideos(ctx, n)

}
