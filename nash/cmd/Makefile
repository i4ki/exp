all: build run

build:
	go build -race pipe.go
	go build gen.go
	go build toupper.go
	go build filter.go

run:
	./pipe
