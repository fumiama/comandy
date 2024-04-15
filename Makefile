# Tested under Apple M1.
# Edit it if you use different platform.

PROJECT_NAME := comandy
BUILD_PATH := build
GOOS := android
GOARCH := arm64
BUILD_MACHINE := darwin
BUILD_ARCH := x86_64
NDK_VERSION := 26.3.11579264
TARGET_SDK := android23
TARGET_ARCH := aarch64 # optional: armv7a i686 x86_64

CGO_ENABLED := 1
GO_SRC := $(shell find . -name '*.go')
NDK_TOOLCHAIN := ~/Library/Android/sdk/ndk/$(NDK_VERSION)/toolchains/llvm/prebuilt/$(BUILD_MACHINE)-$(BUILD_ARCH)
CC := $(NDK_TOOLCHAIN)/bin/$(TARGET_ARCH)-linux-$(TARGET_SDK)-clang
TEST_OUTPUT = '$(shell cd $(BUILD_PATH) && ./test | head -c 12)'
TEST_EXPECTED = '{"code":200,'

all:
	@BUILD_PATH=$(BUILD_PATH)/aarch64 TARGET_ARCH=aarch64 GOARCH=arm64 $(MAKE) -e shared
	@BUILD_PATH=$(BUILD_PATH)/armv7a TARGET_ARCH=armv7a GOARCH=arm TARGET_SDK=androideabi23 $(MAKE) -e shared
	@BUILD_PATH=$(BUILD_PATH)/i686 TARGET_ARCH=i686 GOARCH=amd64 $(MAKE) -e shared
	@BUILD_PATH=$(BUILD_PATH)/x86_64 TARGET_ARCH=x86_64 GOARCH=386 $(MAKE) -e shared
	rm -rf $(BUILD_PATH)/*/*.h
	cd $(BUILD_PATH) && xz -z -9 -k aarch64/* armv7a/* i686/* x86_64/*
	cd $(BUILD_PATH) && zip -r -9 $(BUILD_PATH).zip aarch64 armv7a i686 x86_64
shared: $(GO_SRC) dir tidy
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
tidy:
	go mod tidy
dir:
	@if [ ! -d "$(BUILD_PATH)" ]; then mkdir $(BUILD_PATH); fi
clean:
	@if [ -d "$(BUILD_PATH)" ]; then \
		rm -rf $(BUILD_PATH)/lib$(PROJECT_NAME).*; \
		rm -rf $(BUILD_PATH)/test; \
		rm -rf $(BUILD_PATH)/*.zip; \
		find $(BUILD_PATH) -mindepth 1 -maxdepth 1 -type d -exec rm -rf {} \;; \
	fi
