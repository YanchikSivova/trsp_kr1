package models

var ForbiddenWords = []string{"кринж", "рофл", "вайб"}

type ValidationErrorDetail struct {
	Type     string      `json:"type"`
	Location []string    `json:"loc"`
	Msg      string      `json:"msg"`
	Input    interface{} `json:"input"`
	Ctx      interface{} `json:"ctx,omitempty"`
}

type ValidationError struct {
	Detail []ValidationErrorDetail `json:"detail"`
}
