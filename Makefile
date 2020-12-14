all: test      \
	 benchmark \
     coverage

clean:
	go clean

format: 
	go fmt ./...

debug: build
	go test ./... -run TestTaps2BeatsWithForgetting

debugx: build
	go test ./... -run TestInterpolate

build: format
	mkdir -p bin
	go build -o bin ./...

test: build
	go test ./...

vet: build
	go vet ./...

lint: build
	golint ./...

benchmark: build
	go test -bench ./...

coverage: build
	go test -cover ./...

run: build
	./bin/taps2beats --verbose --precision 1ms --latency 7ms --quantize --interval '*' --shift ./runtime/taps.txt
	./bin/taps2beats --verbose --precision 1ms --latency 7ms --quantize --interval 1s:12s  ./runtime/taps.txt
	./bin/taps2beats --verbose --precision 1ms --latency 7ms --quantize --interval 1s:12s  --json ./runtime/taps.txt

help: build
	./bin/taps2beats 
	./bin/taps2beats --help

