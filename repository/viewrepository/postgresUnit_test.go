package viewrepository

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"testing"
	"view_count/model"

	"github.com/DATA-DOG/go-sqlmock"
)

func Test_db_GetView(t *testing.T) {

	database, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	testRepo := NewPostgresRepo(database)

	defer database.Close()

	videoId := "video0"

	t.Run("New video, should insert and return 0 views", func(t *testing.T) {

		mock.ExpectBegin()

		mock.ExpectQuery("SELECT views FROM videos WHERE id = \\$1").
			WithArgs(videoId).
			WillReturnError(sql.ErrNoRows)

		mock.ExpectExec("INSERT INTO videos \\(id, views\\) VALUES \\(\\$1, 0\\)").
			WithArgs(videoId).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		result, err := testRepo.GetView(context.Background(), videoId)
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

		mock.ExpectBegin()

		mock.ExpectQuery("SELECT views FROM videos WHERE id = \\$1").
			WithArgs(videoId).
			WillReturnRows(sqlmock.NewRows([]string{"views"}).
				AddRow(2))

		mock.ExpectCommit()

		result, err := testRepo.GetView(context.Background(), videoId)
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

	t.Run("Error on transaction begin", func(t *testing.T) {

		mock.ExpectBegin().
			WillReturnError(fmt.Errorf("transaction begin error"))

		result, err := testRepo.GetView(context.Background(), videoId)
		if err == nil {
			t.Fatal("Expected error but got no error.")
		}
		if result != 0 {
			t.Fatalf("Expected 0 views, but got %v", result)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unmet expectations : %v", err)
		}
	})

	t.Run("Error scaning the row", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery("SELECT views FROM videos WHERE id = \\$1").
			WithArgs(videoId).
			WillReturnRows(sqlmock.NewRows([]string{"views"}).
				AddRow(nil))

		mock.ExpectRollback()

		result, err := testRepo.GetView(context.Background(), videoId)
		if err == nil {
			t.Fatal("Expected error but got no error.")
		}
		if result != 0 {
			t.Fatalf("Expected 0 views, but got %v", result)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unmet expectations : %v", err)
		}
	})

	t.Run("New Video, Error inserting the video", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery("SELECT views FROM videos WHERE id = \\$1").
			WithArgs(videoId).
			WillReturnError(sql.ErrNoRows)

		mock.ExpectExec("INSERT INTO videos \\(id, views\\) VALUES \\(\\$1, 0\\)").
			WithArgs(videoId).
			WillReturnError(fmt.Errorf("Error inserting the video"))

		mock.ExpectRollback()

		result, err := testRepo.GetView(context.Background(), videoId)
		if err == nil {
			t.Fatal("Expected error but got no error.")
		}
		if result != 0 {
			t.Fatalf("Expected 0 views, but got %v", result)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unmet expectations : %v", err)
		}
	})
}

func Test_db_GetAllViews(t *testing.T) {
	database, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}

	testRepo := NewPostgresRepo(database)

	videoRows := sqlmock.NewRows([]string{"id", "views"}).AddRow("video1", 2).AddRow("video2", 3)

	defer database.Close()

	t.Run("Get All Views", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, views FROM videos").
			WillReturnRows(videoRows)

		result, err := testRepo.GetAllViews(context.Background())
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
	})

	t.Run("Query throws an error", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, views FROM videos").
			WillReturnError(fmt.Errorf("Custom Error"))

		result, err := testRepo.GetAllViews(context.Background())
		if err == nil {
			t.Fatal("Expected error, but got none")
		}

		if result != nil {
			t.Fatalf("Expected nil, but got %+v", result)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unmet expectations : %v", err)
		}
	})

	t.Run("Error scanning a row", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, views FROM videos").
			WillReturnRows(sqlmock.NewRows([]string{"views"}).
				AddRow(nil))

		result, err := testRepo.GetAllViews(context.Background())
		if err == nil {
			t.Fatal("Expected error, but got none")
		}

		if result != nil {
			t.Fatalf("Expected nil, but got %+v", result)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unmet expectations : %v", err)
		}
	})
}

func Test_db_Increment(t *testing.T) {
	database, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to connect to the databse: %v", err)
	}

	testRepo := NewPostgresRepo(database)

	defer database.Close()

	videoId := "video1"
	// not able to match the query without these "\s*" and "(?i)"
	mock.ExpectExec(`(?i)INSERT INTO videos\s*\(id,\s*views,\s*last_updated\)\s*VALUES\s*\(\$1,\s*1,\s*NOW\(\)\)\s*ON CONFLICT\s*\(id\)\s*DO\s*UPDATE\s*SET\s*views\s*=\s*videos\.views\s*\+\s*1,\s*last_updated\s*=\s*NOW\(\)`).
		WithArgs(videoId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = testRepo.Increment(context.Background(), "video1")
	if err != nil {
		t.Fatalf("Unexpected error while incrementing: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unmet expectations : %v", err)
	}
}

func Test_db_GetTopVideos(t *testing.T) {
	database, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}

	testRepo := NewPostgresRepo(database)
	defer database.Close()

	n := 3

	mockRows := sqlmock.NewRows([]string{"id", "views"}).
		AddRow("video1", 20).
		AddRow("video2", 15).
		AddRow("video3", 10)

	expectedVideos := []model.VideoInfo{
		{Id: "video1", Views: 20},
		{Id: "video2", Views: 15},
		{Id: "video3", Views: 10},
	}

	t.Run("Get top videos", func(t *testing.T) {
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

		for i, video := range videos {
			if video.Id != expectedVideos[i].Id || video.Views != expectedVideos[i].Views {
				t.Errorf("Expected video %+v, got %+v", expectedVideos[i], video)
			}
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unmet expectations: %v", err)
		}
	})

	t.Run("Error executing the query", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, views FROM videos ORDER BY views DESC LIMIT \\$1").
			WithArgs(n).
			WillReturnError(fmt.Errorf("custom error"))

		videos, err := testRepo.GetTopVideos(context.Background(), n)
		if err == nil {
			t.Fatal("Expected error, but got one", err)
		}

		if videos != nil {
			t.Fatalf("expected nil, got %v", videos)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unmet expectations: %v", err)
		}
	})

	t.Run("Error scanning the rows", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, views FROM videos ORDER BY views DESC LIMIT \\$1").
			WithArgs(n).
			WillReturnRows(sqlmock.NewRows([]string{"views"}).
				AddRow(nil))

		videos, err := testRepo.GetTopVideos(context.Background(), n)
		if err == nil {
			t.Fatal("Expected error, but got one", err)
		}

		if videos != nil {
			t.Fatalf("expected nil, got %v", videos)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unmet expectations: %v", err)
		}
	})

}

func Test_db_GetRecentVideos(t *testing.T) {
	database, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	defer database.Close()

	testRepo := NewPostgresRepo(database)

	n := 3

	mockRows := sqlmock.NewRows([]string{"id", "views"}).
		AddRow("video2", 3).
		AddRow("video1", 5).
		AddRow("video3", 10)

	expected := []model.VideoInfo{
		{Id: "video2", Views: 3},
		{Id: "video1", Views: 5},
		{Id: "video3", Views: 10},
	}

	t.Run("Get recent videos", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, views FROM videos ORDER BY last_updated DESC LIMIT \\$1").
			WithArgs(n).
			WillReturnRows(mockRows)

		result, err := testRepo.GetRecentVideos(context.Background(), n)
		if err != nil {
			t.Fatalf("Unexpected error while getting recent videos: %v", err)
		}

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, but got %v", expected, result)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unmet expectations: %v", err)
		}
	})

	t.Run("Error executing the query", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, views FROM videos ORDER BY last_updated DESC LIMIT \\$1").
			WithArgs(n).
			WillReturnError(fmt.Errorf("custom error"))

		result, err := testRepo.GetRecentVideos(context.Background(), n)
		if err == nil {
			t.Fatalf("Expected error but got none")
		}

		if result != nil {
			t.Fatalf("Expected nil, but got %v", result)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unmet expectations: %v", err)
		}
	})

	t.Run("Error scanning the rows", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, views FROM videos ORDER BY last_updated DESC LIMIT \\$1").
			WithArgs(n).
			WillReturnRows(sqlmock.NewRows([]string{"views"}).AddRow(nil))

		result, err := testRepo.GetRecentVideos(context.Background(), n)
		if err == nil {
			t.Fatalf("Expected error but got none")
		}

		if result != nil {
			t.Fatalf("Expected nil, but got %v", result)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("There were unmet expectations: %v", err)
		}
	})
}
