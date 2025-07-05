// main.go
package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

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
	node.PeerInfo.Port = config.ListenPort
	for _, peer := range config.Peers {
		node.Cluster = append(node.Cluster, raft.Peer{
			ID:      peer.ID,
			Address: peer.Address,
			Port:    peer.Port,
		})
	}
	node.ViewNode()

	acceptCounter := 0
	bindAddr := node.PeerInfo.Address + ":" + strconv.Itoa(node.PeerInfo.Port)
	ln, err := net.Listen("tcp", bindAddr)
	ln.Addr()
	if err != nil {
		log.Fatalf("❌ Failed to bind to %s: %v", bindAddr, err)
	}
	fmt.Println("✅ Listening on", bindAddr)
	time.Sleep(1 * time.Second) // Keep the server running for a while to accept connections
	timeout := 2 * time.Second
	for _, peer := range node.Cluster {
		address := peer.Address + ":" + strconv.Itoa(peer.Port)
		for i := 1; i <= 5; i++ {
			conn, err := net.DialTimeout("tcp", address, timeout)
			if err == nil {
				fmt.Printf("✅ Connected to peer %s at %s (attempt %d)\n", peer.ID, address, i)
				acceptCounter++
				conn.Close()
				break
			}
			fmt.Printf("⏳ Retry %d: Failed to connect to peer %s at %s\n", i, peer.ID, address)
			time.Sleep(1 * time.Second)
		}
	}

	time.Sleep(1 * time.Second) // Allow time for connections to be established

	message := raft.Message{
		Type:     raft.Heartbeat,
		Sender:   node.ID,
		Term:     node.CurrentTerm,
		Receiver: "",
		Command:  "",
		Data:     nil,
	}

	go func() {
		for {
			Listenconn, err := ln.Accept()
			if err != nil {
				log.Println("❌ Failed to accept connection:", err)
				continue
			}
			go node.HandleConnection(Listenconn) // Handle each connection concurrently
		}
	}()

	for _, peer := range node.Cluster {
		address := peer.Address + ":" + strconv.Itoa(peer.Port)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			log.Printf("❌ Failed to connect to peer %s at %s: %v", peer.ID, address, err)
			continue
		}
		defer conn.Close()

		if err := node.SendMessage(message, address); err != nil {
			log.Printf("❌ Failed to send message to %s: %v", address, err)
		}
		time.Sleep(100 * time.Millisecond) // Small delay to avoid overwhelming the network
	}

}
