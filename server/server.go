package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"net"
	"os"

	pb "chat-program/routeguide"

	"google.golang.org/grpc"
)

const port = ":50051"

var users = make(map[string]UserEntity)
var timestamp = 0

type server struct {
	pb.UnimplementedChatServiceServer
}

type UserEntity struct {
	server pb.ChatService_JoinChatServerServer
	leave  chan bool
}

func newUserEntity(srv pb.ChatService_JoinChatServerServer) UserEntity {
	return UserEntity{
		server: srv,
		leave:  make(chan bool),
	}
}

func broadcast(sender string, message string) {
	log.Println(sender+":", message)
	timestamp++
	for _, v := range users {
		if err := v.server.Send(&pb.MessageResponse{Sender: sender, Message: message, Timestamp: int32(timestamp)}); err != nil {
			log.Println("Failed to broadcast:", err)
		}
	}
}

func updateTimestamp(inputTimestamp int) {
	timestamp = int(math.Max(float64(timestamp), float64(inputTimestamp))) + 1
}

func (s *server) SendMessage(ctx context.Context, msg *pb.Message) (*pb.Response, error) {
	broadcast(msg.Sender, msg.Message)
	updateTimestamp(int(msg.Timestamp))
	return &pb.Response{Status: 1}, nil
}

func (s *server) LeaveChatServer(ctx context.Context, user *pb.User) (*pb.Response, error) {
	users[user.Uuid].leave <- true
	delete(users, user.Uuid)
	updateTimestamp(int(user.Timestamp))
	broadcast("Server", (user.Name + " left the server"))
	return &pb.Response{Status: 1}, nil
}

func (s *server) JoinChatServer(user *pb.User, srv pb.ChatService_JoinChatServerServer) error {
	users[user.Uuid] = newUserEntity(srv)
	updateTimestamp(int(user.Timestamp))
	broadcast("Server", (user.Name + " has joined the chat room"))

	defer func() {
		if err := recover(); err != nil {
			log.Printf("panic: %v", err)
			os.Exit(1)
		}
	}()
	<-users[user.Uuid].leave
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
