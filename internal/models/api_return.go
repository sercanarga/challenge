package models

type APIReturn struct {
	StatusCode   int
	Response     string
	ResponseTime int64
}

type APIEventReturn struct {
	StatusCode   int
	Success      []Event
	Unsuccess    []Event
	ResponseTime int64
}
