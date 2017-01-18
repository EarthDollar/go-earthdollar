# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

.PHONY: ged android ios ged-cross evm all test clean
.PHONY: ged-linux ged-linux-386 ged-linux-amd64 ged-linux-mips64 ged-linux-mips64le
.PHONY: ged-linux-arm ged-linux-arm-5 ged-linux-arm-6 ged-linux-arm-7 ged-linux-arm64
.PHONY: ged-darwin ged-darwin-386 ged-darwin-amd64
.PHONY: ged-windows ged-windows-386 ged-windows-amd64

GOBIN = build/bin
GO ?= latest

ged:
	build/env.sh go run build/ci.go install ./cmd/ged
	@echo "Done building."
	@echo "Run \"$(GOBIN)/ged\" to launch ged."

evm:
	build/env.sh go run build/ci.go install ./cmd/evm
	@echo "Done building."
	@echo "Run \"$(GOBIN)/evm\" to start the evm."

all:
	build/env.sh go run build/ci.go install

android:
	build/env.sh go run build/ci.go aar --local
	@echo "Done building."
	@echo "Import \"$(GOBIN)/ged.aar\" to use the library."

ios:
	build/env.sh go run build/ci.go xcode --local
	@echo "Done building."
	@echo "Import \"$(GOBIN)/Geth.framework\" to use the library."

test: all
	build/env.sh go run build/ci.go test

clean:
	rm -fr build/_workspace/pkg/ $(GOBIN)/*

# Cross Compilation Targets (xgo)

ged-cross: ged-linux ged-darwin ged-windows ged-android ged-ios
	@echo "Full cross compilation done:"
	@ls -ld $(GOBIN)/ged-*

ged-linux: ged-linux-386 ged-linux-amd64 ged-linux-arm ged-linux-mips64 ged-linux-mips64le
	@echo "Linux cross compilation done:"
	@ls -ld $(GOBIN)/ged-linux-*

ged-linux-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --dest=$(GOBIN) --targets=linux/386 -v ./cmd/ged
	@echo "Linux 386 cross compilation done:"
	@ls -ld $(GOBIN)/ged-linux-* | grep 386

ged-linux-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --dest=$(GOBIN) --targets=linux/amd64 -v ./cmd/ged
	@echo "Linux amd64 cross compilation done:"
	@ls -ld $(GOBIN)/ged-linux-* | grep amd64

ged-linux-arm: ged-linux-arm-5 ged-linux-arm-6 ged-linux-arm-7 ged-linux-arm64
	@echo "Linux ARM cross compilation done:"
	@ls -ld $(GOBIN)/ged-linux-* | grep arm

ged-linux-arm-5:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --dest=$(GOBIN) --targets=linux/arm-5 -v ./cmd/ged
	@echo "Linux ARMv5 cross compilation done:"
	@ls -ld $(GOBIN)/ged-linux-* | grep arm-5

ged-linux-arm-6:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --dest=$(GOBIN) --targets=linux/arm-6 -v ./cmd/ged
	@echo "Linux ARMv6 cross compilation done:"
	@ls -ld $(GOBIN)/ged-linux-* | grep arm-6

ged-linux-arm-7:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --dest=$(GOBIN) --targets=linux/arm-7 -v ./cmd/ged
	@echo "Linux ARMv7 cross compilation done:"
	@ls -ld $(GOBIN)/ged-linux-* | grep arm-7

ged-linux-arm64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --dest=$(GOBIN) --targets=linux/arm64 -v ./cmd/ged
	@echo "Linux ARM64 cross compilation done:"
	@ls -ld $(GOBIN)/ged-linux-* | grep arm64

ged-linux-mips64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --dest=$(GOBIN) --targets=linux/mips64 -v ./cmd/ged
	@echo "Linux MIPS64 cross compilation done:"
	@ls -ld $(GOBIN)/ged-linux-* | grep mips64

ged-linux-mips64le:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --dest=$(GOBIN) --targets=linux/mips64le -v ./cmd/ged
	@echo "Linux MIPS64le cross compilation done:"
	@ls -ld $(GOBIN)/ged-linux-* | grep mips64le

ged-darwin: ged-darwin-386 ged-darwin-amd64
	@echo "Darwin cross compilation done:"
	@ls -ld $(GOBIN)/ged-darwin-*

ged-darwin-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --dest=$(GOBIN) --targets=darwin/386 -v ./cmd/ged
	@echo "Darwin 386 cross compilation done:"
	@ls -ld $(GOBIN)/ged-darwin-* | grep 386

ged-darwin-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --dest=$(GOBIN) --targets=darwin/amd64 -v ./cmd/ged
	@echo "Darwin amd64 cross compilation done:"
	@ls -ld $(GOBIN)/ged-darwin-* | grep amd64

ged-windows: ged-windows-386 ged-windows-amd64
	@echo "Windows cross compilation done:"
	@ls -ld $(GOBIN)/ged-windows-*

ged-windows-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --dest=$(GOBIN) --targets=windows/386 -v ./cmd/ged
	@echo "Windows 386 cross compilation done:"
	@ls -ld $(GOBIN)/ged-windows-* | grep 386

ged-windows-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --dest=$(GOBIN) --targets=windows/amd64 -v ./cmd/ged
	@echo "Windows amd64 cross compilation done:"
	@ls -ld $(GOBIN)/ged-windows-* | grep amd64
