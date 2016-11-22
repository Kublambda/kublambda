SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

all: kublambda-controller kublambda-runner
	docker build -t forjared/kublambda-controller -f ./Dockerfile.controller .
	docker push forjared/kublambda-controller
	docker build -t forjared/kublambda-runner -f ./Dockerfile.runner .
	docker push forjared/kublambda-runner

kublambda-controller: $(SOURCES) 
	env GOOS=linux GOARCH=amd64 go build ./cmd/kublambda-controller

kublambda-runner: $(SOURCES)
	env GOOS=linux GOARCH=amd64 go build ./cmd/kublambda-runner
