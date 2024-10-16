package viewservice

import (
	"context"
	"testing"
	"view_count/model"
	"view_count/repository/viewrepository"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type testCases struct {
	testName       string
	videoId        string
	expectedErr    error
}

func TestGetViews(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := viewrepository.NewMockRepository(ctrl)

	svc := NewService(mockRepo)

	tests := []testCases{
		{
			testName:    "invalid",
			videoId:     "",
			expectedErr: ErrInvalidArgument,
		},
		{
			testName:    "valid",
			videoId:     "video1",
			expectedErr: ErrInvalidArgument,
		},
	}

	views := 99

	for _, test := range tests {
		if test.videoId == "" {
			t.Run("Invalid", func(t *testing.T) {
				_, err := svc.GetView(context.Background(), test.videoId)
				assert.Error(t, err)
				assert.Equal(t, test.expectedErr, err)
			})
		} else {
			t.Run("Valid", func(t *testing.T) {
				mockRepo.EXPECT().GetView(context.Background(), test.videoId).Return(views, nil)
				result, err := svc.GetView(context.Background(), test.videoId)
				assert.NoError(t, err)
				assert.Equal(t, views, result)
			})
		}
	}

}

func TestGetAllViews(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := viewrepository.NewMockRepository(ctrl)

	svc := NewService(mockRepo)

	expectedResult := []model.VideoInfo{
		{
			Id:    "video1",
			Views: 1,
		},
		{
			Id:    "video2",
			Views: 2,
		},
	}

	mockRepo.EXPECT().GetAllViews(context.Background()).Return(expectedResult, nil)

	result, err := svc.GetAllViews(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestIncrement(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := viewrepository.NewMockRepository(ctrl)

	svc := NewService(mockRepo)

	tests := []testCases{
		{
			testName:    "invalid",
			videoId:     "",
			expectedErr: ErrInvalidArgument,
		},
		{
			testName:    "valid",
			videoId:     "video1",
			expectedErr: ErrInvalidArgument,
		},
	}

	for _, test := range tests {
		if test.videoId == "" {
			err := svc.Increment(context.Background(), test.videoId)
			assert.Error(t, err)
			assert.Equal(t, test.expectedErr, err)
		} else {
			mockRepo.EXPECT().Increment(context.Background(), test.videoId).Return(nil)

			err := svc.Increment(context.Background(), test.videoId)
			assert.NoError(t, err)
		}
	}
}

func TestGetTopVideos(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := viewrepository.NewMockRepository(ctrl)

	svc := NewService(mockRepo)

	n := 2

	expectedResult := []model.VideoInfo{
		{
			Id:    "video1",
			Views: 2,
		},
		{
			Id:    "video2",
			Views: 1,
		},
	}

	mockRepo.EXPECT().GetTopVideos(context.Background(), n).Return(expectedResult, nil)

	result, err := svc.GetTopVideos(context.Background(), n)

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestGetRecentVideos(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := viewrepository.NewMockRepository(ctrl)
	svc := NewService(mockRepo)

	n := 2

	expectedResult := []model.VideoInfo{
		{
			Id:    "video1",
			Views: 1,
		},
		{
			Id:    "video2",
			Views: 2,
		},
	}

	mockRepo.EXPECT().GetRecentVideos(context.Background(), n).Return(expectedResult, nil)

	result, err := svc.GetRecentVideos(context.Background(), n)

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
}
