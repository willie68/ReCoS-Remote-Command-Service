package models

// Message is a message from the server to the client to update an action
type Message struct {
	Profile  string `json:"profile"`
	Action   string `json:"action"`
	Page     string `json:"page"`
	ImageURL string `json:"imageurl"`
	Title    string `json:"title"`
	Text     string `json:"text"`
	State    int    `json:"state"`
	Command  string `json:"command"`
}
