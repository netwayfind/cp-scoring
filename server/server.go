package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"github.com/sumwonyuno/cp-scoring/model"
)

func submit(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Write([]byte("HTTP 405\n"))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("ERROR: retrieving body;", err)
	}

	var state model.State
	err = json.Unmarshal(body, &state)
	if err != nil {
		log.Println("ERROR: unmarshalling state;", err)
	}

	response := "Received"
	log.Println(response)
	w.Write([]byte(response))
}

func main() {
	http.HandleFunc("/submit", submit)

	http.ListenAndServe(":8080", nil)
}