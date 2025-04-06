PROJECT_PATH=$(shell pwd)

BUILD_NUM_FILE=build_num.txt
BUILD_NUM=$$(cat ./build_num.txt)
APP_VERSION=0.0
TARGET_VERSION=$(APP_VERSION).$(BUILD_NUM)
TARGET_DIR=bin

API_MODULE_NAME=live-chat-api
API_OUTPUT=$(PROJECT_PATH)/$(TARGET_DIR)/$(API_MODULE_NAME)
API_MAIN=cmd/controller/main.go

WORKER_MODULE_NAME=live-chat-worker
WORKER_OUTPUT=$(PROJECT_PATH)/$(TARGET_DIR)/$(WORKER_MODULE_NAME)
WORKER_MAIN=cmd/worker/main.go

LDFLAGS=-X main.BUILD_TIME=`date -u '+%Y-%m-%d_%H:%M:%S'`
LDFLAGS+=-X main.APP_VERSION=$(TARGET_VERSION)
LDFLAGS+=-X main.GIT_HASH=`git rev-parse HEAD`
LDFLAGS+=-s -w

api: config api-build

worker: config worker-build

config:
	@if [ ! -d $(TARGET_DIR) ]; then mkdir $(TARGET_DIR); fi

api-build:
	GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(API_OUTPUT) $(PROJECT_PATH)/$(API_MAIN)
	cp $(API_OUTPUT) ./live-chat-api

worker-build:
	GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(WORKER_OUTPUT) $(PROJECT_PATH)/$(WORKER_MAIN)
	cp $(WORKER_OUTPUT) ./live-chat-worker

target-version:
	@echo "========================================"
	@echo "APP_VERSION    : $(APP_VERSION)"
	@echo "BUILD_NUM      : $(BUILD_NUM)"
	@echo "TARGET_VERSION : $(TARGET_VERSION)"
	@echo "========================================"

build_num:
	@echo $$(($$(cat $(BUILD_NUM_FILE)) + 1 )) > $(BUILD_NUM_FILE)
	@echo "BUILD_NUM      : $(BUILD_NUM)"

test:
	@go clean -testcache
	@go test ./...
