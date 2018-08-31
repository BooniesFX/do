GOFMT=gofmt
GC=go build
VERSION := $(shell git describe --abbrev=4 --always --tags)
BUILD_DO_PAR = -ldflags "-X github.com/daseinio/do/common/config.VERSION=$(VERSION)"

SRC_FILES = $(shell git ls-files | grep -e .go$ | grep -v _test.go)

do: $(SRC_FILES)
	$(GC)  $(BUILD_DO_PAR) -o do main.go

all: do

do-cross: wdo ldo ddo

wdo:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GC) $(BUILD_DO_PAR) -o do-windows-amd64.exe main.go

ldo:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GC) $(BUILD_DO_PAR) -o do-linux-amd64 main.go

ddo:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GC) $(BUILD_DO_PAR) -o do-darwin-amd64 main.go

format:
	$(GOFMT) -w main.go

clean:
	rm -rf *.8 *.o *.out *.6 *exe
	rm -rf do do-*
