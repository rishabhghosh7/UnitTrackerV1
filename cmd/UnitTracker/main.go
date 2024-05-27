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
	"rg/UnitTracker/svc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/guptarohit/asciigraph"
)

const port = 50051

/*
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
*/

type Project proto.Project
type Unit proto.Unit

func main() {
	log.Println("Hello from UT")

	runServer()

}
/*
func unitFunction(ctx context.Context, store store.Store) {
	// s := &serverImpl{db: store}
	//
	// var i int32
	//
	// for {
	// 	metadata := &proto.Metadata{
	// 		CreatedTs: &timestamppb.Timestamp{
	// 			Seconds: int64(time.Now().Unix()),
	// 		}}
	// 	unit := &proto.Unit{ProjectId: i + 1, Metadata: metadata}
	// 	auReq := &proto.AddUnitRequest{Unit: unit}
	// 	_, err := s.AddUnit(ctx, auReq)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// 	if i == 2 {
	// 		break
	// 	}
	// 	i++
	// }
	//
	// var projectIds []int32
	// projectIds = append(projectIds, 1)
	// projectIds = append(projectIds, 2)
	// projectIds = append(projectIds, 3)
	// guReq := &proto.GetUnitsRequest{ProjectIds: projectIds}
	// _, err := s.GetUnits(ctx, guReq)
	// if err != nil {
	// 	fmt.Println(err)
	// }
}
*/

func runServer() {
	ctx := context.Background()
	var err error

	// get store
	store, err := sqlite.NewSqliteConnector().Connect(ctx)
	if err != nil {
		log.Fatalf("failed to connect to store: %v", err)
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
