package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"rg/UnitTracker/pkg/proto"
	"rg/UnitTracker/store"
	"rg/UnitTracker/utils/fsutils"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const dbFilepath = "./store/sqlite/_sqlite.db"
const migrationDir = "./store/migrations/"

var dbSingleton *sql.DB

type sqliteConnector struct {
	db *sql.DB // never access this directly
}

func NewSqliteConnector() store.Connecter {
	return &sqliteConnector{}
}

// @TODO
// func RunTransaction(store, func(trancsaction) {}) error

func (c *sqliteConnector) Connect(ctx context.Context) (store.Store, error) {
	if dbSingleton == nil {
		var err error
		dbSingleton, err = initDb(ctx)
		if err != nil {
			return nil, err
		}
	}
	return &sqliteConnector{db: dbSingleton}, nil
}

func (c *sqliteConnector) ProjectStore() store.ProjectStore {
	return &projectDb{db: c.db}
}

type projectDb struct {
	db *sql.DB
}

func (c *sqliteConnector) UnitStore() store.UnitStore {
	return &unitDb{db: c.db}
}

type unitDb struct {
	db *sql.DB
}

// ===================== UNIT METHODS ======================
func (u *unitDb) AddUnit(ctx context.Context, in *proto.Unit) (*proto.Unit, error) {
	_, err := u.db.Exec("INSERT INTO unit(project_id, created_ts, updated_ts) VALUES($1, $2, $3)", in.ProjectId, in.Metadata.CreatedTs, in.Metadata.UpdatedTs)
	if err != nil {
		return nil, nil
	}
	return in, nil
}

func (u *unitDb) GetUnits(ctx context.Context, in []int32) ([]*proto.Unit, error) {
	if len(in) == 0 {
		return nil, errors.New("No project ids in the array")
	}
	query := `SELECT * FROM Unit WHERE project_id IN (` + strings.Repeat("?,", len(in))
	query = query[:len(query)-1] + `)`
	args := make([]interface{}, len(in))
	for i, id := range in {
		args[i] = id
	}
	rows, err := u.db.Query(query, args[:]...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	units := make([]*proto.Unit, 1)
	if !rows.Next() {
		return nil, errors.New("No units for the given project ID")
	}
	for rows.Next() {
		var unit proto.Unit
		err := rows.Scan(&unit.Id, &unit.ProjectId, &unit.Metadata.CreatedTs, &unit.Metadata.UpdatedTs)
		if err != nil {
			return nil, err
		}
		units = append(units, &unit)
	}
	return units, nil
}

// ===================== PROJECT METHODS ======================
func (p *projectDb) GetProject(ctx context.Context, projectIds []int32) ([]*proto.Project, error) {
	if len(projectIds) == 0 {
		return nil, errors.New("No project ids sent")
	}
	query := `SELECT * FROM Project WHERE id IN (` + strings.Repeat("?,", len(projectIds))
	query = query[:len(query)-1] + `)`
	args := make([]interface{}, len(projectIds))
	for i, id := range projectIds {
		args[i] = id
	}
	rows, err := p.db.Query(query, args[:]...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	projects := make([]*proto.Project, 1)

	for rows.Next() {
		project := &proto.Project{Metadata: &proto.Metadata{}}
		if err := rows.Scan(&project.Id, &project.Name, &project.Description, &project.Metadata.CreatedTs, &project.Metadata.UpdatedTs); err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}

func (p *projectDb) CreateProject(ctx context.Context, project *proto.Project) (*proto.Project, error) {

	name := strings.TrimSpace(project.Name)
	rows, err := p.db.Query("SELECT * FROM project WHERE name = $1", name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		return nil, errors.New("Project with the given name exists")
	}

	desc := strings.TrimSpace(project.Description)
	_, err = p.db.Exec("INSERT INTO project(name, desc, created_ts, updated_ts) VALUES($1, $2, $3, $4)", name, desc, project.Metadata.CreatedTs, project.Metadata.UpdatedTs)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (p *projectDb) ListProjects(ctx context.Context) ([]*proto.Project, error) {

	rows, err := p.db.Query("SELECT * FROM project")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	projects := make([]*proto.Project, 1)

	for rows.Next() {
		var project proto.Project
		var createdTsUnix int64
		var updatedTsUnix int64
		err := rows.Scan(&project.Id, &project.Name, &project.Description, &createdTsUnix, &updatedTsUnix)
		if err != nil {
			return nil, nil
		}
		project.Metadata.CreatedTs = timestamppb.New(createdTsUnix)
		projects = append(projects, &project)
	}
	return projects, nil
}

// =========================== UTIL FUNCS ===============================

func initDb(ctx context.Context) (*sql.DB, error) {
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
	return migrateDb(ctx, db, migrationDir)
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
