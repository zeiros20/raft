// main.go
package main

import (
	"log"
	"github.com/zeiros20/raft/raft"
)

func main() {
	node := raft.NewRaftNode("node1")
	log.Println("Raft node started:", node.ID)
}
