package productcontroller_test

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
	if _, err := models.DB.Exec(helper.TableProductCreationQuery); err != nil {
		log.Fatal(err)
	}
	if _, err := models.DB.Exec(helper.TableProductIndexingQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	models.DB.Exec("TRUNCATE products;")
	models.DB.Exec("DELETE FROM products;")
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

func TestCreateProductFailAuthorization(t *testing.T) {
	defer clearTable()
	var payload = []byte(`{
		"name": "iniproduk0",
		"price": 0
	}`)

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/product", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	app.R.ServeHTTP(rec, req)

	if rec.Code != 401 {
		t.Errorf("Expected response code to be 401. Got %d", rec.Code)
	}
}

func TestCreateProductFailPayload(t *testing.T) {
	defer clearTable()
	helper.AddUsers(1)

	var payload = []byte(`{
		"name": "iniproduk0",
	}`)

	var accessToken config.TokenPayload
	if err := accessToken.CreateToken("iniusername0", 15); err != nil {
		log.Fatal("can't procced when creating token.")
	}

	var access = fmt.Sprintf("Bearer %s", accessToken.Token)

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/product", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", access)

	app.R.ServeHTTP(rec, req)

	var respPayload map[string]string
	json.Unmarshal(rec.Body.Bytes(), &respPayload)

	if respPayload["error"] != "Invalid request payload" {
		t.Errorf("Expected the 'message' key of resp to be 'Invalid request payload'. Got '%s'", respPayload["error"])
	}
}

func TestCreateProductFailValidation(t *testing.T) {
	defer clearTable()
	helper.AddUsers(1)

	var payload = []byte(`{
		"name": "iniproduk",
		"price": 21
	}`)

	var accessToken config.TokenPayload
	if err := accessToken.CreateToken("iniusername0", 15); err != nil {
		log.Fatal("can't procced when creating token.")
	}

	var access = fmt.Sprintf("Bearer %s", accessToken.Token)
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/product", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", access)

	app.R.ServeHTTP(rec, req)

	var m map[string][]map[string]string
	json.Unmarshal(rec.Body.Bytes(), &m)

	if m["errors"][0]["error"] != "Name value must greater than 10" {
		t.Errorf("Expected the 'error' key of resp to be 'Name value must greater than 10'. Got '%s'", m["errors"][0]["error"])
	}
}

func TestCreateProductSuccess(t *testing.T) {
	defer clearTable()
	helper.AddUsers(1)

	var payload = []byte(`{
		"name": "iniproduk0",
		"price": 21
	}`)

	var accessToken config.TokenPayload
	if err := accessToken.CreateToken("iniusername0", 15); err != nil {
		log.Fatal("can't procced when creating token.")
	}

	var access = fmt.Sprintf("Bearer %s", accessToken.Token)
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/product", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", access)

	app.R.ServeHTTP(rec, req)

	if rec.Code != 201 {
		t.Errorf("Expected the resp code to be 201. Got %d", rec.Code)
	}
}

func TestGetSpesificProductNotFound(t *testing.T) {
	defer clearTable()
	helper.AddUsers(1)
	helper.AddProducts(1)

	var accessToken config.TokenPayload
	if err := accessToken.CreateToken("iniusername0", 15); err != nil {
		log.Fatal("can't procced when creating token.")
	}

	var access = fmt.Sprintf("Bearer %s", accessToken.Token)

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/product/12asda-123asd-gkl75", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", access)

	app.R.ServeHTTP(rec, req)

	if rec.Code != 404 {
		t.Errorf("Expected the resp code to be 404. Got %d", rec.Code)
	}
}

func TestGetSpesificProductSuccess(t *testing.T) {
	defer clearTable()
	helper.AddUsers(1)
	productID := helper.AddProducts(1)

	var accessToken config.TokenPayload
	if err := accessToken.CreateToken("iniusername0", 15); err != nil {
		log.Fatal("can't procced when creating token.")
	}

	var access = fmt.Sprintf("Bearer %s", accessToken.Token)
	var urlPath = fmt.Sprintf("/v1/product/%s", productID)

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", urlPath, nil)
	req.Header.Set("Authorization", access)

	app.R.ServeHTTP(rec, req)

	var m map[string]string
	json.Unmarshal(rec.Body.Bytes(), &m)

	if m["id"] != productID {
		t.Errorf("Expected the resp to be valid. Got %v", m)
	}
}

func TestGetMultipleProduct(t *testing.T) {
	defer clearTable()

	var expectedLength int = 5
	helper.AddUsers(1)
	helper.AddProducts(expectedLength)

	var accessToken config.TokenPayload
	if err := accessToken.CreateToken("iniusername0", 15); err != nil {
		log.Fatal("can't procced when creating token")
	}

	var access = fmt.Sprintf("Bearer %s", accessToken.Token)

	req, _ := http.NewRequest("GET", "/v1/products", nil)
	rec := httptest.NewRecorder()
	req.Header.Set("Authorization", access)

	app.R.ServeHTTP(rec, req)

	var m map[string][]map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &m)

	if len(m["data"]) != expectedLength {
		t.Errorf("Expected resp items to be 5. Got %d", len(m["data"]))
	}
}

func TestUpdateSpesificProduct(t *testing.T) {
	defer clearTable()
	helper.AddUsers(1)
	productID := helper.AddProducts(1)

	var accessToken config.TokenPayload
	if err := accessToken.CreateToken("iniusername0", 15); err != nil {
		log.Fatal("can't procced when creating token")
	}

	var access = fmt.Sprintf("Bearer %s", accessToken.Token)
	var urlPath = fmt.Sprintf("/v1/product/%s", productID)
	var updatedProduct = []byte(`{
		"name": "iniupdatedproduct",
		"price": 100
	}`)

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", urlPath, bytes.NewBuffer(updatedProduct))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", access)

	app.R.ServeHTTP(rec, req)

	var m map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &m)

	if m["name"] == "iniproduk0" {
		t.Errorf("Expected the name to change from 'iniproduk0' to 'iniupdatedproduct'. Got %s", m["name"])
	}

	if m["price"] == 0.0 {
		t.Errorf("Expected the name to change from 0 to 100. Got %v", m["price"])
	}
}

func TestDeleteSpesificProduct(t *testing.T) {
	defer clearTable()
	helper.AddUsers(1)
	productID := helper.AddProducts(1)

	var accessToken config.TokenPayload
	if err := accessToken.CreateToken("iniusername0", 15); err != nil {
		log.Fatal("can't procced when creating token")
	}
	var access = fmt.Sprintf("Bearer %s", accessToken.Token)
	var urlPath = fmt.Sprintf("/v1/product/%s", productID)

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", urlPath, nil)
	req.Header.Set("Authorization", access)

	app.R.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Errorf("Expected the resp code to be 200. Got %d", rec.Code)
	}
}
