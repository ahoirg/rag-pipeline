package models

import "time"

type ApiResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message,omitempty"`
	Query     string      `json:"query,omitempty"`
	Answer    string      `json:"answer,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}
