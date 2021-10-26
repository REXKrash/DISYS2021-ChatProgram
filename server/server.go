package main

import (
	"fmt"
	"log"
	"math"
	"net"
	"os"

	pb "chat-program/routeguide"

	uuid "github.com/nu7hatch/gouuid"
	"google.golang.org/grpc"
)

const port = ":50051"

var users = make(map[string]pb.ChatService_ChatMessageServer)
var timestamp = 0

type server struct {
	pb.UnimplementedChatServiceServer
}

func broadcast(sender string, message string) {
	log.Println(sender+":", message)
	for _, v := range users {
		if err := v.Send(&pb.MessageResponse{Sender: sender, Message: message, Timestamp: int32(timestamp + 1)}); err != nil {
			log.Println("Failed to broadcast:", err)
		}
	}
}

func (s *server) ChatMessage(in pb.ChatService_ChatMessageServer) error {
	uuid, _ := uuid.NewV4()
	users[uuid.String()] = in
	broadcast("Server", "Some user has joined the chat room")

	defer func() {
		if err := recover(); err != nil {
			log.Printf("panic: %v", err)
			os.Exit(1)
		}
	}()

	for {
		input, error := in.Recv()
		if error != nil {
			log.Fatalln("Fatal error:", error)
			break
		}
		timestamp = int(math.Max(float64(timestamp), float64(input.Timestamp))) + 1

		broadcast(input.Sender, input.Message)
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
