all: test      \
	 benchmark \
     coverage

clean:
	go clean

format: 
	go fmt ./...

debug: build
	go test ./... -run ExampleBeats_Interpolate

build: format
	mkdir -p bin
	go build -o bin ./...

wasm: format
	mkdir -p html/wasm
	GOOS=js GOARCH=wasm go build -o ./html/wasm ./...

build-all: format
	mkdir -p bin
	mkdir -p html/wasm
	go build -o bin ./...
	GOOS=js GOARCH=wasm go build -o ./html/wasm ./...

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

godoc: build
	godoc -http=:80     

run: build
	./bin/taps2beats --verbose --precision 1ms --latency 7ms --quantize --interval '*' --shift ./examples/taps.txt
	./bin/taps2beats --verbose --precision 1ms --latency 7ms --quantize --interval 1s:12s  ./examples/taps.txt

json: build
#	./bin/taps2beats --verbose --precision 1ms --latency 7ms --quantize --interval 1s:12s  --json ./examples/taps.txt
	./bin/taps2beats --verbose --precision 1ms --latency 7ms --quantize --interval 1s:12s  --json ./examples/taps.json

pipe: build
	cat examples/taps.txt | ./bin/taps2beats --verbose --precision 1ms --latency 7ms --quantize --interval '*' --shift

stdin: build
	./bin/taps2beats --verbose --precision 1ms --latency 7ms --quantize --interval '*'

help: build
	./bin/taps2beats --help

httpd: build
	go run cmd/httpd/httpd.go	

