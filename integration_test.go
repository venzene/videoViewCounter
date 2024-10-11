package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	pool *dockertest.Pool
)

func TestMain(m *testing.M) {
	var err error
	pool, err = dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}
	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}
	code := m.Run()
	os.Exit(code)
}

func TestPostgres(t *testing.T) {
	resource, err := pool.Run("postgres", "latest", []string{
		"POSTGRES_USER=postgres",
		"POSTGRES_PASSWORD=secret",
		"POSTGRES_DB=view_count",
	})
	require.NoError(t, err)

	assert.NotEmpty(t, resource.GetPort("5432/tcp"))

	err = pool.Retry(func() error {
		conn := fmt.Sprintf("host=localhost port=%s user=postgres dbname=view_count password=secret sslmode=disable", resource.GetPort("5432/tcp"))
		db, err := sql.Open("postgres", conn)
		if err != nil {
			return err
		}
		return db.Ping()
	})
	require.NoError(t, err)

	require.NoError(t, pool.Purge(resource))
}
