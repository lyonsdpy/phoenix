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

# # compile build result
# TARGET_NAME=phoenix


# # support OS
# # mac os "darwin"
# # linux "linux"
# # windows "windows"
# GOOS=linux

# # compile platform, default amd64
# GOARCH=amd64

# win: 

# all: pull format clean build compress package release

# prd: pull clean build compress package

# stg: pull clean build compress package release

# cli: TARGET_NAME=sealion-cli
# cli: BUILD_FILE=cli.go
# cli: VERSION=0.0.1

# cli: pull clean build-cli compress package

# build-cli:
# 	$(eval GIT_COMMIT=$(shell git rev-parse --short HEAD))
# 	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags "-X main.CliCommitId=$(GIT_COMMIT) -X 'main.CliBuilt=$(BUILT)' -X main.CliVersion=$(VERSION) -s -w" -v -o $(TARGET_NAME) $(BUILD_FILE)

# pull:
# 	git pull
# 	git submodule foreach git pull

# format:
# 	gofmt -w .

# build:
# 	$(eval GIT_COMMIT=$(shell git rev-parse --short HEAD))
# 	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -i -ldflags "-X main.CommitId=$(GIT_COMMIT) -X 'main.Built=$(BUILT)' -s -w" -v -o $(TARGET_NAME) $(BUILD_FILE)

# clean:
# 	rm -fr ./objs

# package: 
# 	dos2unix ./scripts/*.sh
# 	./scripts/package.sh package $(TARGET_NAME)

# release:
# 	dos2unix ./scripts/*.sh
# 	./scripts/package.sh release $(TARGET_NAME)

# compress:
# 	upx $(TARGET_NAME)
