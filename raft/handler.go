package raft

import (
	"encoding/json"
	"fmt"
	"net"
)

func handleRequestVote(message Message, conn net.Conn, n *RaftNode) {

}

func handleAppendEntries(message Message, conn net.Conn, n *RaftNode) {
	logMessage := fmt.Sprintf("AppendEntries from %s for term %d", message.Sender, message.Term)
	InfoLogger.Println(logMessage)
}

func handleResponseVote(message Message, conn net.Conn, n *RaftNode) {

}

func handleHeartbeat(message Message, conn net.Conn, n *RaftNode) {

}

func handleHeartbeatResponse(message Message, conn net.Conn, n *RaftNode) {
	logMessage := fmt.Sprintf("HeartbeatResponse from %s for term %d", message.Sender, message.Term)
	InfoLogger.Println(logMessage)
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
		handleRequestVote(message, conn, n)
	case AppendEntries:
		handleAppendEntries(message, conn, n)
	case ResponseVote:
		handleResponseVote(message, conn, n)
	case Heartbeat:
		handleHeartbeat(message, conn, n)
	case HeartbeatResponse:
		handleHeartbeatResponse(message, conn, n)
	default:
		ErrorLogger.Printf("❌ Unknown message type: %s from %s", message.Type, message.Sender)
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
		ErrorLogger.Println("❌ Error encoding message:", err)
		return err
	}

	return nil
}
