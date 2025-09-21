package webhook

// GitHubEvent represents a generic GitHub webhook event
type GitHubEvent struct {
	Action     string                 `json:"action,omitempty"`
	Repository map[string]interface{} `json:"repository,omitempty"`
	Sender     map[string]interface{} `json:"sender,omitempty"`
}