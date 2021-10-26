package main

import (
	"fmt"
	"log"
	"net"
	"os"

	pb "chat-program/routeguide"

	uuid "github.com/nu7hatch/gouuid"
	"google.golang.org/grpc"
)

const port = ":50051"

var users = make(map[string]pb.ChatService_ChatMessageServer)

type server struct {
	pb.UnimplementedChatServiceServer
}

func (s *server) ChatMessage(in pb.ChatService_ChatMessageServer) error {
	uuid, _ := uuid.NewV4()
	users[uuid.String()] = in

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

		for _, v := range users {
			if err := v.Send(&pb.MessageResponse{Sender: input.Sender, Message: input.Message}); err != nil {
				log.Printf("broadcast err: %v", err)
			}
		}
	}
	return nil
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
