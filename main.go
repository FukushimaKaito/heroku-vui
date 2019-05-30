/*-----------
[ref] json  https://qiita.com/nayuneko/items/2ec20ba69804e8bf7ca3
[ref] alexa https://github.com/yamaryu0508/alexa-skills-kit-color-expert-go
2019/05/29 FUKUSHIMA Kaito
--------*/

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"./alexa"
	"bytes"
	"errors"
	"net/http"

	"log"
)

// Value is JSONデコード用の構造体
type Value struct {
	Light float32 `json:"d1"`
	Vib   float32 `json:"d2"`
	DT    string  `json:"created"`
}

var (
	// ErrInvalidIntent is error-object
	ErrInvalidIntent = errors.New("Invalid intent")
)

/*
 * Functions that control the skill's behavior
 */

// GetWelcomeResponse is function-type
func GetWelcomeResponse() alexa.Response {
	sessionAttributes := make(map[string]interface{})
	cardTitle := "Welcome"
	speechOutput := "Welcome to the Alexa Skills Kit sample. Please tell me your favorite color by saying, my favorite color is red"
	repromptText := "Please tell me your favorite color by saying, my favorite color is red."
	shouldEndSession := false
	return alexa.BuildResponse(sessionAttributes, alexa.BuildSpeechletResponse(cardTitle, speechOutput, repromptText, shouldEndSession))
}

// HandleSessionEndRequest is function-type
func HandleSessionEndRequest() alexa.Response {
	sessionAttributes := make(map[string]interface{})
	cardTitle := "Session Ended"
	speechOutput := "Thank you for trying the Alexa Skills Kit sample. Have a nice day! "
	repromptText := ""
	shouldEndSession := true
	return alexa.BuildResponse(sessionAttributes, alexa.BuildSpeechletResponse(cardTitle, speechOutput, repromptText, shouldEndSession))
}

// CreateFavoriteColorAttributes is function-type
func CreateFavoriteColorAttributes(favoriteColor string) alexa.SessionAttributes {
	sessionAttributes := make(map[string]interface{})
	sessionAttributes["favoriteColor"] = favoriteColor
	return sessionAttributes
}

// SetColorInSession is function-type
func SetColorInSession(intent alexa.RequestIntent, session alexa.Session) alexa.Response {
	cardTitle := intent.Name
	sessionAttributes := make(map[string]interface{})
	shouldEndSession := false
	speechOutput := ""
	repromptText := ""

	if color, ok := intent.Slots["Color"]; ok {
		favoriteColor := color.Value
		sessionAttributes = CreateFavoriteColorAttributes(favoriteColor)
		speechOutput = "I now know your favorite color is " + favoriteColor +
			". You can ask me your favorite color by saying, " +
			"what's my favorite color?"
		repromptText = "You can ask me your favorite color by saying, " +
			"what's my favorite color?"
	} else {
		speechOutput = "I'm not sure what your favorite color is. " +
			"Please try again."
		repromptText = "I'm not sure what your favorite color is. " +
			"You can tell me your favorite color by saying, " +
			"my favorite color is red."
	}
	return alexa.BuildResponse(sessionAttributes, alexa.BuildSpeechletResponse(cardTitle, speechOutput, repromptText, shouldEndSession))
}

// GetColorFromSession is function-type
func GetColorFromSession(intent alexa.RequestIntent, session alexa.Session) alexa.Response {
	cardTitle := intent.Name
	sessionAttributes := make(map[string]interface{})
	shouldEndSession := false
	speechOutput := ""
	repromptText := ""

	if favoriteColor, ok := session.Attributes["favoriteColor"].(string); ok {
		speechOutput = "Your favorite color is " + favoriteColor + ". Goodbye."
		shouldEndSession = true
	} else {
		speechOutput = "I'm not sure what your favorite color is, you can say, my favorite color is red"
	}
	return alexa.BuildResponse(sessionAttributes, alexa.BuildSpeechletResponse(cardTitle, speechOutput, repromptText, shouldEndSession))
}

// GetNoEntityResponse is function-type
func GetNoEntityResponse() alexa.Response {
	cardTitle := ""
	sessionAttributes := make(map[string]interface{})
	shouldEndSession := false
	speechOutput := ""
	repromptText := ""
	return alexa.BuildResponse(sessionAttributes, alexa.BuildSpeechletResponse(cardTitle, speechOutput, repromptText, shouldEndSession))
}

/*
 * Events
 */

// OnSessionStarted is function-type
func OnSessionStarted(sessionStartedRequest map[string]string, session alexa.Session) (alexa.Response, error) {
	fmt.Println("OnSessionStarted requestId=" + sessionStartedRequest["requestId"] + ", sessionId=" + session.SessionID)
	return GetNoEntityResponse(), nil
}

// OnLaunch is function-type
func OnLaunch(launchRequest alexa.RequestDetail, session alexa.Session) (alexa.Response, error) {
	fmt.Println("OnLaunch requestId=" + launchRequest.RequestID + ", sessionId=" + session.SessionID)
	return GetWelcomeResponse(), nil
}

// OnIntent is function-type
func OnIntent(intentRequest alexa.RequestDetail, session alexa.Session) (alexa.Response, error) {
	fmt.Println("OnLaunch requestId=" + intentRequest.RequestID + ", sessionId=" + session.SessionID)
	intent := intentRequest.Intent
	intentName := intentRequest.Intent.Name

	if intentName == "MyColorIsIntent" {
		return SetColorInSession(intent, session), nil
	} else if intentName == "WhatsMyColorIntent" {
		return GetColorFromSession(intent, session), nil
	} else if intentName == "AMAZON.HelpIntent" {
		return GetWelcomeResponse(), nil
	} else if intentName == "AMAZON.StopIntent" || intentName == "AMAZON.CancelIntent" {
		return HandleSessionEndRequest(), nil
	}
	return alexa.Response{}, ErrInvalidIntent
}

// OnSessionEnded is function-type
func OnSessionEnded(sessionEndedRequest alexa.RequestDetail, session alexa.Session) (alexa.Response, error) {
	fmt.Println("OnSessionEnded requestId=" + sessionEndedRequest.RequestID + ", sessionId=" + session.SessionID)
	return GetNoEntityResponse(), nil
}

// Handler is main
func Handler(event alexa.Request) (alexa.Response, error) {
	fmt.Println("event.session.application.applicationId=" + event.Session.Application.ApplicationID)

	eventRequestType := event.Request.Type
	if event.Session.New {
		return OnSessionStarted(map[string]string{"requestId": event.Request.RequestID}, event.Session)
	} else if eventRequestType == "LaunchRequest" {
		return OnLaunch(event.Request, event.Session)
	} else if eventRequestType == "IntentRequest" {
		return OnIntent(event.Request, event.Session)
	} else if eventRequestType == "SessionEndedRequest" {
		return OnSessionEnded(event.Request, event.Session)
	}
	return alexa.Response{}, ErrInvalidIntent
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

func main() {
	html := httpResponse()
	decodeJSON(html)
	lambda.Start(Handler)
}
