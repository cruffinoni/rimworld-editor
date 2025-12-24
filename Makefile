.PHONY: deps check build build-all build-app build-cycle build-gen-go build-gen-xml test clean

BINS_DIR := bin

APP_BIN := $(BINS_DIR)/app
CYCLE_BIN := $(BINS_DIR)/cycle
GEN_GO_BIN := $(BINS_DIR)/gen_go
GEN_XML_BIN := $(BINS_DIR)/gen_xml

check:
	go build ./...

vendor:
	go mod vendor

build: build-all

build-all: $(APP_BIN) $(CYCLE_BIN) $(GEN_GO_BIN) $(GEN_XML_BIN)

$(BINS_DIR):
	mkdir -p $(BINS_DIR)

$(APP_BIN): $(BINS_DIR)
	go build -o $(APP_BIN) ./cmd/app

$(CYCLE_BIN): $(BINS_DIR)
	go build -o $(CYCLE_BIN) ./cmd/cycle

$(GEN_GO_BIN): $(BINS_DIR)
	go build -o $(GEN_GO_BIN) ./cmd/gen_go

$(GEN_XML_BIN): $(BINS_DIR)
	go build -o $(GEN_XML_BIN) ./cmd/gen_xml

build-app: $(APP_BIN)
build-cycle: $(CYCLE_BIN)
build-gen-go: $(GEN_GO_BIN)
build-gen-xml: $(GEN_XML_BIN)

test:
	go test ./...

clean:
	rm -rf $(BINS_DIR)
