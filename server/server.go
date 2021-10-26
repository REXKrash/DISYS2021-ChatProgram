package main

import (
	"fmt"
	"log"
	"net"
	"os"

	pb "chat-program/routeguide"

	"google.golang.org/grpc"
)

const port = ":50051"

//var users = make(map[string]pb.User)

type server struct {
	pb.UnimplementedChatServiceServer
}

func (s *server) ChatMessage(in pb.ChatService_ChatMessageServer) error {
	//log.Printf("Received Join server request from user: %v and uuid: %v", in.GetName(), in.GetUuid())
	//users[in.GetUuid()] = in
	defer func() {
		if err := recover(); err != nil {
			log.Printf("panic: %v", err)
			os.Exit(1)
		}
	}()
	for {
		input, error := in.Recv()
		if error != nil {
			log.Fatalln("Fatal", error)
			break
		}
		log.Println("Received input", input.Message)

		if err := in.Send(&pb.MessageResponse{Message: "Hello again"}); err != nil {
			log.Printf("broadcast err: %v", err)
		}
	}
	return nil

	//return &pb.MessageResponse{Status: "1", Message: "You successfully joined the chat room"}, nil
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
