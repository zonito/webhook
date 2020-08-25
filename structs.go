package main

// Response
type Response struct {
	Success bool   `json:"success"`
	Reason  string `json:"reason"`
	Handler string `json:"handler"`
}
