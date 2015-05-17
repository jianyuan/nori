package task

import "time"

type Task struct {
	Name       string                 `json:"task"`
	ID         string                 `json:"id"`
	Args       []string               `json:"args,omitempty"`
	KWArgs     map[string]interface{} `json:"kwargs,omitempty"`
	Retries    int                    `json:"retries,omitempty"`
	ETA        *time.Time             `json:"eta,omitempty"`
	Expires    *time.Time             `json:"expires,omitempty"`
	UTC        bool                   `json:"utc,omitempty"`
	Callbacks  []string               `json:"callbacks,omitempty"`
	Errbacks   []string               `json:"errbacks,omitempty"`
	TimeLimits [2]*float64            `json:"timelimit,omitempty"`
	TaskSet    *string                `json:"taskset,omitempty"`
	Chord      *string                `json:"chord,omitempty"`
}
