package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"rg/UnitTracker/pkg/proto"
	"rg/UnitTracker/queries"
	"rg/UnitTracker/utils/timeutils"
	"strings"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type projectDb struct {
	db      *sql.DB
	queries *queries.Queries
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
		createdTs := timeutils.ProtobufTimestampToUnix(project.GetMetadata().GetCreatedTs())
		updatedTs := timeutils.ProtobufTimestampToUnix(project.GetMetadata().GetUpdatedTs())
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
