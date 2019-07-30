package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes/empty"

	bp "github.com/salapao2136/grpc-server-stream/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func main() {
	fmt.Println("gRPC client running")

	conn, err := NewGrpcClient(context.Background())
	if err != nil {
		log.Print(err)
	}
	defer conn.Close()

	userClient := bp.NewServiceUserClient(conn.ClientConn)
	stream, e := userClient.GetAllStream(context.Background(), &empty.Empty{})
	if e != nil {
		log.Print(e)
	}

	for {
		feature, err := stream.Recv()
		if err == io.EOF {
			log.Println(err)
			break
		}
		if err != nil {
			log.Fatalf("%v.GetAllStream(_) = _, %v", conn, err)
		}
		log.Println(feature)
	}

}

type rpcer struct {
	*grpc.ClientConn
	ctx context.Context
}

func NewGrpcClient(ctx context.Context) (*rpcer, error) {
	var kacp = keepalive.ClientParameters{
		Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
		Timeout:             time.Second,      // wait 1 second for ping ack before considering the connection dead
		PermitWithoutStream: true,             // send pings even without active streams
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Unable to connect grpc server %v", r)
		}
	}()

	conn, err := grpc.Dial("localhost:5599", grpc.WithInsecure(), grpc.WithKeepaliveParams(kacp))
	if err != nil {
		log.Printf("Unable to connect grpc server %v", err)
		return nil, err
	}
	return &rpcer{
		ClientConn: conn,
		ctx:        ctx,
	}, nil
}
