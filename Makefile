GOBIN=$(GOPATH)/bin
APP_DIR_LIST=$(shell go list ./... | grep -v /vendor/)

build: verify_gopath
	CGO_ENABLED=0 go install -tags netgo $(APP_DIR_LIST)
	go fmt $(APP_DIR_LIST)
	cp $(GOPATH)/bin/tap-cli ./application/tap

verify_gopath:
	@if [ -z "$(GOPATH)" ] || [ "$(GOPATH)" = "" ]; then\
		echo "GOPATH not set. You need to set GOPATH before run this command";\
		exit 1 ;\
	fi

deps_fetch_specific: bin/govendor
	@if [ "$(DEP_URL)" = "" ]; then\
		echo "DEP_URL not set. Run this comand as follow:";\
		echo " make deps_fetch_specific DEP_URL=github.com/nu7hatch/gouuid";\
	exit 1 ;\
	fi
	@echo "Fetching specific dependency in newest versions"
	$(GOBIN)/govendor fetch -v $(DEP_URL)

deps_update_tap: verify_gopath
	$(GOBIN)/govendor update github.com/trustedanalytics/...
	rm -Rf vendor/github.com/trustedanalytics/tap-cli
	@echo "Done"

bin/govendor: verify_gopath
	go get -v -u github.com/kardianos/govendor

tests: verify_gopath
	go test --cover $(APP_DIR_LIST)

prepare_dirs:
	mkdir -p ./temp/src/github.com/trustedanalytics/tap-cli
	$(eval REPOFILES=$(shell pwd)/*)
	ln -sf $(REPOFILES) temp/src/github.com/trustedanalytics/tap-cli

build_anywhere:
	$(MAKE) prepare_dirs build_anywhere_linux
	$(MAKE) prepare_dirs build_anywhere_win32
	$(MAKE) prepare_dirs build_anywhere_osx

build_anywhere_linux: prepare_dirs
	$(eval GOPATH=$(shell cd ./temp; pwd))
	$(eval APP_DIR_LIST=$(shell GOPATH=$(GOPATH) go list ./temp/src/github.com/trustedanalytics/tap-cli/... | grep -v /vendor/))
	GOPATH=$(GOPATH) CGO_ENABLED=0 go build -tags netgo $(APP_DIR_LIST)
	mkdir -p application && rm -f application/tap-cli-linux-amd64.elf
	cp ./tap-cli ./application/tap-linux-amd64.elf
	cp ./tap-cli ./application/tap
	rm -Rf ./temp

build_anywhere_win32: prepare_dirs
	$(eval GOPATH=$(shell cd ./temp; pwd))
	$(eval APP_DIR_LIST=$(shell GOPATH=$(GOPATH) go list ./temp/src/github.com/trustedanalytics/tap-cli/... | grep -v /vendor/))
	mkdir -p application
	GOPATH=$(GOPATH) CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o ./application/tap.exe -tags netgo $(APP_DIR_LIST)
	rm -Rf ./temp

build_anywhere_osx: prepare_dirs
	$(eval GOPATH=$(shell cd ./temp; pwd))
	$(eval APP_DIR_LIST=$(shell GOPATH=$(GOPATH) go list ./temp/src/github.com/trustedanalytics/tap-cli/... | grep -v /vendor/))
	mkdir -p application
	GOPATH=$(GOPATH) CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ./application/tap.osx -tags netgo $(APP_DIR_LIST)
	rm -Rf ./temp
