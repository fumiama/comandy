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
TEST_INSTALL_PATH := /usr/local/lib

CGO_ENABLED := 1
GO_SRC := $(shell find . -name '*.go' | grep -v '_test.go$$')
NDK_HOME := ~/Library/Android/sdk/ndk/$(NDK_VERSION)
NDK_TOOLCHAIN := $(NDK_HOME)/toolchains/llvm/prebuilt/$(BUILD_MACHINE)-$(BUILD_ARCH)
CC := $(NDK_TOOLCHAIN)/bin/$(TARGET_ARCH)-linux-$(TARGET_SDK)-clang
TEST_EXPECTED := '{"code":200,'

all:
	@BUILD_PATH=$(BUILD_PATH)/aarch64 TARGET_ARCH=aarch64 GOARCH=arm64 $(MAKE) -e shared
	@BUILD_PATH=$(BUILD_PATH)/armv7a TARGET_ARCH=armv7a GOARCH=arm GOARM=7 TARGET_SDK=androideabi23 $(MAKE) -e shared
	@BUILD_PATH=$(BUILD_PATH)/arm TARGET_ARCH=arm GOARCH=arm GOARM=5 TARGET_SDK=androideabi23 $(MAKE) -e shared
	@BUILD_PATH=$(BUILD_PATH)/i686 TARGET_ARCH=i686 GOARCH=amd64 $(MAKE) -e shared
	@BUILD_PATH=$(BUILD_PATH)/x86_64 TARGET_ARCH=x86_64 GOARCH=386 $(MAKE) -e shared
	@BUILD_PATH=$(BUILD_PATH)/mips TARGET_ARCH=mips GOOS=linux GOARCH=mipsel GOMIPS=softfloat $(MAKE) -e shared
	@BUILD_PATH=$(BUILD_PATH)/mips64 TARGET_ARCH=mips64 GOOS=linux GOARCH=mips64el $(MAKE) -e shared
	rm -rf $(BUILD_PATH)/*/*.h
	cd $(BUILD_PATH) && gzip -9 -k aarch64/* armv7a/* i686/* x86_64/*
	find $(BUILD_PATH) -mindepth 1 -maxdepth 1 -type d -exec mv {}/lib$(PROJECT_NAME).so.gz {}_lib$(PROJECT_NAME).so.gz \;; \
	cd $(BUILD_PATH) && zip -r -9 $(BUILD_PATH).zip aarch64 armv7a i686 x86_64
shared: $(GO_SRC) dir tidy
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=$(CGO_ENABLED) NDK_TOOLCHAIN=$(NDK_TOOLCHAIN) CC=$(CC) go build -buildmode=c-shared -ldflags "-s -w" -trimpath -o $(BUILD_PATH)/lib$(PROJECT_NAME).so $(GO_SRC)
test: dir
	@GOOS=$(BUILD_MACHINE) CC=cc NDK_TOOLCHAIN="" $(MAKE) -e shared
	cc -o $(BUILD_PATH)/test $(BUILD_PATH)/test.c -l$(PROJECT_NAME) -L$(BUILD_PATH)
runtest:
	@if [ ! -f "$(BUILD_PATH)/test" ]; then \
		$(MAKE) -e test; \
	fi
	@TEST_OUTPUT=$$(cd $(BUILD_PATH) && ./test | head -c 12); \
	if [ $$TEST_OUTPUT = $(TEST_EXPECTED) ]; then \
		echo "test succeeded."; \
	else \
		echo "test failed, expected:" $(TEST_EXPECTED) "but got:" $$TEST_OUTPUT; \
		exit 1; \
	fi
tidy:
	go mod tidy
dir:
	@if [ ! -d "$(BUILD_PATH)" ]; then mkdir $(BUILD_PATH); fi
install: test
	sudo cp $(BUILD_PATH)/lib$(PROJECT_NAME).so $(TEST_INSTALL_PATH)/
	sudo ldconfig
clean:
	@if [ -d "$(BUILD_PATH)" ]; then \
		rm -rf $(BUILD_PATH)/lib$(PROJECT_NAME).*; \
		rm -rf $(BUILD_PATH)/test; \
		rm -rf $(BUILD_PATH)/*.zip; \
		rm -rf $(BUILD_PATH)/*.gz; \
		find $(BUILD_PATH) -mindepth 1 -maxdepth 1 -type d -exec rm -rf {} \;; \
	fi
