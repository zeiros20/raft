package raft

import (
	"context"
	"fmt"
	"net"
	"sync/atomic"
	"time"
)

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

	// Clock for timing
	Clock int64

	// Buffered channels
	HeartbeatTimeout chan bool
	HeartbeatRecived chan bool
	BecameLeader     chan bool
	ClockStopped     chan struct{}

	// Votes
	Votes     int
	VotedFrom []Peer
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
		PeerInfo: Peer{
			ID:      id,
			Address: "localhost",
			Port:    8080,
		},
		Clock: 0,

		HeartbeatTimeout: make(chan bool, 1),     // Buffered channel to avoid blocking
		HeartbeatRecived: make(chan bool, 1),     // Buffered channel to avoid blocking
		BecameLeader:     make(chan bool, 1),     // Buffered channel to signal leader election
		ClockStopped:     make(chan struct{}, 1), // Buffered channel to signal clock stop

		Votes:     1,
		VotedFrom: []Peer{},
	}
}

func (n *RaftNode) ViewNode() {
	fmt.Printf("Node ID: %s\n", n.ID)
	fmt.Printf("Address: %s\n", n.PeerInfo.Address)
	fmt.Printf("Current Leader: %s\n", n.CurrentLeader)
	fmt.Println("Cluster Nodes:")
	for _, node := range n.Cluster {
		fmt.Printf("Node ID: %s, Address: %s, Port %d\n", node.ID, node.Address, node.Port)
	}
	fmt.Printf("Current Term: %d\n", n.CurrentTerm)
	fmt.Printf("Voted For: %s\n", n.VotedFor)
	fmt.Printf("State: %s\n", n.State)
	fmt.Println("Log Entries:")
	for i, entry := range n.Log {
		fmt.Printf("  Entry %d: Term %d, Command %s\n", i, entry.Term, entry.Command)
	}
}

func (n *RaftNode) ConnectionLoop() {
	address := fmt.Sprintf("%s:%d", n.PeerInfo.Address, n.PeerInfo.Port)
	InfoMessage(fmt.Sprintf("Node %s listening on %s", n.ID, address))
	// Here you would implement the logic to accept connections and handle messages
	listener, err := net.Listen("tcp", address)
	if HandleError(err, fmt.Sprintf("❌ Error starting listener on %s", address)) {
		return
	}

	for {
		conn, err := listener.Accept()
		if HandleError(err, "❌ Error accepting connection") {
			continue
		}
		go n.HandleConnection(conn) // Handle each connection in a separate goroutine
	}
}

func (n *RaftNode) Heartbeat() {

	for _, peer := range n.Cluster {
		address := fmt.Sprintf("%s:%d", peer.Address, peer.Port)
		message := NewHeartbeatMessage(n.PeerInfo.ID, fmt.Sprintf("%s:%d", n.PeerInfo.Address, n.PeerInfo.Port), n.CurrentTerm, peer.ID, address)
		if HandleError(n.SendMessage(message, address), fmt.Sprintf("❌ Failed to send Heartbeat message to %s", address)) {
			continue
		} else {
			InfoMessage(fmt.Sprintf("✅ Sent Heartbeat message to %s", address))
		}
	}
	time.Sleep(100 * time.Millisecond) // Small delay to avoid overwhelming the network
}

func (n *RaftNode) StartElection() {

}

func (n *RaftNode) NodeClockStart(ctx context.Context) {
	go func() {
		start := time.Now()
		ticker := time.NewTicker(1 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				select {
				case n.ClockStopped <- struct{}{}: // signal clean shutdown
				default:
				}
				return
			case <-ticker.C:
				elapsed := time.Since(start).Milliseconds()
				// InfoLogger.Printf("Node %s %s clock: %d ms", n.ID, n.State, elapsed)
				if n.State == Follower && elapsed > 500 {
					select {
					case n.HeartbeatTimeout <- true:
						InfoLogger.Printf("Node %s heartbeat timeout: %d ms", n.ID, elapsed)
					default:
					}
				}
				atomic.StoreInt64(&n.Clock, elapsed)
			}
		}
	}()
}

func (n *RaftNode) GetClock() int64 {
	return atomic.LoadInt64(&n.Clock)
}

func (n *RaftNode) Run() {
	for {
		switch n.State {
		case Follower:
			n.FollowerLoop()
		case Candidate:
			n.CandidateLoop()
		case Leader:
			n.LeaderLoop()
		default:
			ErrorLogger.Printf("❌ Unknown state: %s", n.State)
		}
	}

}

func (n *RaftNode) FollowerLoop() {
	ctx, cancel := context.WithCancel(context.Background())
	n.NodeClockStart(ctx)

	for {
		select {
		case <-n.HeartbeatTimeout:
			cancel()         // Cancel the clock context to stop the clock
			<-n.ClockStopped // ✅ wait for the clock to shut down safely
			InfoMessage(fmt.Sprintf("Node %s received heartbeat timeout, transitioning to Candidate state", n.ID))
			n.State = Candidate
			n.CurrentLeader = ""
			return
		case <-n.HeartbeatRecived:
			InfoMessage(fmt.Sprintf("Node %s received heartbeat, remaining in Follower state", n.ID))
			cancel()         // Cancel the clock context to stop the clock
			<-n.ClockStopped // ✅ wait for the clock to shut down safely
			ctx, cancel = context.WithCancel(context.Background())
			n.NodeClockStart(ctx) // Restart the clock for the follower
		default:
			// Continue in follower state, waiting for heartbeats or election timeouts
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func (n *RaftNode) CandidateLoop() {

}

func (n *RaftNode) LeaderLoop() {
	for {
		n.Heartbeat()
		time.Sleep(400 * time.Millisecond) // Leader sends heartbeats
	}
	return
}
