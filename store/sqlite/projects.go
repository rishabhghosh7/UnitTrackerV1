package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"rg/UnitTracker/pkg/proto"
	"strings"
)


type projectDb struct {
	db *sql.DB
}

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


   var projects []*proto.Project
	for rows.Next() {
		var project proto.Project
		var createdTsUnix int64
		var updatedTsUnix int64
		err := rows.Scan(&project.Id, &project.Name, &project.Description, &createdTsUnix, &updatedTsUnix)
		if err != nil {
			return nil, nil
		}
		// project.Metadata.CreatedTs = timestamppb.New(createdTsUnix)
		projects = append(projects, &project)
	}
	return projects, nil
}

