package config

import "sync"

type BoardCastCfg struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type SecureBoardCast struct {
	Lock sync.RWMutex
	Data map[string]string
}
