SHELL = /bin/bash
GO-VER = go1.17

default: build

# #### GO Binary Management ####

deps-go-binary:
	echo "Expect: $(GO-VER)" && \
		echo "Actual: $$(go version)" && \
	 	go version | grep $(GO-VER) > /dev/null

HAS_GO_IMPORTS := $(shell command -v goimports;)

deps-goimports: deps-go-binary
ifndef HAS_GO_IMPORTS
	go get -u golang.org/x/tools/cmd/goimports
endif


# #### CLEAN ####

clean: deps-go-binary 
	rm -rf build/*
	go clean --modcache

# #### DEPS ####

deps-modules: deps-goimports deps-go-binary
	go mod download

deps-counterfeiter: deps-modules
	command -v counterfeiter >/dev/null 2>&1 || go get -u github.com/maxbrunsfeld/counterfeiter/v6

deps-ginkgo: deps-go-binary
	command -v ginkgo >/dev/null 2>&1 || go get -u github.com/onsi/ginkgo/ginkgo github.com/onsi/gomega

deps: deps-modules deps-counterfeiter deps-ginkgo


# #### BUILD ####

SRC = $(shell find . -name "*.go" | grep -v "_test\." )
VERSION := $(or $(VERSION), "dev")
LDFLAGS="-X github.com/cf-platform-eng/mrlog/version.Version=$(VERSION)"

build/mrlog: $(SRC)
	go build -o build/mrlog -ldflags ${LDFLAGS} ./cmd/mrlog/main.go

build: deps build/mrlog

build/mrlog-linux: $(SRC)
	GOARCH=amd64 GOOS=linux go build -o build/mrlog-linux -ldflags ${LDFLAGS} ./cmd/mrlog/main.go

build-linux: deps build/mrlog-linux

build/mrlog-darwin: $(SRC) $(GENERATE_ARTIFACTS)
	GOARCH=amd64 GOOS=darwin go build -o build/mrlog-darwin -ldflags ${LDFLAGS} ./cmd/mrlog/main.go

build-darwin: deps build/mrlog-darwin

build-all: build-linux build-darwin

build-image: build/mrlog-linux
	docker build --tag cfplatformeng/mrlog:${VERSION} --file Dockerfile .

# #### TESTS ####

units: deps 
	ginkgo -r -skipPackage features .

features: deps $(GENERATE_ARTIFACTS)
	ginkgo -r -tags=feature features

test: deps lint units features

lint: deps-goimports
	git ls-files | grep '.go$$' | xargs goimports -l -w
