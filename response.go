package gosoap

import (
	"encoding/xml"
	"fmt"
)

// Response Soap Response
type Response struct {
	Body    []byte
	Headers []byte
	Payload []byte
}

type ResponseWithCustomEnvelope struct {
	Response []byte
	Payload  []byte
}

// FaultError implements error interface
type FaultError struct {
	fault *Fault
}

func (e FaultError) Error() string {
	if e.fault != nil {
		return e.fault.String()
	}

	return ""
}

// IsFault returns whether the given error is a fault error or not.
//
// IsFault will return false when the error could not be typecasted to FaultError, because
// every fault error should have it's dynamic type as FaultError.
func IsFault(err error) bool {
	if _, ok := err.(FaultError); !ok {
		return false
	}

	return true
}

// Unmarshal get the body and unmarshal into the interface
func (r *Response) Unmarshal(v interface{}) error {
	if len(r.Body) == 0 {
		return fmt.Errorf("Body is empty")
	}

	var fault Fault
	err := xml.Unmarshal(r.Body, &fault)
	if err != nil {
		return fmt.Errorf("error unmarshalling the body to Fault: %v", err.Error())
	}

	if fault.Code != "" {
		return FaultError{fault: &fault}
	}

	return xml.Unmarshal(r.Body, v)
}

func (r *ResponseWithCustomEnvelope) CustomUnmarshal(v interface{}) error {
	if len(r.Response) == 0 {
		return fmt.Errorf("Response is empty")
	}

	var fault Fault
	err := xml.Unmarshal(r.Response, &fault)
	if err != nil {
		return fmt.Errorf("error unmarshalling the body to Fault: %v", err.Error())
	}

	if fault.Code != "" {
		return FaultError{fault: &fault}
	}

	return xml.Unmarshal(r.Response, v)
}
