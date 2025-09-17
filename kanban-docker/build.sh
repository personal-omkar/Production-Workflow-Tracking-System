#!/bin/bash

echo -------------VERSION-BUILD-INCREASE------------------
echo `awk -F',' '{printf("%d",$1+1)}' ./version-build` > ./version-build

VERSION=$(cat ./version)
BUILD=$(cat ./version-build)

echo Version $VERSION
echo Build $BUILD

Flag="-X main.Version=$VERSION -X main.Build=$BUILD"

echo --------------------COMMONS------------------
cd "../kanban-commons"
go mod tidy
cd - > /dev/null

echo --------------------Services------------------
cd "../kanban-commons/services"

CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "$Flag -X main.CFG=$CFG" -o "../../kanban-docker/scripts/service" "bash-service.go"

rm -rf /RUBBER/scripts/*
cp -r ../../kanban-docker/scripts/* /RUBBER/scripts/
cd - > /dev/null

goProjects=("../kanban-dao" "../kanban-rest" "../kanban-web")

for i in "${goProjects[@]}"
do
   echo --------------------$i------------------ | tr [a-z] [A-Z]
   cd "$i"
   
   go clean
   go mod tidy
   CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "$Flag" -buildvcs=false -o ${i##*/} .
   cd - > /dev/null

done