package viewrepository

import (
	"context"
	"reflect"
	"testing"
	"view_count/model"
)

type testCase struct {
	testName       string
	vid            string
	expectedViews  int
	nParams        int
	testInput      []model.VideoInfo
	expectedResult []model.VideoInfo
	expectedErr    error
}

func Test_IM_GetView(t *testing.T) {

	tests := []testCase{
		{
			testName:      "Get Views",
			vid:           "video1",
			expectedViews: 3,
			expectedErr:   nil,
		},
		{
			testName:      "Get Views",
			vid:           "video2",
			expectedViews: 1,
			expectedErr:   nil,
		},
		{
			testName:      "Get Views",
			vid:           "video3",
			expectedViews: 10,
			expectedErr:   nil,
		},
		{
			testName:      "Get Views",
			vid:           "video4",
			expectedViews: 0,
			expectedErr:   nil,
		},
		// TODO how to write multiple testcases for above : Done
	}

	for _, test := range tests {

		testRepo := NewInmemoryRepo()

		for i := 0; i < test.expectedViews; i++ {
			testRepo.Increment(context.Background(), test.vid)
		}

		t.Run(test.testName, func(t *testing.T) {
			result, err := testRepo.GetView(context.Background(), test.vid)
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

	tests := []testCase{
		{
			testName: "Get all videos",
			expectedResult: []model.VideoInfo{
				{Id: "video1", Views: 3},
				{Id: "video2", Views: 1},
			},
			expectedErr: nil,
		},
		{
			testName: "Get all videos",
			expectedResult: []model.VideoInfo{
				{Id: "video1", Views: 5},
				{Id: "video2", Views: 7},
				{Id: "video3", Views: 9},
			},
			expectedErr: nil,
		},
	}

	for _, test := range tests {
		testRepo := NewInmemoryRepo()

		for i := 0; i < len(test.expectedResult); i++ {
			for j := 1; j <= test.expectedResult[i].Views; j++ {
				testRepo.Increment(context.Background(), test.expectedResult[i].Id)

			}
		}

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

	case1 := "video does not exists"
	case2 := "video already exists"

	tests := []testCase{
		{
			vid:           "video1",
			testName:      case1,
			expectedViews: 1,
			expectedErr:   nil,
		},
		{
			vid:           "video2",
			testName:      case2,
			expectedViews: 4,
			expectedErr:   nil,
		},
	}

	for _, test := range tests {

		testRepo := NewInmemoryRepo()

		if test.testName == case1 {
			t.Run(test.testName, func(t *testing.T) {
				err := testRepo.Increment(context.Background(), test.vid)
				if err != test.expectedErr {
					t.Fatalf("Expected %v, got %v", test.expectedErr, err)
				}
				result, _ := testRepo.GetView(context.Background(), test.vid)
				if result != test.expectedViews {
					t.Fatalf("Expected %v, got %v", test.expectedViews, result)
				}
			})
		} else {
			t.Run(test.testName, func(t *testing.T) {
				for i := 0; i < test.expectedViews; i++ {
					err := testRepo.Increment(context.Background(), test.vid)
					if err != test.expectedErr {
						t.Fatalf("Expected %v, got %v", test.expectedErr, err)
					}
				}
				result, _ := testRepo.GetView(context.Background(), test.vid)
				if result != test.expectedViews {
					t.Fatalf("Expected %v, got %v", test.expectedViews, result)
				}
			})
		}
	}
}

func Test_IM_GetTopVideos(t *testing.T) {

	tests := []testCase{
		{
			testName: "Get all videos",
			nParams:  2,
			expectedResult: []model.VideoInfo{
				{Id: "video3", Views: 3},
				{Id: "video2", Views: 2},
				{Id: "video1", Views: 1},
			},
			expectedErr: nil,
		},
		{
			testName: "Get all videos",
			nParams:  3,
			expectedResult: []model.VideoInfo{
				{Id: "video3", Views: 10},
				{Id: "video1", Views: 3},
			},
			expectedErr: nil,
		},
	}

	for _, test := range tests {

		testRepo := NewInmemoryRepo()

		for i := 0; i < len(test.expectedResult); i++ {
			for j := 0; j < test.expectedResult[i].Views; j++ {
				testRepo.Increment(context.Background(), test.expectedResult[i].Id)
			}
		}

		t.Run(test.testName, func(t *testing.T) {
			result, err := testRepo.GetTopVideos(context.Background(), test.nParams)

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

	tests := []testCase{
		{
			testName: "Get all videos",
			nParams:  4,
			expectedResult: []model.VideoInfo{
				{Id: "video2", Views: 2},
				{Id: "video4", Views: 1},
				{Id: "video3", Views: 1},
				{Id: "video1", Views: 1},
			},
			expectedErr: nil,
		},
		{
			testName: "Get all videos",
			nParams:  10,
			expectedResult: []model.VideoInfo{
				{Id: "video1", Views: 1},
				{Id: "video2", Views: 1},
				{Id: "video3", Views: 1},
				{Id: "video4", Views: 1},
			},
			expectedErr: nil,
		},
	}

	for _, test := range tests {
		testRepo := NewInmemoryRepo()

		for i := len(test.expectedResult) - 1; i >= 0; i-- {
			for j := 0; j < test.expectedResult[i].Views; j++ {
				testRepo.Increment(context.Background(), test.expectedResult[i].Id)
			}
		}

		t.Run(test.testName, func(t *testing.T) {
			result, err := testRepo.GetRecentVideos(context.Background(), test.nParams)

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
