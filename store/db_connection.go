package store

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"rg/UnitTracker/utils/fsutils"

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
const migrationDir = "./store/migrations/"

var dbSingleton *sql.DB

type sqliteConnectionImpl struct {
}

func (c *sqliteConnectionImpl) GetDb() (*sql.DB, error) {
	if dbSingleton == nil {
		var err error
		dbSingleton, err = initDb()
		if err != nil {
			return nil, err
		}
	}
	return dbSingleton, nil
}

func initDb() (*sql.DB, error) {
   if !fsutils.FileExists(dbFilepath) {
      log.Printf("Db not found, creating %s...", dbFilepath)
   }
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
	return migrateDb(context.TODO(), db, migrationDir)
}

func migrateDb(ctx context.Context, db *sql.DB, migrationDir string) (*sql.DB, error) {
	if err := goose.RunContext(ctx, "status", db,migrationDir); err != nil {
		return nil, fmt.Errorf("failed to get goose status: %v", err)
	}

	if err := goose.RunContext(ctx, "up", db, migrationDir); err != nil {
		return nil, fmt.Errorf("failed to get goose up: %v", err)
	}
	return db, nil
}
