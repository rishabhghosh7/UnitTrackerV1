package sqlite_test

import (
	"context"
	"rg/UnitTracker/pkg/proto"
	"rg/UnitTracker/store/sqlite"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProjectsCRUD(t *testing.T) {
	db, err := sqlite.NewSqliteConnector().Connect(context.TODO()) // @TODO : testing db
	require.NoError(t, err)

	projectNoNameNoDesc := &proto.Project{}
	projectNoName := &proto.Project{Description: "test desc"}
	legalProject := &proto.Project{}

	projectStore := db.ProjectStore()
	_, err = projectStore.CreateProject(context.TODO(), projectNoNameNoDesc)
	// require.ErrorIs() // @TODO : check error
	require.Error(t, err)

	_, err = projectStore.CreateProject(context.TODO(), projectNoName)
	// require.ErrorIs() // @TODO : check error
	require.Error(t, err)

	_, err = projectStore.CreateProject(context.TODO(), legalProject)
	// require.ErrorIs() // @TODO : check error
	require.NoError(t, err)
}
