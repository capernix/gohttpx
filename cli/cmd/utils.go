package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Note struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type APIResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

func HTTPClient() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
	}
}

func GetBaseURL() string {
	config := LoadConfig()
	if config.BaseURL == "" {
		return "http://localhost"
	}
	return config.BaseURL
}

func MakeRequest(method, endpoint string, body interface{}) (*http.Response, error) {
	client := HTTPClient()
	baseURL := GetBaseURL()
	url := baseURL + endpoint

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %v", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return client.Do(req)
}

func ParseAPIResponse(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	bodyStr := strings.TrimSpace(string(body))
	fmt.Println("DEBUG: Response body:", bodyStr)

	if strings.HasPrefix(bodyStr, "[") {
		fmt.Println("DEBUG: Detected array response")
		if err := json.Unmarshal(body, target); err != nil {
			return fmt.Errorf("failed to parse array response: %v", err)
		}
		return nil
	}

	if strings.HasPrefix(bodyStr, "{") {
		fmt.Println("DEBUG: Detected object response, trying direct unmarshal")
		if err := json.Unmarshal(body, target); err == nil {
			fmt.Println("DEBUG: Direct unmarshal succeeded")
			return nil
		}
		fmt.Println("DEBUG: Direct unmarshal failed, trying wrapped")
		var apiResp APIResponse
		if err := json.Unmarshal(body, &apiResp); err != nil {
			return fmt.Errorf("failed to parse response: %v", err)
		}
		if apiResp.Status != "" && apiResp.Status != "success" {
			return fmt.Errorf("API returned status: %s", apiResp.Status)
		}
		if apiResp.Data != nil {
			dataBytes, err := json.Marshal(apiResp.Data)
			if err != nil {
				return fmt.Errorf("failed to marshal data: %v", err)
			}
			if err := json.Unmarshal(dataBytes, target); err != nil {
				return fmt.Errorf("failed to unmarshal data: %v", err)
			}
		}
		return nil
	}

	return fmt.Errorf("unexpected response format")
}

func PrintUsers(users []User, jsonOutput bool) {
	if jsonOutput {
		output, _ := json.MarshalIndent(users, "", "  ")
		fmt.Println(string(output))
		return
	}

	if len(users) == 0 {
		fmt.Println("No users found")
		return
	}

	fmt.Println("┌────┬──────────────────────────────────┐")
	fmt.Println("│ ID │ Name                             │")
	fmt.Println("├────┼──────────────────────────────────┤")

	for _, user := range users {
		name := user.Name
		if len(name) > 32 {
			name = name[:29] + "..."
		}
		fmt.Printf("│ %-2d │ %-32s │\n", user.ID, name)
	}

	fmt.Println("└────┴──────────────────────────────────┘")
}

func PrintUser(user User, jsonOutput bool) {
	if jsonOutput {
		output, _ := json.MarshalIndent(user, "", "  ")
		fmt.Println(string(output))
		return
	}

	fmt.Printf("ID:   %d\n", user.ID)
	fmt.Printf("Name: %s\n", user.Name)
}

func PrintNotes(notes []Note, jsonOutput bool) {
	if jsonOutput {
		output, _ := json.MarshalIndent(notes, "", "  ")
		fmt.Println(string(output))
		return
	}

	if len(notes) == 0 {
		fmt.Println("No notes found")
		return
	}

	// Print table header
	fmt.Println("┌────┬────────────────────┬──────────────────────────────────┐")
	fmt.Println("│ ID │ Title              │ Content                          │")
	fmt.Println("├────┼────────────────────┼──────────────────────────────────┤")

	// Print notes
	for _, note := range notes {
		title := note.Title
		if len(title) > 18 {
			title = title[:15] + "..."
		}
		content := note.Content
		if len(content) > 32 {
			content = content[:29] + "..."
		}
		fmt.Printf("│ %-2d │ %-18s │ %-32s │\n", note.ID, title, content)
	}

	fmt.Println("└────┴────────────────────┴──────────────────────────────────┘")
}

func PrintNote(note Note, jsonOutput bool) {
	if jsonOutput {
		output, _ := json.MarshalIndent(note, "", "  ")
		fmt.Println(string(output))
		return
	}

	fmt.Printf("ID:      %d\n", note.ID)
	fmt.Printf("Title:   %s\n", note.Title)
	fmt.Printf("Content: %s\n", note.Content)
}

func ValidateID(idStr string) (int, error) {
	if idStr == "" {
		return 0, fmt.Errorf("ID is required")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("ID must be a valid number")
	}

	if id <= 0 {
		return 0, fmt.Errorf("ID must be a positive number")
	}

	return id, nil
}

func PrintSuccess(message string) {
	fmt.Printf("Success: %s\n", message)
}

func PrintError(message string) {
	fmt.Printf("Error: %s\n", message)
}

func ConfirmAction(prompt string) bool {
	fmt.Printf("%s [y/N]: ", prompt)
	var response string
	fmt.Scanln(&response)
	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}
