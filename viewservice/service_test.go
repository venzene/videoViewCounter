package viewservice

import (
	"context"
	"testing"
	"view_count/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetViews(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)

	svc := NewService(mockRepo)

	videoId := "video1"
	mockRepo.EXPECT().GetView(context.Background(), videoId).Return(7, nil)

	views, err := svc.GetView(context.Background(), videoId)

	assert.NoError(t, err)
	assert.Equal(t, 7, views)
}

func TestGetAllViews(t *testing.T) {

}

func TestIncrement(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)

	svc := NewService(mockRepo)

	videoId := "video1"
	mockRepo.EXPECT().Increment(context.Background(), videoId).Return(nil)

	err := svc.Increment(context.Background(), videoId)
	assert.NoError(t, err)
}
