package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/netwayfind/cp-scoring/model"
	"github.com/netwayfind/cp-scoring/processing"
	"golang.org/x/crypto/openpgp"
)

// APIHandler asdf
type APIHandler struct {
	BackingStore backingStore
	jwtSecret    []byte
	dirResults   string
	entities     openpgp.EntityList
}

func (handler APIHandler) middlewareLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method + " " + r.URL.String())

		next.ServeHTTP(w, r)
	})
}

func (handler APIHandler) middlewareAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwtCookie, err := r.Cookie(model.AuthCookieName)
		if err != nil {
			httpErrorNotAuthenticated(w)
			return
		}

		token, err := getJwtToken(handler.jwtSecret, jwtCookie.Value)
		claims := token.Claims.(jwt.MapClaims)

		// TODO: check roles
		if len(claims) == 0 {
			httpErrorNotAuthenticated(w)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (handler APIHandler) middlewareTeam(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwtCookie, err := r.Cookie(model.TeamCookieName)
		if err != nil {
			httpErrorNotAuthenticated(w)
			return
		}

		token, err := getJwtToken(handler.jwtSecret, jwtCookie.Value)
		claims := token.Claims.(jwt.MapClaims)

		// TODO: check roles
		if len(claims) == 0 {
			httpErrorNotAuthenticated(w)
			return
		}

		ctx := context.WithValue(r.Context(), model.TeamCookieName, uint64(claims["TeamID"].(float64)))

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getJwtToken(jwtSecret []byte, jwtStr string) (*jwt.Token, error) {
	if len(jwtStr) < 5 {
		return nil, errors.New("invalid jwt length")
	}

	token, err := jwt.Parse(jwtStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			errMsg := fmt.Sprintf("%s", token.Header["alg"])
			return nil, errors.New(errMsg)
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid jwt")
	}

	return token, nil
}

func getRequestID(r *http.Request) (uint64, error) {
	vars := mux.Vars(r)

	return strconv.ParseUint(vars["id"], 10, 64)
}

func getSourceIP(r *http.Request) string {
	conn := r.RemoteAddr
	ips := r.Header.Get("X-Forwarded-For")
	if ips != "" {
		conn = strings.Split(ips, ",")[0]
	}
	return strings.Split(conn, ":")[0]
}

func httpErrorBadRequest(w http.ResponseWriter) {
	msg := "ERROR: bad request;"
	http.Error(w, msg, http.StatusBadRequest)
}

func httpErrorDatabase(w http.ResponseWriter, err error) {
	msg := "ERROR: database query;"
	log.Println(msg, err)
	http.Error(w, msg, http.StatusInternalServerError)
}

func httpErrorForbidden(w http.ResponseWriter) {
	msg := "ERROR: forbidden;"
	http.Error(w, msg, http.StatusForbidden)
}

func httpErrorInternal(w http.ResponseWriter, err error) {
	msg := "ERROR: internal server error;"
	log.Println(msg, err)
	http.Error(w, msg, http.StatusInternalServerError)
}

func httpErrorInvalidID(w http.ResponseWriter) {
	msg := "ERROR: cannot parse scenario id;"
	log.Println(msg)
	http.Error(w, msg, http.StatusBadRequest)
}

func httpErrorMarshall(w http.ResponseWriter, err error) {
	msg := "ERROR: cannot marshall;"
	log.Println(msg, err)
	http.Error(w, msg, http.StatusInternalServerError)
}

func httpErrorNotFound(w http.ResponseWriter) {
	msg := "ERROR: not found;"
	log.Println(msg)
	http.Error(w, msg, http.StatusNotFound)
}

func httpErrorReadRequestBody(w http.ResponseWriter, err error) {
	msg := "ERROR: cannot read request body;"
	log.Println(msg, err)
	http.Error(w, msg, http.StatusInternalServerError)
}

func httpErrorNotAuthenticated(w http.ResponseWriter) {
	msg := "ERROR: not authenticated;"
	log.Println(msg)
	http.Error(w, msg, http.StatusUnauthorized)
}

func httpErrorUnmarshall(w http.ResponseWriter, err error) {
	msg := "ERROR: cannot unmarshall;"
	log.Println(msg, err)
	http.Error(w, msg, http.StatusBadRequest)
}

func httpNotModified(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotModified)
}

func readRequestBody(w http.ResponseWriter, r *http.Request, o interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		httpErrorReadRequestBody(w, err)
		return err
	}

	err = json.Unmarshal(body, &o)
	if err != nil {
		httpErrorUnmarshall(w, err)
		return err
	}

	return err
}

func sendResponse(w http.ResponseWriter, o interface{}) {
	b, err := json.Marshal(o)
	if err != nil {
		httpErrorMarshall(w, err)
		return
	}
	w.Write(b)
}

func (handler APIHandler) readAPIRoot(w http.ResponseWriter, r *http.Request) {
	sendResponse(w, "OK")
}

func (handler APIHandler) readAPIVersion(w http.ResponseWriter, r *http.Request) {
	sendResponse(w, version)
}

func (handler APIHandler) redirectToUI(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/ui", http.StatusPermanentRedirect)
}

func (handler APIHandler) audit(w http.ResponseWriter, r *http.Request) {
	log.Println("audit")

	source := getSourceIP(r)
	timestamp := time.Now().Unix()
	var bs []byte
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		httpErrorInternal(w, errors.New("ERROR: unable to read request body"))
		return
	}
	result, err := processing.FromBytes(bs, handler.entities)
	if err != nil {
		httpErrorBadRequest(w)
		return
	}
	entry := model.AuditQueueEntry{
		Timestamp: timestamp,
		Source:    source,
		Body:      result,
	}
	err = handler.BackingStore.auditQueueInsert(entry)
	if err != nil {
		httpErrorInternal(w, errors.New("ERROR: Unable to save to audit queue"))
	}
}

func (handler APIHandler) auditEntries(entries []model.AuditQueueEntry) error {
	log.Println("audit entries")

	for _, entry := range entries {
		err := handler.auditEntry(entry)
		if err != nil {
			log.Println("ERROR: unable to audit entry;", err)
			err2 := handler.BackingStore.auditQueueUpdateStatusFailed(entry.ID)
			if err2 != nil {
				return err2
			}
			return err
		}
		err = handler.BackingStore.auditQueueDelete(entry.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (handler APIHandler) auditEntry(entry model.AuditQueueEntry) error {
	log.Printf("audit entry %d", entry.ID)

	auditCheckResults := entry.Body
	timestamp := entry.Timestamp
	source := entry.Source

	scenario, err := handler.BackingStore.scenarioSelect(auditCheckResults.ScenarioID)
	if err != nil {
		log.Println("ERROR: unable to read scenario;", err)
		return err
	}
	if scenario.ID == 0 {
		return errors.New("ERROR: scenario not found;")
	}

	if len(auditCheckResults.HostToken) == 0 {
		return err
	}

	hostname, err := handler.BackingStore.hostTokenSelectHostname(auditCheckResults.HostToken)
	if err != nil {
		log.Println("ERROR: unable to read hostname from host token;", err)
		return err
	}
	if len(hostname) == 0 {
		return errors.New("ERROR: hostname not found;")
	}

	lastModified, err := handler.BackingStore.scenarioHostsSelectLastModified(scenario.ID, hostname)
	lastModifiedStr := time.Unix(lastModified, 0).Format(model.JavascriptDateFormat)
	if auditCheckResults.ChecksLastModified != lastModifiedStr {
		return fmt.Errorf("ERROR: expected last modified %s, received %s", lastModifiedStr, auditCheckResults.ChecksLastModified)
	}

	teamID, err := handler.BackingStore.hostTokenSelectTeamID(auditCheckResults.HostToken)
	if err != nil {
		return err
	}
	if teamID == 0 {
		return errors.New("ERROR: team not found;")
	}

	checkResultsID, err := handler.BackingStore.auditCheckResultsInsert(auditCheckResults, teamID, timestamp, source)
	if err != nil {
		return err
	}

	answers, err := handler.BackingStore.scenarioHostsSelectAnswers(auditCheckResults.ScenarioID, hostname)
	if err != nil {
		return err
	}

	if len(answers) != len(auditCheckResults.CheckResults) {
		return errors.New("ERROR: result count, answer count mismatch;")
	}

	answerResults := make([]model.AnswerResult, len(answers))
	score := 0
	for i, answer := range answers {
		checkResult := auditCheckResults.CheckResults[i]
		points := 0
		if answer.Operator == model.OperatorTypeEqual {
			if answer.Value == checkResult {
				points = answer.Points
				score += points
			}
		} else if answer.Operator == model.OperatorTypeNotEqual {
			if answer.Value != checkResult {
				points = answer.Points
				score += points
			}
		}
		answerResults[i] = model.AnswerResult{
			Description: answer.Description,
			Points:      points,
		}
	}

	auditAnswerResults := model.AuditAnswerResults{
		ScenarioID:     scenario.ID,
		TeamID:         teamID,
		HostToken:      auditCheckResults.HostToken,
		Timestamp:      auditCheckResults.Timestamp,
		CheckResultsID: checkResultsID,
		Score:          score,
		AnswerResults:  answerResults,
	}

	err = handler.BackingStore.auditAnswerResultsInsert(auditAnswerResults)
	if err != nil {
		return err
	}

	err = handler.BackingStore.scoreboardUpdate(scenario.ID, teamID, hostname, score, auditCheckResults.Timestamp)
	if err != nil {
		return err
	}

	return nil
}

func (handler APIHandler) requestHostToken(w http.ResponseWriter, r *http.Request) {
	log.Println("request host token")

	var hostTokenRequest model.HostTokenRequest
	err := readRequestBody(w, r, &hostTokenRequest)
	if err != nil {
		return
	}

	hostToken := randHexStr(16)
	hostname := hostTokenRequest.Hostname
	scenarioID := hostTokenRequest.ScenarioID
	timestamp := time.Now().Unix()
	sourceIP := getSourceIP(r)

	// make sure scenario + hostname exists
	s, err := handler.BackingStore.scenarioHostsSelectChecks(scenarioID, hostname)
	if err != nil {
		httpErrorDatabase(w, err)
		return
	}
	if s == nil {
		httpErrorNotFound(w)
		return
	}

	err = handler.BackingStore.hostTokenInsert(hostToken, hostname, timestamp, sourceIP)
	if err != nil {
		httpErrorDatabase(w, err)
		return
	}

	sendResponse(w, hostToken)
}

func (handler APIHandler) registerHostToken(w http.ResponseWriter, r *http.Request) {
	log.Println("register host token")

	var hostTokenRegistration model.HostTokenRegistration
	err := readRequestBody(w, r, &hostTokenRegistration)
	if err != nil {
		return
	}

	team, err := handler.BackingStore.teamSelectByKey(hostTokenRegistration.TeamKey)
	if err != nil {
		httpErrorDatabase(w, err)
		return
	}
	if team.ID == 0 {
		httpErrorNotFound(w)
		return
	}

	hostToken := hostTokenRegistration.HostToken
	timestamp := time.Now().Unix()

	err = handler.BackingStore.teamHostTokenInsert(team.ID, hostToken, timestamp)
	if err != nil {
		httpErrorDatabase(w, err)
		return
	}
}

func (handler APIHandler) checkLoginUser(w http.ResponseWriter, r *http.Request) {
	log.Println("check login user")

	jwtCookie, err := r.Cookie(model.AuthCookieName)
	if err != nil {
		httpErrorNotAuthenticated(w)
		return
	}

	if len(jwtCookie.Value) < 5 {
		httpErrorNotAuthenticated(w)
		return
	}

	token, err := getJwtToken(handler.jwtSecret, jwtCookie.Value)
	if err != nil {
		httpErrorNotAuthenticated(w)
		return
	}
	if !token.Valid {
		httpErrorNotAuthenticated(w)
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := uint64(claims["UserID"].(float64))
	user, err := handler.BackingStore.userSelect(userID)
	if err != nil {
		httpErrorDatabase(w, err)
		return
	}

	sendResponse(w, user.Username)
}

func (handler APIHandler) checkLoginTeam(w http.ResponseWriter, r *http.Request) {
	log.Println("check login team")

	jwtCookie, err := r.Cookie(model.TeamCookieName)
	if err != nil {
		httpErrorNotAuthenticated(w)
		return
	}

	if len(jwtCookie.Value) < 5 {
		httpErrorNotAuthenticated(w)
		return
	}

	token, err := getJwtToken(handler.jwtSecret, jwtCookie.Value)
	if err != nil {
		httpErrorNotAuthenticated(w)
		return
	}
	if !token.Valid {
		httpErrorNotAuthenticated(w)
	}

	claims := token.Claims.(jwt.MapClaims)
	teamID := uint64(claims["TeamID"].(float64))
	team, err := handler.BackingStore.teamSelect(teamID)
	if err != nil {
		httpErrorDatabase(w, err)
		return
	}

	sendResponse(w, team.Name)
}

func (handler APIHandler) loginUser(w http.ResponseWriter, r *http.Request) {
	log.Println("login user")

	var loginUser model.LoginUser
	err := readRequestBody(w, r, &loginUser)
	if err != nil {
		return
	}

	user, err := handler.BackingStore.userSelectByUsername(loginUser.Username)
	if err != nil {
		httpErrorDatabase(w, err)
		return
	}

	if !checkPasswordHash(loginUser.Password, user.Password) {
		log.Println("user authentication failed: " + loginUser.Username)
		httpErrorNotAuthenticated(w)
		return
	}

	log.Println("user authentication successful: " + loginUser.Username)

	roles, err := handler.BackingStore.userRolesSelect(user.ID)
	if err != nil {
		httpErrorDatabase(w, err)
		return
	}

	claims := model.ClaimsAuth{
		UserID: user.ID,
		Roles:  roles,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(handler.jwtSecret)
	if err != nil {
		httpErrorInternal(w, err)
		return
	}

	cookie := &http.Cookie{
		Name:     model.AuthCookieName,
		Value:    signedToken,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().AddDate(0, 0, 1),
		SameSite: http.SameSiteLaxMode,
		// TODO: Secure when enforcing HTTPS
	}
	http.SetCookie(w, cookie)

	sendResponse(w, user.Username)
}

func (handler APIHandler) loginTeam(w http.ResponseWriter, r *http.Request) {
	log.Println("login team")

	var loginTeam model.LoginTeam
	err := readRequestBody(w, r, &loginTeam)
	if err != nil {
		return
	}

	team, err := handler.BackingStore.teamSelectByKey(loginTeam.TeamKey)
	if err != nil {
		httpErrorDatabase(w, err)
		return
	}
	if team.ID == 0 {
		log.Printf("team authentication failed " + loginTeam.TeamKey)
		httpErrorNotAuthenticated(w)
		return
	}

	log.Printf("team authentication successful: %d", team.ID)

	claims := model.ClaimsTeam{
		TeamID: team.ID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(handler.jwtSecret)
	if err != nil {
		httpErrorInternal(w, err)
		return
	}

	cookie := &http.Cookie{
		Name:     model.TeamCookieName,
		Value:    signedToken,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().AddDate(0, 0, 1),
		SameSite: http.SameSiteLaxMode,
		// TODO: Secure when enforcing HTTPS
	}
	http.SetCookie(w, cookie)

	sendResponse(w, team.Name)
}

func (handler APIHandler) logoutUser(w http.ResponseWriter, r *http.Request) {
	log.Println("logout user")

	cookie := &http.Cookie{
		Name:     model.AuthCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
		SameSite: http.SameSiteLaxMode,
		// TODO: Secure when enforcing HTTPS
	}
	http.SetCookie(w, cookie)
	return
}

func (handler APIHandler) logoutTeam(w http.ResponseWriter, r *http.Request) {
	log.Println("logout team")

	cookie := &http.Cookie{
		Name:     model.TeamCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
		SameSite: http.SameSiteLaxMode,
		// TODO: Secure when enforcing HTTPS
	}
	http.SetCookie(w, cookie)
	return
}

func (handler APIHandler) createScenario(w http.ResponseWriter, r *http.Request) {
	log.Println("create scenario")

	var scenario model.Scenario
	err := readRequestBody(w, r, &scenario)
	if err != nil {
		return
	}

	s, err := handler.BackingStore.scenarioInsert(scenario)
	if err != nil {
		httpErrorDatabase(w, err)
		return
	}

	sendResponse(w, s)
}

func (handler APIHandler) deleteScenario(w http.ResponseWriter, r *http.Request) {
	log.Println("delete scenario")

	id, err := getRequestID(r)
	if err != nil {
		httpErrorInvalidID(w)
		return
	}

	scenario, err := handler.BackingStore.scenarioSelect(id)
	if err != nil {
		httpErrorDatabase(w, err)
		return
	}
	if scenario.ID == 0 {
		httpErrorNotFound(w)
		return
	}

	err = handler.BackingStore.scenarioDelete(id)
	if err != nil {
		httpErrorDatabase(w, err)
	}
}

func (handler APIHandler) readScenario(w http.ResponseWriter, r *http.Request) {
	log.Println("read scenario")

	id, err := getRequestID(r)
	if err != nil {
		httpErrorInvalidID(w)
		return
	}

	s, err := handler.BackingStore.scenarioSelect(id)
	if err != nil {
		httpErrorDatabase(w, err)
		return
	}
	if s.ID == 0 {
		httpErrorNotFound(w)
		return
	}

	sendResponse(w, s)
}

func (handler APIHandler) readScenarios(w http.ResponseWriter, r *http.Request) {
	log.Println("read scenarios")

	s, err := handler.BackingStore.scenarioSelectAll()
	if err != nil {
		httpErrorDatabase(w, err)
		return
	}

	sendResponse(w, s)
}

func (handler APIHandler) updateScenario(w http.ResponseWriter, r *http.Request) {
	log.Println("update scenario")

	id, err := getRequestID(r)
	if err != nil {
		httpErrorInvalidID(w)
		return
	}

	var scenario model.Scenario
	err = readRequestBody(w, r, &scenario)
	if err != nil {
		return
	}

	s, err := handler.BackingStore.scenarioUpdate(id, scenario)
	if err != nil {
		if err.Error() == model.ErrorDBUpdateNoChange {
			httpErrorNotFound(w)
			return
		}
		httpErrorDatabase(w, err)
		return
	}

	sendResponse(w, s)
}

func (handler APIHandler) readScenarioChecks(w http.ResponseWriter, r *http.Request) {
	log.Println("read scenario checks")

	id, err := getRequestID(r)
	if err != nil {
		httpErrorInvalidID(w)
		return
	}

	hostnameParam, present := r.URL.Query()["hostname"]
	if !present || len(hostnameParam) != 1 {
		httpErrorBadRequest(w)
		return
	}
	hostname := hostnameParam[0]

	modifiedSince := r.Header.Get("If-Modified-Since")
	if len(modifiedSince) == 0 {
		modifiedSince = "Thu, 01 Jan 1970 00:00:00 GMT"
	}
	t, err := time.Parse(model.JavascriptDateFormat, modifiedSince)
	if err != nil {
		log.Println(err)
		httpErrorBadRequest(w)
		return
	}
	lastModified, err := handler.BackingStore.scenarioHostsSelectLastModified(id, hostname)
	if err != nil {
		httpErrorDatabase(w, err)
		return
	}
	if lastModified <= t.Unix() {
		httpNotModified(w)
		return
	}

	s, err := handler.BackingStore.scenarioHostsSelectChecks(id, hostname)
	if err != nil {
		httpErrorDatabase(w, err)
		return
	}
	if s == nil {
		httpErrorNotFound(w)
		return
	}

	w.Header().Set("Last-Modified", time.Unix(lastModified, 0).Format(model.JavascriptDateFormat))
	sendResponse(w, s)
}

func (handler APIHandler) readScenarioConfig(w http.ResponseWriter, r *http.Request) {
	log.Println("read scenario config")

	id, err := getRequestID(r)
	if err != nil {
		httpErrorInvalidID(w)
		return
	}

	hostnameParam, present := r.URL.Query()["hostname"]
	if !present || len(hostnameParam) != 1 {
		httpErrorBadRequest(w)
		return
	}
	hostname := hostnameParam[0]

	s, err := handler.BackingStore.scenarioHostsSelectConfig(id, hostname)
	if err != nil {
		httpErrorDatabase(w, err)
		return
	}
	if s == nil {
		httpErrorNotFound(w)
		return
	}

	sendResponse(w, s)
}

func (handler APIHandler) readScenarioHosts(w http.ResponseWriter, r *http.Request) {
	log.Println("read scenario hosts")

	id, err := getRequestID(r)
	if err != nil {
		httpErrorInvalidID(w)
		return
	}

	s, err := handler.BackingStore.scenarioHostsSelectAll(id)
	if err != nil {
		httpErrorDatabase(w, err)
		return
	}

	sendResponse(w, s)
}

func (handler APIHandler) updateScenarioHosts(w http.ResponseWriter, r *http.Request) {
	log.Println("update scenario hosts")

	id, err := getRequestID(r)
	if err != nil {
		httpErrorInvalidID(w)
		return
	}

	var hostMap map[string]model.ScenarioHost
	err = readRequestBody(w, r, &hostMap)
	if err != nil {
		return
	}

	err = handler.BackingStore.scenarioHostsUpdate(id, hostMap)
	if err != nil {
		if err.Error() == model.ErrorDBUpdateNoChange {
			httpErrorNotFound(w)
			return
		}
		httpErrorDatabase(w, err)
		return
	}
}

func (handler APIHandler) readScenarioReport(w http.ResponseWriter, r *http.Request) {
	log.Println("read scenario report")

	id, err := getRequestID(r)
	if err != nil {
		httpErrorInvalidID(w)
		return
	}

	teamIDValue := r.Context().Value(model.TeamCookieName)
	if teamIDValue == nil {
		httpErrorNotAuthenticated(w)
		return
	}
	teamID := teamIDValue.(uint64)

	team, err := handler.BackingStore.teamSelect(teamID)
	if err != nil {
		httpErrorDatabase(w, err)
		return
	}

	hostnameParam, present := r.URL.Query()["hostname"]
	if !present || len(hostnameParam) != 1 {
		httpErrorBadRequest(w)
		return
	}
	hostname := hostnameParam[0]

	s, err := handler.BackingStore.auditAnswerResultsReport(id, team.ID, hostname)
	if err != nil {
		httpErrorDatabase(w, err)
		return
	}
	if s.AnswerResults == nil {
		httpErrorNotFound(w)
		return
	}

	// filter out 0 points to not hint on answers
	filtered := make([]model.AnswerResult, 0)
	for _, answerResult := range s.AnswerResults {
		if answerResult.Points != 0 {
			filtered = append(filtered, answerResult)
		}
	}
	s.AnswerResults = filtered

	sendResponse(w, s)
}

func (handler APIHandler) readScenarioReportHostnames(w http.ResponseWriter, r *http.Request) {
	log.Println("read scenario report hostnames")

	id, err := getRequestID(r)
	if err != nil {
		httpErrorInvalidID(w)
		return
	}

	teamIDValue := r.Context().Value(model.TeamCookieName)
	if teamIDValue == nil {
		httpErrorNotAuthenticated(w)
		return
	}
	teamID := teamIDValue.(uint64)

	team, err := handler.BackingStore.teamSelect(teamID)
	if err != nil {
		httpErrorDatabase(w, err)
		return
	}

	s, err := handler.BackingStore.auditAnswerResultsSelectHostnames(id, team.ID)
	if err != nil {
		httpErrorDatabase(w, err)
		return
	}

	sendResponse(w, s)
}

func (handler APIHandler) readScenarioReportTimeline(w http.ResponseWriter, r *http.Request) {
	log.Println("read scenario report timeline")

	id, err := getRequestID(r)
	if err != nil {
		httpErrorInvalidID(w)
		return
	}

	teamIDValue := r.Context().Value(model.TeamCookieName)
	if teamIDValue == nil {
		httpErrorNotAuthenticated(w)
		return
	}
	teamID := teamIDValue.(uint64)

	team, err := handler.BackingStore.teamSelect(teamID)
	if err != nil {
		httpErrorDatabase(w, err)
		return
	}

	hostnameParam, present := r.URL.Query()["hostname"]
	if !present || len(hostnameParam) != 1 {
		httpErrorBadRequest(w)
		return
	}
	hostname := hostnameParam[0]

	s, err := handler.BackingStore.auditAnswerResultsReportTimeline(id, team.ID, hostname)
	if err != nil {
		httpErrorDatabase(w, err)
		return
	}

	sendResponse(w, s)
}

func (handler APIHandler) readScoreboardForScenario(w http.ResponseWriter, r *http.Request) {
	log.Println("read scoreboard for scenarios")

	id, err := getRequestID(r)
	if err != nil {
		httpErrorInvalidID(w)
		return
	}

	s, err := handler.BackingStore.scoreboardSelectByScenarioID(id)
	if err != nil {
		httpErrorDatabase(w, err)
		return
	}

	sendResponse(w, s)
}

func (handler APIHandler) readScoreboardScenarios(w http.ResponseWriter, r *http.Request) {
	log.Println("read scoreboard scenarios")

	s, err := handler.BackingStore.scoreboardSelectScenarios()
	if err != nil {
		httpErrorDatabase(w, err)
		return
	}

	sendResponse(w, s)
}

func (handler APIHandler) createTeam(w http.ResponseWriter, r *http.Request) {
	log.Println("create team")

	var team model.Team
	err := readRequestBody(w, r, &team)
	if err != nil {
		return
	}

	t, err := handler.BackingStore.teamInsert(team)
	if err != nil {
		httpErrorDatabase(w, err)
		return
	}

	sendResponse(w, t)
}

func (handler APIHandler) deleteTeam(w http.ResponseWriter, r *http.Request) {
	log.Println("delete team")

	id, err := getRequestID(r)
	if err != nil {
		httpErrorInvalidID(w)
		return
	}

	team, err := handler.BackingStore.teamSelect(id)
	if err != nil {
		httpErrorDatabase(w, err)
		return
	}
	if team.ID == 0 {
		httpErrorNotFound(w)
		return
	}

	err = handler.BackingStore.teamDelete(id)
	if err != nil {
		httpErrorDatabase(w, err)
	}
}

func (handler APIHandler) readTeam(w http.ResponseWriter, r *http.Request) {
	log.Println("read team")

	id, err := getRequestID(r)
	if err != nil {
		httpErrorInvalidID(w)
		return
	}

	t, err := handler.BackingStore.teamSelect(id)
	if err != nil {
		httpErrorDatabase(w, err)
		return
	}
	if t.ID == 0 {
		httpErrorNotFound(w)
		return
	}

	sendResponse(w, t)
}

func (handler APIHandler) readTeams(w http.ResponseWriter, r *http.Request) {
	log.Println("read teams")

	t, err := handler.BackingStore.teamSelectAll()
	if err != nil {
		httpErrorDatabase(w, err)
		return
	}

	sendResponse(w, t)
}

func (handler APIHandler) updateTeam(w http.ResponseWriter, r *http.Request) {
	log.Println("update team")

	id, err := getRequestID(r)
	if err != nil {
		httpErrorInvalidID(w)
		return
	}

	var team model.Team
	err = readRequestBody(w, r, &team)
	if err != nil {
		return
	}

	s, err := handler.BackingStore.teamUpdate(id, team)
	if err != nil {
		if err.Error() == model.ErrorDBUpdateNoChange {
			httpErrorNotFound(w)
			return
		}
		httpErrorDatabase(w, err)
		return
	}

	sendResponse(w, s)
}

func (handler APIHandler) createUser(w http.ResponseWriter, r *http.Request) {
	log.Println("create user")

	var user model.User
	err := readRequestBody(w, r, &user)
	if err != nil {
		return
	}

	t, err := handler.BackingStore.userInsert(user)
	if err != nil {
		httpErrorDatabase(w, err)
		return
	}

	sendResponse(w, t)
}

func (handler APIHandler) deleteUser(w http.ResponseWriter, r *http.Request) {
	log.Println("delete user")

	id, err := getRequestID(r)
	if err != nil {
		httpErrorInvalidID(w)
		return
	}

	user, err := handler.BackingStore.userSelect(id)
	if err != nil {
		httpErrorDatabase(w, err)
		return
	}
	if user.ID == 0 {
		httpErrorNotFound(w)
		return
	}

	err = handler.BackingStore.userDelete(id)
	if err != nil {
		httpErrorDatabase(w, err)
	}
}

func (handler APIHandler) readUser(w http.ResponseWriter, r *http.Request) {
	log.Println("read user")

	id, err := getRequestID(r)
	if err != nil {
		httpErrorInvalidID(w)
		return
	}

	s, err := handler.BackingStore.userSelect(id)
	if err != nil {
		httpErrorDatabase(w, err)
		return
	}
	if s.ID == 0 {
		httpErrorNotFound(w)
		return
	}

	sendResponse(w, s)
}

func (handler APIHandler) readUsers(w http.ResponseWriter, r *http.Request) {
	log.Println("read users")

	s, err := handler.BackingStore.userSelectAll()
	if err != nil {
		httpErrorDatabase(w, err)
		return
	}

	sendResponse(w, s)
}

func (handler APIHandler) updateUser(w http.ResponseWriter, r *http.Request) {
	log.Println("update user")

	id, err := getRequestID(r)
	if err != nil {
		httpErrorInvalidID(w)
		return
	}

	var user model.User
	err = readRequestBody(w, r, &user)
	if err != nil {
		return
	}

	s, err := handler.BackingStore.userUpdate(id, user)
	if err != nil {
		if err.Error() == model.ErrorDBUpdateNoChange {
			httpErrorNotFound(w)
			return
		}
		httpErrorDatabase(w, err)
		return
	}

	sendResponse(w, s)
}

func (handler APIHandler) readUserRoles(w http.ResponseWriter, r *http.Request) {
	log.Println("read user roles")

	id, err := getRequestID(r)
	if err != nil {
		httpErrorInvalidID(w)
		return
	}

	s, err := handler.BackingStore.userRolesSelect(id)
	if err != nil {
		httpErrorDatabase(w, err)
		return
	}

	sendResponse(w, s)
}

func (handler APIHandler) updateUserRoles(w http.ResponseWriter, r *http.Request) {
	log.Println("update user roles")

	id, err := getRequestID(r)
	if err != nil {
		httpErrorInvalidID(w)
		return
	}

	var roles []model.Role
	err = readRequestBody(w, r, &roles)
	if err != nil {
		return
	}

	err = handler.BackingStore.userRolesUpdate(id, roles)
	if err != nil {
		if err.Error() == model.ErrorDBUpdateNoChange {
			httpErrorNotFound(w)
			return
		}
		httpErrorDatabase(w, err)
		return
	}
}
