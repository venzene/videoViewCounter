package viewservice

import (
	"context"
	"time"
	"view_count/model"

	"github.com/go-kit/log"
)

type ServiceLogging struct {
	logger log.Logger
	Service
}

func NewServiceLogging(logger log.Logger, s Service) Service {
	return &ServiceLogging{logger, s}
}

func (s *ServiceLogging) GetAllViews(ctx context.Context) (info []model.VideoInfo, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"Method", "GetAllViews",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.GetAllViews(ctx)
}

func (s *ServiceLogging) GetView(ctx context.Context, videoId string) (view int, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"Method", "GetView",
			"videoId", videoId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.GetView(ctx, videoId)
}

func (s *ServiceLogging) Increment(ctx context.Context, videoId string) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"Method", "Increment",
			"videoId", videoId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Increment(ctx, videoId)
}

func (s *ServiceLogging) TopVideos(ctx context.Context, num int) (info []model.VideoInfo, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"Method", "TopVideos",
			"Params", num,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.GetTopVideos(ctx, num)
}

func (s *ServiceLogging) RecentViews(ctx context.Context, num int) (info []model.VideoInfo, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"Method", "RecentViews",
			"Params", num, "took",
			time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.GetRecentVideos(ctx, num)
}
