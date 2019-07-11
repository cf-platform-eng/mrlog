package mrl

import (
	"time"
)

type MachineReadableLog struct {
	Type     string    `json:"type"`
	Filename string    `json:"filename"`
	Hash     string    `json:"hash"`
	Version  string    `json:"version"`
	Time     time.Time `json:"time"`
}
