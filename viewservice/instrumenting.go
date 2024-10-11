package viewservice

import (
	"context"
	"time"
	"view_count/model"

	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/log"
)

type instrumentingService struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	Service
	logger log.Logger
}

// NewInstrumentingService returns an instance of an instrumenting Service.
func NewInstrumentingService(counter metrics.Counter, latency metrics.Histogram, logger log.Logger, s Service) Service {
	return &instrumentingService{
		requestCount:   counter,
		requestLatency: latency,
		logger:         logger,
		Service:        s,
	}
}

func (s *instrumentingService) GetAllViews(ctx context.Context) ([]model.VideoInfo, error) {
	defer func(begin time.Time) {
		requestLatency := time.Since(begin)
		s.requestCount.With("method", "GetAllViews").Add(1)
		s.requestLatency.With("method", "GetAllViews").Observe(requestLatency.Seconds())
		s.logger.Log(
			"method", "GetAllViews",
			"requestLatency", requestLatency.Microseconds(),
		)
	}(time.Now())
	return s.Service.GetAllViews(ctx)
}

func (s *instrumentingService) Increment(ctx context.Context, videoId string) (err error) {
	defer func(begin time.Time) {
		requestLatency := time.Since(begin)
		s.requestCount.With("method", "Increment").Add(1)
		s.requestLatency.With("method", "Increment").Observe(requestLatency.Seconds())
		s.logger.Log(
			"method", "Increment",
			"requestLatency", requestLatency.Microseconds(),
		)
	}(time.Now())
	return s.Service.Increment(ctx, videoId)
}

func (s *instrumentingService) GetView(ctx context.Context, videoId string) (view int, err error) {
	defer func(begin time.Time) {
		requestLatency := time.Since(begin)
		s.requestCount.With("method", "GetView").Add(1)
		s.requestLatency.With("method", "GetView").Observe(requestLatency.Seconds())
		s.logger.Log(
			"method", "GetViews",
			"requestLatency", requestLatency.Microseconds(),
		)
	}(time.Now())
	return s.Service.GetView(ctx, videoId)
}

func (s *instrumentingService) GetRecentVideos(ctx context.Context, n int) ([]model.VideoInfo, error) {
	defer func(begin time.Time) {
		requestLatency := time.Since(begin)
		s.requestCount.With("method", "GetRecentVideos").Add(1)
		s.requestLatency.With("method", "GetRecentVideos").Observe(requestLatency.Seconds())
		s.logger.Log(
			"method", "GetRecentVideos",
			"requestLatency", requestLatency.Microseconds(),
		)
	}(time.Now())
	return s.Service.GetRecentVideos(ctx, n)
}

func (s *instrumentingService) GetTopVideos(ctx context.Context, n int) ([]model.VideoInfo, error) {
	defer func(begin time.Time) {
		requestLatency := time.Since(begin)
		s.requestCount.With("method", "GetTopVideos").Add(1)
		s.requestLatency.With("method", "GetTopVideos").Observe(requestLatency.Seconds())
		s.logger.Log(
			"method", "GetTopVideos",
			"requestLatency", requestLatency.Microseconds(),
		)
	}(time.Now())
	return s.Service.GetTopVideos(ctx, n)
}
