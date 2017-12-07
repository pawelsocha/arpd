REPO_URI ?= github.com/pawelsocha/
REPO_PATH ?= $(REPO_URI)/arpd

prepare:
	@echo "Preapre GOPATH"
	test -h gopath/src/$(REPO_PATH) || \
		( mkdir -p gopath/src/$(REPO_URI); \
		ln -s ../../../.. gopath/src/$(REPO_PATH) )

build: prepare
	@echo "Building arpd for $(GOOS)/$(GOARCH) $(GOPATH)"
	cd gopath/src/${REPO_PATH}; \
	go build -o arpd

linux-amd64:
	export GOOS="linux"; \
	export GOARCH="amd64"; \
	export GOPATH="$(PWD)/gopath"; \
	$(MAKE) build 

darwin-amd64:
	export GOOS="darwin"; \
	export GOARCH="amd64"; \
	export GOPATH="$(PWD)/gopath"; \
	$(MAKE) build 

clean:
	rm -rf ${REPO_PATH}/arpd

test: prepare
	export GOPATH="$(PWD)/gopath"; \
	go test

all: linux-amd64