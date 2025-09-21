package webhook

// GitHubEvent represents a generic GitHub webhook event
type GitHubEvent struct {
	Action     string                 `json:"action,omitempty"`
	Repository map[string]interface{} `json:"repository,omitempty"`
	Sender     map[string]interface{} `json:"sender,omitempty"`
}

// SupportedEventTypes contains the event types we want to store in the database
var SupportedEventTypes = map[string]bool{
	"push":          true,
	"issue_comment": true,
	"pull_request":  true,
}

// IsSupportedEvent checks if an event type should be stored in the database
func IsSupportedEvent(eventType string) bool {
	return SupportedEventTypes[eventType]
}