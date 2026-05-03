package main

import (
	"encoding/json"
	"os"

	"github.com/openai/openai-go/v3"
)

// Tools is exported so it can be used from other packages/files (e.g., app.Tools).
// If you only need it within the same package, keep it unexported as `tools`.
var Tools = []openai.ChatCompletionToolUnionParam{
	openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
		Name:        "Read",
		Description: openai.String("Read and return the contents of a file"),
		Parameters: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]any{
				"file_path": map[string]any{
					"type":        "string",
					"description": "The path to the file to read",
				},
			},
			"required": []string{"file_path"},
		},
	}),
}

type ToolHandler func(args string) string

var toolHandlers = map[string]ToolHandler{
	"Read": handleReadTool,
}

func handleToolCall(toolCalls []openai.ChatCompletionMessageToolCallUnion) string {
	if len(toolCalls) == 0 {
		return ""
	}
	tc := toolCalls[0]
	if handler, exists := toolHandlers[tc.Function.Name]; exists {
		return handler(tc.Function.Arguments)
	}
	return ""
}

func handleReadTool(args string) string {
	var parsedArgs map[string]string
	err := json.Unmarshal([]byte(args), &parsedArgs)
	if err != nil {
		return ""
	}
	filePath := parsedArgs["file_path"]
	content, err := os.ReadFile(filePath)
	if err != nil {
		return ""
	}
	return string(content)
}
