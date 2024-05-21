package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"time"

	"rg/UnitTracker/pkg/proto"
	"rg/UnitTracker/store"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const port = 50051

type serverImpl struct {
   db *sql.DB
	proto.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *serverImpl) SayHello(ctx context.Context, in *proto.HelloRequest) (*proto.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &proto.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func (s *serverImpl) CreateProject(ctx context.Context, in *proto.Project) (*proto.Project, error) {
   err := createProject(s.db, (*Project)(in))
   if err != nil {
      return in, nil
   }
   log.Printf("Creating project: %s", in.String())
   return nil, err
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

func (p *Project) String() string {
	if p == nil {
		return ""
	}
	return fmt.Sprintf("%d\t %s\t %s\n", p.Id, p.Name, p.Description)
}

// createProject inserts a new project into the Project table.
func createProject(db *sql.DB, project *Project) error {
	query := "INSERT INTO Project (name, desc) VALUES (?, ?)"
	_, err := db.Exec(query, project.Name, project.Description)
	return err
}

func getProjects(db *sql.DB) ([]*Project, error) {
	rows, err := db.Query("SELECT id, name, desc FROM Project")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*Project
	for rows.Next() {
		var project Project
		if err := rows.Scan(&project.Id, &project.Name, &project.Description); err != nil {
			return nil, err
		}
		projects = append(projects, &project)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}

func main() {
	log.Println("Hello from UT")

	go sampleClientCall()

	db, err := store.NewDbConnectionGetter().GetDb()
	if err != nil {
		log.Fatalf("failed to get db: %v", err)
	}

	projects, err := getProjects(db)
	if err != nil {
		log.Fatalf("failed to get projects: %v", err)
	}
	for _, prj := range projects {
		fmt.Println(prj.String())
	}

	// setup server
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	proto.RegisterGreeterServer(s, &serverImpl{db: db})
	log.Printf("Listening at %v", listen.Addr())
	if err := s.Serve(listen); err != nil {
		log.Fatalf("failed to serve :%v", err)
	}

}
