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
		msg := "HTTP 405"
		log.Println(msg)
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
	err = DBInsertState(string(body))
	if err != nil {
		msg := "ERROR: cannot insert state;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}

	log.Println("Auditing state")
	templates, err := DBSelectTemplatesForHostname(state.Hostname)
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

func templates(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// get all templates
		templates, err := DBSelectTemplates()
		if err != nil {
			msg := "ERROR: cannot retrieve templates;"
			log.Println(msg)
			w.Write([]byte(msg))
			return
		}
		b, err := json.Marshal(templates)
		if err != nil {
			msg := "ERROR: cannot marshal templates;"
			log.Println(msg)
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
	
		var template model.Template
		err = json.Unmarshal(body, &template)
		if err != nil {
			msg := "ERROR: cannot unmarshal template;"
			log.Println(msg, err)
			w.Write([]byte(msg))
			return
		}

		err = DBInsertTemplate(string(body))
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
	} else {
		msg := "HTTP 405"
		log.Println(msg)
		w.Write([]byte(msg))
		return
	}
}

func template(w http.ResponseWriter, r *http.Request) {
	// parse out int64 id
	// remove /templates/ from URL
	id, err := strconv.ParseInt(r.URL.Path[11:], 10, 64)
	if err != nil {
		msg := "ERROR: cannot parse template id;"
		log.Println(msg, err)
		w.Write([]byte(msg))
		return
	}

	template, err := DBSelectTemplate(id)
	if err != nil {
		msg := "ERROR: cannot retrieve template;"
		log.Println(msg)
		w.Write([]byte(msg))
		return
	}
	w.Write([]byte(template))
}

func hosts(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// get all hosts
		hosts, err := DBSelectHosts()
		if err != nil {
			msg := "ERROR: cannot retrieve hosts;"
			log.Println(msg)
			w.Write([]byte(msg))
			return
		}
		b, err := json.Marshal(hosts)
		if err != nil {
			msg := "ERROR: cannot marshal hosts;"
			log.Println(msg)
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

		err = DBInsertHost(host)
		if err != nil {
			msg := "ERROR: cannot insert host;"
			log.Println(msg, err)
			w.Write([]byte(msg))
			return
		}

		// new host
		msg := "Saved host"
		log.Println(msg)
		w.Write([]byte(msg))
	} else {
		msg := "HTTP 405"
		log.Println(msg)
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

	host, err := DBSelectHost(id)
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

func hosts_templates(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// get all hosts templates
		hosts_templates, err := DBSelectHostsTemplates()
		if err != nil {
			msg := "ERROR: cannot retrieve hosts templates;"
			log.Println(msg)
			w.Write([]byte(msg))
			return
		}
		b, err := json.Marshal(hosts_templates)
		if err != nil {
			msg := "ERROR: cannot marshal hosts templates;"
			log.Println(msg)
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

		err = DBInsertHostsTemplates(hostsTemplates.HostId, hostsTemplates.TemplateId)
		if err != nil {
			msg := "ERROR: cannot insert hosts templates;"
			log.Println(msg, err)
			w.Write([]byte(msg))
			return
		}

		// new host
		msg := "Saved hosts templates"
		log.Println(msg)
		w.Write([]byte(msg))
	} else {
		msg := "HTTP 405"
		log.Println(msg)
		w.Write([]byte(msg))
		return
	}
}

func main() {
	DBInit()

	http.HandleFunc("/audit", audit)
	http.HandleFunc("/templates", templates)
	http.HandleFunc("/templates/", template)
	http.HandleFunc("/hosts", hosts)
	http.HandleFunc("/hosts/", host)
	http.HandleFunc("/hosts_templates", hosts_templates)

	http.ListenAndServe(":8080", nil)
}