package store

import (
	"database/sql"
	"sync/atomic"
)

type DbConnectionGetter interface {
	GetDb() (*sql.DB, error)
}

func NewDbConnectionGetter() DbConnectionGetter {
	return &SqliteConnectionImpl{}
}

type SqliteConnectionImpl struct {
	db *sql.DB

	/*
		when this is 0, we know its safe to close
      VS. depending on user of this package to call something like db.Close()
	*/

	referencesHeld int32
}

func (c *SqliteConnectionImpl) GetDb() (*sql.DB, error) {
	// @TODO
	atomic.AddInt32(&c.referencesHeld, 1) // we just gave someone a reference
	return nil, nil
}
