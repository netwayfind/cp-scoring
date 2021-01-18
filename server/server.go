package main

import (
	"flag"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/netwayfind/cp-scoring/model"
)

// to be set by build
var version string

func main() {
	// set seed
	rand.Seed(time.Now().UTC().UnixNano())

	// default path
	ex, err := os.Executable()
	if err != nil {
		log.Fatal("ERROR: unable to get executable", err)
	}
	dirWork := filepath.Dir(ex)

	// program arguments
	var askVersion bool
	flag.BoolVar(&askVersion, "version", false, "get version number")
	flag.StringVar(&dirWork, "dir_work", dirWork, "working directory path")
	flag.Parse()

	// version
	if askVersion {
		log.Println("Version: " + version)
		os.Exit(0)
	}

	dirConfig := path.Join(dirWork, "config")
	// dirData := path.Join(dirWork, "data")
	dirUI := path.Join(dirWork, "ui")

	// read config file
	fileConfig := path.Join(dirConfig, "server.conf")
	log.Printf("Reading config file %s\n", fileConfig)
	bytesConfig, err := ioutil.ReadFile(fileConfig)
	if err != nil {
		log.Fatal("ERROR: unable to read config file;", err)
	}
	var port string
	var dbURL string
	var bytesJwtSecret []byte
	for _, line := range strings.Split(string(bytesConfig), "\n") {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}

		tokens := strings.Split(line, " ")
		if tokens[0] == "port" {
			port = tokens[1]
		} else if tokens[0] == "db_url" {
			dbURL = tokens[1]
		} else if tokens[0] == "jwt_secret" {
			bytesJwtSecret = []byte(tokens[1])
		} else {
			log.Fatalf("ERROR: unknown config file setting %s\n", tokens[0])
		}
	}

	backingStore, err := getBackingStore("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	apiHandler := APIHandler{
		BackingStore: backingStore,
		jwtSecret:    bytesJwtSecret,
	}

	// generate default user if no users
	users, err := apiHandler.BackingStore.userSelectAll()
	if err != nil {
		log.Fatal("Could not get users list;", err)
	}
	if len(users) == 0 {
		log.Println("Creating default user")
		if err != nil {
			log.Fatal("ERROR: cannot generate password hash;", err)
		}
		user := model.User{
			Username: "admin",
			Password: "admin",
			Enabled:  true,
			Email:    "",
		}
		_, err := apiHandler.BackingStore.userInsert(user)
		if err != nil {
			log.Fatal("ERROR: cannot save default user;", err)
		}
	}

	// API routing
	r := mux.NewRouter().StrictSlash(true)
	r.Use(apiHandler.middlewareLog)
	r.HandleFunc("/", apiHandler.redirectToUI).Methods("GET")
	r.PathPrefix("/public").Handler(http.FileServer(http.Dir(dirWork)))
	r.PathPrefix("/static").Handler(http.FileServer(http.Dir(dirUI)))
	r.PathPrefix("/ui").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, dirUI+"/index.html")
	})

	apiRouter := r.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/", apiHandler.readAPIRoot).Methods("GET")

	// audit, no auth
	auditRouter := apiRouter.PathPrefix("/audit").Subrouter()
	auditRouter.HandleFunc("/", apiHandler.audit).Methods("POST")

	// host-token, no auth
	hostTokenRouter := apiRouter.PathPrefix("/host-token").Subrouter()
	hostTokenRouter.HandleFunc("/request", apiHandler.requestHostToken).Methods("POST")
	hostTokenRouter.HandleFunc("/register", apiHandler.registerHostToken).Methods("POST")

	// login, no auth
	loginRouter := apiRouter.PathPrefix("/login").Subrouter()
	loginRouter.HandleFunc("/", apiHandler.checkLogin).Methods("GET")
	loginRouter.HandleFunc("/", apiHandler.loginUser).Methods("POST")
	loginRouter.HandleFunc("/team", apiHandler.loginTeam).Methods("POST")

	// logout, auth required
	logoutRouter := apiRouter.PathPrefix("/logout").Subrouter()
	logoutRouter.Use(apiHandler.middlewareAuth)
	logoutRouter.HandleFunc("/", apiHandler.logout).Methods("POST")

	// scenarios, auth required
	scenarioRouter := apiRouter.PathPrefix("/scenarios").Subrouter()
	scenarioRouter.Use(apiHandler.middlewareAuth)
	scenarioRouter.HandleFunc("/", apiHandler.readScenarios).Methods("GET")
	scenarioRouter.HandleFunc("/", apiHandler.createScenario).Methods("POST")
	scenarioRouter.HandleFunc("/{id:[0-9]+}", apiHandler.deleteScenario).Methods("DELETE")
	scenarioRouter.HandleFunc("/{id:[0-9]+}", apiHandler.readScenario).Methods("GET")
	scenarioRouter.HandleFunc("/{id:[0-9]+}", apiHandler.updateScenario).Methods("PUT")
	scenarioRouter.HandleFunc("/{id:[0-9]+}/hosts", apiHandler.readScenarioHosts).Methods("GET")
	scenarioRouter.HandleFunc("/{id:[0-9]+}/hosts", apiHandler.updateScenarioHosts).Methods("PUT")

	// report, no auth
	reportRouter := apiRouter.PathPrefix("/report").Subrouter()
	reportRouter.HandleFunc("/{id:[0-9]+}", apiHandler.readScenarioReport).Methods("GET")
	reportRouter.HandleFunc("/{id:[0-9]+}/hostnames", apiHandler.readScenarioReportHostnames).Methods("GET")
	reportRouter.HandleFunc("/{id:[0-9]+}/timeline", apiHandler.readScenarioReportTimeline).Methods("GET")

	// scenario-desc, no auth
	scenarioDescRouter := apiRouter.PathPrefix("/scenario-desc").Subrouter()
	scenarioDescRouter.HandleFunc("/{id:[0-9]+}", apiHandler.readScenario).Methods("GET")

	// scenario-checks, no auth
	scenarioChecksRouter := apiRouter.PathPrefix("/scenario-checks").Subrouter()
	scenarioChecksRouter.HandleFunc("/{id:[0-9]+}", apiHandler.readScenarioChecks).Methods("GET")

	// scoreboard, no auth
	scoreboardRouter := apiRouter.PathPrefix("/scoreboard").Subrouter()
	scoreboardRouter.HandleFunc("/scenarios", apiHandler.readScoreboardScenarios).Methods("GET")
	scoreboardRouter.HandleFunc("/scenarios/{id:[0-9]+}", apiHandler.readScoreboardForScenario).Methods("GET")

	// teams, auth required
	teamRouter := apiRouter.PathPrefix("/teams").Subrouter()
	teamRouter.Use(apiHandler.middlewareAuth)
	teamRouter.HandleFunc("/", apiHandler.readTeams).Methods("GET")
	teamRouter.HandleFunc("/", apiHandler.createTeam).Methods("POST")
	teamRouter.HandleFunc("/{id:[0-9]+}", apiHandler.deleteTeam).Methods("DELETE")
	teamRouter.HandleFunc("/{id:[0-9]+}", apiHandler.readTeam).Methods("GET")
	teamRouter.HandleFunc("/{id:[0-9]+}", apiHandler.updateTeam).Methods("PUT")

	// users, auth required
	userRouter := apiRouter.PathPrefix("/users").Subrouter()
	userRouter.Use(apiHandler.middlewareAuth)
	userRouter.HandleFunc("/", apiHandler.readUsers).Methods("GET")
	userRouter.HandleFunc("/", apiHandler.createUser).Methods("POST")
	userRouter.HandleFunc("/{id:[0-9]+}", apiHandler.deleteUser).Methods("DELETE")
	userRouter.HandleFunc("/{id:[0-9]+}", apiHandler.readUser).Methods("GET")
	userRouter.HandleFunc("/{id:[0-9]+}", apiHandler.updateUser).Methods("PUT")
	userRouter.HandleFunc("/{id:[0-9]+}/roles", apiHandler.readUserRoles).Methods("GET")
	userRouter.HandleFunc("/{id:[0-9]+}/roles", apiHandler.updateUserRoles).Methods("PUT")

	log.Println("Ready to serve requests")
	addr := "0.0.0.0:" + port
	l, err := net.Listen("tcp4", addr)
	if err != nil {
		log.Fatal(err)
	}

	err = http.Serve(l, r)
	if err != nil {
		log.Fatal("ERROR: cannot start server;", err)
	}

}
