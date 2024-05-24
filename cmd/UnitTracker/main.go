package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"rg/UnitTracker/pkg/proto"
	"rg/UnitTracker/store"
	"rg/UnitTracker/store/sqlite"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
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
  log.Printf("Creating new project")
  projectStore := s.db.ProjectStore()
  _, err := projectStore.CreateProject(ctx, in.Project)
  if err != nil{
    fmt.Println(err)
    return nil, err
  }
  fmt.Printf("Created new project: %s\n", in.String())
	return &proto.CreateProjectResponse{Project: in.Project}, nil
}
func (s *serverImpl) ListProjects(ctx context.Context, data *proto.ListProjectRequest)(*proto.ListProjectResponse, error){
  log.Printf("Listing projects for user")
  projectStore := s.db.ProjectStore()
  projects, err := projectStore.ListProjects(ctx)
  if err != nil{
    fmt.Println(err)
    return nil, err
  }
  log.Printf("Returning project(s) %s\n:", projects)
  return &proto.ListProjectResponse{Project: projects}, nil;
}

func (s *serverImpl) AddUnit(ctx context.Context, unit *proto.AddUnitRequest)(*proto.AddUnitResponse,error){
  log.Printf("Adding unit: %s\n", unit.Unit)
  unitStore:=s.db.UnitStore()
  _, err:=unitStore.AddUnit(ctx, unit.Unit)
  if err != nil{
    fmt.Println(err)
    return nil, err
  }
  log.Printf("Added unit: %s\n", unit.Unit)
  return &proto.AddUnitResponse{}, nil
}

func (s *serverImpl) GetUnits(ctx context.Context, data *proto.GetUnitsRequest)(*proto.GetUnitsResponse, error){
  log.Printf("Getting units for projectId(s): %d\n", data.ProjectIds)
  unitStore:=s.db.UnitStore()
  units, err:=unitStore.GetUnits(ctx, data.ProjectIds)  
  if err != nil{
    fmt.Println(err)
    return nil, err
  }
  log.Printf("Units for projectId(s): %s\n", units)
  return &proto.GetUnitsResponse{Units: units}, nil
}

func sampleClientCall() {
	time.Sleep(1000 * time.Millisecond)

	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := proto.NewGreeterClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &proto.HelloRequest{Name: "Randy"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())

  c.CreateProject(context.TODO(), &proto.CreateProjectRequest{Project: &proto.Project{Name: "Created Project", Description: "Project desc"}})
  c.CreateProject(context.TODO(), &proto.CreateProjectRequest{Project: &proto.Project{Name: "Another Project", Description: "Project desc"}})

}

type Project proto.Project

type Unit proto.Unit

func main() {
	log.Println("Hello from UT")

	runServer()
}


func unitFunction(ctx context.Context, store store.Store){
  
  s:=&serverImpl{db: store} 

  var i int32;

  for{
    metadata:=&proto.Metadata{
      CreatedTs: &timestamppb.Timestamp{
      Seconds: int64(time.Now().Unix()), 
    }}
    unit:=&proto.Unit{ProjectId: i+1,  Metadata: metadata}
    auReq:=&proto.AddUnitRequest{Unit: unit}
    _, err:=s.AddUnit(ctx, auReq)
    if err != nil{
      fmt.Println(err)
    }
    if i==2{
      break
    }
    i++;
  }
  
  var projectIds []int32
  projectIds=append(projectIds, 1)
  projectIds=append(projectIds, 2)
  projectIds=append(projectIds, 3)
  guReq:=&proto.GetUnitsRequest{ProjectIds: projectIds}
  _, err:=s.GetUnits(ctx, guReq)
  if err != nil{
    fmt.Println(err)
  }
}


func runServer() {
	ctx := context.TODO()
	// get store
	store, err := sqlite.NewSqliteConnector().Connect(ctx)
	if err != nil {
		log.Fatalf("failed to connect to store: %v", err)
	}
  
   
	projectDb := store.ProjectStore()
  var projectIds []int32
  projectIds=append(projectIds, 1)
  projectIds=append(projectIds, 2)
  projectIds=append(projectIds, 4)
  project1, err := projectDb.GetProject(ctx, projectIds)
	if err != nil {
    log.Println(err)
		log.Fatalf("could not get p1")
	}
  fmt.Println(project1)

  unitFunction(ctx, store)

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
