package main

import (
	"net/http"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"io/ioutil"
	"log"
)

// Value for ambient JSON decode
type Value struct {
	Light float32 `json:"d1"`
	Vib   float32 `json:"d2"`
	DT    string  `json:"created"`
}

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
	req, err := DecodeInput(r)
	if err != nil {
		log.Println(err)
		return
	}

	var res *Response
	intent := req.Result.Intent.DisplayName

	switch intent {
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

//asklightIntent is mk asklight msg
func asklightIntent(r *Request) (*Response, error) {
	//ambient
	url := "http://ambidata.io/api/v2/channels/10905/data?readKey=7e7df40858ef249c&n=1440"
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	buf := bytes.NewBuffer(body)
	html := buf.String()
	fmt.Println(html)
	bytes := []byte(html)

	// //JSONデコード
	var values []Value
	if err := json.Unmarshal(bytes, &values); err != nil {
		log.Fatal(err)
	}

	//json count
	high:=0
	mid:=0
	low:=0
    for i:=0;i< 1440;i++{
        if values[i].Light > 1000{
			high ++
		} else if values[i].Light < 300{
			low ++
		} else{
			mid ++
		}
	}
	//mkmsg
	lightHighest := "光が強すぎます．"
	lightJust := "光合成にちょうどよい照度継続時間です．"
	lightHigher := "もう少し光を強くてもいいかもしれません．"
	lightLack := "光が足りていません．"
	lightMissing := "は登録されていません．🙇"

	msg:=""
    if r.Result.Parameters.Vegelight == "ミニトマト"{// positive class
        if high > 360{
            msg = lightJust
		} else if high + mid > 360{
            msg = lightHigher
		} else{
            msg = lightLack
		}
	} else if r.Result.Parameters.Vegelight == "ジャガイモ"{// negative class
        if high > 30 || mid > 180{
            msg = lightHighest
		}else if high + mid > 60{
            msg = lightJust
		}else{
            msg = lightLack
		}
	}else if r.Result.Parameters.Vegelight =="大葉"{// half class
        if high > 120 || mid > 180{
			msg = lightHighest
		} else if high + mid > 300{
            msg = lightJust
		} else{
            msg = lightLack
		}
	}else{
		msg = r.Result.Parameters.Vegelight + lightMissing
	}
	return NewResponse(msg).SetDisplayText(msg), nil
}

// asknowIntent is mk asknowIntent msg 
func asknowIntent(r *Request) (*Response, error) {
	//ambient
	url := "http://ambidata.io/api/v2/channels/10905/data?readKey=7e7df40858ef249c&n=1"
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	buf := bytes.NewBuffer(body)
	html := buf.String()
	fmt.Println(html)

	//JSON
	bytes := []byte(html)
	// //JSONデコード
	var values []Value
	if err := json.Unmarshal(bytes, &values); err != nil {
		log.Fatal(err)
	}
	//デコードデータの表示
	fmt.Printf("%f : %f\n", values[0].Light, values[0].Vib)
	template:="%s現在の振動値は%.2fGal 、明るさは%.2flxです．"
	msg :=fmt.Sprintf(template,values[0].DT,values[0].Vib,values[0].Light)
	return NewResponse(msg).SetDisplayText(msg), nil
}

//httpResponse
func httpResponse() string {
	//http----------------
	url := "http://ambidata.io/api/v2/channels/10905/data?readKey=7e7df40858ef249c&n=1"
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	buf := bytes.NewBuffer(body)
	html := buf.String()
	fmt.Println(html)
	return html
}

//decodeJSON
func decodeJSON(html string) {
	//JSON---------------
	bytes := []byte(html)
	// //JSONデコード
	var values []Value
	if err := json.Unmarshal(bytes, &values); err != nil {
		log.Fatal(err)
	}
	//デコードデータの表示
	for _, d := range values {
		fmt.Printf("%f : %f\n", d.Light, d.Vib)
	}
}