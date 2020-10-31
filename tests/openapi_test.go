package controllers

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
	var a = []byte(w.Body.String())
	var out bytes.Buffer
	json.Indent(&out, a, "", "  ")
	ioutil.WriteFile("../openapi.json", out.Bytes(), 0644)
}
