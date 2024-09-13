package main

type Todo struct {
	Id   string `json:"id"`
	Text string `json:"text"`
}

type ResBody struct {
	Message string `json:"message"`
}
