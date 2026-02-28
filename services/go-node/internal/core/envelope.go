package core

import "encoding/json"

// Envelope is a standard response envelope for internal and external APIs.
// It keeps responses consistent across services and provides a single point
// to extend metadata (trace id, timestamps, etc).
type Envelope struct {
	OK      bool        `json:"ok"`
	Payload interface{} `json:"payload,omitempty"`
	Error   *ErrBody    `json:"error,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

// ErrBody represents structured error information to transport in responses.
type ErrBody struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

// MarshalEnvelope is a helper to produce JSON bytes of the envelope.
func MarshalEnvelope(ok bool, payload interface{}, err *ErrBody, meta interface{}) ([]byte, error) {
	e := Envelope{
		OK:      ok,
		Payload: payload,
		Error:   err,
		Meta:    meta,
	}
	return json.Marshal(e)
}
