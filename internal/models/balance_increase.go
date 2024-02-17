package models

type Event struct {
	App        string     `json:"app"`
	Type       string     `json:"type"`
	Time       string     `json:"time"`
	Meta       Meta       `json:"meta"`
	Wallet     string     `json:"wallet"`
	Attributes Attributes `json:"attributes"`
}

type Meta struct {
	User string `json:"user"`
}

type Attributes struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

type EventList struct {
	Events []Event `json:"events"`
}
