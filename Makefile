.PHONY: all

VERSION=1.0.0
CommitId=$(shell git rev-parse --short HEAD)

# compile file
BUILD_FILE=main.go

BUILD_TIME=`date`

win:TARGET_NAME=phoenix.exe
win:GOOS=windows
win:GOARCH=amd64
win: build-cli

linux:TARGET_NAME=phoenix
linux:GOOS=linux
linux:GOARCH=amd64
linux: build-cli

mac:TARGET_NAME=phoenix
mac:GOOS=darwin
mac:GOARCH=amd64
mac: build-cli

build-cli:
	$(eval GIT_COMMIT=$(shell git rev-parse --short HEAD))
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags "-X main.CommitID=$(GIT_COMMIT) -X 'main.BuildTime=$(BUILD_TIME)' -X main.Version=$(VERSION) -s -w" -v -o $(TARGET_NAME) $(BUILD_FILE)