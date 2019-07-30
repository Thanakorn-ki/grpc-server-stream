package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	bp "github.com/salapao2136/grpc-server-stream/proto"

	"google.golang.org/grpc"
)

const (
	port = ":5599"
)

func main() {
	lis, err := net.Listen(("tcp"), port)
	if err != nil {
		log.Fatalf("Unable to listen grpc server: %v", err)
	}

	grpcServer := grpc.NewServer()
	bp.RegisterServiceUserServer(grpcServer, newUserService())

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Unable to start grpc server: %v", err)
	}

	log.Printf("Started gRPC server port: %v", port)

}

type ServiceUserServer interface {
	GetAll(context.Context, *empty.Empty) (*bp.UserResponse, error)
	GetAllStream(_ *empty.Empty, stream bp.ServiceUser_GetAllStreamServer) error
}

type userService struct {
}

func newUserService() ServiceUserServer {
	return &userService{}
}

func (us userService) GetAll(context.Context, *empty.Empty) (*bp.UserResponse, error) {
	return &bp.UserResponse{
		Name: "aum",
	}, nil
}

func (us userService) GetAllStream(_ *empty.Empty, stream bp.ServiceUser_GetAllStreamServer) error {
	for i := 0; i < 10; i++ {
		if err := stream.Send(&bp.UserResponse{Name: "AUM stream"}); err != nil {
			return err
		}
		time.Sleep(2 * time.Second)
	}

	return nil
}
