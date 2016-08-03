all: build run

build:
	go build pipe.go
	go build gen.go
	go build toupper.go
	go build -race filter.go

run:
	./pipe
