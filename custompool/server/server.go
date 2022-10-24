package main

import (
	context "context"
	"flag"
	"fmt"
	"github.com/rpcxio/rpcx-examples/custompool/pb"
	"github.com/smallnest/rpcx/server"
	"github.com/smallnest/rpcx/share"
	"log"
	"time"
)

var addr = flag.String("addr", "localhost:8972", "server address")

type Worker struct {
	funChan chan func()
}

func NewWorker() *Worker {
	p := new(Worker)
	p.funChan = make(chan func(), 10240)
	return p
}

func (w *Worker) Submit(task func()) {
	w.funChan <- task
}

func (w *Worker) Stop() {
	close(w.funChan)
}

func (w *Worker) StopAndWait() {

}

func (w *Worker) StopAndWaitFor(deadline time.Duration) {

}

func (w *Worker) Len() int { return len(w.funChan) }

func (p *Worker) Run() {
	for f := range p.funChan {
		log.Printf("recv msg\n")
		f()
	}

}

type Greeter struct{}

func (*Greeter) SayHello(ctx context.Context, args *pb.HelloRequest, reply *pb.HelloReply) (err error) {
	*reply = pb.HelloReply{
		Message: fmt.Sprintf("ok hello %s!", args.Name),
	}
	return nil
}

func main() {
	flag.Parse()

	share.Trace = true

	worker := NewWorker()

	s := server.NewServer(
		server.WithCustomPool(worker))
	greeter := &Greeter{}
	pb.RegisterGreeterServe(s, greeter, "")
	go func() {
		s.Serve("tcp", *addr)
	}()

	worker.Run()
}
