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

//===================== UNIT METHODS ======================
func (u *unitDb) AddUnitToProject(ctx context.Context, in *proto.Unit)(error){
  _, err := u.db.Exec("INSERT INTO unit(project_id, create_ts) VALUES($1, $2)", in.ProjectId, in.CreateTs)  
  if err != nil {
    return err
  }
  return nil 
}

func (u *unitDb) GetUnitsForProject(ctx context.Context, in int32)([]*proto.Unit, error){
  rows, err := u.db.Query("SELECT * FROM unit WHERE project_id = $1", in)
  if err != nil{
    return nil, err
  }
  defer rows.Close()
  units := make([]*proto.Unit, 1)
  if !rows.Next(){
    return nil, errors.New("No units for the given project ID")
  }
  for rows.Next(){
    var unit proto.Unit
    err := rows.Scan(&unit.ProjectId, &unit.CreateTs)
    if err != nil{
      return nil, err
    }
    units=append(units, &unit)
  }
  return units, nil
}


//===================== PROJECT METHODS ======================
func (p *projectDb) GetProject(ctx context.Context, id int) (*proto.Project, error) {
	rows, err := p.db.Query("SELECT id, name, desc FROM Project WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	project := &proto.Project{}
	for rows.Next() {
		if err := rows.Scan(&project.Id, &project.Name, &project.Description); err != nil {
			return nil, err
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return project, nil
}

func (p *projectDb) CreateProject(ctx context.Context, project *proto.Project) (*proto.Project, error) {

  name:=strings.TrimSpace(project.Name)
  rows, err := p.db.Query("SELECT * FROM project WHERE name = $1", name)
  if err != nil{
    return nil, err
  }
  defer rows.Close()

  if rows.Next(){
    return nil, errors.New("Project with the given name exists")
  }

  desc:= strings.TrimSpace(project.Description)
  _, err = p.db.Exec("INSERT INTO project(name, desc) VALUES($1, $2)", name, desc);
  if err != nil{
    return nil, err
  }
	return nil, nil
}


func (p *projectDb) ListProjects(ctx context.Context)([]*proto.Project, error){
  rows, err := p.db.Query("SELECT * FROM project")
  if err != nil{
    return nil, err
  }
  defer rows.Close()

  projects := make([]*proto.Project, 1)

  for rows.Next(){
    var project proto.Project
    err := rows.Scan(&project.Id, &project.Name, &project.Description)
    if err != nil{
      return nil, nil
    }
    projects=append(projects, &project)
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
