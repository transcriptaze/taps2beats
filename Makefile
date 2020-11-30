all: test      \
	 benchmark \
     coverage

clean:
	go clean

format: 
	go fmt ./...

debug: build
	go test ./... -run TestTaps2Beats

build: format
	mkdir -p bin
	go build -o bin ./...

test: build
	go test ./...

run: build
	./bin/taps2beats --debug --precision 1ms --latency 7ms --quantize ./runtime/taps.txt

vet: build
	go vet ./...

lint: build
	golint ./...

benchmark: build
	go test -bench ./...

coverage: build
	go test -cover ./...


