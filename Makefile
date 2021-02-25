all:
	mkdir -p build
	go build -o build/rs-iudx main.go

docker:
	docker build -t dataspacein/rs-iudx:latest . 