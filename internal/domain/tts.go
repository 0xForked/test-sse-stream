package domain

import (
	"log"
)

type TTS struct {
	Sequence int    `json:"sequence"`
	Locket   int    `json:"locket"`
	Message  string `json:"message"`
}

type TTSEvent struct {
	// Events are pushed to this channel by the main events-gathering routine
	Message chan *TTS

	// New client connections
	NewClients chan chan string

	// Closed client connections
	ClosedClients chan chan string

	// Total client connections
	TotalClients map[chan string]bool
}

type ClientChan chan string

func (stream *TTSEvent) Listen() {
	for {
		select {
		// Add new available client
		case client := <-stream.NewClients:
			stream.TotalClients[client] = true
			log.Printf("Client added. %d registered clients", len(stream.TotalClients))

		// Remove closed client
		case client := <-stream.ClosedClients:
			delete(stream.TotalClients, client)
			close(client)
			log.Printf("Removed client. %d registered clients", len(stream.TotalClients))

		// Broadcast message to client
		case eventMsg := <-stream.Message:
			for clientMessageChan := range stream.TotalClients {
				select {
				case clientMessageChan <- eventMsg.Message:
				default:
					// If the client channel is full, skip sending the message to that client
					log.Println("Client channel is full. Skipping message.")
				}
			}
		}
	}
}
