package raft

type LogEntry struct {
	Term    int
	Command string
}
