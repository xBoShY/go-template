MODNAME := $(shell go list -m)
LICENCE := AGPLv3.0
REPO := https://$(MODNAME)

UNAME := $(shell uname)

GOTAGSLIST  := ${GOTAGSCUSTOM}

ifeq ($(UNAME), Linux)
GOTAGSLIST	+= osusergo netgo static_build
GOBUILDMODE := -buildmode pie
EXTLDFLAGS	:= -static-libstdc++ -static-libgcc
# the following predicate is abit misleading; it tests if we're not in centos.
ifeq (,$(wildcard /etc/centos-release))
EXTLDFLAGS  += -static
endif
endif

# If build number already set, use it - to ensure same build number across multiple platforms being built
BUILDNUMBER		?= $(shell echo 9)
COMMITHASH		:= $(shell echo 3)
BRANCH			:= $(shell echo main)

GOLDFLAGS_BASE	:= -X $(MODNAME)/config.BuildNumber=$(BUILDNUMBER)
GOLDFLAGS_BASE	+= -X $(MODNAME)/config.CommitHash=$(COMMITHASH)
GOLDFLAGS_BASE	+= -X $(MODNAME)/config.Branch=$(BRANCH)
GOLDFLAGS_BASE	+= -X $(MODNAME)/config.License=$(LICENCE)
GOLDFLAGS_BASE	+= -X $(MODNAME)/config.Repo=$(REPO)
GOLDFLAGS_BASE	+= -extldflags \"$(EXTLDFLAGS)\"

GOMOD_DIRS := 

GOTRIMPATH	:= $(shell GOPATH=$(GOPATH) && go help build | grep -q .-trimpath && echo -trimpath)
GOTAGS      := --tags "$(GOTAGSLIST)"
GOLDFLAGS 	:= $(GOLDFLAGS_BASE)

REST_DAEMONS := rest-server

default: build

build: clean tidy rest-api
	go build -o ./build/ $(GOTRIMPATH) $(GOTAGS) $(GOBUILDMODE) -ldflags="$(GOLDFLAGS)" ./...

clean:
	go clean -i ./...
	rm -rf ./build

tidy:
	@echo "Tidying"
	go mod tidy
	@for dir in $(GOMOD_DIRS); do \
		echo "Tidying $$dir" && \
		(cd $$dir && go mod tidy); \
	done

rest-api: oapi-codegen
	@for daemon in $(REST_DAEMONS); do \
		echo "Building REST-API for $$daemon" && \
		(cd ./daemon/$$daemon/api && make); \
	done

oapi-codegen:
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
