package raft

import (
	"os"

	"gopkg.in/yaml.v3"
)

const (
	RequestVote       string = "RequestVote"
	AppendEntries     string = "AppendEntries"
	ResponseVote      string = "ResponseVote"
	Heartbeat         string = "Heartbeat"
	HeartbeatResponse string = "HeartbeatResponse"
)

type Peer struct {
	ID      string `yaml:"id"`
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
}

type Config struct {
	NodeID        string `yaml:"node_id"`
	ListenAddress string `yaml:"listen_address"`
	ListenPort    int    `yaml:"listen_port"`
	Peers         []Peer `yaml:"peers"`
}

type Message struct {
	Type     string `json:"type"`
	Sender   string `json:"sender"`
	Term     int    `json:"term"`
	Receiver string `json:"receiver"`
	Command  string `json:"command"`
	Data     []byte `json:"data"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
