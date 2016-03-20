package protocol

import "time"

type CeleryTask struct {
	Name       string                 `json:"task"`
	ID         string                 `json:"id"`
	Args       []interface{}          `json:"args,omitempty"`
	KWArgs     map[string]interface{} `json:"kwargs,omitempty"`
	Retries    int                    `json:"retries,omitempty"`
	ETA        *time.Time             `json:"eta,omitempty"`
	ExpiresAt  *time.Time             `json:"expires,omitempty"`
	IsUTC      bool                   `json:"utc,omitempty"`
	Callbacks  []string               `json:"callbacks,omitempty"`
	Errbacks   []string               `json:"errbacks,omitempty"`
	TimeLimits [2]*float64            `json:"timelimit,omitempty"`
	TaskSet    *string                `json:"taskset,omitempty"`
	Chord      *string                `json:"chord,omitempty"`
}

type CeleryResult struct {
	Status    string      `json:"status"`
	Traceback *string     `json:"traceback"`
	Result    interface{} `json:"result"`
	TaskID    string      `json:"task_id"`
	Children  []string    `json:"children"`
}

type CeleryExceptionResult struct {
	Message string `json:"exc_message"`
	Type    string `json:"exc_type"`
}
