MAINS = examples/%/main.go
EXAMPLES = $(wildcard examples/*/main.go)
BINS = $(patsubst $(MAINS),bin/%,$(EXAMPLES))

default: all

all: $(BINS)

bin/%: $(MAINS)
	@echo ">> Building $@ from $< ..."
	@mkdir -p ./bin
	@go build -o ./$@ $$(dirname $<)/*.go
