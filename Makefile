# Copyright (c) 2016 Intel Corporation
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
GOBIN=$(GOPATH)/bin
APP_DIR_LIST=$(shell go list ./... | grep -v /vendor/)

build: verify_gopath
	go fmt $(APP_DIR_LIST)
	CGO_ENABLED=0 go install -tags netgo $(APP_DIR_LIST)
	mkdir -p application && cp -f $(GOBIN)/tap-cli ./application/tap

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
	$(GOBIN)/govendor remove github.com/trustedanalytics/tap-cli/...
	@echo "Done"

bin/govendor: verify_gopath
	go get -v -u github.com/kardianos/govendor

install_mockgen:
	scripts/install_mockgen.sh

mock_update:
	$(GOPATH)/bin/mockgen -source=vendor/github.com/trustedanalytics/tap-api-service/client/client.go -package=api -destination=api/api_service_client_mock.go
	$(GOPATH)/bin/mockgen -source=vendor/github.com/trustedanalytics/tap-api-service/client/login.go -package=api -destination=api/api_service_login_mock.go
	./add_license.sh

test: verify_gopath
	go test --cover $(APP_DIR_LIST)

prepare_temp:
	mkdir -p temp/src/github.com/trustedanalytics/tap-cli
	$(eval REPOFILES=$(shell pwd)/*)
	ln -sf $(REPOFILES) temp/src/github.com/trustedanalytics/tap-cli
	mkdir -p application

clear_temp:
	rm -Rf ./temp

build_anywhere:
	$(MAKE) build_anywhere_linux
	$(MAKE) build_anywhere_win32
	$(MAKE) build_anywhere_osx

build_anywhere_linux:
	$(MAKE) build_anywhere_linux64
	$(MAKE) build_anywhere_linux32

build_anywhere_linux64:
	$(call build,linux,amd64,tap-linux64)
	ln -sf tap-linux64 application/tap

build_anywhere_linux32:
	$(call build,linux,386,tap-linux32)

build_anywhere_win32:
	$(call build,windows,386,tap-windows32.exe)

build_anywhere_osx:
	$(call build,darwin,amd64,tap-macosx64.osx)

define build
	$(MAKE) prepare_temp
	$(eval GOPATH=$(shell readlink -f temp))
	$(eval APP_DIR_LIST=$(shell GOPATH=$(GOPATH) go list ./temp/src/github.com/trustedanalytics/tap-cli/... | grep -v /vendor/))
	GOPATH=$(GOPATH) CGO_ENABLED=0 GOOS=$(1) GOARCH=$(2) go build -o ./application/$(3) -tags netgo $(APP_DIR_LIST)
	$(MAKE) clear_temp
endef
