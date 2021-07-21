package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"gitlab.com/banter-bus/banter-bus-management-api/internal/api"
	"gitlab.com/banter-bus/banter-bus-management-api/internal/core"
)

func TestGetOpenAPI(_ *testing.T) {
	os.Setenv("BANTER_BUS_CONFIG_PATH", "../tests/config.test.yml")
	config, err := core.NewConfig()
	if err != nil {
		fmt.Printf("unable to load config %v", err)
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

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("unexpected error while serving HTTP: %s", err)
		}
	}()
	req, _ := http.NewRequest("GET", "/openapi", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	err = ioutil.WriteFile("openapi.yaml", w.Body.Bytes(), 0600)
	if err != nil {
		fmt.Printf("Failed to save file %v", err)
	}
}
