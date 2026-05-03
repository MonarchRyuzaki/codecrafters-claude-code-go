package main

import (
	"encoding/json"
	"os"
	"os/exec"

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
	openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
		Name:        "Write",
		Description: openai.String("Write content to a file"),
		Parameters: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]any{
				"file_path": map[string]any{
					"type":        "string",
					"description": "The path of the file to write to",
				},
				"content": map[string]any{
					"type":        "string",
					"description": "The content to write to the file",
				},
			},
			"required": []string{"file_path", "content"},
		},
	}),
	openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
		Name:        "Bash",
		Description: openai.String("Execute a shell command"),
		Parameters: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]any{
				"command": map[string]any{
					"type":        "string",
					"description": "The command to execute",
				},
			},
			"required": []string{"command"},
		},
	}),
}

type ToolHandler func(args string) string

var toolHandlers = map[string]ToolHandler{
	"Read":  handleReadTool,
	"Write": handleWriteTool,
	"Bash":  handleBashTool,
}

func handleToolCall(toolCalls []openai.ChatCompletionMessageToolCallUnion, messages []openai.ChatCompletionMessageParamUnion) []openai.ChatCompletionMessageParamUnion {
	if len(toolCalls) == 0 {
		return messages
	}
	for _, tc := range toolCalls {
		if handler, exists := toolHandlers[tc.Function.Name]; exists {
			content := handler(tc.Function.Arguments)
			messages = append(messages, openai.ToolMessage(content, tc.ID))
		}
	}
	return messages
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

func handleWriteTool(args string) string {
	var parsedArgs map[string]string
	err := json.Unmarshal([]byte(args), &parsedArgs)
	if err != nil {
		return "Error parsing args"
	}
	filePath, okPath := parsedArgs["file_path"]
	content, okContent := parsedArgs["content"]
	if !okPath || !okContent {
		return "Missing required arguments"
	}
	err = os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return "Error writing file: " + err.Error()
	}
	return "File written successfully"
}

func handleBashTool(args string) string {
	var parsedArgs map[string]string
	err := json.Unmarshal([]byte(args), &parsedArgs)
	if err != nil {
		return "Error parsing args"
	}
	command, ok := parsedArgs["command"]
	if !ok {
		return "Missing required argument: command"
	}
	
	cmd := exec.Command("sh", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output) + "\nError: " + err.Error()
	}
	
	if len(output) == 0 {
		return "Command executed successfully with no output"
	}
	
	return string(output)
}
