/*-----------
https://github.com/ymotongpoo/go-dialogflow-fulfillment
2019/05/29 FUKUSHIMA Kaito
--------*/

package main

import (
	"bytes"
	"encoding/json"

	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

// Value is JSONデコード用の構造体
type Value struct {
	Light float32 `json:"d1"`
	Vib   float32 `json:"d2"`
	DT    string  `json:"created"`
}

func main() {
	http.HandleFunc("/", mainHandler)
	http.ListenAndServe("35.192.87.192", nil)
}

const (
	//WelcomeIntent is welcome intent message
	WelcomeIntent = "input.welcome"
	//AskLightIntent is AskLightIntent
	AskLightIntent = "input.asklight"
)

func mainHandler(w http.ResponseWriter, r *http.Request) {
	req, err := DecodeInput(r)
	if err != nil {
		log.Println(err)
		return
	}

	var res *Response
	intent := req.GetIntent()

	switch intent {
	case WelcomeIntent:
		res, err = welcomeIntent(req)
	case AskLightIntent:
		res, err = askLightIntent(req)
	}

	if err != nil {
		log.Println(err)
	}

	if err = EncodeOutput(w, res); err != nil {
		log.Println(err)
	}
}

// DecodeInput is Decode Input.
func DecodeInput(r *http.Request) (*Request, error) {
	var req Request
	var buf bytes.Buffer
	tee := io.TeeReader(r.Body, &buf)
	defer r.Body.Close()
	err := json.NewDecoder(tee).Decode(&req)
	if err != nil {
		// return nil, fmt.Errorf("decode error: %v", err)
		b, err := ioutil.ReadAll(&buf)
		if err != nil {
			return nil, fmt.Errorf("ioutil error: %v", err)
		}
		log.Printf("%s\n", b)
	}
	return &req, nil
}

// EncodeOutput is Encode Output
func EncodeOutput(w http.ResponseWriter, res *Response) error {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Printf("encode error: %v\n", err)
	}
	return nil
}

func welcomeIntent(r *Request) (*Response, error) {
	template := `家庭菜園支援VUI APPです．`
	voice := fmt.Sprintf(template)
	return NewResponse(voice).SetDisplayText(voice), nil
}

func askLightIntent(r *Request) (*Response, error) {
	template := `askLightIntentがcallされました．`
	voice := fmt.Sprintf(template)
	return NewResponse(voice).SetDisplayText(voice), nil
}

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
