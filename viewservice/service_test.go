package viewservice

import (
	"context"
	"database/sql"
	"reflect"
	"testing"
	"view_count/model"
	"view_count/repository/viewrepository"

	"github.com/DATA-DOG/go-sqlmock"
)

type testCase struct {
	testName       string
	expectedViews  int
	expectedResult []model.VideoInfo
	expectedErr    error
}

func Test_IM_GetView(t *testing.T) {
	testRepo := viewrepository.NewInmemoryRepo()

	svc := NewService(testRepo)

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
			result, err := svc.GetView(context.Background(), "video1")
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

	testRepo := viewrepository.NewInmemoryRepo()
	svc := NewService(testRepo)

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
			result, err := svc.GetAllViews(context.Background())

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
	testRepo := viewrepository.NewInmemoryRepo()
	svc := NewService(testRepo)

	tests := []testCase{
		{
			testName:    "Increment",
			expectedErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			err := svc.Increment(context.Background(), "video1")
			if err != test.expectedErr {
				t.Fatalf("Expected %v, got %v", test.expectedErr, err)
			}
		})
	}

}

func Test_IM_GetTopVideos(t *testing.T) {
	testRepo := viewrepository.NewInmemoryRepo()
	svc := NewService(testRepo)

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
			result, err := svc.GetTopVideos(context.Background(), 3)

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
	testRepo := viewrepository.NewInmemoryRepo()
	svc := NewService(testRepo)

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
			result, err := svc.GetRecentVideos(context.Background(), 4)

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

// Database testing code

func Test_DB_GetView(t *testing.T) {

	database, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	testRepo := viewrepository.NewPostgresRepo(database)

	svc := NewService(testRepo)

	_ = mock.NewRows([]string{"id", "views"}).AddRow("video1", 1).AddRow("video2", 2)

	defer database.Close()

	t.Run("New video, should insert and return 0 views", func(t *testing.T) {
		videoID := "video0"

		mock.ExpectBegin()

		mock.ExpectQuery("SELECT views FROM videos WHERE id = \\$1").
			WithArgs(videoID).
			WillReturnError(sql.ErrNoRows)

		mock.ExpectExec("INSERT INTO videos \\(id, views\\) VALUES \\(\\$1, 0\\)").
			WithArgs(videoID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		result, err := svc.GetView(context.Background(), videoID)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if result != 0 {
			t.Errorf("Expected 0, but got %v", result)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unmet expectations: %v", err)
		}
	})

	t.Run("Existing Video, should return its view count", func(t *testing.T) {
		videoId := "video1"

		mock.ExpectBegin()

		mock.ExpectQuery("SELECT views FROM videos WHERE id = \\$1").
			WithArgs(videoId).
			WillReturnRows(sqlmock.NewRows([]string{"views"}).AddRow(2))

		mock.ExpectCommit()

		result, err := svc.GetView(context.Background(), videoId)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if result != 2 {
			t.Errorf("Expected 2, but got %v", result)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unmet expectations : %v", err)
		}
	})

}

func Test_DB_GetAllViews(t *testing.T) {
	database, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}

	testRepo := viewrepository.NewPostgresRepo(database)
	svc := NewService(testRepo)

	videoRows := sqlmock.NewRows([]string{"id", "views"}).AddRow("video1", 2).AddRow("video2", 3)

	defer database.Close()

	mock.ExpectQuery("SELECT id, views FROM videos").
		WillReturnRows(videoRows)

	result, err := svc.GetAllViews(context.Background())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedResult := []model.VideoInfo{
		{Id: "video1", Views: 2},
		{Id: "video2", Views: 3},
	}

	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("Expected %+v, but got %+v", expectedResult, result)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unmet expectations : %v", err)
	}
}

func Test_DB_Increment(t *testing.T) {
	database, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to connect to the databse: %v", err)
	}

	testRepo := viewrepository.NewPostgresRepo(database)

	svc := NewService(testRepo)

	defer database.Close()

	videoId := "video1"
	// not able to match the query without these "\s*" and "(?i)"
	mock.ExpectExec(`(?i)INSERT INTO videos\s*\(id,\s*views,\s*last_updated\)\s*VALUES\s*\(\$1,\s*1,\s*NOW\(\)\)\s*ON CONFLICT\s*\(id\)\s*DO\s*UPDATE\s*SET\s*views\s*=\s*videos\.views\s*\+\s*1,\s*last_updated\s*=\s*NOW\(\)`).
		WithArgs(videoId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = svc.Increment(context.Background(), "video1")
	if err != nil {
		t.Fatalf("Unexpected error while incrementing: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unmet expectations : %v", err)
	}
}

func Test_DB_GetTopVideos(t *testing.T) {
	database, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}

	testRepo := viewrepository.NewPostgresRepo(database)
	defer database.Close()

	n := 3

	mockRows := sqlmock.NewRows([]string{"id", "views"}).
		AddRow("video1", 20).
		AddRow("video2", 15).
		AddRow("video3", 10)

	mock.ExpectQuery("SELECT id, views FROM videos ORDER BY views DESC LIMIT \\$1").
		WithArgs(n).
		WillReturnRows(mockRows)

	videos, err := testRepo.GetTopVideos(context.Background(), n)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(videos) != 3 {
		t.Errorf("Expected 3 videos, got %d", len(videos))
	}

	expectedVideos := []model.VideoInfo{
		{Id: "video1", Views: 20},
		{Id: "video2", Views: 15},
		{Id: "video3", Views: 10},
	}

	for i, video := range videos {
		if video.Id != expectedVideos[i].Id || video.Views != expectedVideos[i].Views {
			t.Errorf("Expected video %+v, got %+v", expectedVideos[i], video)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unmet expectations: %v", err)
	}
}

func Test_DB_GetRecentVideos(t *testing.T) {
	database, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	defer database.Close()

	testRepo := viewrepository.NewPostgresRepo(database)
	svc := NewService(testRepo)

	n := 3

	mockRows := sqlmock.NewRows([]string{"id", "views"}).
		AddRow("video2", 3).
		AddRow("video1", 5).
		AddRow("video3", 10)

	mock.ExpectQuery("SELECT id, views FROM videos ORDER BY last_updated DESC LIMIT \\$1").
		WithArgs(n).
		WillReturnRows(mockRows)

	result, err := svc.GetRecentVideos(context.Background(), n)
	if err != nil {
		t.Fatalf("Unexpected error while getting recent videos: %v", err)
	}

	expected := []model.VideoInfo{
		{Id: "video2", Views: 3},
		{Id: "video1", Views: 5},
		{Id: "video3", Views: 10},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unmet expectations: %v", err)
	}
}
