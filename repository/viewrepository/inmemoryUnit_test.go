package viewrepository

import (
	"context"
	"reflect"
	"testing"
	"view_count/model"
)

type testCase struct {
	testName       string
	expectedViews  int
	expectedResult []model.VideoInfo
	expectedErr    error
}

func Test_IM_GetView(t *testing.T) {
	testRepo := NewInmemoryRepo()

	_ = testRepo.Increment(context.Background(), "video1")
	_ = testRepo.Increment(context.Background(), "video1")
	_ = testRepo.Increment(context.Background(), "video1")

	tests := []testCase{
		{
			testName:      "Get Views",
			expectedViews: 3,
			expectedErr:   nil,
		},
		// TODO how to write multiple testcases for above
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			result, err := testRepo.GetView(context.Background(), "video1")
			if err != test.expectedErr {
				t.Errorf("Expected %v, got %v", test.expectedErr, err)
			}

			if result != test.expectedViews {
				t.Fatalf("Expected %v, got %v", test.expectedViews, result)
			}

		})
	}

}

func Test_IM_GetAllViews(t *testing.T) {

	testRepo := NewInmemoryRepo()

	_ = testRepo.Increment(context.Background(), "video1")
	_ = testRepo.Increment(context.Background(), "video1")
	_ = testRepo.Increment(context.Background(), "video1")
	_ = testRepo.Increment(context.Background(), "video2")

	tests := []testCase{
		{
			testName: "Get all videos",
			expectedResult: []model.VideoInfo{
				{Id: "video1", Views: 3},
				{Id: "video2", Views: 1},
			},
			expectedErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			result, err := testRepo.GetAllViews(context.Background())

			if err != test.expectedErr {
				t.Fatalf("Expected error %v, got %v", test.expectedErr, err)
			}

			resultMap := make(map[string]int)
			for _, video := range result {
				resultMap[video.Id] = video.Views
			}

			expectedMap := make(map[string]int)
			for _, video := range test.expectedResult {
				expectedMap[video.Id] = video.Views
			}

			if !reflect.DeepEqual(resultMap, expectedMap) {
				t.Fatalf("Expected %v, but got %v", expectedMap, resultMap)
			}
		})
	}
}

func Test_IM_Increment(t *testing.T) {
	testRepo := NewInmemoryRepo()

	tests := []testCase{
		{
			testName:    "Increment",
			expectedErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			err := testRepo.Increment(context.Background(), "video1")
			if err != test.expectedErr {
				t.Fatalf("Expected %v, got %v", test.expectedErr, err)
			}
		})
	}

}

func Test_IM_GetTopVideos(t *testing.T) {
	testRepo := NewInmemoryRepo()

	_ = testRepo.Increment(context.Background(), "video1")
	_ = testRepo.Increment(context.Background(), "video2")
	_ = testRepo.Increment(context.Background(), "video3")
	_ = testRepo.Increment(context.Background(), "video3")
	_ = testRepo.Increment(context.Background(), "video3")

	tests := []testCase{
		{
			testName: "Get all videos",
			expectedResult: []model.VideoInfo{
				{Id: "video3", Views: 3},
				{Id: "video2", Views: 1},
				{Id: "video1", Views: 1},
			},
			expectedErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			result, err := testRepo.GetTopVideos(context.Background(), 3)

			if err != test.expectedErr {
				t.Fatalf("Expected error %v, got %v", test.expectedErr, err)
			}

			for i, v := range result {
				if v.Id != test.expectedResult[i].Id || v.Views != test.expectedResult[i].Views {
					t.Fatalf("Expected result %v, got %v", test.expectedResult[i], v)
				}
			}
		})
	}
}

func Test_IM_GetRecentVideos(t *testing.T) {
	testRepo := NewInmemoryRepo()

	_ = testRepo.Increment(context.Background(), "video1")
	_ = testRepo.Increment(context.Background(), "video2")
	_ = testRepo.Increment(context.Background(), "video3")
	_ = testRepo.Increment(context.Background(), "video4")
	_ = testRepo.Increment(context.Background(), "video2")

	tests := []testCase{
		{
			testName: "Get all videos",
			expectedResult: []model.VideoInfo{
				{Id: "video2", Views: 2},
				{Id: "video4", Views: 1},
				{Id: "video3", Views: 1},
				{Id: "video1", Views: 1},
			},
			expectedErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			result, err := testRepo.GetRecentVideos(context.Background(), 4)

			if err != test.expectedErr {
				t.Fatalf("Expected error %v, got %v", test.expectedErr, err)
			}

			for i, v := range result {
				if v.Id != test.expectedResult[i].Id || v.Views != test.expectedResult[i].Views {
					t.Fatalf("Expected result %v, got %v", test.expectedResult[i], v)
				}
			}
		})
	}
}
