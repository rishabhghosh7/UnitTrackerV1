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
   log.Printf("Getting project for id: %d\n", in.GetId())

   projectStore := s.db.ProjectStore()
   project, err := projectStore.GetProject(ctx, int(in.GetId()))
   if err != nil {
      return nil, err
   }

   log.Printf("Returning project: %s \n", project.String())
   return &proto.GetProjectResponse{Project: project}, nil
}

func (s *serverImpl) CreateProject(ctx context.Context, in *proto.Project) (*proto.Project, error) {

  projectStore := s.db.ProjectStore()
  _, err := projectStore.CreateProject(ctx, in)
  if err != nil{
    fmt.Println(err)
    return nil, err
  }
  fmt.Println("Created new project")
	return in, nil
}
func (s *serverImpl) ListProjects(ctx context.Context)([]*proto.Project, error){

  projectStore := s.db.ProjectStore()
  projects, err := projectStore.ListProjects(ctx)
  if err != nil{
    fmt.Println(err)
    return nil, err
  }
  return projects, nil;
}

func (s *serverImpl) AddUnit(ctx context.Context, projectId int32)error{
  unitStore:=s.db.UnitStore()
  unit:=&proto.Unit{CreateTs: int32(time.Now().Unix()), ProjectId: int32(projectId)}
  err:=unitStore.AddUnitToProject(ctx, unit)
  if err != nil{
    fmt.Println(err)
    return err
  }
  return nil
}

func (s *serverImpl) GetUnitsForProject(ctx context.Context, projectId int32)([]*proto.Unit, error){
  unitStore:=s.db.UnitStore()
  units, err:=unitStore.GetUnitsForProject(ctx, projectId)  
  if err != nil{
    fmt.Println(err)
    return nil, err
  }
  return units, nil
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

	c.CreateProject(context.TODO(), &proto.Project{Name: "Created Project", Description: "Project desc"})
	c.CreateProject(context.TODO(), &proto.Project{Name: "Another Project", Description: "Project desc"})

}

type Project proto.Project

type Unit proto.Unit

func (p *Project) String() string {
	if p == nil {
		return ""
	}
	return fmt.Sprintf("%d\t %s\t %s\n", p.Id, p.Name, p.Description)
}

func (u *Unit) String() string {
	if u == nil {
		return ""
	}
	return fmt.Sprintf("%d\t %d\t", u.ProjectId, u.CreateTs)
}

func main() {
	log.Println("Hello from UT")

	runServer()
}


func runServer() {
	ctx := context.TODO()
	// get store
	store, err := sqlite.NewSqliteConnector().Connect(ctx)
	if err != nil {
		log.Fatalf("failed to connect to store: %v", err)
	}
  
   
	projectDb := store.ProjectStore()
	project1, err := projectDb.GetProject(ctx, 1)
	if err != nil {
		log.Fatalf("could not get p1")
	}
	project2, err := projectDb.GetProject(ctx, 2)
	if err != nil {
		log.Fatalf("could not get p1")
	}
	fmt.Println("from store", project1.String(), project2.String())
  
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
