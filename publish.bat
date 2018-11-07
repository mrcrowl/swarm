@echo off

cd scripts

echo Getting go packages...
go get
go run publish.go

cd..
echo.
echo Done!