package main

import (
	"context"
	"database/sql"

	"cos316.princeton.edu/assignment5/conn_pool"

	_ "github.com/mattn/go-sqlite3"
)

type SQLConn struct {
	c *sql.Conn
}

func (s *SQLConn) Query(query string) (*sql.Rows, error) {
	ctx, _ := context.WithCancel(context.Background())
	return s.c.QueryContext(ctx, query)
}

type ConnFunc = func() (conn_pool.Conn, error)

func Connector() ConnFunc {
	driver, err := sql.Open("sqlite3", "file:MovieLens.db")
	return (func() (conn_pool.Conn, error) {
		if err != nil {
			return nil, err
		}

		ctx, _ := context.WithCancel(context.Background())

		conn, err := driver.Conn(ctx)
		if err != nil {
			return nil, err
		}

		return &SQLConn{c: conn}, nil
	})

}
