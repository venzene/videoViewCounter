package viewrepository

// TODO add integration build tag

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"
	"view_count/model"

	"github.com/docker/go-connections/nat"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var testSqlDB *sql.DB

func cleanupDB(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM videos;")
	return err
}

func createContainer() (testcontainers.Container, *sql.DB, error) {
	port := "5432/tcp"
	info := map[string]string{
		"POSTGRES_USER":     "postgres",
		"POSTGRES_PASSWORD": "password",
		"POSTGRES_DB":       "test",
	}
	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:latest",
			ExposedPorts: []string{port},
			Env:          info,
			WaitingFor: wait.ForSQL(nat.Port(port), "postgres", func(host string, port nat.Port) string {
				return fmt.Sprintf("host=localhost port=%s user=postgres password=password dbname=test sslmode=disable", port.Port())
			}),
		},
		Started: true,
	}
	container, err := testcontainers.GenericContainer(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}

	mappedPort, err := container.MappedPort(context.Background(), nat.Port(port))
	if err != nil {
		log.Fatal(err)
	}

	url := fmt.Sprintf("host=localhost port=%s user=postgres password=password dbname=test sslmode=disable", mappedPort.Port())
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, nil, err
	}

	return container, db, nil
}

func TestMain(m *testing.M) {
	container, testDB, err := createContainer()
	if err != nil {
		log.Fatal(err)
	}
	defer testDB.Close()
	defer container.Terminate(context.Background())

	// Setting up the database schema
	_, err = testDB.Exec(`CREATE TABLE IF NOT EXISTS videos (
            id TEXT PRIMARY KEY,
            views INT NOT NULL,
            last_updated TIMESTAMP NOT NULL
        );`)
	if err != nil {
		log.Fatal(err)
	}

	testSqlDB = testDB
	// TODO findout why using Exit?
	os.Exit(m.Run())
}

func Test_DB_GetView(t *testing.T) {

	tests := []testCase{
		{
			testName:      "Get Views",
			vid:           "video1",
			expectedViews: 3,
			expectedErr:   nil,
		},
		{
			testName:      "Get Views",
			vid:           "video3",
			expectedViews: 1,
			expectedErr:   nil,
		},
		{
			testName:      "Get Views",
			vid:           "video3",
			expectedViews: 10,
			expectedErr:   nil,
		},
	}

	for _, test := range tests {

		testRepo := NewPostgresRepo(testSqlDB)

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
		if err := cleanupDB(testSqlDB); err != nil {
			t.Fatalf("Error cleaning up database: %v", err)
		}
	}

}

func Test_DB_GetAllViews(t *testing.T) {

	tests := []testCase{
		{
			testName: "Get all videos",
			testInput: []model.VideoInfo{
				{Id: "video2", Views: 1},
				{Id: "video1", Views: 2},
				{Id: "video1", Views: 1},
			},
			expectedResult: []model.VideoInfo{
				{Id: "video1", Views: 3},
				{Id: "video2", Views: 1},
			},
			expectedErr: nil,
		},
		{
			testName: "Get all videos",
			testInput: []model.VideoInfo{
				{Id: "video1", Views: 3},
				{Id: "video2", Views: 4},
				{Id: "video1", Views: 2},
				{Id: "video2", Views: 3},
				{Id: "video3", Views: 9},
			},
			expectedResult: []model.VideoInfo{
				{Id: "video1", Views: 5},
				{Id: "video2", Views: 7},
				{Id: "video3", Views: 9},
			},
			expectedErr: nil,
		},
	}

	for _, test := range tests {
		testRepo := NewPostgresRepo(testSqlDB)

		for i := 0; i < len(test.testInput); i++ {
			for j := 1; j <= test.testInput[i].Views; j++ {
				testRepo.Increment(context.Background(), test.testInput[i].Id)
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
		if err := cleanupDB(testSqlDB); err != nil {
			t.Fatalf("Error cleaning up database: %v", err)
		}
	}
}

func Test_DB_Increment(t *testing.T) {

	tests := []testCase{
		{
			testName:    "Increment",
			expectedErr: nil,
		},
	}

	for _, test := range tests {

		testRepo := NewPostgresRepo(testSqlDB)

		t.Run(test.testName, func(t *testing.T) {
			err := testRepo.Increment(context.Background(), "video1")
			if err != test.expectedErr {
				t.Fatalf("Expected %v, got %v", test.expectedErr, err)
			}
		})
		if err := cleanupDB(testSqlDB); err != nil {
			t.Fatalf("Error cleaning up database: %v", err)
		}
	}

}

func Test_DB_GetTopVideos(t *testing.T) {

	tests := []testCase{
		{
			testName: "Get all videos",
			testInput: []model.VideoInfo{
				{Id: "video3", Views: 2},
				{Id: "video2", Views: 2},
				{Id: "video3", Views: 1},
				{Id: "video1", Views: 1},
			},
			expectedResult: []model.VideoInfo{
				{Id: "video3", Views: 3},
				{Id: "video2", Views: 2},
				{Id: "video1", Views: 1},
			},
			expectedErr: nil,
		},
		{
			testName: "Get all videos",
			testInput: []model.VideoInfo{
				{Id: "video3", Views: 6},
				{Id: "video1", Views: 2},
				{Id: "video3", Views: 4},
				{Id: "video1", Views: 1},
			},
			expectedResult: []model.VideoInfo{
				{Id: "video3", Views: 10},
				{Id: "video1", Views: 3},
			},
			expectedErr: nil,
		},
	}

	for _, test := range tests {

		testRepo := NewPostgresRepo(testSqlDB)

		for i := 0; i < len(test.testInput); i++ {
			for j := 1; j <= test.testInput[i].Views; j++ {
				testRepo.Increment(context.Background(), test.testInput[i].Id)
			}
		}

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
		if err := cleanupDB(testSqlDB); err != nil {
			t.Fatalf("Error cleaning up database: %v", err)
		}
	}
}

func Test_DB_GetRecentVideos(t *testing.T) {

	tests := []testCase{
		{
			testName: "Get all videos",
			testInput: []model.VideoInfo{
				{Id: "video1", Views: 1},
				{Id: "video3", Views: 1},
				{Id: "video4", Views: 1},
				{Id: "video2", Views: 2},
			},
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
			testInput: []model.VideoInfo{
				{Id: "video4", Views: 1},
				{Id: "video3", Views: 1},
				{Id: "video2", Views: 1},
				{Id: "video1", Views: 1},
			},
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
		testRepo := NewPostgresRepo(testSqlDB)

		for i := 0; i < len(test.testInput); i++ {
			for j := 1; j <= test.testInput[i].Views; j++ {
				testRepo.Increment(context.Background(), test.testInput[i].Id)
			}
		}

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
			if err := cleanupDB(testSqlDB); err != nil {
				t.Fatalf("Error cleaning up database: %v", err)
			}
		})
	}
}
