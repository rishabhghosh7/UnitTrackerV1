package main

import (
	"fmt"
   "rg/UnitTracker/pkg/proto/appPb"
)

func main() {
   fmt.Println("Hello from UT")
   rpcVar := appPb.HelloReply{Message: "Randy Orton"}
   fmt.Println("From rpc defns : ", rpcVar.GetMessage())
}
