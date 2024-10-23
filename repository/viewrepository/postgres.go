package viewrepository

import (
	"context"
	"database/sql"
	"view_count/model"
)

type postgresRepo struct {
	*sql.DB
}

func NewPostgresRepo(db *sql.DB) *postgresRepo {
	return &postgresRepo{
		DB: db,
	}
}

// TODO: write docker integration test cases. @Abhishek Gupta/Abhishek AK

func (db *postgresRepo) GetView(ctx context.Context, videoId string) (view int, err error) {
	// TODO impl : done
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	// TODO: check for err than rollback : done
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	row := tx.QueryRow("SELECT views FROM videos WHERE id = $1", videoId)
	err = row.Scan(&view)

	if err == sql.ErrNoRows {
		_, errs := tx.Exec("INSERT INTO videos (id, views) VALUES ($1, 0)", videoId)
		if errs != nil {
			return 0, err
		}
		view = 0
	} else if err != nil {
		return 0, err
	}

	return view, tx.Commit()
}

func (db *postgresRepo) GetAllViews(ctx context.Context) (info []model.VideoInfo, err error) {
	rows, err := db.Query("SELECT id, views FROM videos")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var video model.VideoInfo

		if err = rows.Scan(&video.Id, &video.Views); err != nil {
			return nil, err
		}
		info = append(info, video)
	}

	return info, nil
}

func (db *postgresRepo) Increment(ctx context.Context, videoId string) (err error) {
	_, err = db.Exec(`INSERT INTO videos (id, views, last_updated) VALUES ($1, 1, NOW()) ON CONFLICT (id) DO UPDATE SET views = videos.views + 1, last_updated = NOW()`, videoId)
	return err
}

func (db *postgresRepo) GetTopVideos(ctx context.Context, n int) (info []model.VideoInfo, err error) {
	rows, err := db.Query("SELECT id, views FROM videos ORDER BY views DESC LIMIT $1", n)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var video model.VideoInfo

		if err = rows.Scan(&video.Id, &video.Views); err != nil {
			return nil, err
		}
		info = append(info, video)
	}

	return info, nil
}

func (db *postgresRepo) GetRecentVideos(ctx context.Context, n int) (info []model.VideoInfo, err error) {
	rows, err := db.Query("SELECT id, views FROM videos ORDER BY last_updated DESC LIMIT $1", n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var video model.VideoInfo
		if err := rows.Scan(&video.Id, &video.Views); err != nil {
			return nil, err
		}
		info = append(info, video)
	}
	return info, nil
}
