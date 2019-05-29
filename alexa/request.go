package alexa

import "time"

/*
 * define structs for Alexa Request
 */

// Session is object-type
type Session struct {
	New         bool   `json:"new"`
	SessionID   string `json:"sessionId"`
	Application struct {
		ApplicationID string `json:"applicationId"`
	} `json:"application"`
	Attributes map[string]interface{} `json:"attributes"`
	User       struct {
		UserID      string `json:"userId"`
		Permissions struct {
			ConsentToken string `json:"consentToken"`
		} `json:"permissions"`
		AccessToken string `json:"accessToken"`
	} `json:"user"`
}

// Slot is object-type
type Slot struct {
	Name               string                 `json:"name"`
	Value              string                 `json:"value"`
	ConfirmationStatus string                 `json:"confirmationStatus"`
	Resolutions        map[string]interface{} `json:"resolutions"`
}

// RequestIntent is object-type
type RequestIntent struct {
	Name               string `json:"name"`
	ConfirmationStatus string `json:"confirmationStatus"`
	Slots              map[string]Slot
}

// RequestDetail is map-type
type RequestDetail struct {
	Locale    string        `json:"locale"`
	Timestamp time.Time     `json:"timestamp"`
	Type      string        `json:"type"`
	RequestID string        `json:"requestId"`
	Intent    RequestIntent `json:"intent"`
}

// Request is object-type
type Request struct {
	Version string  `json:"version"`
	Session Session `json:"session"`
	Context struct {
		System struct {
			Application struct {
				ApplicationID string `json:"applicationId"`
			} `json:"application"`
			User struct {
				UserID      string `json:"userId"`
				Permissions struct {
					ConsentToken string `json:"consentToken"`
				} `json:"permissions"`
				AccessToken string `json:"accessToken"`
			} `json:"user"`
			Device struct {
				DeviceID            string `json:"deviceId"`
				SupportedInterfaces struct {
					AudioPlayer struct {
					} `json:"AudioPlayer"`
				} `json:"supportedInterfaces"`
			} `json:"device"`
			APIEndpoint string `json:"apiEndpoint"`
		} `json:"System"`
		AudioPlayer struct {
			Token                string  `json:"token"`
			OffsetInMilliseconds float32 `json:"offsetInMilliseconds"`
			PlayerActivity       string  `json:"playerActivity"`
		} `json:"AudioPlayer"`
	} `json:"context"`
	Request RequestDetail `json:"request"`
}
