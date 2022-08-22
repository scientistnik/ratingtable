mockgen:
	mockgen -source=./internal/app/domain/ports.go -destination=./internal/app/domain/tests/mocks/ports.go

tests: mockgen
	go test ./... -test.v

build:
	go  build -o build/test -tags json1 cmd/test/test.go 

clean:
	rm -rf build