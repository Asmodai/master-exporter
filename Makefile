PACKAGE = github.com/Asmodai/master-exporter
VERSION ?= 0.0.0

CMD_DIR  = $(PACKAGE)/cmd
ROOT_DIR = $(PWD)

APP = master-exporter
CFG = exporters.conf

.phony: configs clean

all: deps build

.PHONY: configs
configs:
	@echo Copying configuration files
	@mkdir -p $(ROOT_DIR)/etc
	@cp configs/* $(ROOT_DIR)/etc

deps:
	@echo Getting dependencies
	@go mod vendor

tidy:
	@echo Tidying dependencies
	@go mod tidy

listdeps:
	@echo Listing dependencies
	@go list -m all

build: deps
	@echo Building $(APP)
	$(eval GIT_VERSION = $(shell scripts/tag2semver.sh $(VERSION) 2>/dev/null))
	go build                                \
		-tags=go_json                   \
		-o $(ROOT_DIR)/bin/$(APP)       \
		-ldflags "-s -w $(GIT_VERSION)" \
		$(CMD_DIR)/$(APP)

test: deps
	@echo Running tests
	@go test --tags=testing $$(go list ./...) -coverprofile=tests.out
	@go tool cover -html=tests.out -o coverage.html

run: deps
	@echo Running $(APP)
	@go run $(CMD_DIR)/$(APP) \
		-config `pwd`/config/$(CFG)

debug: deps
	@echo Running $(APP) in debug mode
	@go run                              \
		-tags=go_json                \
		$(CMD_DIR)/$(APP)            \
		-debug                       \
		-config `pwd`/config/$(CFG)

.PHONY: clean
clean:
	@echo Cleaning $(APP)
	@-rm $(ROOT_DIR)/bin/* 2>/dev/null
	@-rm tests.out 2>/dev/null
	@-rm coverage.html 2>/dev/null
