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

type LogEntry struct {
	Term    int
	Command string
}

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
	Type            string `json:"type"`
	Sender          string `json:"sender"`
	SenderAddress   string `json:"sender_address"`
	Term            int    `json:"term"`
	Receiver        string `json:"receiver"`
	ReceiverAddress string `json:"receiver_address"`
	Command         string `json:"command"`
	Data            []byte `json:"data"`
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

func NewRequestVoteMessage(sender string, senderAddress string, term int, receiver string, receiverAddress string) Message {
	return Message{
		Type:            RequestVote,
		Sender:          sender,
		SenderAddress:   senderAddress,
		Term:            term,
		Receiver:        receiver,
		ReceiverAddress: receiverAddress,
		Command:         "",
		Data:            nil,
	}
}

func NewAppendEntriesMessage(sender string, senderAddress string, term int, receiver string, receiverAddress string, command string, data []byte) Message {
	return Message{
		Type:            AppendEntries,
		Sender:          sender,
		SenderAddress:   senderAddress,
		Term:            term,
		Receiver:        receiver,
		ReceiverAddress: receiverAddress,
		Command:         command,
		Data:            data,
	}
}

func NewResponseVoteMessage(sender string, senderAddress string, term int, receiver string, receiverAddress string) Message {
	return Message{
		Type:            ResponseVote,
		Sender:          sender,
		SenderAddress:   senderAddress,
		Term:            term,
		Receiver:        receiver,
		ReceiverAddress: receiverAddress,
		Command:         "",
		Data:            nil,
	}
}

func NewHeartbeatMessage(sender string, senderAddress string, term int, receiver string, receiverAddress string) Message {
	return Message{
		Type:            Heartbeat,
		Sender:          sender,
		SenderAddress:   senderAddress,
		Term:            term,
		Receiver:        receiver,
		ReceiverAddress: receiverAddress,
		Command:         "",
		Data:            nil,
	}
}

func NewHeartbeatResponseMessage(sender string, senderAddress string, term int, receiver string, receiverAddress string) Message {
	return Message{
		Type:            HeartbeatResponse,
		Sender:          sender,
		SenderAddress:   senderAddress,
		Term:            term,
		Receiver:        receiver,
		ReceiverAddress: receiverAddress,
		Command:         "",
		Data:            nil,
	}
}

func (m *Message) IsHeartbeat() bool {
	return m.Type == Heartbeat || m.Type == HeartbeatResponse
}

func (m *Message) IsValidMessage() bool {
	return (m.Type == RequestVote || m.Type == AppendEntries || m.Type == ResponseVote || m.Type == Heartbeat || m.Type == HeartbeatResponse) &&
		m.Sender != "" && m.Term >= 0 && m.Receiver != ""
}
