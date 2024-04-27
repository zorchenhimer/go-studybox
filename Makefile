.PHONY: all

COMMANDS: bin/script-decode

all: $(COMMANDS)

bin/%: cmd/%.go script/*.go
	go build -o $@ $<
