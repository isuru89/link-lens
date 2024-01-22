package server

type AnalyzeRequest struct {
	Url string `json:"url"`
}

type ErrorResponse struct {
	Message string `json:"error"`
}

type HealthResponse struct {
	Alive bool `json:"alive"`
}
