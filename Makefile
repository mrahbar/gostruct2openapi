GOCMD:=$(shell which go)
GOLINT:=$(shell which golint)
GOIMPORT:=$(shell which goimports)
GOFMT:=$(shell which gofmt)
GOBUILD:=$(GOCMD) build
GOCLEAN:=$(GOCMD) clean
GOTEST:=$(GOCMD) test
GOGET:=$(GOCMD) get
GOLIST:=$(GOCMD) list
GOPATH:=$(shell $(GOCMD) env GOPATH)
u := $(if $(update),-u)

BINARY_NAME:=gostruct2openapi
PACKAGES:=$(shell $(GOLIST) github.com/mrahbar/gostruct2openapi)
GOFILES:=$(shell find . -name "*.go" -type f)

export GO111MODULE := on

all: test build

.PHONY: build
build: deps
	$(GOBUILD) -o $(BINARY_NAME) ./cmd/doc

.PHONY: test
test:
	echo "mode: count" > coverage.out
	for PKG in $(PACKAGES); do \
		$(GOCMD) test -v -covermode=count -coverprofile=profile.out $$PKG > tmp.out; \
		cat tmp.out; \
		if grep -q "^--- FAIL" tmp.out; then \
			rm tmp.out; \
			exit 1; \
		elif grep -q "build failed" tmp.out; then \
			rm tmp.out; \
			exit; \
		fi; \
		if [ -f profile.out ]; then \
			cat profile.out | grep -v "mode:" >> coverage.out; \
			rm profile.out; \
		fi; \
	done

.PHONY: clean
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

.PHONY: deps
deps:
	echo "no deps to install"

.PHONY: devel-deps
devel-deps:
	GO111MODULE=off $(GOGET) -v -u golang.org/x/lint/golint

.PHONY: fmt
fmt:
	$(GOFMT) -s -w $(GOFILES)

.PHONY: fmt-check
fmt-check:
	@diff=$$($(GOFMT) -s -d $(GOFILES)); \
	if [ -n "$$diff" ]; then \
		echo "Please run 'make fmt' and commit the result:"; \
		echo "$${diff}"; \
		exit 1; \
	fi;

.PHONY: view-covered
view-covered:
	$(GOTEST) -coverprofile=cover.out $(TARGET)
	$(GOCMD) tool cover -html=cover.out