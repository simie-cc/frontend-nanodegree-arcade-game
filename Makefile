all: build

build:
	GOOS=js GOARCH=wasm go build -v -o main.wasm