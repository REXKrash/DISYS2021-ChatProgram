package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "chat-program/routeguide"

	"google.golang.org/grpc"
)

const port = ":50051"

var users = make(map[string]*pb.User)

type server struct {
	pb.UnimplementedChatServiceServer
}

func (s *server) JoinServer(ctx context.Context, in *pb.User) (*pb.Empty, error) {
	log.Printf("Received Join server request from user: %v and uuid: %v", in.GetName(), in.GetUuid())
	users[in.GetUuid()] = in

	log.Println("Sending empty response...")
	return &pb.Empty{}, nil
}

func main() {
	fmt.Println("--- SERVER APP ---")

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterChatServiceServer(s, &server{})
	log.Printf("Server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
