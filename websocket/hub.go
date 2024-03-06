package main

import (
	"context"
	"encoding/json"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"server/proto"
)

type Hub struct {
	clients map[*Client]bool

	broadcast chan Expedition

	register chan *Client

	unregister chan *Client
}

type Contenu struct {
	Destinataire string `json:"destinataire"`

	Message string `json:"message"`
}

type Expedition struct {
	Expediteur string

	Content []byte
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan Expedition),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	name_expediteur := ""
	for {
		select {
		case client := <-h.register:
			name_expediteur = client.name
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case content := <-h.broadcast:
			var sender Contenu
			err := json.Unmarshal(content.Content, &sender)
			if err != nil {
				log.Println(err)
			}

			for client := range h.clients {
				if sender.Destinataire == client.name {
					select {
					case client.send <- []byte("destinataire : " + name_expediteur + " / " + sender.Message):
						conn, err2 := grpc.Dial("localhost:3000", grpc.WithTransportCredentials(insecure.NewCredentials()))
						if err2 != nil {
							log.Println(err2)
						}

						greeterClient := proto.NewGreeterClient(conn)

						_, err2 = greeterClient.Catch(context.TODO(), &proto.Content{Destinataire: sender.Destinataire, Expediteur: name_expediteur, Message: sender.Message})
						if err2 != nil {
							log.Println(err2)
							return
						}

					default:
						close(client.send)
						delete(h.clients, client)
					}
				}

			}
		}
	}
}
