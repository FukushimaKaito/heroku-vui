package main

// Request V2 API
// https://dialogflow.com/docs/fulfillment/how-it-works
type Request struct {
	ID        string                 `json:"responseId"`
	SessionID string                 `json:"session"`
	Result    *RequestResult         `json:"queryResult"`
	ODIR      map[string]interface{} `json:"originalDetectIntentRequest"`
}

// GetIntent is Request get
func (req *Request) GetIntent() string {
	return req.Result.Intent
}

// RequestResult into Request V2 API
type RequestResult struct {
	Text                      string                 `json:"queryText"`
	Parameters                *Parameter `json:"parameters"`
	AllRequiredParamsPresent  bool                   `json:"allRequiredParamsPresent"`
	FulfillmentText           string                 `json:"fulfillmentText"`
	FulfillmentMessages       map[string]interface{} `json:"fulfillmentMessages"`
	OutputContexts            map[string]interface{} `json:"outputContexts"`
	Intent                    string  `json:"intent"`
	IntentDetectionConfidence int8                   `json:"intentDetectionConfidence"`
	DiagnosticInfo            map[string]interface{} `json:"diagnosticInfo"`
	LanguageCode              string                 `json:"languageCode"`
}

type Parameter struct{
	name string `json:"parameter_name"`
	value int `json:"parameter_value"`
}

type Response struct {
	fulfillmentText string `json:fulfillmentText`
	fulfillmentMessages[]  map[string]interface{} `json:"fulfillmentMessages"`
	Source        string         `json:"source"`
	payload map[string]interface{} `json:"payload"`
	outputContexts[]	map[string]interface{} `json:"outputContexts"`
	followupEventInput	map[string]interface{} `json:"followupEventInput"`
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

// func (res *Response) AddContext(ctx *OutputContexts) *Response {
// 	res.OutputContexts = append(res.OutputContexts, ctx)
// 	return res
// }

type GoogleData struct {
	ExpectUserResponse bool              `json:"expectUserResponse"`
	IsSSML             bool              `json:"isSsml"`
	NoInputPrompts     []*SimpleResponse `json:"noInputPrompts"`
	RichResponse       *RichResponse     `json:"richResponse"`
	//SystemIntent       *SystemIntent     `json:"systemIntent"`
}

type Item interface {
	isItem() bool
}

type SimpleResponse struct {
	TextToSpeech string `json:"textToSpeech"`
	SSML         string `json:"ssml"`
	DisplayText  string `json:"displayText"`
}

func (s *SimpleResponse) isItem() bool {
	return true
}

type RichResponse struct {
	Items             []*Item            `json:"items"`
	Suggestions       []*Suggestion      `json:"suggestions"`
	LinkOutSuggestion *LinkOutSuggestion `json:"linkOutSuggestion"`
}

func (r *RichResponse) isItem() bool {
	return true
}

type StructuredResponse struct {
	OrderUpdate *OrderUpdate `json:"orderUpdate"`
}

func (s *StructuredResponse) isItem() bool {
	return true
}

type Suggestion struct {
	Title string `json:"title"`
}

type LinkOutSuggestion struct {
	DestinationName string `json:"destinationName"`
	URL             string `json:"url"`
}

type SystemIntent struct {
}

type OrderUpdate struct {
	GoogleOrderID          string                     `json:"googleOrderId"`
	ActionOrderID          string                     `json:"actionOrderId"`
	OrderState             *OrderState                `json:"orderState"`
	OrderManagementActions []*Action                  `json:"orderManagementActions"`
	Receipt                *Receipt                   `json:"receipt"`
	UpdateTime             string                     `json:"updateTime"`
	TotalPrice             *Price                     `json:"totalPrice"`
	LineItemUpdates        map[string]*LineItemUpdate `json:"lineItemUpdates"`
	UserNotification       *UserNotification          `json:"userNotification"`
	//InfoExtension          *InfoExtension             `json:"infoExtension"`
	// TODO: add extra information here.
}

type OrderState struct {
	State string `json:"state"`
	Label string `json:"label"`
}

type Action struct {
	Type   ActionType `json:"type"`
	Button *Button    `json:"button"`
}

type ActionType string

const (
	UNKNOWN          ActionType = "UNKNOWN"
	VIEW_DETAILS                = "VIEW_DETAILS"
	MODIFY                      = "MODIFY"
	CANCEL                      = "CANCEL"
	RETURN                      = "RETURN"
	EXCHANGE                    = "EXCHANGE"
	EMAIL                       = "EMAIL"
	CALL                        = "CALL"
	REORDER                     = "REORDER"
	REVIEW                      = "REVIEW"
	CUSTOMER_SERVICE            = "CUSTOMER_SERVICE"
)

type Button struct {
	Title         string         `json:"title"`
	OpenURLAction *OpenURLAction `json:"openUrlAction"`
}

type OpenURLAction struct {
	URL string `json:"url"`
}

type Receipt struct {
	UserVisibleOrderID string `json:"userVisibleOrderId"`
}

type Price struct {
	Type   PriceType `json:"type"`
	Amount Money     `json:"amount"`
}

type PriceType string

const (
	UNKNOWNP PriceType = "UNKNOWN"
	ESTIMATE           = "ESTIMATE"
	ACTUAL             = "ACTUAL"
)

type Money struct {
	CurrencyCode string `json:"currencyCode"`
	Units        string `json:"units"`
	Nanos        int64  `json:"nanos"`
}

type LineItemUpdate struct {
	OrderState *OrderState `json:"orderState"`
	Price      *Price      `json:"price"`
	Reason     string      `json:"reason"`
	//Extension  *InfoExtension `json:"extension"`
}

type UserNotification struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type FollowupEvent struct {
	EventName string                 `json:"name"`
	Data      map[string]interface{} `json:"data"`
}
