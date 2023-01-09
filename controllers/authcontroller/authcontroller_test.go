package authcontroller_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/bayudha2/go-test-0/app"
	"github.com/bayudha2/go-test-0/config"
	"github.com/bayudha2/go-test-0/helper"
	"github.com/bayudha2/go-test-0/models"
)

type errorVal struct {
	Errors []map[string]string
}

func ensureTableExist() {
	if _, err := models.DB.Exec(helper.TableUserCreationQuery); err != nil {
		log.Fatal(err)
	}
	if _, err := models.DB.Exec(helper.TableSessionCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	models.DB.Exec("TRUNCATE users;")
	models.DB.Exec("DELETE FROM users;")
	models.DB.Exec("TRUNCATE sessions;")
	models.DB.Exec("DELETE FROM sessions;")
}

func TestMain(m *testing.M) {
	models.ConnectDatabase(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_TEST_NAME"))

	app.Initialize()

	ensureTableExist()
	code := m.Run()
	clearTable()

	os.Exit(code)
}

var payload = []byte(`{
	"username": "iniusername0",
	"fullname": "inifullname0",
	"email": "iniemail@0.com",
	"password": "inipassword0"
}`)

func TestRegisterUsernameExist(t *testing.T) {
	defer clearTable()
	helper.AddUsers(1)

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	app.R.ServeHTTP(rec, req)

	var m map[string]string
	json.Unmarshal(rec.Body.Bytes(), &m)

	if m["error"] != "Username already used!" {
		t.Errorf("Expected the 'error' key of resp to be 'Username already used!'. Got '%s'", m["error"])
	}
}

func TestRegisterFail(t *testing.T) {
	defer clearTable()
	var jsonStr = []byte(`{
		"fullname": "inifullname",
		"email": "iniemail@email.com"
	}`)

	var expectedVal errorVal
	expectedVal.Errors = append(expectedVal.Errors, map[string]string{"error": "Username is required"})
	expectedVal.Errors = append(expectedVal.Errors, map[string]string{"error": "Password is required"})

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	var errors errorVal
	app.R.ServeHTTP(rec, req)
	decoder := json.NewDecoder(rec.Body)
	if err := decoder.Decode(&errors); err != nil {
		t.Log("cannot retrieve JSON value")
		return
	}

	if expectedVal.Errors[0]["error"] != errors.Errors[0]["error"] || expectedVal.Errors[1]["error"] != errors.Errors[1]["error"] {
		t.Errorf("Failed => Got %v wanted %v", errors, expectedVal)
	}
}

func TestRegisterSuccess(t *testing.T) {
	defer clearTable()

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	app.R.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Errorf("Expected response code tobe 200. Got %v", rec.Code)
	}
}

func TestLoginFailRequirePassword(t *testing.T) {
	defer clearTable()

	var payloadFailLogin = []byte(`{
		"username": "inifullname",
		"email": "iniemail@email.com"
	}`)

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/signin", bytes.NewBuffer(payloadFailLogin))
	req.Header.Set("Content-Type", "application/json")

	app.R.ServeHTTP(rec, req)
	var m map[string][]map[string]string
	json.Unmarshal(rec.Body.Bytes(), &m)

	if m["errors"][0]["error"] != "Password is required" {
		t.Errorf("Expected the 'error' key of resp to be 'Password is required'. Got '%s'", m["errors"][0]["error"])
	}
}

func TestLoginWrongPassword(t *testing.T) {
	defer clearTable()
	helper.AddUsers(1)

	var payloadLoginFail = []byte(`{
		"username": "iniusername",
		"password": "inipassword321"
	}`)

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/signin", bytes.NewBuffer(payloadLoginFail))
	req.Header.Set("Content-Type", "application/json")

	app.R.ServeHTTP(rec, req)
	var m map[string]string
	json.Unmarshal(rec.Body.Bytes(), &m)

	if m["error"] != "Username or password is incorrect" {
		t.Errorf("Expected the 'error' key of resp to be 'Username or password is incorrect'. Got '%s'", m["error"])
	}
}

func TestLoginSuccess(t *testing.T) {
	defer clearTable()
	helper.AddUsers(1)

	var payloadLogin = []byte(`{
		"username": "iniusername0",
		"password": "inipassword0"
	}`)

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/signin", bytes.NewBuffer(payloadLogin))
	req.Header.Set("Content-Type", "application/json")

	app.R.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Errorf("Expected the resp code to be 200. Got %d", rec.Code)
	}
}

func TestRefreshTokenWithoutBearer(t *testing.T) {
	defer clearTable()

	var refresh = fmt.Sprintf(`{"refresh_token": "asdasdadasd"}`)
	var refreshPayload = []byte(refresh)

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/refresh", bytes.NewBuffer(refreshPayload))
	req.Header.Set("Content-Type", "application/json")
	app.R.ServeHTTP(rec, req)

	var respPayload map[string]string
	json.Unmarshal(rec.Body.Bytes(), &respPayload)

	if respPayload["error"] != "Not Authorized!" {
		t.Errorf("Expected the 'error' key of resp to be 'Not Authorized!'. Got %s", respPayload["error"])
	}
}

func TestRefreshFailRandomToken(t *testing.T) {
	defer clearTable()
	var accessToken config.TokenPayload
	if err := accessToken.CreateToken("iniusername0", 15); err != nil {
		t.Errorf("can't procced when creating token.")
	}

	var refresh = fmt.Sprintf(`{"refresh_token": "asdasdadasd"}`)
	var access = fmt.Sprintf("Bearer %s", accessToken.Token)

	var refreshPayload = []byte(refresh)
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/refresh", bytes.NewBuffer(refreshPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", access)
	app.R.ServeHTTP(rec, req)

	var respPayload map[string]string
	json.Unmarshal(rec.Body.Bytes(), &respPayload)

	if respPayload["error"] != "token contains an invalid number of segments" {
		t.Errorf("Expected the 'error' key of resp to be 'token contains an invalid number of segments'. Got %s", respPayload["error"])
	}
}

func TestRefreshSuccess(t *testing.T) {
	defer clearTable()

	var accessToken config.TokenPayload
	if err := accessToken.CreateToken("iniusername0", 15); err != nil {
		t.Errorf("can't procced when creating token.")
	}
	var refreshToken config.TokenPayload
	if err := refreshToken.CreateToken("iniusername0", 30); err != nil {
		t.Errorf("can't procced when creating token.")
	}

	helper.AddUsers(1)
	helper.AddSession(refreshToken.Token, int(refreshToken.ExpTime.Unix()))

	var access = fmt.Sprintf("Bearer %s", accessToken.Token)
	var refresh = fmt.Sprintf(`{"refresh_token": "%s"}`, refreshToken.Token)

	var refreshPayload = []byte(refresh)
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/refresh", bytes.NewBuffer(refreshPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("authorization", access)
	app.R.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Errorf("Expected the resp code to be 200. Got %d", rec.Code)
	}
}

func TestLogoutSuccess(t *testing.T) {
	defer clearTable()

	var accessToken config.TokenPayload
	if err := accessToken.CreateToken("iniusername0", 15); err != nil {
		t.Errorf("can't procced when creating token.")
	}
	var refreshToken config.TokenPayload
	if err := refreshToken.CreateToken("iniusername0", 30); err != nil {
		t.Errorf("can't procced when creating token.")
	}

	helper.AddUsers(1)
	helper.AddSession(refreshToken.Token, int(refreshToken.ExpTime.Unix()))
	var access = fmt.Sprintf("Bearer %s", accessToken.Token)

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/signout", nil)
	req.Header.Set("authorization", access)
	app.R.ServeHTTP(rec, req)

	var respPayload map[string]string
	json.Unmarshal(rec.Body.Bytes(), &respPayload)

	if respPayload["message"] != "Logout successfully" {
		t.Errorf("Expected the 'message' key of resp to be 'Logout Successfully'. Got '%s'", respPayload["message"])
	}
}
