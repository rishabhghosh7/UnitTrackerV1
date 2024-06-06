package main

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net"
	"rg/UnitTracker/pkg/proto"
	"rg/UnitTracker/store"
	"rg/UnitTracker/store/sqlite"
	"rg/UnitTracker/svc"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/guptarohit/asciigraph"
	// "google.golang.org/protobuf/types/known/timestamppb"
)

const port = 50051

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

	// c.CreateProject(context.TODO(), &proto.CreateProjectRequest{Project: &proto.Project{Name: "Created Project", Description: "Project desc"}})
	// c.CreateProject(context.TODO(), &proto.CreateProjectRequest{Project: &proto.Project{Name: "Another Project", Description: "Project desc"}})

	data := []float64{3, 4, 5, 1, 2, 3, 7, 8, 9}
	data1 := []float64{12, 12, 12, 12, 12, 12, 12, 12, 12, 12}
	graph := asciigraph.PlotMany([][]float64{data, data1})

	fmt.Println(graph)
}

type Project proto.Project
type Unit proto.Unit

func main() {
	log.Println("Hello from UT")

	runServer()
}

func testSqliteFns(ctx context.Context, store store.Store) {

	/*
		s := &serverImpl{db: store}

		project := &proto.Project{
			Name:        "Marathon",
			Description: "",
			Metadata: &proto.Metadata{
				CreatedTs: timestamppb.Now(),
				UpdatedTs: timestamppb.Now(),
			},
		}
		resp, err := s.CreateProject(ctx, &proto.CreateProjectRequest{Project: project})
		if err != nil {
			fmt.Println("Encountered an error")
		}
		fmt.Println(resp)
	*/
}

func runServer() {
	ctx := context.Background()
	var err error

	// get store
	store, err := sqlite.NewTestSqliteConnector().Connect(ctx)
	if err != nil {
		log.Fatalf("failed to connect to store: %v", err)
	}

	projectStore := store.ProjectStore()
	project := &proto.Project{
		Name:        "blah",
		Description: "blah blah",
		Metadata: &proto.Metadata{
			CreatedTs: timestamppb.Now(),
			UpdatedTs: timestamppb.Now(),
		},
	}
	_, err = projectStore.CreateProject(ctx, project)
	if err != nil {
		log.Fatal("failed to connect to store: ", err)
	}

	// setup server
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	service := svc.NewService(store)
	proto.RegisterGreeterServer(s, service)
	log.Printf("Listening at %v", listen.Addr())
	if err := s.Serve(listen); err != nil {
		log.Fatalf("failed to serve :%v", err)
	}
}
