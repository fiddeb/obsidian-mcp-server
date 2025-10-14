package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// ObsidianAPI represents the Obsidian REST API client
type ObsidianAPI struct {
	baseURL string
	token   string
	client  *http.Client
}

// Note represents an Obsidian note
type Note struct {
	Path     string                 `json:"path"`
	Content  string                 `json:"content"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// VaultInfo represents vault information
type VaultInfo struct {
	Name         string `json:"name"`
	Path         string `json:"path"`
	NotesCount   int    `json:"notes_count"`
	FoldersCount int    `json:"folders_count"`
}

// SearchResult represents a search result
type SearchResult struct {
	Path    string   `json:"path"`
	Matches []string `json:"matches"`
	Score   float64  `json:"score"`
	Content string   `json:"content"`
}

// NewObsidianAPI creates a new Obsidian API client
func NewObsidianAPI(baseURL, token string) *ObsidianAPI {
	return &ObsidianAPI{
		baseURL: strings.TrimRight(baseURL, "/"),
		token:   token,
		client:  &http.Client{},
	}
}

// normalizeNotePath ensures the path has .md extension if it's a note
func normalizeNotePath(path string) string {
	// Don't add .md if path is empty or already has an extension
	if path == "" || strings.Contains(path, ".") {
		return path
	}
	// Add .md extension
	return path + ".md"
}

func (api *ObsidianAPI) makeRequest(method, endpoint string, body interface{}) (*http.Response, error) {
	var bodyReader io.Reader

	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %v", err)
		}
		bodyReader = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, api.baseURL+endpoint, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if api.token != "" {
		req.Header.Set("Authorization", "Bearer "+api.token)
	}

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}

	return resp, nil
}

func (api *ObsidianAPI) makeTextRequest(method, endpoint string, content string) (*http.Response, error) {
	var bodyReader io.Reader

	if content != "" {
		bodyReader = strings.NewReader(content)
	}

	req, err := http.NewRequest(method, api.baseURL+endpoint, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "text/markdown")
	if api.token != "" {
		req.Header.Set("Authorization", "Bearer "+api.token)
	}

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}

	return resp, nil
}

// GetNote retrieves a note by its path
func (api *ObsidianAPI) GetNote(path string) (string, error) {
	// Normalize path to ensure .md extension
	path = normalizeNotePath(path)
	endpoint := fmt.Sprintf("/vault/%s", url.PathEscape(path))

	resp, err := api.makeRequest("GET", endpoint, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get note: %s", resp.Status)
	}

	// Read the plain text content
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	result := fmt.Sprintf("# Note: %s\n\n%s", path, string(content))

	return result, nil
}

// CreateNote creates a new note
func (api *ObsidianAPI) CreateNote(path, content string) (string, error) {
	endpoint := fmt.Sprintf("/vault/%s", url.PathEscape(path))

	resp, err := api.makeTextRequest("PUT", endpoint, content)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		return "", fmt.Errorf("failed to create note: %s", resp.Status)
	}

	return fmt.Sprintf("Successfully created note: %s", path), nil
}

// UpdateNote updates an existing note
func (api *ObsidianAPI) UpdateNote(path, content string) (string, error) {
	endpoint := fmt.Sprintf("/vault/%s", url.PathEscape(path))

	resp, err := api.makeTextRequest("PUT", endpoint, content)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return "", fmt.Errorf("failed to update note: %s", resp.Status)
	}

	return fmt.Sprintf("Successfully updated note: %s", path), nil
}

// DeleteNote deletes a note
func (api *ObsidianAPI) DeleteNote(path string) (string, error) {
	endpoint := fmt.Sprintf("/vault/%s", url.PathEscape(path))

	resp, err := api.makeRequest("DELETE", endpoint, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return "", fmt.Errorf("failed to delete note: %s", resp.Status)
	}

	return fmt.Sprintf("Successfully deleted note: %s", path), nil
}

// ListNotes lists all notes in the vault or a specific folder
func (api *ObsidianAPI) ListNotes(folder string) (string, error) {
	endpoint := "/vault/"
	if folder != "" {
		endpoint = fmt.Sprintf("/vault/%s/", url.PathEscape(folder))
	}

	resp, err := api.makeRequest("GET", endpoint, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Handle 404 for empty/non-existent folders - return empty list instead of error
	if resp.StatusCode == http.StatusNotFound {
		if folder != "" {
			return fmt.Sprintf("Folder '%s' is empty or does not exist yet. No notes found.", folder), nil
		}
		return "No notes found in vault.", nil
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to list notes: %s", resp.Status)
	}

	var files map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&files); err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	var notes []string
	if filesList, ok := files["files"].([]interface{}); ok {
		for _, file := range filesList {
			// Handle both string and object formats
			if fileStr, ok := file.(string); ok {
				if strings.HasSuffix(fileStr, ".md") {
					notes = append(notes, fileStr)
				}
			} else if fileInfo, ok := file.(map[string]interface{}); ok {
				if path, ok := fileInfo["path"].(string); ok {
					if strings.HasSuffix(path, ".md") {
						notes = append(notes, path)
					}
				}
			}
		}
	}

	if len(notes) == 0 {
		return "No notes found.", nil
	}

	result := fmt.Sprintf("Found %d notes:\n", len(notes))
	for _, note := range notes {
		result += fmt.Sprintf("- %s\n", note)
	}

	return result, nil
}

// SearchNotes searches for notes containing specific text
func (api *ObsidianAPI) SearchNotes(query string) (string, error) {
	// Use simple search endpoint
	endpoint := fmt.Sprintf("/search/simple/?query=%s&contextLength=100", url.QueryEscape(query))

	req, err := http.NewRequest("POST", api.baseURL+endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	if api.token != "" {
		req.Header.Set("Authorization", "Bearer "+api.token)
	}

	resp, err := api.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("search request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to search notes: %s - %s", resp.Status, string(bodyBytes))
	}

	var searchResults []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&searchResults); err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	if len(searchResults) == 0 {
		return fmt.Sprintf("No notes found matching \"%s\".", query), nil
	}

	result := fmt.Sprintf("Found %d notes matching \"%s\":\n\n", len(searchResults), query)
	for i, item := range searchResults {
		if filename, ok := item["filename"].(string); ok {
			result += fmt.Sprintf("%d. **%s**\n", i+1, filename)

			// Add context from matches
			if matches, ok := item["matches"].([]interface{}); ok && len(matches) > 0 {
				for _, match := range matches {
					if matchMap, ok := match.(map[string]interface{}); ok {
						if context, ok := matchMap["context"].(string); ok {
							result += fmt.Sprintf("   > %s\n", strings.TrimSpace(context))
						}
					}
				}
			}
			result += "\n"
		}
	}

	return result, nil
}

// GetVaultInfo gets information about the vault
func (api *ObsidianAPI) GetVaultInfo() (string, error) {
	resp, err := api.makeRequest("GET", "/", nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get vault info: %s", resp.Status)
	}

	var info map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	result := "# Vault Information\n\n"

	if authenticated, ok := info["authenticated"].(bool); ok {
		result += fmt.Sprintf("**Authenticated:** %v\n", authenticated)
	}

	if service, ok := info["service"].(string); ok {
		result += fmt.Sprintf("**Service:** %s\n", service)
	}

	if versions, ok := info["versions"].(map[string]interface{}); ok {
		result += "\n## Versions\n"
		for key, value := range versions {
			result += fmt.Sprintf("- **%s:** %v\n", key, value)
		}
	}

	// Get additional vault statistics
	notesResp, err := api.makeRequest("GET", "/vault/", nil)
	if err == nil && notesResp.StatusCode == http.StatusOK {
		var files map[string]interface{}
		if err := json.NewDecoder(notesResp.Body).Decode(&files); err == nil {
			if filesList, ok := files["files"].([]interface{}); ok {
				noteCount := 0
				folderCount := 0
				for _, file := range filesList {
					if fileInfo, ok := file.(map[string]interface{}); ok {
						if path, ok := fileInfo["path"].(string); ok {
							if strings.HasSuffix(path, ".md") {
								noteCount++
							}
						}
						if isFolder, ok := fileInfo["is_folder"].(bool); ok && isFolder {
							folderCount++
						}
					}
				}
				result += "\n## Statistics\n"
				result += fmt.Sprintf("- **Notes:** %d\n", noteCount)
				result += fmt.Sprintf("- **Folders:** %d\n", folderCount)
			}
		}
		notesResp.Body.Close()
	}

	return result, nil
}
