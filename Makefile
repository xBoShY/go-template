MODNAME := $(shell go list -m)
LICENCE := AGPLv3.0
REPO := https://$(MODNAME)

UNAME := $(shell uname)

GOTAGSLIST  := ${GOTAGSCUSTOM}

ifeq ($(UNAME), Linux)
GOTAGSLIST  += osusergo netgo static_build
GOBUILDMODE := -buildmode pie
EXTLDFLAGS  := -static-libstdc++ -static-libgcc
# the following predicate is abit misleading; it tests if we're not in centos.
ifeq (,$(wildcard /etc/centos-release))
EXTLDFLAGS  += -static
endif
endif

# If version and build number are already set, use them
# - to ensure same build number across multiple platforms being built
V_MAJOR     ?= 0
V_MINOR     ?= 0
BUILDNUMBER ?= $(shell echo 9)
COMMITHASH  := $(shell echo 8dd03144ffdc6c0d486d6b705f9c7fba871ee7c3)
BRANCH      := $(shell echo main)

GOLDFLAGS_BASE  := -w -s
GOLDFLAGS_BASE  += -X $(MODNAME)/config.Major=$(V_MAJOR)
GOLDFLAGS_BASE  += -X $(MODNAME)/config.Minor=$(V_MINOR)
GOLDFLAGS_BASE  += -X $(MODNAME)/config.BuildNumber=$(BUILDNUMBER)
GOLDFLAGS_BASE  += -X $(MODNAME)/config.CommitHash=$(COMMITHASH)
GOLDFLAGS_BASE  += -X $(MODNAME)/config.Branch=$(BRANCH)
GOLDFLAGS_BASE  += -X $(MODNAME)/config.License=$(LICENCE)
GOLDFLAGS_BASE  += -X $(MODNAME)/config.Repo=$(REPO)
GOLDFLAGS_BASE  += -extldflags \"$(EXTLDFLAGS)\"

GOMOD_DIRS  := 
GOTRIMPATH  := $(shell GOPATH=$(GOPATH) && go help build | grep -q .-trimpath && echo -trimpath)
GOTAGS      := --tags "$(GOTAGSLIST)"
GOLDFLAGS   := $(GOLDFLAGS_BASE)

API_DIRS := stringd

docker_entrypoint = $(error Please set docker_entrypoint flag)


default: build

build: clean tidy api
	@go build -o ./build/ $(GOTRIMPATH) $(GOTAGS) $(GOBUILDMODE) -ldflags="$(GOLDFLAGS)" ./...

clean:
	@go clean -i ./...
	@rm -rf ./build

docker: clean
	@docker build --no-cache --build-arg DOCKER_ENTRYPOINT=$(docker_entrypoint) -t "$(docker_entrypoint):$(V_MAJOR).$(V_MINOR).$(BUILDNUMBER)" .

tidy:
	@echo "Tidying"
	@go mod tidy
	@for dir in $(GOMOD_DIRS); do \
    echo "Tidying $$dir" && \
    (cd $$dir && go mod tidy); \
  done

api: oapi-codegen
	@for api in $(API_DIRS); do \
    echo "Building API for $$api" && \
    (cd ./daemon/$$api/api && [ -f ./Makefile ] && make); \
  done

oapi-codegen:
	@go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
