PROJECT_NAME := comandy
BUILD_PATH := build
GOOS := android
GOARCH := arm64
BUILD_MACHINE := darwin
BUILD_ARCH := x86_64
NDK_VERSION := 26.3.11579264
TARGET_SDK := android23

CGO_ENABLED := 1
GO_SRC := $(shell find . -name '*.go')
NDK_TOOLCHAIN := ~/Library/Android/sdk/ndk/$(NDK_VERSION)/toolchains/llvm/prebuilt/$(BUILD_MACHINE)-$(BUILD_ARCH)
CC := $(NDK_TOOLCHAIN)/bin/aarch64-linux-$(TARGET_SDK)-clang
TEST_OUTPUT = '$(shell cd $(BUILD_PATH) && ./test)'
TEST_EXPECTED = '{"code":500,"data":"aW52YWxpZCB1cmwgJyc="}'

all: shared

shared: $(GO_SRC) dir
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=$(CGO_ENABLED) NDK_TOOLCHAIN=$(NDK_TOOLCHAIN) CC=$(CC) go build -buildmode=c-shared -o $(BUILD_PATH)/lib$(PROJECT_NAME).so $(GO_SRC)
test: dir
	@GOOS=$(BUILD_MACHINE) CC=cc NDK_TOOLCHAIN="" $(MAKE) -e shared
	cc -o $(BUILD_PATH)/test $(BUILD_PATH)/test.c -l$(PROJECT_NAME) -L$(BUILD_PATH)
runtest: test
	@if [ $(TEST_OUTPUT) = $(TEST_EXPECTED) ]; then \
		echo "test succeeded."; \
	else \
		echo "test failed, expected:" $(TEST_EXPECTED) "but got:" $(TEST_OUTPUT); \
	fi
dir:
	@if [ ! -d "$(BUILD_PATH)" ]; then mkdir $(BUILD_PATH); fi
clean:
	@if [ -d "$(BUILD_PATH)" ]; then rm -rf $(BUILD_PATH)/lib$(PROJECT_NAME).*; fi
