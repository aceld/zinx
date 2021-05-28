
.PHONY: build

SERVICE := zinx
CUR_PWD := $(shell pwd)

SERVER_DEMO_PATH := $(CUR_PWD)/examples/zinx_server
CLIENT_DEMO_PATH := $(CUR_PWD)/examples/zinx_client
SERVER_DEMO_BIN := $(SERVER_DEMO_PATH)/server
CLIENT_DEMO_BIN := $(CLIENT_DEMO_PATH)/client

AUTHOR := $(shell git log --pretty=format:"%an"|head -n 1)
VERSION := $(shell git rev-list HEAD | head -1)
BUILD_INFO := $(shell git log --pretty=format:"%s" | head -1)
BUILD_DATE := $(shell date +%Y-%m-%d\ %H:%M:%S)

export GO111MODULE=on

LD_FLAGS='-X "$(SERVICE)/version.TAG=$(TAG)" -X "$(SERVICE)/version.VERSION=$(VERSION)" -X "$(SERVICE)/version.AUTHOR=$(AUTHOR)" -X "$(SERVICE)/version.BUILD_INFO=$(BUILD_INFO)" -X "$(SERVICE)/version.BUILD_DATE=$(BUILD_DATE)"'

default: build

build:
	go build  -ldflags $(LD_FLAGS) -gcflags "-N"  -o $(SERVER_DEMO_BIN) $(SERVER_DEMO_PATH)/main.go
	go build  -ldflags $(LD_FLAGS) -gcflags "-N"  -o $(CLIENT_DEMO_BIN) $(CLIENT_DEMO_PATH)/main.go
clean:
	rm $(SERVER_DEMO_BIN)
	rm $(CLIENT_DEMO_BIN)
