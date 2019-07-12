package mrl

import (
	"time"
)

type MachineReadableLog struct {
	Type    string    `json:"type"`
	Hash    string    `json:"hash,omitempty"`
	Version string    `json:"version,omitempty"`
	Name    string    `json:"name"`
	Time    time.Time `json:"time"`
}
