package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"google.golang.org/protobuf/types/known/timestamppb"
	"rg/UnitTracker/pkg/proto"
	"rg/UnitTracker/queries"
	"rg/UnitTracker/utils/timeutils"
	"time"
)

type unitDb struct {
	db      *sql.DB
	queries *queries.Queries
}

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
