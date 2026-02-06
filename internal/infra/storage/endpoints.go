package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Octrafic/octrafic-cli/internal/core/parser"
)

const endpointsFile = "endpoints.json"

// SaveEndpoints saves parsed endpoints to endpoints.json
func SaveEndpoints(projectID string, endpoints []parser.Endpoint, isTemporary bool) error {
	projectPath, err := GetProjectPathByType(projectID, isTemporary)
	if err != nil {
		return fmt.Errorf("failed to get project path: %w", err)
	}

	filePath := filepath.Join(projectPath, endpointsFile)
	data, err := json.MarshalIndent(endpoints, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal endpoints: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write endpoints file: %w", err)
	}

	return nil
}

// LoadEndpoints loads endpoints from endpoints.json
func LoadEndpoints(projectID string, isTemporary bool) ([]parser.Endpoint, error) {
	projectPath, err := GetProjectPathByType(projectID, isTemporary)
	if err != nil {
		return nil, fmt.Errorf("failed to get project path: %w", err)
	}

	filePath := filepath.Join(projectPath, endpointsFile)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read endpoints file: %w", err)
	}

	var endpoints []parser.Endpoint
	if err := json.Unmarshal(data, &endpoints); err != nil {
		return nil, fmt.Errorf("failed to unmarshal endpoints: %w", err)
	}

	return endpoints, nil
}

// GetEndpointsList returns comma-separated list of endpoints for system prompt
// Format: "GET /users, POST /users, GET /users/{id}, PUT /users/{id}, DELETE /users/{id}, GET /health"
func GetEndpointsList(endpoints []parser.Endpoint) string {
	if len(endpoints) == 0 {
		return "No endpoints available"
	}

	var result string
	for i, ep := range endpoints {
		if i > 0 {
			result += ", "
		}
		result += fmt.Sprintf("%s %s", ep.Method, ep.Path)
	}

	return result
}

// HasEndpoints checks if endpoints.json exists for a project
func HasEndpoints(projectID string, isTemporary bool) bool {
	projectPath, err := GetProjectPathByType(projectID, isTemporary)
	if err != nil {
		return false
	}

	filePath := filepath.Join(projectPath, endpointsFile)
	_, err = os.Stat(filePath)
	return err == nil
}
