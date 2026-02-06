package agent

import (
	"github.com/Octrafic/octrafic-cli/internal/llm/common"
	"fmt"
)

type ChatMessage struct {
	Role             string                `json:"role"`
	Content          string                `json:"content"`
	ReasoningContent string                `json:"reasoning_content,omitempty"`
	FunctionCalls    []ToolCall            `json:"function_calls,omitempty"`
	FunctionResponse *FunctionResponseData `json:"function_response,omitempty"`
}

type FunctionResponseData struct {
	ID       string                 `json:"id"` // tool_use_id from original tool use block
	Name     string                 `json:"name"`
	Response map[string]any `json:"response"`
}

type ChatResponse struct {
	Message      string     `json:"message"`
	Reasoning    string     `json:"reasoning,omitempty"`
	ToolCalls    []ToolCall `json:"tool_calls,omitempty"`
	TokensUsed   int64      `json:"tokens_used,omitempty"`   // Total tokens (input + output)
	InputTokens  int64      `json:"input_tokens,omitempty"`  // Input tokens used
	OutputTokens int64      `json:"output_tokens,omitempty"` // Output tokens used
}

type ToolCall struct {
	ID               string                 `json:"id,omitempty"`
	Name             string                 `json:"name"`
	Arguments        map[string]any `json:"arguments"`
	ThoughtSignature string                 `json:"thought_signature,omitempty"` // base64 encoded []byte (not used by Claude)
}

// ReasoningCallback is called for each chunk as it's streamed
// isThought=true for reasoning chunks, false for text chunks
type ReasoningCallback func(chunk string, isThought bool)

// getMainAgentTools returns the tools for the main agent
func getMainAgentTools() []common.Tool {
	return []common.Tool{
		{
			Name:        "get_endpoints_details",
			Description: "Get detailed information about specified endpoints including description, parameters, security, request body, and responses.",
			InputSchema: map[string]any{
				"type":                 "object",
				"additionalProperties": false,
				"properties": map[string]any{
					"endpoints": map[string]any{
						"type": "array",
						"items": map[string]any{
							"type":                 "object",
							"additionalProperties": false,
							"properties": map[string]any{
								"path": map[string]any{
									"type":        "string",
									"description": "Endpoint path (e.g., /users)",
								},
								"method": map[string]any{
									"type":        "string",
									"description": "HTTP method (GET, POST, PUT, DELETE, PATCH)",
								},
							},
							"required": []string{"path", "method"},
						},
					},
				},
				"required": []string{"endpoints"},
			},
		},
		{
			Name:        "GenerateTestPlan",
			Description: "Generate test cases for API endpoints. Describe endpoints with all relevant details from get_endpoints_details.",
			InputSchema: map[string]any{
				"type":                 "object",
				"additionalProperties": false,
				"properties": map[string]any{
					"what": map[string]any{
						"type":        "string",
						"description": "Detailed endpoint description including: method, path, what it does, authentication requirements (Security field), request body schema, expected responses, parameters. Be thorough!",
					},
					"focus": map[string]any{
						"type":        "string",
						"description": "Testing focus: 'happy path' (basic success), 'authentication' (with/without auth), 'error handling' (validation, 404, etc), 'all aspects' (comprehensive)",
					},
				},
				"required": []string{"what", "focus"},
			},
		},
		{
			Name:        "ExecuteTestGroup",
			Description: "Execute a group of tests against the API. Tests are run locally by the CLI and results are returned. Call this AFTER GenerateTestPlan to actually run the tests.",
			InputSchema: map[string]any{
				"type":                 "object",
				"additionalProperties": false,
				"properties": map[string]any{
					"tests": map[string]any{
						"type": "array",
						"items": map[string]any{
							"type":                 "object",
							"additionalProperties": false,
							"properties": map[string]any{
								"method": map[string]any{
									"type":        "string",
									"description": "HTTP method (GET, POST, PUT, DELETE, etc)",
								},
								"endpoint": map[string]any{
									"type":        "string",
									"description": "API endpoint path (e.g., /api/health)",
								},
								"headers": map[string]any{
									"type":                 []any{"object", "null"},
									"additionalProperties": false,
									"description":          "Optional HTTP headers",
								},
								"body": map[string]any{
									"type":        []any{"string", "null"},
									"description": "Optional request body (JSON string)",
								},
								"requires_auth": map[string]any{
									"type":        "boolean",
									"description": "Whether authentication is required for this test",
								},
							},
							"required": []string{"method", "endpoint", "headers", "body", "requires_auth"},
						},
					},
				},
				"required": []string{"tests"},
			},
		},
	}
}

func (a *Agent) Chat(messages []ChatMessage, thinkingEnabled bool, endpointsList ...string) (*ChatResponse, error) {
	systemPrompt := buildSystemPrompt(a.baseURL, endpointsList...)
	tools := getMainAgentTools()
	return a.baseAgent.Chat(systemPrompt, tools, messages, thinkingEnabled)
}

func buildSystemPrompt(baseURL string, endpointsList ...string) string {
	endpointsInfo := ""
	if len(endpointsList) > 0 && endpointsList[0] != "" {
		endpointsInfo = fmt.Sprintf(`

# Available Endpoints
%s
`, endpointsList[0])
	}

	return fmt.Sprintf(`Role: API testing assistant
Base URL: %s%s

# Rules
1. Be proactive - don't over-ask for clarification
2. When user mentions endpoints (e.g., "users", "auth"), automatically fetch details and show/test them
3. Default to "happy path" tests unless user specifies otherwise
4. List endpoints from above, use get_endpoints_details for technical details
5. Use ONE tool per response

# Tools

## get_endpoints_details
Fetch detailed specs (params, auth, schemas). Use when:
- Need technical details for response/tests
- User asks about specific endpoint behavior

## GenerateTestPlan
Generate tests. Parameters:
- what: endpoint details from get_endpoints_details
- focus: default "happy path", or user's choice

## ExecuteTestGroup
Run tests after GenerateTestPlan.

# Behavior
- User says "users" → fetch details, show info OR generate tests
- User says "test X" → fetch details, generate & run tests
- User says "list endpoints" → show list from above (no tool call)
- Default focus: "happy path"
- requires_auth=true → send auth, requires_auth=false → no auth`, baseURL, endpointsInfo)
}

func (a *Agent) ChatStream(messages []ChatMessage, thinkingEnabled bool, callback ReasoningCallback, endpointsList ...string) (*ChatResponse, error) {
	systemPrompt := buildSystemPrompt(a.baseURL, endpointsList...)
	tools := getMainAgentTools()
	return a.baseAgent.ChatStream(systemPrompt, tools, messages, thinkingEnabled, callback)
}
