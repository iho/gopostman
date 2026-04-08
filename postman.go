package main

import "encoding/json"

// Collection represents a Postman collection v2.1
type Collection struct {
	Info Info   `json:"info"`
	Item []Item `json:"item"`
}

type Info struct {
	Name string `json:"name"`
}

// Item is either a folder (has Item children) or a request (has Request).
type Item struct {
	Name    string   `json:"name"`
	Item    []Item   `json:"item"`
	Request *Request `json:"request"`
}

type Request struct {
	Method string   `json:"method"`
	URL    URL      `json:"url"`
	Header []Header `json:"header"`
	Body   *Body    `json:"body"`
}

// URL can be a plain string or an object with a "raw" field.
type URL struct {
	Raw string
}

func (u *URL) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		u.Raw = s
		return nil
	}
	var obj struct {
		Raw string `json:"raw"`
	}
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	u.Raw = obj.Raw
	return nil
}

type Header struct {
	Key      string `json:"key"`
	Value    string `json:"value"`
	Disabled bool   `json:"disabled"`
}

type Body struct {
	Mode string `json:"mode"`
	Raw  string `json:"raw"`
}
