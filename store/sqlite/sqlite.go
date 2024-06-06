package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"rg/UnitTracker/queries"
	"rg/UnitTracker/store"
	"rg/UnitTracker/utils/fsutils"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

const testDbFilepath = "./_testSqlite.db"
const testMigrationDir = "../migrations/"

const mainDbFilepath = "./store/sqlite/_sqlite.db"
const mainMigrationDir = "./store/migrations/"

type sqliteConnector struct {
	db      *sql.DB // never access this directly
	queries *queries.Queries
}

type testSqliteConnector struct {
	db      *sql.DB // never access this directly
	queries *queries.Queries
	store.Store
}

func NewSqliteConnector() store.Connecter {
	return &sqliteConnector{}
}

func NewTestSqliteConnector() store.Connecter {
	return &testSqliteConnector{}
}

func (c *testSqliteConnector) Connect(ctx context.Context) (store.Store, error) {
	if c.db == nil {
		var err error
		c.db, c.queries, err = initDb(ctx, testDbFilepath, testMigrationDir)
		if err != nil {
			return nil, err
		}
	}
	return c, nil
}

func (c *sqliteConnector) Connect(ctx context.Context) (store.Store, error) {
	if c.db == nil {
		var err error
		c.db, c.queries, err = initDb(ctx, mainDbFilepath, mainMigrationDir)
		if err != nil {
			return nil, err
		}
	}
	return c, nil
}

func (c *sqliteConnector) ProjectStore() store.ProjectStore {
	return &projectDb{db: c.db, queries: c.queries}
}

func (c *sqliteConnector) UnitStore() store.UnitStore {
	return &unitDb{db: c.db, queries: c.queries}
}

// =========================== UTIL FUNCS ===============================
func initDb(ctx context.Context, dbFile string, migrationDir string) (*sql.DB, *queries.Queries, error) {
	if !fsutils.FileExists(dbFile) {
		log.Printf("Db not found, creating %s...", dbFile)
	}
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open database: %v", err)
	}

	queries := queries.New(db)

	if err = db.Ping(); err != nil {
		return nil, nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	if err = goose.SetDialect("sqlite3"); err != nil {
		return nil, nil, fmt.Errorf("failed to set goose dialect: %v", err)
	}
	// run migrations
	db, err = migrateDb(ctx, db, migrationDir)
	if err != nil {
		dir, err := os.Getwd()
		if err != nil {
			return nil, nil, fmt.Errorf("Failed to run dir")
		}
		fmt.Println("printing from directory: ", dir)
		return nil, nil, fmt.Errorf("Failed to run migrations: %v", err)
	}
	return db, queries, err
}

func migrateDb(ctx context.Context, db *sql.DB, migrationDir string) (*sql.DB, error) {
	if err := goose.RunContext(ctx, "status", db, migrationDir); err != nil {
		return nil, fmt.Errorf("failed to get goose status: %v", err)
	}

	if err := goose.RunContext(ctx, "up", db, migrationDir); err != nil {
		return nil, fmt.Errorf("failed to get goose up: %v", err)
	}
	return db, nil
}
