package main

import (
	"fmt"
	"rg/UnitTracker/pkg/proto/rg/UnitTracker/pkg/proto"
)

func main() {
   fmt.Println("Hello from UT")
   rpcVar := proto.HelloReply{Message: "Randy Orton"}
   fmt.Println("From rpc defns : ", rpcVar.GetMessage())
}
