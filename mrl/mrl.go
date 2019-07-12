package mrl

import (
	"time"
)

type MachineReadableLog struct {
	Type     string    `json:"type"`
	Filename string    `json:"filename"`
	Hash     string    `json:"hash"`
	Version  string    `json:"version"`
	Name     string    `json:"name"`
	Time     time.Time `json:"time"`
}
