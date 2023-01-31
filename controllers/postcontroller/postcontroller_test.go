package postcontroller_test

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

func ensureTableExist() {
	if _, err := models.DB.Exec(helper.TableUserCreationQuery); err != nil {
		log.Fatal(err.Error())
	}

	if _, err := models.DB.Exec(helper.TablePostCreationQuery); err != nil {
		log.Fatal(err.Error())
	}
}

func clearTable() {
	models.DB.Exec("TRUNCATE users;")
	models.DB.Exec("DELETE FROM users;")
	models.DB.Exec("TRUNCATE posts;")
	models.DB.Exec("DELETE FROM posts;")
}

func TestMain(m *testing.M) {
	models.ConnectDatabase(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_TEST_NAME"),
	)

	app.Initialize()

	ensureTableExist()
	code := m.Run()
	clearTable()

	os.Exit(code)
}

func TestCreatePostFail(t *testing.T) {
	defer clearTable()
	helper.AddUsers(1)

	var accessToken config.TokenPayload
	if err := accessToken.CreateToken("iniuserid0", "iniusername0", 15); err != nil {
		log.Fatal("can't procced when creating token.")
	}

	var access = fmt.Sprintf("Bearer %s", accessToken.Token)
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/post", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", access)

	app.R.ServeHTTP(rec, req)

	var err map[string]string
	json.Unmarshal(rec.Body.Bytes(), &err)

	if err["error"] != "Invalid Request Payload" {
		t.Errorf("Expected the 'error' key of resp to be 'Invalid Request payload'. Got %s", err["error"])
	}
}

func TestCreatePostSuccess(t *testing.T) {
	defer clearTable()
	helper.AddUsers(1)

	var accessToken config.TokenPayload
	if err := accessToken.CreateToken("iniuserid0", "iniusername0", 15); err != nil {
		log.Fatal("can't procced when creating token.")
	}

	var jsonPayload = []byte(`{
		"description": "ini postingan pertamaku :)"
	}`)

	var access = fmt.Sprintf("Bearer %s", accessToken.Token)

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/post", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", access)

	app.R.ServeHTTP(rec, req)

	if rec.Code != 201 {
		t.Errorf("Expected the resp code to be 200. Got %d", rec.Code)
	}
}

func TestGetAllPost(t *testing.T) {
	defer clearTable()
	helper.AddUsers(1)
	helper.AddPost(10, "iniuserid0")

	var expectedLengthData = 10

	var accessToken config.TokenPayload
	if err := accessToken.CreateToken("iniuserid0", "iniusername0", 15); err != nil {
		log.Fatal("can't procced when creating token.")
	}

	var access = fmt.Sprintf("Bearer %s", accessToken.Token)
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/posts", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", access)

	app.R.ServeHTTP(rec, req)

	var m models.PayloadPosts
	json.Unmarshal(rec.Body.Bytes(), &m)

	if m.TotalData != expectedLengthData {
		t.Errorf("expected total resp data to be %d. Got %d", expectedLengthData, m.TotalData)
	}
}

func TestGetSpesificPost(t *testing.T) {
	defer clearTable()
	helper.AddUsers(1)
	helper.AddPost(4, "iniuserid0")

	var accessToken config.TokenPayload
	if err := accessToken.CreateToken("iniuserid0", "iniusername0", 15); err != nil {
		log.Fatal("can't procced when creating token.")
	}

	var access = fmt.Sprintf("Bearer %s", accessToken.Token)
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/post/inipostid3", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", access)

	app.R.ServeHTTP(rec, req)

	var m models.Post
	json.Unmarshal(rec.Body.Bytes(), &m)

	if m.ID != "inipostid3" {
		t.Errorf("Expected the resp to be valid. Got %v", m)
	}
}

func TestGetSpesificPostNotFound(t *testing.T) {
	defer clearTable()
	helper.AddUsers(1)
	helper.AddPost(1, "iniuserid0")

	var accessToken config.TokenPayload
	if err := accessToken.CreateToken("iniuserid0", "iniusername0", 15); err != nil {
		log.Fatal("can't procced when creating token.")
	}

	var access = fmt.Sprintf("Bearer %s", accessToken.Token)
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/post/inipostid12", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", access)

	app.R.ServeHTTP(rec, req)

	var m map[string]string
	json.Unmarshal(rec.Body.Bytes(), &m)

	if m["error"] != "Post not found" {
		t.Errorf("Expected the resp value to be 'Post not found'. Got %s", m["error"])
	}
}

func TestUpdatePostFail(t *testing.T) {
	defer clearTable()
	helper.AddUsers(1)
	helper.AddPost(1, "iniuserid0")

	var accessToken config.TokenPayload
	if err := accessToken.CreateToken("iniuserid0", "iniusername0", 15); err != nil {
		log.Fatal("can't procced when creating token")
	}

	expectedDesc := "halo aku adalah post terbaru"
	payloadByte := fmt.Sprintf(`{"description": "%s"}`, expectedDesc)
	updatedPayload := []byte(payloadByte)

	var access = fmt.Sprintf("Bearer %s", accessToken.Token)
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/v1/post/inipostid33", bytes.NewBuffer(updatedPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", access)

	app.R.ServeHTTP(rec, req)

	var m map[string]string
	json.Unmarshal(rec.Body.Bytes(), &m)

	if m["error"] != "Unauthorized request!" {
		t.Errorf("Expected the key resp error value to be 'Unauthorized request!'. Got %s ", m["error"])
	}
}

func TestUpdatePostSuccess(t *testing.T) {
	defer clearTable()
	helper.AddUsers(1)
	helper.AddPost(3, "iniuserid0")

	var accessToken config.TokenPayload
	if err := accessToken.CreateToken("iniuserid0", "iniusername0", 15); err != nil {
		log.Fatal("can't procced when creating token")
	}

	expectedDesc := "halo aku adalah post terbaru"
	payloadByte := fmt.Sprintf(`{"description": "%s"}`, expectedDesc)
	updatedPayload := []byte(payloadByte)

	var access = fmt.Sprintf("Bearer %s", accessToken.Token)
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/v1/post/inipostid2", bytes.NewBuffer(updatedPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", access)

	app.R.ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Errorf("Expected the resp code to be 200. Got %d", rec.Code)
	}
}

func TestDeletePostSuccess(t *testing.T) {
	defer clearTable()
	helper.AddUsers(1)
	helper.AddPost(3, "iniuserid0")

	var accessToken config.TokenPayload
	if err := accessToken.CreateToken("iniuserid0", "iniusername0", 15); err != nil {
		log.Fatal("can't procced when creating token")
	}

	var access = fmt.Sprintf("Bearer %s", accessToken.Token)
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/v1/post/inipostid1", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", access)

	app.R.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Errorf("Expected the resp code to be 200. Got %d", rec.Code)
	}
}
