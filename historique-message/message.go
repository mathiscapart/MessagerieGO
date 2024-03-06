package main

import (
	"context"
	"fmt"
	"historique-message/db"
	"historique-message/proto"
)

type MessageServer struct {
	proto.UnimplementedGreeterServer
}

func NewMessageServer() *MessageServer {
	return &MessageServer{}
}

func (m MessageServer) Catch(ctx context.Context, req *proto.Content) (*proto.Reply, error) {
	fmt.Println("message : ", req.GetMessage())
	var content db.Message
	content.Exp = req.Expediteur
	content.Dest = req.Destinataire
	content.Message = req.Message
	db.InsertUser(content)
	return &proto.Reply{}, nil
}
