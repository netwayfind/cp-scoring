package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"github.com/sumwonyuno/cp-scoring/auditor"
	"github.com/sumwonyuno/cp-scoring/model"
)

func audit(w http.ResponseWriter, r *http.Request) {
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

	log.Println("Saving state")
	DBInsertState(string(body))

	log.Println("Auditing state")
	templates := getTemplates(state.Hostname)
	auditor.Audit(state, templates)

	response := "Received and saved"
	log.Println(response)
	w.Write([]byte(response))
}

func getTemplates(hostname string) []model.Template {
	templates := DBSelectTemplatesForHostname(hostname)
	templateObjs := make([]model.Template, len(templates))
	for i := 0; i < len(templates); i++ {
		var template model.Template
		err := json.Unmarshal([]byte(templates[i]), &template)
		if err != nil {
			log.Println("ERROR: unmarshalling template;", err)
			continue
		}
		templateObjs[i] = template
	}
	return templateObjs
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

func template(w http.ResponseWriter, r *http.Request) {
	// parse out int64 id
	// remove /templates/ from URL
	id, err := strconv.ParseInt(r.URL.Path[11:], 10, 64)
	if err != nil {
		log.Println("ERROR: cannot parse template id;", err)
		return
	}

	template := DBSelectTemplate(id)
	w.Write([]byte(template))
}

func hosts(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// get all hosts
		hosts := DBSelectHosts()
		b, err := json.Marshal(hosts)
		if err != nil {
			msg := "ERROR: returning hosts"
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
	
		var host model.Host
		err = json.Unmarshal(body, &host)
		if err != nil {
			log.Println("ERROR: unmarshalling template;", err)
		}

		DBInsertHost(host)

		// new host
		msg := "Saved host"
		log.Println(msg)
		w.Write([]byte(msg))
	} else {
		w.Write([]byte("HTTP 405\n"))
		return
	}
}

func host(w http.ResponseWriter, r *http.Request) {
	// parse out int64 id
	// remove /hosts/ from URL
	id, err := strconv.ParseInt(r.URL.Path[7:], 10, 64)
	if err != nil {
		log.Println("ERROR: cannot parse host id;", err)
		return
	}

	host, err := DBSelectHost(id)
	if err == nil {
		log.Println(err)
	}
	out, err := json.Marshal(host)
	if err != nil {
		log.Println("ERROR: cannot marshal host;", err)
		return
	}
	w.Write([]byte(out))
}

func main() {
	DBInit()

	http.HandleFunc("/audit", audit)
	http.HandleFunc("/templates", templates)
	http.HandleFunc("/templates/", template)
	http.HandleFunc("/hosts", hosts)
	http.HandleFunc("/hosts/", host)

	http.ListenAndServe(":8080", nil)
}