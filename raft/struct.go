package raft

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Peer struct {
	ID      string `yaml:"id"`
	Address string `yaml:"address"`
}

type Config struct {
	NodeID        string `yaml:"node_id"`
	ListenAddress string `yaml:"listen_address"`
	Peers         []Peer `yaml:"peers"`
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
