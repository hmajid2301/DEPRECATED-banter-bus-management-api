#!/bin/bash

printf "Generating new openapi specification from the app. \n"
go test ./tests/openapi_test.go -v > /dev/null

if cmp -s "openapi.json" "specs/openapi.json"; then
    printf "The openapi files are the same, deleting new openapi.json. \n"
    rm openapi.json
else
    printf "Replacing current openapi.json file. \n"
    mv openapi.json specs/openapi.json
fi
