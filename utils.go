package main

import "encoding/json"

func createResponseBody(msg string) []byte {
	var body ResBody = ResBody{Message: msg}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	return bodyBytes
}
