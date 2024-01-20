package main

type AnalyzeRequest struct {
	Url string `json:"url"`
}

type ErrorResponse struct {
	Message string `json:"error"`
}
