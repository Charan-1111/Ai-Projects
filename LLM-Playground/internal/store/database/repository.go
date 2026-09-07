package database

import "context"

type Repository interface {
	Create(ctx context.Context) error
}

func (db *DataBaseStore) Create(ctx context.Context) error {
	for tableName, createQuery := range db.Queries.Create {
		_, err := db.Db.Exec(ctx, createQuery)
		if err != nil {
			return err
		} else {
			db.Log.Log.Info().Msgf("Table %s created successfully", tableName)
		}
	}

	return nil
}
