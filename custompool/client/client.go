package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/rpcxio/rpcx-examples/custompool/pb"
	"github.com/smallnest/rpcx/client"
)

var (
	addr = flag.String("addr", "localhost:8972", "server address")
)

func main() {
	flag.Parse()

	d, _ := client.NewPeer2PeerDiscovery("tcp@"+*addr, "")
	xclient := client.NewXClient("Greeter", client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer xclient.Close()

	args := &pb.HelloRequest{
		Name: "world",
	}

	for {
		reply := &pb.HelloReply{}
		err := xclient.Call(context.Background(), "SayHello", args, reply)
		if err != nil {
			log.Fatalf("failed to call: %v", err)
		}

		log.Printf("reply:%s", reply.Message)
		time.Sleep(time.Millisecond * 100)
	}

}
