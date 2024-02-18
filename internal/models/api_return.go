package models

type APIReturn struct {
	StatusCode   int
	Response     string
	ResponseTime int64
}

type EventResult struct {
	StatusCode   int    `json:"statusCode,omitempty"`
	ErrorDetails string `json:"errorDetails,omitempty"`
}

type APIEventReturn struct {
	Data   Event       `json:"data"`
	Result EventResult `json:"result"`
}
