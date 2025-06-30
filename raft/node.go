package raft

import "fmt"

type Role string

const (
	Follower  Role = "Follower"
	Candidate Role = "Candidate"
	Leader    Role = "Leader"
)

type RaftNode struct {
	// Nodes structure
	ID            string
	CurrentLeader string
	Cluster       []Peer

	// Raft state
	CurrentTerm int
	VotedFor    string
	State       Role

	// Log entries
	Log []LogEntry

	// Network communication
	PeerInfo Peer
}

func NewRaftNode(id string) *RaftNode {
	return &RaftNode{
		ID:            id,
		CurrentLeader: "",
		Cluster:       []Peer{},
		CurrentTerm:   0,
		VotedFor:      "",
		Log:           []LogEntry{},
		State:         Follower,
	}
}

func (n *RaftNode) ViewNode() {
	fmt.Printf("Node ID: %s\n", n.ID)
	fmt.Printf("Address: %s\n", n.PeerInfo.Address)
	fmt.Printf("Current Leader: %s\n", n.CurrentLeader)
	fmt.Println("Cluster Nodes:")
	for _, node := range n.Cluster {
		fmt.Printf("  Node ID: %s, Address: %s\n", node.ID, node.Address)
	}
	fmt.Printf("Current Term: %d\n", n.CurrentTerm)
	fmt.Printf("Voted For: %s\n", n.VotedFor)
	fmt.Printf("State: %s\n", n.State)
	fmt.Println("Log Entries:")
	for i, entry := range n.Log {
		fmt.Printf("  Entry %d: Term %d, Command %s\n", i, entry.Term, entry.Command)
	}
}
