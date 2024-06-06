package sqlite_test

import (
	"context"
	"fmt"
	"os"
	"rg/UnitTracker/pkg/proto"
	"rg/UnitTracker/store"
	"rg/UnitTracker/store/sqlite"
	"testing"

	google_proto "github.com/golang/protobuf/proto"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

var db store.Store

func TestMain(m *testing.M) {
	testStore, err := sqlite.NewTestSqliteConnector().Connect(context.TODO())
	if err != nil {
		fmt.Printf("could not get db :%s\n", err)
		os.Exit(1)
	}
	db = testStore
	os.Exit(m.Run())
}

func TestProjectCRUD(t *testing.T) {
	ctx := context.Background()
	require.NotNil(t, db)

	// Test that we are able to create projects

	// var err error
	projectStore := db.ProjectStore()
	require.NotNil(t, projectStore)

	testCases := []struct {
		// name    string
		project *proto.Project
	}{
		{project: &proto.Project{
			Name:        "N",
			Description: "Valid Desc",
		}},
		{project: &proto.Project{
			Name:        "Valid Name",
			Description: "",
		}},
		{project: &proto.Project{
			Name:        "Valid Name",
			Description: "Valid Desc",
		}},
	}

	for _, testCase := range testCases {
		projectFromCreate, err := projectStore.CreateProject(ctx, testCase.project)
		if err != nil {
			t.Fatalf("could not create project :%s \n", err)
		}
		projectsFromStore, err := projectStore.GetProject(ctx, []int64{projectFromCreate.Id})
		if err != nil {
			t.Fatalf("could not get project :%s \n", err)
		}

		if len(projectsFromStore) != 1 {
			t.Fatalf("expected 1 project, got %d \n", len(projectsFromStore))
		}

		if !google_proto.Equal(projectsFromStore[0], testCase.project) {
			// @TODO : utils proto compare with diff
			t.Fatalf("project from store not equal \nProject : %s\nProject(S) : %s\n",
				testCase.project, projectsFromStore)
		}
	}

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
