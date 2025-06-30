package raft

import "fmt"

type Role string

const (
	Follower  Role = "Follower"
	Candidate Role = "Candidate"
	Leader    Role = "Leader"
)

type ClusterStructure struct {
	NodeID  string
	Address string
	Port    int
}

type RaftNode struct {
	// Nodes structure
	ID            string
	CurrentLeader string
	Cluster       []ClusterStructure

	// Raft state
	CurrentTerm int
	VotedFor    string
	State       Role

	// Log entries
	Log []LogEntry
}

func NewRaftNode(id string) *RaftNode {
	return &RaftNode{
		ID:            id,
		CurrentLeader: "",
		Cluster:       []ClusterStructure{},
		CurrentTerm:   0,
		VotedFor:      "",
		Log:           []LogEntry{},
		State:         Follower,
	}
}

func (n *RaftNode) ViewNode() {
	fmt.Printf("Node ID: %s\n", n.ID)
	fmt.Printf("Current Leader: %s\n", n.CurrentLeader)
	fmt.Println("Cluster Nodes:")
	for _, node := range n.Cluster {
		fmt.Printf("  Node ID: %s, Address: %s, Port: %d\n", node.NodeID, node.Address, node.Port)
	}
	fmt.Printf("Current Term: %d\n", n.CurrentTerm)
	fmt.Printf("Voted For: %s\n", n.VotedFor)
	fmt.Printf("State: %s\n", n.State)
	fmt.Println("Log Entries:")
	for i, entry := range n.Log {
		fmt.Printf("  Entry %d: Term %d, Command %s\n", i, entry.Term, entry.Command)
	}
}
