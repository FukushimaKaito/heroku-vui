package main

import (
	"net/http"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	//"io/ioutil"
	"log"
)

//Request is request from filfullment
type Request struct {
	ID  			    string `json:"responseId"`
	Result struct{
		Parameters struct{
			Vegelight   string `json:"Vegelight"`
		} 					   `json:"parameters"`
		Intent struct{
			DisplayName string `json:"displayName"`
		} 					   `json:"intent"`
	}						   `json:"queryResult"`
	SessionID			string `json:"session"`
}

// Response is response for filfullment
type Response struct{
	FulfillmentText string `json:"fulfillmentText"`
}

//NewResponse is new response
func NewResponse(speech string) *Response {
	return &Response{
	//	Speech: speech,
	}
}

//SetDisplayText is sendmsg
func (res *Response) SetDisplayText(text string) *Response {
	res.FulfillmentText = text
	return res
}


func main(){
	http.HandleFunc("/",handler)
	port := os.Getenv("PORT")
	if port == ""{
		port = "8080"
	}
	http.ListenAndServe(":"+port,nil)
}

func handler(w http.ResponseWriter,r *http.Request){
	//Validate request
	// if r.Method != "POST" {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	// if r.Header.Get("Content-Type") != "application/json" {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	// //To allocate slice for request body
	// length, err := strconv.Atoi(r.Header.Get("Content-Length"))
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	// //Read body data to parse json
	// body := make([]byte, length)
	// length, err = r.Body.Read(body)
	// if err != nil && err != io.EOF {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	// //parse json
	// var jsonBody map[string]interface{}
	// err = json.Unmarshal(body[:length], &jsonBody)
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }
	// fmt.Printf("%v\n", jsonBody)

	// w.WriteHeader(http.StatusOK)  


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

// DecodeInput is decode recvmsg
func DecodeInput(r *http.Request) (*Request, error) {
	var req Request
	var buf bytes.Buffer
	tee := io.TeeReader(r.Body, &buf)
	defer r.Body.Close()
	err := json.NewDecoder(tee).Decode(&req)
	if err != nil {
		return nil, fmt.Errorf("decode error: %v", err)
		// b, err := ioutil.ReadAll(&buf)
		// if err != nil {
		// 	return nil, fmt.Errorf("ioutil error: %v", err)
		// }
		// log.Printf("%s\n", b)
	}
	return &req, nil
}

//EncodeOutput is mk sendmsg
func EncodeOutput(w http.ResponseWriter, res *Response) error {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Printf("encode error: %v\n", err)
	}
	return nil
}

//welcomeIntent is mk welcome msg
func welcomeIntent(r *Request) (*Response, error) {
	template := `こんにちは。`
	msg := fmt.Sprintf(template)
	return NewResponse(msg).SetDisplayText(msg), nil
}

//asklightIntent is mk asklight msg
func asklightIntent(r *Request) (*Response, error) {
	msg :=fmt.Sprintf("msg")
	return NewResponse(msg).SetDisplayText(msg), nil
}

// asknowIntent is mk asknowIntent msg 
func asknowIntent(r *Request) (*Response, error) {
	msg :=fmt.Sprintf("msg")
	return NewResponse(msg).SetDisplayText(msg), nil
}