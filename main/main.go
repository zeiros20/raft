// main.go
package main

import (
	"github.com/zeiros20/raft/raft"
)

func main() {
	node := raft.NewRaftNode("node1")
	node.ViewNode()
}
