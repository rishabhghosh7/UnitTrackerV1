package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"rg/UnitTracker/pkg/proto"
	"rg/UnitTracker/queries"
	"rg/UnitTracker/store"
	"rg/UnitTracker/utils/fsutils"
	"rg/UnitTracker/utils/timeutils"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const dbFilepath = "./store/sqlite/_sqlite.db"
const migrationDir = "./store/migrations/"

var dbSingleton *sql.DB

type sqliteConnector struct {
	db      *sql.DB // never access this directly
	queries *queries.Queries
}

func NewSqliteConnector() store.Connecter {
	return &sqliteConnector{}
}

// @TODO
// func RunTransaction(store, func(trancsaction) {}) error

func (c *sqliteConnector) Connect(ctx context.Context) (store.Store, error) {
	var q *queries.Queries
	if dbSingleton == nil {
		var err error
		dbSingleton, q, err = initDb(ctx)
		if err != nil {
			return nil, err
		}
		if q == nil {
			return nil, errors.New("Unable to get SQLC generated queries")
		}
	}
	return &sqliteConnector{db: dbSingleton, queries: q}, nil
}

func (c *sqliteConnector) ProjectStore() store.ProjectStore {
	return &projectDb{db: c.db, queries: c.queries}
}

type projectDb struct {
	db      *sql.DB
	queries *queries.Queries
}

func (c *sqliteConnector) UnitStore() store.UnitStore {
	return &unitDb{db: c.db, queries: c.queries}
}

type unitDb struct {
	db      *sql.DB
	queries *queries.Queries
}

// ===================== UNIT METHODS ======================
func (u *unitDb) AddUnit(ctx context.Context, unit *proto.Unit) (*proto.Unit, error) {
	createdTs := timeutils.ProtobufTimestampToUnix(unit.Metadata.CreatedTs)
	updatedTs := timeutils.ProtobufTimestampToUnix(unit.Metadata.UpdatedTs)

	err := u.queries.AddUnit(ctx, queries.AddUnitParams{ProjectID: unit.ProjectId, CreatedTs: createdTs, UpdatedTs: updatedTs})
	//_, err := u.db.Exec("INSERT INTO unit(project_id, created_ts, updated_ts) VALUES($1, $2, $3)", unit.ProjectId, createdTs, updatedTs)
	if err != nil {
		return nil, err
	}
	return unit, nil
}

func (u *unitDb) GetUnits(ctx context.Context, projectIds []int64) ([]*proto.Unit, error) {
	if len(projectIds) == 0 {
		return nil, errors.New("No project ids in the array")
	}
	units := make([]*proto.Unit, 0)
	rows, err := u.queries.GetUnits(ctx, projectIds)
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range rows {
		unit := &proto.Unit{
			Metadata: &proto.Metadata{
				CreatedTs: timestamppb.New(time.Unix(v.CreatedTs, 0)),
				UpdatedTs: timestamppb.New(time.Unix(v.UpdatedTs, 0)),
			},
			ProjectId: v.ProjectID,
			Id:        v.ID,
		}
		units = append(units, unit)
	}
	return units, nil
}

// ===================== PROJECT METHODS ======================
func (p *projectDb) GetProject(ctx context.Context, projectIds []int64) ([]*proto.Project, error) {
	if len(projectIds) == 0 {
		return nil, errors.New("No project ids sent")
	}

	rows, err := p.queries.GetProject(ctx, projectIds)
	if err != nil {
		return nil, err
	}
	projects := make([]*proto.Project, 0)
	for _, v := range rows {
		if !v.Desc.Valid {
			v.Desc.String = ""
		}
		project := &proto.Project{
			Metadata: &proto.Metadata{
				CreatedTs: timestamppb.New(time.Unix(v.CreatedTs, 0)),
				UpdatedTs: timestamppb.New(time.Unix(v.UpdatedTs, 0)),
			},
			Id:          v.ID,
			Name:        v.Name,
			Description: v.Desc.String,
		}
		projects = append(projects, project)
	}
	return projects, nil
}

func (p *projectDb) CreateProject(ctx context.Context, project *proto.Project) (*proto.Project, error) {
	name := strings.TrimSpace(project.Name)
	_, err := p.queries.GetProjectByName(ctx, name)
	if err == nil {
		return nil, errors.New("Project with the given name already exists")
	}
	if err == sql.ErrNoRows {
		desc := sql.NullString{String: strings.TrimSpace(project.Description)}
		createdTs := timeutils.ProtobufTimestampToUnix(project.Metadata.CreatedTs)
		updatedTs := timeutils.ProtobufTimestampToUnix(project.Metadata.UpdatedTs)
		err := p.queries.CreateProject(ctx, queries.CreateProjectParams{Name: name, Desc: desc, CreatedTs: createdTs, UpdatedTs: updatedTs})
		if err != nil {
			return nil, err
		}
		return project, nil
	}
	return nil, err
}

func (p *projectDb) UpdateProject(ctx context.Context, in *proto.Project) (*proto.Project, error) {
	name := strings.TrimSpace(in.Name)
	_, err := p.queries.GetProjectByName(ctx, name)
	if err == nil {
		return nil, errors.New("Project with the given name already exists")
	}
	desc := sql.NullString{String: strings.TrimSpace(in.Description)}
	err = p.queries.UpdateProject(ctx, queries.UpdateProjectParams{ID: in.Id, Desc: desc})
	if err != nil {
		return nil, err
	}
	return in, nil
}

func (p *projectDb) ListProjects(ctx context.Context) ([]*proto.Project, error) {
	var projects []*proto.Project
	rows, err := p.queries.ListProjects(ctx)
	if err != nil {
		return nil, err
	}
	for _, v := range rows {
		if !v.Desc.Valid {
			v.Desc.String = ""
		}
		project := &proto.Project{
			Description: v.Desc.String,
			Name:        v.Name,
			Id:          v.ID,
			Metadata: &proto.Metadata{
				CreatedTs: timestamppb.New(time.Unix(v.CreatedTs, 0)),
				UpdatedTs: timestamppb.New(time.Unix(v.UpdatedTs, 0)),
			},
		}
		projects = append(projects, project)
	}
	return projects, nil
}

// =========================== UTIL FUNCS ===============================

func initDb(ctx context.Context) (*sql.DB, *queries.Queries, error) {
	if !fsutils.FileExists(dbFilepath) {
		log.Printf("Db not found, creating %s...", dbFilepath)
	}
	db, err := sql.Open("sqlite3", dbFilepath)
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
