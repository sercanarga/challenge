package models

type Event struct {
	App        string          `json:"app"`
	Type       string          `json:"type"`
	Time       string          `json:"time"`
	Meta       EventMeta       `json:"meta"`
	Wallet     string          `json:"wallet"`
	Attributes EventAttributes `json:"attributes"`
}

type EventMeta struct {
	User string `json:"user"`
}

type EventAttributes struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

type EventList struct {
	Events []Event `json:"events"`
}
