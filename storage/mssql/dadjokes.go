package mssql

import (
	"context"
	"database/sql"
)
type DadJoke struct {
	ID   string `db:"id"`
	Joke string `db:"joke"`
}

const getDadJokeSQL = `select d.id, d.joke from dadjokes.dadjokes d where id = @ID`

func (c Client) GetDadJoke(ctx context.Context, id string) (*DadJoke, bool, error) {
	row := c.db.QueryRowContext(ctx, getDadJokeSQL, sql.Named("ID", id))
	var d DadJoke
	err := row.Scan(&d.ID, &d.Joke)
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
