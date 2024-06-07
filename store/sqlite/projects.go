package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"rg/UnitTracker/pkg/proto"
	"rg/UnitTracker/queries"
	"rg/UnitTracker/utils/sliceutils"
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

func ProjectModelToProto(project *queries.Project) *proto.Project {
	if project == nil {
		return nil
	}

	return &proto.Project{
		Id:          project.ID,
		Name:        project.Name,
		Description: project.Desc.String,
		// @TODO : metadata
	}
}

func ProjectProtoToModel(project *proto.Project) *queries.Project {
	if project == nil {
		return nil
	}

	return &queries.Project{
		ID:   project.Id,
		Name: project.Name,
		Desc: sql.NullString{String: project.Description, Valid: true},
		// @TODO : metadata
	}
}

func (p *projectDb) GetProject(ctx context.Context, projectIds []int64) ([]*proto.Project, error) {
	if len(projectIds) == 0 {
		return []*proto.Project{}, nil
	}

	rows, err := p.queries.GetProject(ctx, projectIds)
	if err != nil {
		return nil, err
	}
	fmt.Printf("db %d rows \n", len(rows))
	protoRows := sliceutils.Map(rows, func(p queries.Project) *proto.Project {
		return ProjectModelToProto(&p)
	})

	fmt.Printf("returning %d rows \n", len(protoRows))
	return protoRows, nil
}

func (p *projectDb) CreateProject(ctx context.Context, project *proto.Project) (*proto.Project, error) {
	name := strings.TrimSpace(project.Name)
	_, err := p.queries.GetProjectByName(ctx, name)
	if err == nil {
		return nil, errors.New("Project with the given name already exists")
	}
	if err == sql.ErrNoRows {
    trimmedDesc := strings.TrimSpace(project.Description)
		desc := sql.NullString{String: trimmedDesc, Valid: true}
		createdTs := timeutils.ProtobufTimestampToUnix(project.GetMetadata().GetCreatedTs())
		updatedTs := timeutils.ProtobufTimestampToUnix(project.GetMetadata().GetUpdatedTs())
		row, err := p.queries.CreateProject(ctx, queries.CreateProjectParams{Name: name, Desc: desc, CreatedTs: createdTs, UpdatedTs: updatedTs})
		if err != nil {
			return nil, err
		}
    var projects []queries.Project
    projects = append(projects, row)
    protoRows := sliceutils.Map(projects, func(p queries.Project) *proto.Project {
      return ProjectModelToProto(&p)
    })
		return protoRows[0], nil
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
  row, err := p.queries.UpdateProject(ctx, queries.UpdateProjectParams{ID: in.Id, Desc: desc})
  if err != nil {
    return nil, err
  }
  var projects []queries.Project
  projects = append(projects, row)
  protoRows := sliceutils.Map(projects, func(p queries.Project) *proto.Project {
    return ProjectModelToProto(&p)
  })
	return protoRows[0], nil
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
