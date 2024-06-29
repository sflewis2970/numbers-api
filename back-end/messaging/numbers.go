package messaging

type NumbersAPIData struct {
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
	Message   string `json:"message,omitempty"`
}
