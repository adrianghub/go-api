package resources

type Resource struct {
	ID             int    `json:"id"`
	Title          string `json:"title"`
	Category       string `json:"category"`
	Description    string `json:"description"`
	URL            string `json:"url"`
	DateAdded      string `json:"date_added"`
	ResourceType   string `json:"resource_type"`
	CompletionTime string `json:"completion_time"`
}
