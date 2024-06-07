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

//HEMANT-Created project has id which passed project does not have, resulting in test case failure
func validProjectMessage(project *proto.Project) google_proto.Message {
  return &proto.Project{Name: project.Name, Description: project.Description}
}

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
			Name:        "Valid Name 2",
			Description: "Valid Desc",
		}},
		{project: &proto.Project{
			Description: "Project without name",
		}},
	}

  //@TODO Add a check for metadata dates to be successfully converted back to timestamps for validity

	for _, testCase := range testCases {
		projectFromCreate, err := projectStore.CreateProject(ctx, testCase.project)
		if err != nil && testCase.project.Description=="Project without name" {
		  t.Fatalf("could not create project :%s \n", err)
		}

		projectsFromStore, err := projectStore.GetProject(ctx, []int64{projectFromCreate.Id})
		if err != nil {
			t.Fatalf("could not get project :%s \n", err)
		}

		if len(projectsFromStore) != 1 {
			t.Fatalf("expected 1 project, got %d \n", len(projectsFromStore))
		}
		if !google_proto.Equal(validProjectMessage(testCase.project) , validProjectMessage(projectsFromStore[0])) {
			// @TODO : utils proto compare with diff
			t.Fatalf("project from store not equal \nProject : %s\nProject(S) : %s\n",
				testCase.project, projectsFromStore[0])
		}else {
      fmt.Println("Test case passed")
    }
	}
}

/*
func TestUnitCRUD(t *testing.T) {

	ctx := context.Background()
	require.NotNil(t, db)

	// Test that we are able to create projects

	// var err error
	unitStore := db.UnitStore()
	require.NotNil(t, unitStore)

	testCases := []struct { unit *proto.Unit }{
    { unit: &proto.Unit{ ProjectId: 0, }, },
    { unit: &proto.Unit{ ProjectId: 1, }, },
    { unit: &proto.Unit{ }, },
	}

  for _, test := range testCases{
    unitFromCreate, err := unitStore.AddUnit(ctx, test.unit)    
    if err != nil && test.unit.ProjectId!=0{
		  t.Fatalf("could not add unit :%s \n", err)
    }
    var projectIds []int64
    projectIds = append(projectIds, unitFromCreate.ProjectId)
    unitFromStore, err := unitStore.GetUnits(ctx, projectIds)
    if err != nil{
		  t.Fatalf("could not get unit :%s \n", err)
    }
    if !google_proto.Equal()
  } 


}
*/
