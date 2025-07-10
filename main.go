// main.go
package main

import (
	"flag"

	"github.com/zeiros20/raft/raft"
)

func main() {
	ConfigPath := flag.String("config", "config/node1.yaml", "Path to the configuration file")
	flag.Parse()
	node := raft.NewRaftNode(" ")
	config, err := raft.LoadConfig(*ConfigPath)
	if err != nil {
		panic(err)
	}
	node.ID = config.NodeID
	node.PeerInfo.ID = config.NodeID
	node.PeerInfo.Address = config.ListenAddress
	node.PeerInfo.Port = config.ListenPort
	for _, peer := range config.Peers {
		node.Cluster = append(node.Cluster, raft.Peer{
			ID:      peer.ID,
			Address: peer.Address,
			Port:    peer.Port,
		})
	}
	node.ViewNode()

	go node.ConnectionLoop()

	node.Run()
}
