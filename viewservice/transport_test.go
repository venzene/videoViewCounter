package viewservice

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"view_count/model"

	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// type MockEndpoint struct {
// 	GetView       endpoint.Endpoint
// 	GetAllViews     endpoint.Endpoint
// 	Increment       endpoint.Endpoint
// 	GetTopVideos    endpoint.Endpoint
// 	GetRecentVideos endpoint.Endpoint
// }

func TestTransport(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := kitlog.NewNopLogger()

	endpoints := Endpoints{
		GetView:         MockGetViewsEndpoint(),
		GetAllViews:     MockGetAllViewsEndpoint(),
		Increment:       MockIncrementEndpoint(),
		GetTopVideos:    MockGetTopVideosEndpoint(),
		GetRecentVideos: MockGetRecentVideosEndpoint(),
	}

	handler := MakeHandler(endpoints, mockLogger)

	t.Run("GetView", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/views/vishal", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		res := rec.Result()
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("GetAllViews", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		res := rec.Result()
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("Increment", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/increment/vishal", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		res := rec.Result()
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("GetTopVideos", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/top/2", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		res := rec.Result()
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("GetRecentVideos", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/recent/2", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		res := rec.Result()
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("Encoding throws an status 500", func(t *testing.T) {
		rec := httptest.NewRecorder()
		customError := errors.New("custom error")
		err := encodeResponse(context.Background(), rec, customError)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Result().StatusCode)
	})
}

func MockGetViewsEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		return 1, nil
	}
}

func MockGetAllViewsEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		return []model.VideoInfo{
			{
				Id:    "video0",
				Views: 1,
			},
			{
				Id:    "video1",
				Views: 2,
			},
		}, nil
	}
}

func MockIncrementEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		return nil, nil
	}
}

func MockGetTopVideosEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		return []model.VideoInfo{
			{
				Id:    "video0",
				Views: 2,
			},
			{
				Id:    "video1",
				Views: 1,
			},
		}, nil
	}
}

func MockGetRecentVideosEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		return []model.VideoInfo{
			{
				Id:    "video0",
				Views: 1,
			},
			{
				Id:    "video1",
				Views: 2,
			},
		}, nil
	}
}
