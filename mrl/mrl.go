package mrl

import (
	"time"
)

type MachineReadableLog struct {
	Type     string      `json:"type"`
	Hash     string      `json:"hash,omitempty"`
	Version  string      `json:"version,omitempty"`
	Name     string      `json:"name,omitempty"`
	Metadata interface{} `json:"metadata,omitempty"`
	Result   int         `json:"result,omitempty"`
	Time     time.Time   `json:"time"`
	Message  string      `json:"message,omitempty"`
}
