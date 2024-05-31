package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"rg/UnitTracker/pkg/proto"
	"rg/UnitTracker/store"
	"rg/UnitTracker/store/sqlite"

	"google.golang.org/grpc"
)

const port = 50051

type serverImpl struct {
	db store.Store
	proto.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *serverImpl) SayHello(ctx context.Context, in *proto.HelloRequest) (*proto.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &proto.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func (s *serverImpl) GetProject(ctx context.Context, in *proto.GetProjectRequest) (*proto.GetProjectResponse, error) {
	log.Printf("Getting project for id(s): %d\n", in.ProjectIds)

	projectStore := s.db.ProjectStore()
	project, err := projectStore.GetProject(ctx, in.ProjectIds)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning project: %s \n", project)
	return &proto.GetProjectResponse{Project: project}, nil
}

func (s *serverImpl) CreateProject(ctx context.Context, in *proto.CreateProjectRequest) (*proto.CreateProjectResponse, error) {
	log.Printf("Creating new project name %s description %s \n", in.GetProject().GetName(), in.GetProject().GetDescription())

	projectStore := s.db.ProjectStore()
	projectFromDb, err := projectStore.CreateProject(ctx, in.Project)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Created new project: %s\n", in.String())
	return &proto.CreateProjectResponse{Project: projectFromDb}, nil
}
func (s *serverImpl) ListProjects(ctx context.Context, data *proto.ListProjectRequest) (*proto.ListProjectResponse, error) {
	log.Printf("Listing projects for user \n")

	projectStore := s.db.ProjectStore()
	projects, err := projectStore.ListProjects(ctx)
	if err != nil {
		return nil, err
	}
	log.Printf("Returning project(s) %v\n:", projects)
	return &proto.ListProjectResponse{Project: projects}, nil
}

func (s *serverImpl) AddUnit(ctx context.Context, in *proto.AddUnitRequest) (*proto.AddUnitResponse, error) {
	log.Printf("Adding unit: %s\n", in.Unit)
	unitStore := s.db.UnitStore()
	_, err := unitStore.AddUnit(ctx, in.Unit)
	if err != nil {
		return nil, err
	}
	log.Printf("Added unit: %s\n", in.Unit)
	return &proto.AddUnitResponse{}, nil
}

func (s *serverImpl) GetUnits(ctx context.Context, data *proto.GetUnitsRequest) (*proto.GetUnitsResponse, error) {
	log.Printf("Getting units for projectId(s): %d\n", data.ProjectIds)
	unitStore := s.db.UnitStore()
	units, err := unitStore.GetUnits(ctx, data.ProjectIds)
	if err != nil {
		return nil, err
	}
	log.Printf("Units for projectId(s): %s\n", units)
	return &proto.GetUnitsResponse{Units: units}, nil
}

type Project proto.Project
type Unit proto.Unit

func main() {
	log.Println("Hello from UT")

	runServer()
}

func testSqliteFns(ctx context.Context, store store.Store) {

	s := &serverImpl{db: store}
	projectIds := make([]int32, 0)
	projectIds = append(projectIds, 1)
	projectIds = append(projectIds, 2)
	projectIds = append(projectIds, 3)
	projectIds = append(projectIds, 4)
	/*
		projects, err := s.GetProject(ctx, &proto.GetProjectRequest{ProjectIds: projectIds})
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("==============PROJECTS=============")
		fmt.Println(projects)
		fmt.Println("===========================")
	*/
	units, err := s.GetUnits(ctx, &proto.GetUnitsRequest{ProjectIds: projectIds})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("==============UNITS=============")
	fmt.Println(units)
	fmt.Println("===========================")
}

func runServer() {
	ctx := context.TODO()
	// get store
	store, err := sqlite.NewSqliteConnector().Connect(ctx)
	if err != nil {
		log.Fatalf("failed to connect to store: %v", err)
	}

	testSqliteFns(ctx, store)

	// setup server
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	proto.RegisterGreeterServer(s, &serverImpl{db: store})
	log.Printf("Listening at %v", listen.Addr())
	if err := s.Serve(listen); err != nil {
		log.Fatalf("failed to serve :%v", err)
	}
}
