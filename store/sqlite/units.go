package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"rg/UnitTracker/pkg/proto"
	"strings"
)


type unitDb struct {
	db *sql.DB
}

// ===================== UNIT METHODS ======================
func (u *unitDb) AddUnit(ctx context.Context, unit *proto.Unit) (*proto.Unit, error) {
	_, err := u.db.Exec("INSERT INTO unit(project_id, created_ts, updated_ts) VALUES($1, $2, $3)", unit.ProjectId, unit.Metadata.CreatedTs, unit.Metadata.UpdatedTs)
	if err != nil {
		return nil, nil
	}
	return unit, nil
}

func (u *unitDb) GetUnits(ctx context.Context, projectIds []int32) ([]*proto.Unit, error) {
	if len(projectIds) == 0 {
		return nil, errors.New("No project ids in the array")
	}
	query := `SELECT * FROM Unit WHERE project_id IN (` + strings.Repeat("?,", len(projectIds))
	query = query[:len(query)-1] + `)`
	args := make([]interface{}, len(projectIds))
	for i, id := range projectIds {
		args[i] = id
	}
	rows, err := u.db.Query(query, args[:]...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	units := make([]*proto.Unit, len(projectIds))
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

