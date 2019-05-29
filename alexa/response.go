package alexa

/*
 * define structs for Alexa Response
 */

// OutputSpeech is object-type
type OutputSpeech struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// Card is object-type
type Card struct {
	Type    string `json:"type"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// Reprompt is object-type
type Reprompt struct {
	OutputSpeech OutputSpeech `json:"outputSpeech"`
}

// SpeechletResponse is object-type
type SpeechletResponse struct {
	ShouldEndSession bool         `json:"shouldEndSession"`
	OutputSpeech     OutputSpeech `json:"outputSpeech"`
	Card             Card         `json:"card"`
	Reprompt         Reprompt     `json:"reprompt"`
}

// SessionAttributes is map-type
type SessionAttributes map[string]interface{}

// Response is object-type
type Response struct {
	Version           string            `json:"version"`
	SessionAttributes SessionAttributes `json:"sessionAttributes"`
	Response          SpeechletResponse `json:"response"`
}

// BuildSpeechletResponse is function-type
func BuildSpeechletResponse(title string, output string, repromptText string, shouldEndSession bool) SpeechletResponse {
	return SpeechletResponse{
		ShouldEndSession: shouldEndSession,
		OutputSpeech: OutputSpeech{
			Type: "PlainText",
			Text: output,
		},
		Card: Card{
			Type:    "Simple",
			Title:   "SessionSpeechlet - " + title,
			Content: "SessionSpeechlet - " + output,
		},
		Reprompt: Reprompt{
			OutputSpeech: OutputSpeech{
				Type: "PlainText",
				Text: repromptText,
			},
		},
	}
}

// BuildResponse is function-type
func BuildResponse(sessionAttributes SessionAttributes, speechletResponse SpeechletResponse) Response {
	return Response{
		Version:           "1.0",
		Response:          speechletResponse,
		SessionAttributes: sessionAttributes,
	}
}