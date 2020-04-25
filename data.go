package main

type sendResponse struct {
	ID         uint
	Message    string
	StatusCode int
}

type statusResponse struct {
	Message    string
	StatusCode int
}
