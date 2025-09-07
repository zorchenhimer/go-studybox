.PHONY: all

all: bin/script-decode bin/sbutil bin/just-stats

bin/script-decode: cmd/script-decode.go script/*.go
	go build -o $@ $<

bin/sbutil: cmd/sbutil.go rom/*.go
	go build -o $@ $<

bin/just-stats: cmd/just-stats.go script/*.go
	go build -o $@ $<
