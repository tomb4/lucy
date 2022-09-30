package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	pb "lucy/proto/party"
	"net"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.PartyServiceServer
}

func (s *server) JoinMetaParty(ctx context.Context, req *pb.JoinMetaPartyReq) (*pb.MetaParty, error) {
	log.Printf("Received: %v", req.PartyId)
	return &pb.MetaParty{
		PartyId: req.PartyId,
		Name:    "meta",
	}, nil
}

func (s *server) ExitMetaParty(ctx context.Context, req *pb.ExitMetaPartyReq) (*pb.Nil, error) {
	log.Println("exit:", req.UserId)
	return &pb.Nil{}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterPartyServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
