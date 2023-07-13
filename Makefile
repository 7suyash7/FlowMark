.PHONY: run build

run:
	./FlowMark start

build: 
	go build -o FlowMark ./src/main.go

test:
	go test ./...

clean:
	go clean
	rm FlowBench
