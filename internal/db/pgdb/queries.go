package pgdb

import (
	"context"
	"fmt"

	"github.com/labstack/gommon/log"
)

func (handler *PostgresQuery) InsertPolls(ctx context.Context, data []Poll) (int64, error) {

	res, err := handler.Pgdb.NewInsert().Model(&data).Exec(ctx)

	if err != nil {
		log.Error(fmt.Sprintf("pgdb::InsertOne err %v", err))
		return 0, err
	}

	rows_count, _ := res.RowsAffected()

	return rows_count, nil
}

func (handler *PostgresQuery) InsertComments(ctx context.Context, data []Comment) (int64, error) {

	res, err := handler.Pgdb.NewInsert().Model(&data).Exec(ctx)

	if err != nil {
		log.Error(fmt.Sprintf("pgdb::InsertOne err %v", err))
		return 0, err
	}

	rows_count, _ := res.RowsAffected()
	return rows_count, nil
}
