package json

import (
	"encoding/json"
	"io/ioutil"
)

func New(opts ...string) *Json {

	var indent string
	var prefix string

	if len(opts) == 1 {
		indent = opts[0]
		prefix = ""
	}
	if len(opts) == 2 {
		indent = opts[0]
		prefix = opts[1]
	}

	return &Json{
		Indent: indent,
		Prefix: prefix,
	}
}

type StringifyParser interface {
	Stringifyer
	Parser
}

type Stringifyer interface {
	Stringify(config interface{}) ([]byte, error)
}

type Parser interface {
	Parse(file string, s interface{}) interface{}
}

type Json struct {
	Indent string
	Prefix string
}

func (j Json) Stringify(config interface{}) ([]byte, error) {

	if j.Indent == "" {
		j.Indent = "\t"
	}
	if j.Prefix == "" {
		j.Prefix = ""
	}

	return json.MarshalIndent(config, j.Prefix, j.Indent)
}

func (j Json) Parse(file string, s interface{}) interface{} {

	raw, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	json.Unmarshal(raw, &s)
	return s
}
