package main

import "github.com/openai/openai-go/v3"

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
