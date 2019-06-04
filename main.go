package main

import (
	"net/http"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
)

type Request struct {
	ID        string                 `json:"responseId"`
	Result    struct{
		Parameters struct{
			Vegelight string `json:"Vegelight"`
		} `json:"parameters"`
		Intent struct{
			DisplayName string `json:"displayName"`
		} `json:"intent"`
	}`json:"queryResult"`
	SessionID string                 `json:"session"`
}

type Response struct{
	fulfillmentText string `json:fulfillmentText`
}

func NewResponse(speech string) *Response {
	return &Response{
	//	Speech: speech,
	}
}

func (res *Response) SetDisplayText(text string) *Response {
	res.fulfillmentText = text
	return res
}


func main(){
	http.HandleFunc("/",handler)
	http.ListenAndServe(":8080",nil)
}

func handler(w http.ResponseWriter,r *http.Request){
	req, err := DecodeInput(r)
	if err != nil {
		log.Println(err)
		return
	}

	var res *Response
	intent := req.Result.Intent.DisplayName

	switch intent {
	case "Default Welcome Intent":
		res, err = welcomeIntent(req)
	case "AskLightIntent":
		res,err = asklightIntent(req)
	case "AskNowdata":
		res,err=asknowIntent(req)
	}
	if err != nil {
		log.Println(err)
	}

	if err = EncodeOutput(w, res); err != nil {
		log.Println(err)
	}}

// DecodeInput
func DecodeInput(r *http.Request) (*Request, error) {
	var req Request
	var buf bytes.Buffer
	tee := io.TeeReader(r.Body, &buf)
	defer r.Body.Close()
	err := json.NewDecoder(tee).Decode(&req)
	if err != nil {
		return nil, fmt.Errorf("decode error: %v\n", err)
		b, err := ioutil.ReadAll(&buf)
		if err != nil {
			return nil, fmt.Errorf("ioutil error: %v\n", err)
		}
		log.Printf("%s\n", b)
	}
	return &req, nil
}

func EncodeOutput(w http.ResponseWriter, res *Response) error {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Printf("encode error: %v\n", err)
	}
	return nil
}

func welcomeIntent(r *Request) (*Response, error) {
	template := `こんにちは。`
	msg := fmt.Sprintf(template)
	return NewResponse(msg).SetDisplayText(msg), nil
}

func asklightIntent(r *Request) (*Response, error) {
	msg :=fmt.Sprintf("msg")
	return NewResponse(msg).SetDisplayText(msg), nil
}

func asknowIntent(r *Request) (*Response, error) {
	msg :=fmt.Sprintf("msg")
	return NewResponse(msg).SetDisplayText(msg), nil
}