package models

// Error ...
type Error struct {
	Error error `json:"error"`
	Code int `json:"code"`
	Description string `json:"description"`
}

// StandardErrorModel ...
type StandardErrorModel struct {
	Error Error `json:"error"`
}
