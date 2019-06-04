/*-----------
https://github.com/ymotongpoo/go-dialogflow-fulfillment
2019/05/29 FUKUSHIMA Kaito
--------*/

package main

import (
	"bytes"
	"encoding/json"

	"fmt"
<<<<<<< HEAD
	"io"
	"io/ioutil"
	"log"
=======

	//"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"bytes"
	"encoding/json"
>>>>>>> 87f400a80f2c964eefe5482e5c0e05ef35eca2e8
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
	http.ListenAndServe("0.0.0.0:8080", nil)
}

const (
	//WelcomeIntent is welcome intent message
	WelcomeIntent  = "input.welcome"
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
<<<<<<< HEAD
=======

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s -project-id <PROJECT ID> -session-id <SESSION ID> -language-code <LANGUAGE CODE> <OPERATION> <INPUTS>\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "<PROJECT ID> must be your Google Cloud Platform project id\n")
		fmt.Fprintf(os.Stderr, "<SESSION ID> must be a Dialogflow session ID\n")
		fmt.Fprintf(os.Stderr, "<LANGUAGE CODE> must be a language code from https://dialogflow.com/docs/reference/language; defaults to en\n")
		fmt.Fprintf(os.Stderr, "<OPERATION> must be one of text, audio, stream\n")
		fmt.Fprintf(os.Stderr, "<INPUTS> can be a series of text inputs if <OPERATION> is text, or a path to an audio file if <OPERATION> is audio or stream\n")
	}
	var projectID, sessionID, languageCode string
	flag.StringVar(&projectID, "project-id", "", "Google Cloud Platform project ID")
	flag.StringVar(&sessionID, "session-id", "", "Dialogflow session ID")
	flag.StringVar(&languageCode, "language-code", "jp", "Dialogflow language code from https://dialogflow.com/docs/reference/language; defaults to en")

	flag.Parse()

	args := flag.Args()

	if len(args) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	operation := args[0]
	inputs := args[1:]

	switch operation {
	case "text":
		fmt.Printf("Responses:\n")
		for _, query := range inputs {
			fmt.Printf("\nInput: %s\n", query)
			response, err := DetectIntentText(projectID, sessionID, query, languageCode)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Output: %s\n", response)
		}
	default:
		flag.Usage()
		os.Exit(1)
	}
	html := httpResponse()
	decodeJSON(html)
}

// DetectIntentText : dialogflow_detect_intent_text]
func DetectIntentText(projectID, sessionID, text, languageCode string) (string, error) {
	ctx := context.Background()

	sessionClient, err := dialogflow.NewSessionsClient(ctx)
	if err != nil {
		return "", err
	}
	defer sessionClient.Close()

	if projectID == "" || sessionID == "" {
		return "", fmt.Errorf(fmt.Sprintf("Received empty project (%s) or session (%s)", projectID, sessionID))
	}

	sessionPath := fmt.Sprintf("projects/%s/agent/sessions/%s", projectID, sessionID)
	textInput := dialogflowpb.TextInput{Text: text, LanguageCode: languageCode}
	queryTextInput := dialogflowpb.QueryInput_Text{Text: &textInput}
	queryInput := dialogflowpb.QueryInput{Input: &queryTextInput}
	request := dialogflowpb.DetectIntentRequest{Session: sessionPath, QueryInput: &queryInput}

	response, err := sessionClient.DetectIntent(ctx, &request)
	if err != nil {
		return "", err
	}

	queryResult := response.GetQueryResult()
	fulfillmentText := queryResult.GetFulfillmentText()
	return fulfillmentText, nil
}
>>>>>>> 87f400a80f2c964eefe5482e5c0e05ef35eca2e8
