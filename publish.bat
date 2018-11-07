@echo off

pushd
cd scripts

echo Getting go packages...
go get
go run publish.go

popd
echo.
echo Done!