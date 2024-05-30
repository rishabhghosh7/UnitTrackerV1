package svc

import (
	"context"
	"fmt"
	"log"
	"rg/UnitTracker/pkg/proto"
	"rg/UnitTracker/store"
)

type Service struct {
	db store.Store
	// proto.UnimplementedGreeterServer
	proto.UnsafeGreeterServer
}

func NewService(db store.Store) *Service {
	return &Service{db: db}
}

// SayHello implements helloworld.GreeterServer
func (s *Service) SayHello(ctx context.Context, in *proto.HelloRequest) (*proto.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &proto.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func (s *Service) GetProject(ctx context.Context, in *proto.GetProjectRequest) (*proto.GetProjectResponse, error) {
	log.Printf("Getting project for id(s): %d\n", in.ProjectIds)

	projectStore := s.db.ProjectStore()
	project, err := projectStore.GetProject(ctx, in.ProjectIds)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning project: %s \n", project)
	return &proto.GetProjectResponse{Project: project}, nil
}

func (s *Service) CreateProject(ctx context.Context, in *proto.CreateProjectRequest) (*proto.CreateProjectResponse, error) {
	log.Printf("Creating new project name %s description %s \n", in.GetProject().GetName(), in.GetProject().GetDescription())

	projectStore := s.db.ProjectStore()
	projectFromDb, err := projectStore.CreateProject(ctx, in.Project)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Created new project: %s\n", in.String())
	return &proto.CreateProjectResponse{Project: projectFromDb}, nil
}
func (s *Service) ListProjects(ctx context.Context, data *proto.ListProjectRequest) (*proto.ListProjectResponse, error) {
	log.Printf("Listing projects for user \n")

	projectStore := s.db.ProjectStore()
	projects, err := projectStore.ListProjects(ctx)
	if err != nil {
		return nil, err
	}
	log.Printf("Returning project(s) %d\n", len(projects))
	return &proto.ListProjectResponse{Project: projects}, nil
}

func (s *Service) AddUnit(ctx context.Context, in *proto.AddUnitRequest) (*proto.AddUnitResponse, error) {
	log.Printf("Adding unit: %s\n", in.Unit)
	unitStore := s.db.UnitStore()
	_, err := unitStore.AddUnit(ctx, in.Unit)
	if err != nil {
		return nil, err
	}
	log.Printf("Added unit: %s\n", in.Unit)
	return &proto.AddUnitResponse{}, nil
}

func (s *Service) GetUnits(ctx context.Context, data *proto.GetUnitsRequest) (*proto.GetUnitsResponse, error) {
	log.Printf("Getting units for projectId(s): %d\n", data.ProjectIds)
	unitStore := s.db.UnitStore()
	units, err := unitStore.GetUnits(ctx, data.ProjectIds)
	if err != nil {
		return nil, err
	}
	log.Printf("Units for projectId(s): %s\n", units)
	return &proto.GetUnitsResponse{Units: units}, nil
}
