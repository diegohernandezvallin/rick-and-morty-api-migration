package model

import "time"

type HttpClientResponse struct {
	ResponseBody []byte
	StatusCode   int
}

type Info struct {
	Count int    `json:"count"`
	Pages int    `json:"pages"`
	Next  string `json:"next"`
	Prev  string `json:"prev"`
}

type Location struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Dimension string    `json:"dimension"`
	Residents []string  `json:"residents"`
	URL       string    `json:"url"`
	Created   time.Time `json:"created"`
}

type Character struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Status   string    `json:"status"`
	Species  string    `json:"species"`
	Type     string    `json:"type"`
	Gender   string    `json:"gender"`
	Origin   Location  `json:"origin"`
	Location Location  `json:"location"`
	Image    string    `json:"image"`
	Episode  []string  `json:"episode"`
	URL      string    `json:"url"`
	Created  time.Time `json:"created"`
}

type CharacterResponse struct {
	Info    Info        `json:"info"`
	Results []Character `results:"info"`
}

type LocationResponse struct {
	Info    Info       `json:"info"`
	Results []Location `results:"info"`
}

type Message struct {
	Payload interface{}
	Headers map[string]string
	Key     string
}
