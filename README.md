# Porfoliotracker database access interface golang microservice

A simple golang microservice with minimal json config. 


## Usage 

```bash
# cd to project directory and build executable
$ go build -o microservice .

```

## Docker build

```bash
docker build -t <your-registry-id>/myportfolio-dbinterface:1.13.1 .

```

## Curl timing usage
```
curl -w "@curl-timing.txt" -o /dev/null -s "http://site-to-test

```

## Executing tests
```bash
# clear the cache - this is optional
go clean -testcache
go test -v schema.go handlers.go middleware.go middleware_test.go handlers_test.go -coverprofile tests/results/cover.out
go tool cover -html=tests/results/cover.out -o tests/results/cover.html
# run sonarqube scanner (assuming sonarqube server is running)
# NB the SonarQube host and login will differ - please update it accordingly 
 ~/Programs/sonar-scanner-3.3.0.1492-linux/bin/sonar-scanner  -Dsonar.projectKey=portfoliotracker-stocks-dbinterface  -Dsonar.sources=.   -Dsonar.host.url=http://localhost:9009   -Dsonar.login=3b172e408d048820bc6a633b1c3f0097523e89f4 -Dsonar.go.coverage.reportPaths=tests/results/cover.out -Dsonar.exclusions=vendor/**,*_test.go,main.go,connectors.go,tests/**

```
## Testing container 
```bash

# start the container
# curl the isalive endpoint
curl -k -H 'Token: xxxxx' -w '@curl-timing.txt'  http://127.0.0.1:9000/api/v1/sys/info/isalive

```
