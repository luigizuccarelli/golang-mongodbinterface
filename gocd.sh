!/bin/sh

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
    GOCACHE=off go test -v config.go config_test.go schema.go handlers.go middleware.go middleware_test.go handlers_test.go -coverprofile tests/results/cover.out
    go tool cover -html=tests/results/cover.out -o tests/results/cover.html
fi

if [ "$1" = "sonarqube" ]
then
    echo -e "\nSonarqube scanning project"
    sonarqube/bin/sonar-scanner  -Dsonar.projectKey=portfoliotracker-stocks-dbinterface  -Dsonar.sources=.   -Dsonar.host.url=http://sonarqube-service:9009   -Dsonar.login=c24cdce3d9d3a38680f16a0f069ef90421884eeb -Dsonar.go.coverage.reportPaths=tests/results/cover.out -Dsonar.exclusions=vendor/**,*_test.go,main.go,connectors.go,tests/**
fi
