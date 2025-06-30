// main.go
package main

import (
	"flag"

	"github.com/zeiros20/raft/raft"
)

func main() {
	ConfigPath := flag.String("config", "config/node1.yaml", "Path to the configuration file")
	flag.Parse()
	node := raft.NewRaftNode("temp")
	config, err := raft.LoadConfig(*ConfigPath)
	if err != nil {
		panic(err)
	}
	node.ID = config.NodeID
	node.PeerInfo.ID = config.NodeID
	node.PeerInfo.Address = config.ListenAddress
	for _, peer := range config.Peers {
		node.Cluster = append(node.Cluster, raft.Peer{
			ID:      peer.ID,
			Address: peer.Address,
		})
	}
	node.ViewNode()

}
