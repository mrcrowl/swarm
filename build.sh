#!/bin/bash

echo Getting go packages...
go get

echo ""
echo Building TypeScript files...
tsc -p ./assets/static/tsconfig.json

echo ... and embedding as resources...
go get github.com/omeid/go-resources/cmd/resources
go run github.com/omeid/go-resources/cmd/resources -output assets/assets_prod.go -declare -var FS -package assets -tag=embed assets/static/*.js

echo ""
echo Building 'swarm' for macOS...
go build -o swarm -ldflags "-s -w" -tags=embed

echo ""
echo Done!