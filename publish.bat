@echo off

pushd .
cd ..\..\..\swarm-tools\scripts

echo Getting go packages...
go get
go run publish.go

popd
echo.
echo Done!