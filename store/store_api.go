package store

import (
	"database/sql"
	"sync/atomic"
  "sync"
)

type DbConnectionGetter interface {
	GetDb() (*sql.DB, error)
}

var dbSingletonInstance *SqliteConnectionImpl

func NewDbConnectionGetter() DbConnectionGetter {
  if dbSingletonInstance==nil{
    dbSingletonInstance=&SqliteConnectionImpl{}
    dbSingletonInstance.db,_=sql.Open("sqlite3", "./unitTracker.db")
    go dbSingletonInstance.autoCloseDb()
  }
	return dbSingletonInstance 
}

type SqliteConnectionImpl struct {
	db *sql.DB

	/*
				when this is 0, we know its safe to close
		      VS. depending on user of this package to call something like db.Close()
	*/

	connectionCount int32
  wg sync.WaitGroup
}

func (c *SqliteConnectionImpl ) autoCloseDb(){
  c.wg.Wait()
  if atomic.LoadInt32(&c.connectionCount)==0{
    c.db.Close()
  }
}

func (c *SqliteConnectionImpl) CloseDb(){
  c.wg.Done()
  atomic.AddInt32(&c.connectionCount, -1)
  c.db.Close()
}

func (c *SqliteConnectionImpl) GetDb() (*sql.DB, error) {
  c.wg.Add(1)
	atomic.AddInt32(&c.connectionCount, 1) // we just gave someone a reference
	return c.db, nil
}
