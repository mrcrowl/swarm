@echo off
go get
set GOARCH=amd64
set GOOS=darwin
echo Building Mac...
go build -o swarm -ldflags "-s -w-x"
set GOOS=windows
echo Building Windows...
go build -o swarm.exe -ldflags "-s -w"
