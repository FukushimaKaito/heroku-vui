package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

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

//asklightIntent is mk asklight msg
func asklightIntent(r *Request) (*Response, error) {
	url := "http://ambidata.io/api/v2/channels/10905/data?readKey=7e7df40858ef249c&n=1440"
	var values []Value
	values = decodeAmbient(url)

	//json count
	high := 0
	mid := 0
	low := 0
	for i := 0; i < 1440; i++ {
		if values[i].Light > 1000 {
			high++
		} else if values[i].Light < 300 {
			low++
		} else {
			mid++
		}
	}

	msg := ""
	//if detectWord(r.Result.Parameters.Vegelight, 0)
	if r.Result.Parameters.Vegelight == "トマト" { // positive class
		if high > 360 {
			msg = fmt.Sprintf(lightJust)
		} else if high+mid > 360 {
			msg = fmt.Sprintf(lightHigher)
		} else {
			msg = fmt.Sprintf(lightLack)
		}
		//}else if detectWord(r.Result.Parameters.Vegelight, 2){
	} else if r.Result.Parameters.Vegelight == "シソ" { // negative class
		if high > 30 || mid > 180 {
			msg = fmt.Sprintf(lightHighest)
		} else if high+mid > 60 {
			msg = fmt.Sprintf(lightJust)
		} else {
			msg = fmt.Sprintf(lightLack)
		}
		//	}else if detectWord(r.Result.Parameters.Vegelight, 1){
	} else if r.Result.Parameters.Vegelight == "ジャガイモ" { // half class
		if high > 120 || mid > 180 {
			msg = fmt.Sprintf(lightHighest)
		} else if high+mid > 300 {
			msg = fmt.Sprintf(lightJust)
		} else {
			msg = fmt.Sprintf(lightLack)
		}
	} else {
		msg = fmt.Sprintf(r.Result.Parameters.Vegelight + lightMissing)
	}
	return NewResponse(msg).SetDisplayText(msg), nil
}

// asknowIntent is mk asknowIntent msg
func asknowIntent(r *Request) (*Response, error) {
	url := "http://ambidata.io/api/v2/channels/10905/data?readKey=7e7df40858ef249c&n=1"
	var values []Value
	values = decodeAmbient(url)

	//デコードデータの表示
	fmt.Printf("%f : %f\n", values[0].Light, values[0].Vib)
	msg := fmt.Sprintf(nowresponse, values[0].DT, values[0].Vib, values[0].Light)
	return NewResponse(msg).SetDisplayText(msg), nil
}

func countCheckIntent(r *Request) (*Response, error) {
	url := "http://ambidata.io/api/v2/channels/10905/data?readKey=7e7df40858ef249c&n=1440"
	var values []Value
	values = decodeAmbient(url)

	//json count
	high := 0
	mid := 0
	low := 0
	for i := 0; i < 1440; i++ {
		if values[i].Light > 1000 {
			high++
		} else if values[i].Light < 300 {
			low++
		} else {
			mid++
		}
	}
	msg := fmt.Sprintf(countresponse, high, mid, low)
	return NewResponse(msg).SetDisplayText(msg), nil
}

//detectWord is vegiClass
func detectWord(word string, class int) bool {
	//read csv
	file, err := os.Open("./lightVegiClass.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var line []string

	for {
		line, err = reader.Read()
		if err != nil {
			break
		}
	}
	if word == line[class] {
		return true
	}
	return false
}

//decodeAmbient
func decodeAmbient(url string) []Value {
	//url
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

	//JSON---------------
	bytes := []byte(html)
	var values []Value
	if err := json.Unmarshal(bytes, &values); err != nil {
		log.Fatal(err)
	}
	//print
	// for _, d := range values {
	// 	fmt.Printf("%f : %f\n", d.Light, d.Vib)
	// }
	return values
}
