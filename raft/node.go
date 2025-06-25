// raft/node.go
package raft

import "fmt"

type Role string

const (
	Follower  Role = "Follower"
	Candidate Role = "Candidate"
	Leader    Role = "Leader"
)

type RaftNode struct {
	ID          string
	CurrentTerm int
	VotedFor    string
	Log         []LogEntry
	State       Role
}

func NewRaftNode(id string) *RaftNode {
	return &RaftNode{
		ID:          id,
		CurrentTerm: 0,
		VotedFor:    "",
		Log:         []LogEntry{},
		State:       Follower,
	}
}

func (rn *RaftNode) String() string {
	return fmt.Sprintf("Node %s [%s] term=%d", rn.ID, rn.State, rn.CurrentTerm)
}
