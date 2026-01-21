// pkg/jupyter_client/message_types.go

package jupyterclient

import (
	"encoding/json"
	"time"
)

type Header struct {
	MsgID    string    `json:"msg_id"`
	MsgType  string    `json:"msg_type"`
	Username string    `json:"username"`
	Session  string    `json:"session"`
	Date     time.Time `json:"date"`
	Version  string    `json:"version"`
}

type Message struct {
	Header       Header          `json:"header"`
	ParentHeader Header          `json:"parent_header"`
	Metadata     json.RawMessage `json:"metadata"`
	Content      json.RawMessage `json:"content"`
	Buffers      []json.RawMessage `json:"buffers"`
}

type StreamContent struct {
	Name string `json:"name"` // "stdout" or "stderr"
	Text string `json:"text"`
}

type DisplayDataContent struct {
	Data      map[string]any `json:"data"`
	Metadata  map[string]any `json:"metadata"`
	Transient map[string]any `json:"transient"`
}

type ExecuteResultContent struct {
	ExecutionCount int                    `json:"execution_count"`
	Data           map[string]any `json:"data"`
	Metadata       map[string]any `json:"metadata"`
}

type ErrorContent struct {
	Ename     string   `json:"ename"`
	Evalue    string   `json:"evalue"`
	Traceback []string `json:"traceback"`
}
