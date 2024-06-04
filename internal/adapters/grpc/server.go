package grpc

import (
	"fmt"
	"log"
	"net"

	"github.com/3Davydov/ms-order/internal/ports"
	"google.golang.org/grpc"
)

type Adapter struct {
	api  ports.APIPort
	port int
}

func NewAdapter(api ports.APIPort, port int) *Adapter {
	return &Adapter{
		api:  api,
		port: port,
	}
}

func (a Adapter) Run() {
	var err error

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		log.Fatalf("failed to listen on port %d, error %v", a.port, err)
	}

	grpcServer := grpc.NewServer()
	order.RegisterOrderServer
}
