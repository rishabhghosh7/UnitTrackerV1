package store

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

type DbConnectionGetter interface {
	GetDb() (*sql.DB, error)
}

func NewDbConnectionGetter() DbConnectionGetter {
   return &sqliteConnectionImpl{}
}

const dbFilepath = "./store/sqlite.db"
var dbSingleton *sql.DB

type sqliteConnectionImpl struct {
}

func (c *sqliteConnectionImpl) GetDb() (*sql.DB, error) {
   if dbSingleton == nil {
      var err error
      dbSingleton, err = c.initDb()
      if err != nil {
         return nil,  err
      }
   }
   return dbSingleton, nil
}

func (c *sqliteConnectionImpl) initDb() (*sql.DB, error) {
   db, err := sql.Open("sqlite3", dbFilepath)
   if err != nil {
      return nil, fmt.Errorf("failed to open database: %v", err)
   }


   if err = db.Ping(); err != nil {
      return nil, fmt.Errorf("failed to connect to database: %v", err)
   }

   if err = goose.SetDialect("sqlite3"); err != nil {
      return nil, fmt.Errorf("failed to set goose dialect: %v", err)
   }

   // run migrations
   ctx := context.TODO()
   if err := goose.RunContext(ctx, "status", db, "store/migrations"); err != nil {
      return nil, fmt.Errorf("failed to get goose status: %v", err)
   }

   if err := goose.RunContext(ctx, "up", db, "store/migrations"); err != nil {
      return nil, fmt.Errorf("failed to get goose up: %v", err)
   }

   return db, nil
}

