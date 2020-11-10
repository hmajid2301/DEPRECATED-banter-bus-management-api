package controllers_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"banter-bus-server/src/server"
)

func TestGetOpenAPI(t *testing.T) {
	router, _ := server.NewRouter()
	req, _ := http.NewRequest("GET", "/openapi.json", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	var out bytes.Buffer

	err := json.Indent(&out, w.Body.Bytes(), "", "  ")
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("../openapi.json", out.Bytes(), 0600)
	if err != nil {
		panic(err)
	}
}
