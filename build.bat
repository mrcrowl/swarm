@echo off


echo Getting go packages...
go get

echo.
echo Building TypeScript files...
call tsc -p ./assets/static/tsconfig.json

echo ... and embedding as resources...
resources -output assets/assets_prod.go -declare -var FS -package assets -tag=embed assets/static/*.js

echo.
echo Building 'swarm' for macOS...
set GOARCH=amd64
set GOOS=darwin
go build -o swarm -ldflags "-s -w-x" -tags=embed

echo.
echo Building 'swarm.exe' for Windows...
set GOOS=windows
go build -o swarm.exe -ldflags "-s -w" -tags=embed

echo.
echo Done!