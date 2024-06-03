package sqlite_test

import (
	"context"
	"fmt"
	"os"
	"rg/UnitTracker/pkg/proto"
	"rg/UnitTracker/store"
	"rg/UnitTracker/store/sqlite"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

var db store.Store

func TestMain(m *testing.M) {
	testStore, err := sqlite.NewTestSqliteConnector().Connect(context.TODO())
	if err != nil {
		dir, _ := os.Getwd()
		fmt.Printf("%s could not get db at %s", err, dir)
		os.Exit(1)
	}
	db = testStore
	os.Exit(m.Run())
}

func TestProjectsCRUD(t *testing.T) {
	ctx := context.Background()
	// var err error

	require.NotNil(t, db)
	pdb := db.ProjectStore()
	require.NotNil(t, pdb)

	pdb.GetProject(ctx, []int32{1})
	pdb.CreateProject(ctx, &proto.Project{})

	// projectNoNameNoDesc := &proto.Project{}
	// projectNoName := &proto.Project{Description: "test desc"}
	// legalProject := &proto.Project{}
	//
	// projectStore := db.ProjectStore()
	// _, err = projectStore.CreateProject(context.TODO(), projectNoNameNoDesc)
	// // require.ErrorIs() // @TODO : check error
	// require.Error(t, err)
	//
	// _, err = projectStore.CreateProject(context.TODO(), projectNoName)
	// // require.ErrorIs() // @TODO : check error
	// require.Error(t, err)
	//
	// _, err = projectStore.CreateProject(context.TODO(), legalProject)
	// // require.ErrorIs() // @TODO : check error
	// require.NoError(t, err)
}
