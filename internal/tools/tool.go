package tools

import "encoding/json"

type Tool interface {
	Name() string
	Description() string
	Execute(args json.RawMessage) (string, error)
	Schema() interface{} // Return a struct that can be marshaled to JSON Schema
}
