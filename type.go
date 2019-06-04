package main

// https://dialogflow.com/docs/fulfillment#request
type Request struct {
	Lang            string           `json:"lang"`
	Status          *RequestStatus   `json:"status"`
	Timestamp       string           `json:"timestamp"`
	SessionID       string           `json:"sessionId"`
	Result          *RequestResult   `json:"result"`
	ID              string           `json:"id"`
	OriginalRequest *OriginalRequest `json:"originalRequest"`
}

func (req *Request) GetIntent() string {
	return req.Result.Action
}

type RequestStatus struct {
	ErrorType string `json:"errorType"`
	Code      int    `json:"code"`
}

type RequestResult struct {
	Parameters       map[string]interface{} `json:"parameters"`
	Contexts         []*Context             `json:"contexts"`
	ResolvedQuery    string                 `json:"resolvedQuery"`
	Source           string                 `json:"source"`
	Score            float32                `json:"score"`
	Speech           string                 `json:"speech"`
	Fulfillment      *RequestFulfillment    `json:"fulfillment"`
	ActionIncomplete bool                   `json:"actionIncomplete"`
	Action           string                 `json:"action"`
	Metadata         *Metadata              `json:"metadata"`
}

type Context struct {
	ContextName string                 `json:"name"`
	LifeSpan    int                    `json:"lifespan"`
	Parameters  map[string]interface{} `json:"parameters"`
}

type RequestFulfillment struct {
	Messages []*Message `json:"messages"`
	Speech   string     `json:"speech"`
}

type Message struct {
	Speech string `json:"speech"`
	Type   int    `json:"type"`
}

// TODO: check if there are other values than "true" and "false" for webhookUsed and webhookForSlogFillingUsed.
type Metadata struct {
	IntentID                  string `json:"intentId"`
	WebhookForSlotFillingUsed string `json:"webhookForSlotFillingUsed"`
	IntentName                string `json:"intentName"`
	WebhookUsed               string `json:"webhookUsed"`
}

type OriginalRequest struct {
	Source string               `json:"source"`
	Data   *OriginalRequestData `json:"data"`
}

type OriginalRequestData struct {
	Inputs       []*InputData  `json:"inputs"`
	User         *User         `json:"user"`
	Conversation *Conversation `json:"conversation"`
}

type InputData struct {
	RawInputs []*RawInput  `json:"raw_inputs"`
	Intent    string       `json:"intent"`
	Arguments []*Arguments `json:"arguments"`
}

type RawInput struct {
	Query     string `json:"query"`
	InputType int    `json:"input_type"`
}

type Arguments struct {
	TextValue string `json:"text_value"`
	RawText   string `json:"raw_text"`
	Name      string `json:"name"`
}

type User struct {
	UserID string `json:"user_id"`
}

type Conversation struct {
	ConversationID    string      `json:"conversation_id"`
	Type              interface{} `json:"type"`
	ConversationToken string      `json:"conversation_token"`
}

// https://dialogflow.com/docs/fulfillment#response
// https://developers.google.com/actions/dialogflow/webhook#response
// https://developers.google.com/actions/reference/rest/Shared.Types/AppResponse
type Response struct {
	Speech        string         `json:"speech"`
	DisplayText   string         `json:"displayText"`
	ContextOut    []*Context     `json:"contextOut"`
	Source        string         `json:"source"`
	FollowupEvent *FollowupEvent `json:"followupEvent"`
	Data          struct {
		Google *GoogleData `json:"google"`
	} `json:"data"`
}

func NewResponse(speech string) *Response {
	return &Response{
		Speech: speech,
	}
}

func (res *Response) SetDisplayText(text string) *Response {
	res.DisplayText = text
	return res
}

func (res *Response) AddContext(ctx *Context) *Response {
	res.ContextOut = append(res.ContextOut, ctx)
	return res
}

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