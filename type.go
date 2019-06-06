package main

const (
	//return message
	lightHighest = "光が強すぎます．"
	lightJust    = "光合成にちょうどよい照度継続時間です．"
	lightHigher  = "もう少し光を強くてもいいかもしれません．"
	lightLack    = "光が足りていません．"
	lightMissing = "は登録されていません．🙇"

	nowresponse   = "%s現在の振動値は%.2fGal 、明るさは%.2flxです．"
	countresponse = "強い光が%d分，明るい光が%d分，暗い状態が%d分です．"
)

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