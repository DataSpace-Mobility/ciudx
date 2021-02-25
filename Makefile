all:
	mkdir -p build
	go build -o build/ciudx main.go

docker:
	docker build -t dataspacein/ciudx:latest . 