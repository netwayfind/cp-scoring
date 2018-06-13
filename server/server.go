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

	DBInsertState(string(body))

	response := "Received and saved"
	log.Println(response)
	w.Write([]byte(response))
}

func templates(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// get all templates
		templates := DBSelectTemplates()
		b, err := json.Marshal(templates)
		if err != nil {
			msg := "ERROR: returning templates"
			log.Println(msg)
			w.Write([]byte(msg))
			return
		}
		w.Write(b)
	} else if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("ERROR: retrieving body;", err)
		}
	
		var template model.Template
		err = json.Unmarshal(body, &template)
		if err != nil {
			log.Println("ERROR: unmarshalling template;", err)
		}

		DBInsertTemplate(string(body))

		// new template
		msg := "Saved template"
		log.Println(msg)
		w.Write([]byte(msg))
	} else {
		w.Write([]byte("HTTP 405\n"))
		return
	}
}

func main() {
	DBInit()

	http.HandleFunc("/submit", submit)
	http.HandleFunc("/templates", templates)

	http.ListenAndServe(":8080", nil)
}