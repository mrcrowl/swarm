#!/bin/bash

pushd . >/dev/null 2>&1
cd scripts

echo Getting go packages...
go get
go run publish.go

popd >/dev/null 2>&1

echo Done!