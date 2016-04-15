# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

.PHONY: ged ged-cross evm all test travis-test-with-coverage xgo clean
.PHONY: ged-linux ged-linux-arm ged-linux-386 ged-linux-amd64
.PHONY: ged-darwin ged-darwin-386 ged-darwin-amd64
.PHONY: ged-windows ged-windows-386 ged-windows-amd64
.PHONY: ged-android ged-android-16 ged-android-21

GOBIN = build/bin

MODE ?= default
GO ?= latest

ged:
	build/env.sh go install -v $(shell build/flags.sh) ./cmd/ged
	@echo "Done building."
	@echo "Run \"$(GOBIN)/ged\" to launch ged."

ged-cross: ged-linux ged-darwin ged-windows ged-android
	@echo "Full cross compilation done:"
	@ls -l $(GOBIN)/ged-*

ged-linux: xgo ged-linux-arm ged-linux-386 ged-linux-amd64
	@echo "Linux cross compilation done:"
	@ls -l $(GOBIN)/ged-linux-*

ged-linux-386: xgo
	build/env.sh $(GOBIN)/xgo --go=$(GO) --buildmode=$(MODE) --dest=$(GOBIN) --targets=linux/386 -v $(shell build/flags.sh) ./cmd/ged
	@echo "Linux 386 cross compilation done:"
	@ls -l $(GOBIN)/ged-linux-* | grep 386

ged-linux-amd64: xgo
	build/env.sh $(GOBIN)/xgo --go=$(GO) --buildmode=$(MODE) --dest=$(GOBIN) --targets=linux/amd64 -v $(shell build/flags.sh) ./cmd/ged
	@echo "Linux amd64 cross compilation done:"
	@ls -l $(GOBIN)/ged-linux-* | grep amd64

ged-linux-arm: ged-linux-arm-5 ged-linux-arm-6 ged-linux-arm-7 ged-linux-arm64
	@echo "Linux ARM cross compilation done:"
	@ls -l $(GOBIN)/ged-linux-* | grep arm

ged-linux-arm-5: xgo
	build/env.sh $(GOBIN)/xgo --go=$(GO) --buildmode=$(MODE) --dest=$(GOBIN) --targets=linux/arm-5 -v $(shell build/flags.sh) ./cmd/ged
	@echo "Linux ARMv5 cross compilation done:"
	@ls -l $(GOBIN)/ged-linux-* | grep arm-5

ged-linux-arm-6: xgo
	build/env.sh $(GOBIN)/xgo --go=$(GO) --buildmode=$(MODE) --dest=$(GOBIN) --targets=linux/arm-6 -v $(shell build/flags.sh) ./cmd/ged
	@echo "Linux ARMv6 cross compilation done:"
	@ls -l $(GOBIN)/ged-linux-* | grep arm-6

ged-linux-arm-7: xgo
	build/env.sh $(GOBIN)/xgo --go=$(GO) --buildmode=$(MODE) --dest=$(GOBIN) --targets=linux/arm-7 -v $(shell build/flags.sh) ./cmd/ged
	@echo "Linux ARMv7 cross compilation done:"
	@ls -l $(GOBIN)/ged-linux-* | grep arm-7

ged-linux-arm64: xgo
	build/env.sh $(GOBIN)/xgo --go=$(GO) --buildmode=$(MODE) --dest=$(GOBIN) --targets=linux/arm64 -v $(shell build/flags.sh) ./cmd/ged
	@echo "Linux ARM64 cross compilation done:"
	@ls -l $(GOBIN)/ged-linux-* | grep arm64

ged-darwin: ged-darwin-386 ged-darwin-amd64
	@echo "Darwin cross compilation done:"
	@ls -l $(GOBIN)/ged-darwin-*

ged-darwin-386: xgo
	build/env.sh $(GOBIN)/xgo --go=$(GO) --buildmode=$(MODE) --dest=$(GOBIN) --targets=darwin/386 -v $(shell build/flags.sh) ./cmd/ged
	@echo "Darwin 386 cross compilation done:"
	@ls -l $(GOBIN)/ged-darwin-* | grep 386

ged-darwin-amd64: xgo
	build/env.sh $(GOBIN)/xgo --go=$(GO) --buildmode=$(MODE) --dest=$(GOBIN) --targets=darwin/amd64 -v $(shell build/flags.sh) ./cmd/ged
	@echo "Darwin amd64 cross compilation done:"
	@ls -l $(GOBIN)/ged-darwin-* | grep amd64

ged-windows: xgo ged-windows-386 ged-windows-amd64
	@echo "Windows cross compilation done:"
	@ls -l $(GOBIN)/ged-windows-*

ged-windows-386: xgo
	build/env.sh $(GOBIN)/xgo --go=$(GO) --buildmode=$(MODE) --dest=$(GOBIN) --targets=windows/386 -v $(shell build/flags.sh) ./cmd/ged
	@echo "Windows 386 cross compilation done:"
	@ls -l $(GOBIN)/ged-windows-* | grep 386

ged-windows-amd64: xgo
	build/env.sh $(GOBIN)/xgo --go=$(GO) --buildmode=$(MODE) --dest=$(GOBIN) --targets=windows/amd64 -v $(shell build/flags.sh) ./cmd/ged
	@echo "Windows amd64 cross compilation done:"
	@ls -l $(GOBIN)/ged-windows-* | grep amd64

ged-android: xgo
	build/env.sh $(GOBIN)/xgo --go=$(GO) --buildmode=$(MODE) --dest=$(GOBIN) --targets=android/* -v $(shell build/flags.sh) ./cmd/ged
	@echo "Android cross compilation done:"
	@ls -l $(GOBIN)/ged-android-*

ged-ios: ged-ios-arm-7 ged-ios-arm64
	@echo "iOS cross compilation done:"
	@ls -l $(GOBIN)/ged-ios-*

ged-ios-arm-7: xgo
	build/env.sh $(GOBIN)/xgo --go=$(GO) --buildmode=$(MODE) --dest=$(GOBIN) --targets=ios/arm-7 -v $(shell build/flags.sh) ./cmd/ged
	@echo "iOS ARMv7 cross compilation done:"
	@ls -l $(GOBIN)/ged-ios-* | grep arm-7

ged-ios-arm64: xgo
	build/env.sh $(GOBIN)/xgo --go=$(GO) --buildmode=$(MODE) --dest=$(GOBIN) --targets=ios-7.0/arm64 -v $(shell build/flags.sh) ./cmd/ged
	@echo "iOS ARM64 cross compilation done:"
	@ls -l $(GOBIN)/ged-ios-* | grep arm64

evm:
	build/env.sh $(GOROOT)/bin/go install -v $(shell build/flags.sh) ./cmd/evm
	@echo "Done building."
	@echo "Run \"$(GOBIN)/evm to start the evm."

all:
	build/env.sh go install -v $(shell build/flags.sh) ./...

test: all
	build/env.sh go test ./...

travis-test-with-coverage: all
	build/env.sh build/test-global-coverage.sh

xgo:
	build/env.sh go get github.com/karalabe/xgo

clean:
	rm -fr build/_workspace/pkg/ Godeps/_workspace/pkg $(GOBIN)/*
