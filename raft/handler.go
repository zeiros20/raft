package raft

import (
	"encoding/json"
	"log"
	"net"
)

func handleRequestVote(message Message, conn net.Conn) {
	log.Printf("Received RequestVote from %s for term %d", message.Sender, message.Term)
}

func handleAppendEntries(message Message, conn net.Conn) {
	log.Printf("Received AppendEntries from %s for term %d", message.Sender, message.Term)
}

func handleResponseVote(message Message, conn net.Conn) {
	log.Printf("Received ResponseVote from %s for term %d", message.Sender, message.Term)
}

func handleHeartbeat(message Message, conn net.Conn) {
	log.Printf("Received Heartbeat from %s for term %d", message.Sender, message.Term)
}

func handleHeartbeatResponse(message Message, conn net.Conn) {
	log.Printf("Received HeartbeatResponse from %s for term %d", message.Sender, message.Term)
}

func (n *RaftNode) HandleConnection(conn net.Conn) {
	defer conn.Close()
	decoder := json.NewDecoder(conn)
	var message Message

	if err := decoder.Decode(&message); err != nil {
		return
	}

	switch message.Type {
	case RequestVote:
		handleRequestVote(message, conn)
	case AppendEntries:
		handleAppendEntries(message, conn)
	case ResponseVote:
		handleResponseVote(message, conn)
	case Heartbeat:
		handleHeartbeat(message, conn)
	case HeartbeatResponse:
		handleHeartbeatResponse(message, conn)
	default:
		log.Printf("❌ Unknown message type: %s from %s", message.Type, message.Sender)
	}
}

func (n *RaftNode) SendMessage(message Message, address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	defer conn.Close()

	encoder := json.NewEncoder(conn)
	if err := encoder.Encode(message); err != nil {
		log.Println("❌ Error encoding message:", err)
		return err
	}

	log.Printf("✅ Sent message of type %s to %s", message.Type, address)
	return nil
}
