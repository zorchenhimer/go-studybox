.PHONY: all

all: bin/script-decode bin/sbutil bin/just-stats bin/extract-imgs bin/sbx2wav

bin/script-decode: script/*.go
bin/sbutil: rom/*.go
bin/just-stats: script/*.go
bin/sbx2wav: rom/*.go audio/*.go

bin/%: cmd/%.go
	go build -o $@ $<
