#!/bin/sh

if [ "$1" = "compile" ]
then
    echo -e "\nExecuting golang compile"
    go version
    go get github.com/gorilla/mux
    go get github.com/microlib/simple
    go get github.com/go-redis/redis
    go get github.com/imdario/mergo
    go get gopkg.in/mgo.v2
    go build -o bin/microservice .

fi

if [ "$1" = "test" ]
then
    echo -e "\nExecuting golang unit tests"
    go test -v config.go config_test.go schema.go handlers.go middleware.go middleware_test.go handlers_test.go -coverprofile tests/results/cover.out
    go tool cover -html=tests/results/cover.out -o tests/results/cover.html
fi

if [ "$1" = "sonarqube" ]
then
    echo -e "\nSonarqube scanning project"
    /sonarqube/bin/sonar-scanner  -Dsonar.projectKey=portfoliotracker-stocks-dbinterface  -Dsonar.sources=.   -Dsonar.host.url=http://sonarqube-service:9000  -Dsonar.login=$2 -Dsonar.go.coverage.reportPaths=tests/results/cover.out -Dsonar.exclusions=vendor/**,*_test.go,main.go,connectors.go,tests/**
fi

if [ "$1" = "image" ]
then
    echo -e "\nBuilding container image"
    docker build -f Dockerfile -t $2/$3:$3 .
fi

if [ "$1" = "push" ]
then
    echo -e "\nPushing image to registry"
    docker push $2/$3:$4
fi

if [ "$1" = "deploy" ]
then
    echo -e "\nDeploying to openshift"
    oc login --server=$2 --username=$3 --password=$4 --insecure-skip-tls-verify
    oc project portfoliotracker
    oc delete pod $5
    echo -e "\nPod is re-deploying use oc status or oc get pods -w" 
fi
