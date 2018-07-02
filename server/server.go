package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sumwonyuno/cp-scoring/auditor"
	"github.com/sumwonyuno/cp-scoring/model"
)

func audit(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		msg := "HTTP 405"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := "ERROR: cannot retrieve body;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}

	var state model.State
	err = json.Unmarshal(body, &state)
	if err != nil {
		msg := "ERROR: cannot unmarshal state;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}

	log.Println("Saving state")
	err = dbInsertState(string(body))
	if err != nil {
		msg := "ERROR: cannot insert state;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}

	log.Println("Auditing state")
	templates, err := dbSelectTemplatesForHostname(state.Hostname)
	if err != nil {
		msg := "ERROR: cannot get templates;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	auditor.Audit(state, templates)

	response := "Received and saved"
	log.Println(response)
	w.Write([]byte(response))
}

func hosts(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// get all hosts
		hosts, err := dbSelectHosts()
		if err != nil {
			msg := "ERROR: cannot retrieve hosts;"
			log.Println(msg, err)
			w.Write([]byte(msg))
			return
		}
		b, err := json.Marshal(hosts)
		if err != nil {
			msg := "ERROR: cannot marshal hosts;"
			log.Println(msg, err)
			w.Write([]byte(msg))
			return
		}
		w.Write(b)
	} else if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			msg := "ERROR: cannot retrieve body;"
			log.Println(msg, err)
			w.Write([]byte(msg))
			return
		}

		var host model.Host
		err = json.Unmarshal(body, &host)
		if err != nil {
			msg := "ERROR: cannot unmarshal host;"
			log.Println(msg, err)
			w.Write([]byte(msg))
			return
		}

		err = dbInsertHost(host)
		if err != nil {
			msg := "ERROR: cannot insert host;"
			log.Println(msg, err)
			w.Write([]byte(msg))
			return
		}

		// new host
		msg := "Saved host"
		log.Println(msg, err)
		w.Write([]byte(msg))
	} else {
		msg := "HTTP 405"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
}

func host(w http.ResponseWriter, r *http.Request) {
	// parse out int64 id
	// remove /hosts/ from URL
	id, err := strconv.ParseInt(r.URL.Path[7:], 10, 64)
	if err != nil {
		msg := "ERROR: cannot parse host id;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}

	host, err := dbSelectHost(id)
	if err != nil {
		msg := "ERROR: cannot retrieve host;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	if (model.Host{}) == host {
		w.Write([]byte("Host not found"))
	} else {
		out, err := json.Marshal(host)
		if err != nil {
			msg := "ERROR: cannot marshal host;"
			log.Println(msg, err)
			w.Write([]byte(msg))
			return
		}
		w.Write([]byte(out))
	}
}

func hostsTemplates(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// get all hosts templates
		hts, err := dbSelectHostsTemplates()
		if err != nil {
			msg := "ERROR: cannot retrieve hosts templates;"
			log.Println(msg, err)
			w.Write([]byte(msg))
			return
		}
		b, err := json.Marshal(hts)
		if err != nil {
			msg := "ERROR: cannot marshal hosts templates;"
			log.Println(msg, err)
			w.Write([]byte(msg))
			return
		}
		w.Write(b)
	} else if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			msg := "ERROR: cannot retrieve body;"
			log.Println(msg, err)
			w.Write([]byte(msg))
			return
		}

		var hostsTemplates model.HostsTemplates
		err = json.Unmarshal(body, &hostsTemplates)
		if err != nil {
			msg := "ERROR: cannot unmarshal hosts templates;"
			log.Println(msg, err)
			w.Write([]byte(msg))
			return
		}

		err = dbInsertHostsTemplates(hostsTemplates.HostID, hostsTemplates.TemplateID)
		if err != nil {
			msg := "ERROR: cannot insert hosts templates;"
			log.Println(msg, err)
			w.Write([]byte(msg))
			return
		}

		// new host
		msg := "Saved hosts templates"
		log.Println(msg, err)
		w.Write([]byte(msg))
	} else {
		msg := "HTTP 405"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
}

func getTemplates(w http.ResponseWriter, r *http.Request) {
	log.Println("get all templates")

	// get all templates
	templates, err := dbSelectTemplates()
	if err != nil {
		msg := "ERROR: cannot retrieve templates;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	b, err := json.Marshal(templates)
	if err != nil {
		msg := "ERROR: cannot marshal templates;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	w.Write(b)
}

func getTemplate(w http.ResponseWriter, r *http.Request) {
	log.Println("get a template")

	// parse out int64 id
	// remove /templates/ from URL
	id, err := strconv.ParseInt(r.URL.Path[11:], 10, 64)
	if err != nil {
		msg := "ERROR: cannot parse template id;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	log.Println(id)
	template, err := dbSelectTemplate(id)
	if err != nil {
		msg := "ERROR: cannot retrieve template;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	out, err := json.Marshal(template)
	if err != nil {
		msg := "ERROR: cannot marshal template;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	w.Write(out)
}

func newTemplate(w http.ResponseWriter, r *http.Request) {
	log.Println("new template")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := "ERROR: cannot retrieve body;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}

	var template model.Template
	err = json.Unmarshal(body, &template)
	if err != nil {
		msg := "ERROR: cannot unmarshal template;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}

	err = dbInsertTemplate(template)
	if err != nil {
		msg := "ERROR: cannot insert template;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}

	// new template
	msg := "Saved template"
	log.Println(msg)
	w.Write([]byte(msg))
}

func editTemplate(w http.ResponseWriter, r *http.Request) {
	log.Println("edit template")

	// parse out int64 id
	// remove /templates/ from URL
	id, err := strconv.ParseInt(r.URL.Path[11:], 10, 64)
	if err != nil {
		msg := "ERROR: cannot parse template id;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	log.Println(id)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := "ERROR: cannot retrieve body;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}

	var template model.Template
	err = json.Unmarshal(body, &template)
	if err != nil {
		msg := "ERROR: cannot unmarshal template;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}

	err = dbUpdateTemplate(id, template)
	if err != nil {
		msg := "ERROR: cannot update template;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}

	msg := "Updated template"
	log.Println(msg)
	w.Write([]byte(msg))
}

func deleteTemplate(w http.ResponseWriter, r *http.Request) {
	log.Println("delete template")

	// parse out int64 id
	// remove /templates/ from URL
	id, err := strconv.ParseInt(r.URL.Path[11:], 10, 64)
	if err != nil {
		msg := "ERROR: cannot parse template id;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
	log.Println(id)
	err = dbDeleteTemplate(id)
	if err != nil {
		msg := "ERROR: cannot delete template;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}
}

func main() {
	dbInit()

	r := mux.NewRouter()
	r.PathPrefix("/ui").Handler(http.StripPrefix("/ui", http.FileServer(http.Dir("./ui/"))))

	r.HandleFunc("/audit", audit)
	templatesRouter := r.PathPrefix("/templates").Subrouter()
	templatesRouter.HandleFunc("", getTemplates).Methods("GET")
	templatesRouter.HandleFunc("/", getTemplates).Methods("GET")
	templatesRouter.HandleFunc("", newTemplate).Methods("POST")
	templatesRouter.HandleFunc("/", newTemplate).Methods("POST")
	templatesRouter.HandleFunc("/{id:[0-9]+}", getTemplate).Methods("GET")
	templatesRouter.HandleFunc("/{id:[0-9]+}", editTemplate).Methods("POST")
	templatesRouter.HandleFunc("/{id:[0-9]+}", deleteTemplate).Methods("DELETE")
	r.HandleFunc("/hosts", hosts)
	r.HandleFunc("/hosts/", host)
	r.HandleFunc("/hosts_templates", hostsTemplates)

	http.ListenAndServe(":8080", r)
}
