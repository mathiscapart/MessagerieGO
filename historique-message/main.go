package main

import (
	"google.golang.org/grpc"
	"historique-message/db"
	"historique-message/proto"
	"log"
	"net"
)

func main() {
	db.Database()
	grpcServer := grpc.NewServer()
	msgServer := NewMessageServer()
	proto.RegisterGreeterServer(grpcServer, msgServer)
	//proto.RegisterMessageServer(grpcServer, msgServer)
	//proto.RegisterGreeterServer(grpcServer, msgServer)
	lis, err := net.Listen("tcp", "localhost:3000")
	if err != nil {
		log.Println(err)
	}
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Println(err)
	}
	defer grpcServer.GracefulStop()

}
