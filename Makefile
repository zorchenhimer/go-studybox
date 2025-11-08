.PHONY: all

all: bin/script-decode bin/sbutil bin/just-stats bin/extract-imgs

bin/script-decode: script/*.go
bin/sbutil: rom/*.go
bin/just-stats: script/*.go

bin/%: cmd/%.go
	go build -o $@ $<
