package main

// Value for ambient JSON decode
type Value struct {
	Light float32 `json:"d1"`
	Vib   float32 `json:"d2"`
	DT    string  `json:"created"`
}

//Request is request from filfullment
type Request struct {
	ID     string `json:"responseId"`
	Result struct {
		Parameters struct {
			Vegelight string `json:"Vegelight"`
		} `json:"parameters"`
		Intent struct {
			DisplayName string `json:"displayName"`
		} `json:"intent"`
	} `json:"queryResult"`
	SessionID string `json:"session"`
}

// Response is response for filfullment
type Response struct {
	FulfillmentText string `json:"fulfillmentText"`
}