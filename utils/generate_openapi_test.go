package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/api"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/core"
)

func TestGetOpenAPI() {
	os.Setenv("BANTER_BUS_CONFIG_PATH", "../tests/config.test.yml")
	config, err := core.NewConfig()
	if err != nil {
		fmt.Println("unable to load config")
	}

	env := &api.Env{Logger: nil, Conf: config, DB: nil}
	router, err := api.Setup(env)
	if err != nil {
		fmt.Println("unable to setup webserver")
	}
	srv := http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.Srv.Host, config.Srv.Port),
		Handler: router,
	}
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Printf("unexpected error while serving HTTP: %s", err)
	}
	req, _ := http.NewRequest("GET", "/openapi.json", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	var out bytes.Buffer

	err = json.Indent(&out, w.Body.Bytes(), "", "  ")
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("openapi.json", out.Bytes(), 0600)
	if err != nil {
		panic(err)
	}
}
