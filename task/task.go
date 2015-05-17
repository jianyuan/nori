package task

import "time"

type Task struct {
	Name       string                 `json:"task"`
	ID         string                 `json:"id"`
	Args       []string               `json:"args"`
	KWArgs     map[string]interface{} `json:"kwargs"`
	Retries    int                    `json:"retries"`
	ETA        *time.Time             `json:"eta"`
	Expires    *time.Time             `json:"expires"`
	UTC        bool                   `json:"utc"`
	Callbacks  []string               `json:"callbacks"`
	Errbacks   []string               `json:"errbacks"`
	TimeLimits [2]*float64            `json:"timelimit"`
	TaskSet    *string                `json:"taskset"`
	Chord      *string                `json:"chord"`
}
