package gh

// Comment : The comment associated with a notification
type Comment struct {
	HTMLURL string `json:"html_url"`
	CreatedAt string `json:"created_at"`
	// Those two are don't appear in every Comment, only on PR's. -> Rely on GhNotification.Type
	State string `json:"state"`
	Merged bool `json:"merged"`
	Body string `json:"body"`
}
