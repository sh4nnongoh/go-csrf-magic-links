package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/goccy/go-json"
)

const (
	usageError = errorType(1 << iota)
	decodeError
	internalError
)

type errorType int

type cookieError struct {
	typ   errorType
	msg   string
	cause error
}

func (e cookieError) Cause() error { return e.cause }

func (e cookieError) Error() string {
	parts := []string{"securecookie: "}
	if e.msg == "" {
		parts = append(parts, "error")
	} else {
		parts = append(parts, e.msg)
	}
	if c := e.Cause(); c != nil {
		parts = append(parts, " - caused by: ", c.Error())
	}
	return strings.Join(parts, "")
}

type JSONEncoder struct{}

func (e JSONEncoder) Serialize(src interface{}) ([]byte, error) {
	var data interface{}
	switch v := src.(type) {
	case map[string]interface{}:
		data = src
	case map[interface{}]interface{}:
		converted := make(map[string]interface{})
		for k, val := range v {
			converted[fmt.Sprint(k)] = val
		}
		data = converted
	default:
		data = src
	}
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	if err := enc.Encode(data); err != nil {
		return nil, cookieError{cause: err, typ: usageError}
	}
	return buf.Bytes(), nil
}

func (e JSONEncoder) Deserialize(src []byte, dst interface{}) error {
	dec := json.NewDecoder(bytes.NewReader(src))
	if err := dec.Decode(dst); err != nil {
		return cookieError{cause: err, typ: decodeError}
	}
	return nil
}
