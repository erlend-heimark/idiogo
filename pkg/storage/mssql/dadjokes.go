package mssql

import (
	"context"
	"database/sql"
)

type DadJoke struct {
	ID   string `db:"id"`
	Joke string `db:"joke"`
}

const getDadJokeSQL = `select id, joke from dadjokes.dadjokes`

func (c Client) GetDadJoke(ctx context.Context) (*DadJoke, bool, error) {
	row := c.db.QueryRowContext(ctx, getDadJokeSQL)
	var d DadJoke
	err := row.Scan(&d)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, false, nil
		}
		return nil, false, err
	}
	return &d, true, nil
}

const insertDadJokeSQL = `insert into dadjokes.dadjokes (id, joke) values (@ID, @Joke)`

func (c Client) CreateDadJoke(ctx context.Context, joke DadJoke) error {
	_, err := c.db.ExecContext(ctx, insertDadJokeSQL, sql.Named("ID", joke.ID), sql.Named("Joke", joke.Joke))
	if err != nil {
		return err
	}
	return nil
}
